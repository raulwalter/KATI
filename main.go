package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mileusna/crontab"
)

var (
	// Store will hold all session data
	store     *sessions.CookieStore
	bootstrap map[string]*template.Template
)

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})

	bootstrap = make(map[string]*template.Template)
	bootstrap["home"] = template.Must(template.ParseFiles("templates/index.html", "templates/login.html"))
	bootstrap["report"] = template.Must(template.ParseFiles("templates/index.html", "templates/report.html"))
	bootstrap["diary"] = template.Must(template.ParseFiles("templates/index.html", "templates/diary.html"))
	bootstrap["contact"] = template.Must(template.ParseFiles("templates/index.html", "templates/contact.html"))
	bootstrap["support"] = template.Must(template.ParseFiles("templates/index.html", "templates/support.html"))
	bootstrap["maps"] = template.Must(template.ParseFiles("templates/index.html", "templates/maps.html"))
	bootstrap["api"] = template.Must(template.ParseFiles("templates/index.html", "templates/api.html"))
	bootstrap["privacy"] = template.Must(template.ParseFiles("templates/index.html", "templates/privacy.html"))
	bootstrap["faq"] = template.Must(template.ParseFiles("templates/index.html", "templates/faq.html"))
	bootstrap["dashboard"] = template.Must(template.ParseFiles("templates/index.html", "templates/dashboard.html"))
	bootstrap["contactnetwork"] = template.Must(template.ParseFiles("templates/index.html", "templates/contactnetwork.html"))
	bootstrap["datafeed"] = template.Must(template.ParseFiles("templates/index.html", "templates/datafeed.html"))
	bootstrap["form_radio"] = template.Must(template.ParseFiles("templates/index.html", "templates/form_radio.html"))
	bootstrap["form_input"] = template.Must(template.ParseFiles("templates/index.html", "templates/form_input.html"))
	bootstrap["result"] = template.Must(template.ParseFiles("templates/index.html", "templates/done.html"))

}

// Get current user
func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}

// Render the template by name
func render(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := bootstrap[name]
	if !ok {
		err := errors.New("Template not found: " + name)
		return err
	}
	err := tmpl.ExecuteTemplate(w, "layout", data)
	return err
}

// Is user authenticated
func isAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "kati-session")
	user := getUser(session)
	if auth := user.Authenticated; !auth {
		return false
	}
	return true
}

// Analyse result
func analyseResult(r *http.Request) string {

	session, _ := store.Get(r, "kati-session")
	user := getUser(session)

	message := ""
	diagnoses := getDiagnoses()

	answers, err := json.Marshal(user.getAnswersMap())

	if err != nil {
		// TODO
		fmt.Println(err)
	}

	// Check DiaryUser
	// Replace in case user is different
	diaryUser := user.Username
	if user.Questions[1].Result != "" {
		diaryUser = user.Questions[1].Result
	}

	for _, d := range diagnoses {

		fmt.Println(d.QuestionID, user.Questions[d.QuestionID].Result, d.Result, user.Questions[d.QuestionID].Result == d.Result)
		if user.Questions[d.QuestionID].Result == d.Result {
			fmt.Println(user.Questions[d.QuestionID].Result, d.Result)
			message = d.Message

			err = saveDiaryEntry(&DiaryEntry{
				UserName:  user.Username,
				DiaryUser: diaryUser,
				Answers:   postgres.Jsonb{RawMessage: answers},
				Result:    d.Status,
			})

			if err != nil {
				// TODO
				fmt.Println(err)
			}

			break
		}
	}

	return message
}

// Proccess answerx
func proccessAnswer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	questionIndex, _ := strconv.Atoi(params["question"])
	nextQuestionIndex, _ := strconv.Atoi(r.URL.Query().Get("next_question"))
	userAnswer := r.URL.Query().Get("user_answer")

	fmt.Println("Question:", questionIndex, "Answer:", userAnswer, "Next:", nextQuestionIndex)

	session, _ := store.Get(r, "kati-session")
	user := getUser(session)

	if len(user.Questions) > 0 {

		user.Questions[questionIndex].Result = userAnswer
		session.Values["user"] = user
		err := session.Save(r, w)
		if err != nil {
			// TODO
			fmt.Println(err)
		}
	}

	for i, q := range user.Questions {
		fmt.Println(i, q.Title, q.Result)
	}

	if nextQuestionIndex == -1 {
		http.Redirect(w, r, "/done", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/questionnaire"+fmt.Sprintf("/%v", nextQuestionIndex), http.StatusSeeOther)
}

// Startig page
func getDefault(w http.ResponseWriter, r *http.Request) {

	// Show dashboard in case user is authenticated
	if isAuthenticated(r) {
		dashboard(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Sisene"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "home", page)
	if err != nil {
		log.Fatalln(err)
	}

}

// report
func report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	params := mux.Vars(r)
	reportID := params["id"]

	page := &PageData{}
	page.Title = "Raport"
	page.IsAuthenticated = isAuthenticated(r)

	session, _ := store.Get(r, "kati-session")
	user := getUser(session)
	page.Diary, _ = user.getDiaryEntries()

	for _, report := range page.Diary {
		if strconv.Itoa(int(report.ID)) == reportID {
			page.CurrentReport = report
		}
	}

	err := render(w, "report", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Diary
func diary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Päevik"
	page.IsAuthenticated = isAuthenticated(r)

	session, _ := store.Get(r, "kati-session")
	user := getUser(session)
	page.Diary, _ = user.getDiaryEntries()

	err := render(w, "diary", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Contact
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kontaktid"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "contact", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Support
func support(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Tugi"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "support", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Maps
func maps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kaardirakendus"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "maps", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Api
func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kati API"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "api", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Privacy
func privacy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Privaatsustingimused"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "privacy", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Faq
func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Korduvad küsimused"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "faq", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Dashboard
func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Dashboard"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "dashboard", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Contact Network
func contactNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kontaktvõrgustik"
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "contactnetwork", page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Datafeed
func dataFeed(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	type FeedPage struct {
		Title             string
		IsAuthenticated   bool
		Country           CoVidCountry
		World             map[string]float64
		PositiveByGender  PositiveByGender
		AgeGroupsPositive map[interface{}]int
		ByCountyPositive  map[interface{}]int
	}

	cov, err := getCovCountry()
	if err != nil {
		fmt.Println(err.Error())
	}

	covMap, err := getCoVidMap()
	if err != nil {
		fmt.Println(err.Error())
	}

	var m map[string]float64
	m = make(map[string]float64)

	for _, cd := range covMap {
		_, ok := m[cd.CountryRegion]
		if !ok {
			m[cd.CountryRegion] = cd.Confirmed
		} else {
			m[cd.CountryRegion] = m[cd.CountryRegion] + cd.Confirmed
		}
	}

	// Calculate result per capita
	for i, v := range m {
		m[i] = v / 1000000
	}

	// Analyse Terviseamet
	var positiveGen PositiveByGender
	var ageGroup map[interface{}]int
	var byCounty map[interface{}]int
	ageGroup = make(map[interface{}]int)
	byCounty = make(map[interface{}]int)

	terviseametData, err := getEstoniaCovidTested()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, tp := range terviseametData {
		if tp.Gender == "M" && tp.ResultValue == "P" {
			positiveGen.Men = positiveGen.Men + 1
		}
		if tp.Gender == "N" && tp.ResultValue == "P" {
			positiveGen.Women = positiveGen.Women + 1
		}

		//fmt.Println(tp.AgeGroup, reflect.TypeOf(tp.AgeGroup).Kind(), reflect.TypeOf(tp.AgeGroup).Kind() == reflect.String)

		if tp.AgeGroup != nil {
			if tp.ResultValue == "P" && reflect.TypeOf(tp.AgeGroup).Kind() == reflect.String {
				_, ok := ageGroup[tp.AgeGroup]
				if !ok {
					ageGroup[tp.AgeGroup] = 1
				} else {
					ageGroup[tp.AgeGroup] = ageGroup[tp.AgeGroup] + 1
				}
			}
		}
		if tp.County != nil {
			if tp.ResultValue == "P" && reflect.TypeOf(tp.County).Kind() == reflect.String {

				if tp.County == "Harju maakond" {
					tp.County = "Harjumaa"
				}
				if tp.County == "Hiiu maakond" {
					tp.County = "Hiiumaa"
				}
				if tp.County == "Ida-Viru maakond" {
					tp.County = "Ida-Virumaa"
				}
				if tp.County == "Järva maakond" {
					tp.County = "Järvamaa"
				}
				if tp.County == "Jõgeva maakond" {
					tp.County = "Jõgevamaa"
				}
				if tp.County == "Lääne maakond" {
					tp.County = "Läänemaa"
				}
				if tp.County == "Lääne-Viru maakond" {
					tp.County = "Lääne-Virumaa"
				}
				if tp.County == "Pärnu maakond" {
					tp.County = "Pärnumaa"
				}
				if tp.County == "Põlva maakond" {
					tp.County = "Põlvamaa"
				}
				if tp.County == "Rapla maakond" {
					tp.County = "Raplamaa"
				}
				if tp.County == "Saare maakond" {
					tp.County = "Saaremaa"
				}
				if tp.County == "Tartu maakond" {
					tp.County = "Tartumaa"
				}
				if tp.County == "Valga maakond" {
					tp.County = "Valgamaa"
				}
				if tp.County == "Viljandi maakond" {
					tp.County = "Viljandimaa"
				}
				if tp.County == "Võru maakond" {
					tp.County = "Võrumaa"
				}

				_, ok := byCounty[tp.County]
				if !ok {
					byCounty[tp.County] = 1
				} else {
					byCounty[tp.County] = byCounty[tp.County] + 1
				}
			}
		}

	}

	feed := &FeedPage{}
	feed.Title = "Andmestik"
	feed.IsAuthenticated = isAuthenticated(r)
	feed.Country = cov
	feed.World = m
	feed.PositiveByGender = positiveGen
	feed.AgeGroupsPositive = ageGroup
	feed.ByCountyPositive = byCounty

	err = render(w, "datafeed", feed)
	if err != nil {
		log.Fatalln(err)
	}
}

// Questionnaire
func questionnaire(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	params := mux.Vars(r)
	questionIndex, err := strconv.Atoi(params["id"])
	if err != nil {
		questionIndex = 0
	}

	page := &PageData{}
	page.Title = "Küsimustik"
	page.IsAuthenticated = isAuthenticated(r)
	page.CurrentQuestionIndex = questionIndex

	session, _ := store.Get(r, "kati-session")
	user := getUser(session)

	page.CurrentQuestion = user.Questions[questionIndex]
	questionType := user.Questions[questionIndex].Type

	err = render(w, "form_"+questionType, page)
	if err != nil {
		log.Fatalln(err)
	}
}

// Show result after questions
func resultPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// TODO: Store result to database

	page := &PageData{}
	page.Title = "Test tehtud"
	page.DiagnoseHTML = template.HTML(analyseResult(r))
	page.IsAuthenticated = isAuthenticated(r)

	err := render(w, "result", page)
	if err != nil {
		log.Fatalln(err)
	}

	// Reset questions
	session, _ := store.Get(r, "kati-session")
	user := getUser(session)
	user.Questions = getQuestions()
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
	}

}

// Login
func loginPost(w http.ResponseWriter, r *http.Request) {

	// TODO: TARA login
	// Replace with TARA

	// Random login
	randomNID := []string{"37701130004", "46205124327", "36205129497", "35401020033", "45401020132", "45401020415", "50111110030", "50111111800", "60111113719", "60111114106"}
	userNID := randomNID[rand.Intn(len(randomNID)-1)]

	session, _ := store.Get(r, "kati-session")

	user := &User{
		Username:      userNID,
		Authenticated: true,
		Questions:     getQuestions(),
	}

	session.Values["user"] = user

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Logout
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kati-session")
	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Session handler
var sessionHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "kati-session")

		authRequired := []string{"/questionnaire", "/answer", "/done", "/diary", "/report", "/contactnetwork"}
		requestPath := r.URL.Path

		for _, value := range authRequired {

			if value == requestPath {

				user := getUser(session)

				if auth := user.Authenticated; !auth {
					session.AddFlash("You don't have access!")
					err := session.Save(r, w)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					http.Redirect(w, r, "/forbidden", http.StatusFound)
					return
				}

			}
		}

		next.ServeHTTP(w, r)
	})

}

func main() {

	// Download data on service initialization
	go downloadFromTerviseamet()
	go downloadCovidData()

	// Setup crons
	ctab := crontab.New()
	// Run terviseamet cron once a day
	err := ctab.AddJob("* 0 * * *", downloadFromTerviseamet)
	if err != nil {
		fmt.Println("Cron error:", err)
	}
	// Run Covid cron after every 3 hour
	err = ctab.AddJob("* */3 * * *", downloadCovidData)
	if err != nil {
		fmt.Println("Cron error:", err)
	}

	router := mux.NewRouter()

	serveStatic(router)
	router.HandleFunc("/", getDefault).Methods("GET")
	router.HandleFunc("/datafeed", dataFeed).Methods("GET")
	router.HandleFunc("/done", resultPage).Methods("GET")
	router.HandleFunc("/api", api).Methods("GET")
	router.HandleFunc("/maps", maps).Methods("GET")
	router.HandleFunc("/privacy", privacy).Methods("GET")
	router.HandleFunc("/faq", faq).Methods("GET")
	router.HandleFunc("/diary", diary).Methods("GET")
	router.HandleFunc("/report/{id:[0-9]+}", report).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	router.HandleFunc("/support", support).Methods("GET")
	router.HandleFunc("/contactnetwork", contactNetwork).Methods("GET")
	router.HandleFunc("/questionnaire", questionnaire).Methods("GET")
	router.HandleFunc("/questionnaire/{id:[0-9]+}", questionnaire).Methods("GET")
	router.HandleFunc("/answer/{question:[0-9]+}", proccessAnswer).Methods("GET")
	router.HandleFunc("/login", loginPost).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("GET")

	router.Use(sessionHandler)

	http.ListenAndServe(":8888", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))

}

// Serve static files
func serveStatic(router *mux.Router) {
	staticPaths := map[string]string{
		"/css/":     "/static/css/",
		"/fonts/":   "/static/fonts/",
		"/images/":  "/static/images/",
		"/icons/":   "/static/icons/",
		"/scripts/": "/static/scripts/",
	}
	for pathPrefix, pathValue := range staticPaths {
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir("."+pathValue))))
	}
}

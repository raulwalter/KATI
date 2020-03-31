package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	// Store will hold all session data
	store *sessions.CookieStore
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
	for _, d := range diagnoses {

		fmt.Println(d.QuestionID, user.Questions[d.QuestionID].Result, d.Result, user.Questions[d.QuestionID].Result == d.Result)
		if user.Questions[d.QuestionID].Result == d.Result {
			fmt.Println(user.Questions[d.QuestionID].Result, d.Result)
			message = d.Message
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

	t, _ := template.ParseFiles("templates/index.html", "templates/login.html")
	t.ExecuteTemplate(w, "layout", page)

}

// Diary
func diary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Päevik"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/diary.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Contact
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kontaktid"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/contact.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Support
func support(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Tugi"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/support.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Maps
func maps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kaardirakendus"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/maps.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Api
func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kati API"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/api.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Privacy
func privacy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Privaatsustingimused"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/privacy.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Faq
func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Korduvad küsimused"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/faq.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Dashboard
func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Dashboard"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/dashboard.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Contact Network
func contactNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	page := &PageData{}
	page.Title = "Kontaktvõrgustik"
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/contactnetwork.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Datafeed
func dataFeed(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	type FeedPage struct {
		Title           string
		IsAuthenticated bool
		Country         CoVidCountry
		World           map[string]float64
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

	feed := &FeedPage{}
	feed.Title = "Andmestik"
	feed.IsAuthenticated = isAuthenticated(r)
	feed.Country = cov
	feed.World = m

	t, _ := template.ParseFiles("templates/index.html", "templates/datafeed.html")
	t.ExecuteTemplate(w, "layout", feed)
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

	t, _ := template.ParseFiles("templates/index.html", "templates/form_"+questionType+".html")
	t.ExecuteTemplate(w, "layout", page)
}

// Show result after questions
func resultPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// TODO: Store result to database

	page := &PageData{}
	page.Title = "Test tehtud"
	page.DiagnoseHTML = template.HTML(analyseResult(r))
	page.IsAuthenticated = isAuthenticated(r)

	t, _ := template.ParseFiles("templates/index.html", "templates/done.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Login
func loginPost(w http.ResponseWriter, r *http.Request) {

	//TODO: TARA login

	session, _ := store.Get(r, "kati-session")

	user := &User{
		Username:      "37701130004",
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

		authRequired := []string{"/questionnaire", "/answer", "/done", "/diary", "/contactnetwork"}
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

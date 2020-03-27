package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var currentQuestionID int
var questions []DiagQuestion
var page PageData

var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))

// PageData ...
type PageData struct {
	Title           string
	IsAuthenticated bool
	CurrentQuestion DiagQuestion
	DiagnoseHTML    template.HTML
}

// DiagQuestion ...
type DiagQuestion struct {
	Question    string   `json:"question"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Answers     []Answer `json:"answers"`
	Result      interface{}
}

// Answer ...
type Answer struct {
	Caption string `json:"caption"`
	Value   int    `json:"value"`
	Next    int    `json:"next"`
}

// Diagnose ...
type Diagnose struct {
	QuestionID int    `json:"question"`
	Result     int    `json:"result"`
	Message    string `json:"diagnose"`
}

// CovCountry ...
type CovCountry struct {
	Confirmed struct {
		Value  int    `json:"value"`
		Detail string `json:"detail"`
	} `json:"confirmed"`
	Recovered struct {
		Value  int    `json:"value"`
		Detail string `json:"detail"`
	} `json:"recovered"`
	Deaths struct {
		Value  int    `json:"value"`
		Detail string `json:"detail"`
	} `json:"deaths"`
	LastUpdate time.Time `json:"lastUpdate"`
}

// CoVidMap ...
type CoVidMap struct {
	ProvinceState interface{} `json:"provinceState"`
	CountryRegion string      `json:"countryRegion"`
	LastUpdate    int64       `json:"lastUpdate"`
	Lat           float64     `json:"lat"`
	Long          float64     `json:"long"`
	Confirmed     float64     `json:"confirmed"`
	Recovered     int         `json:"recovered"`
	Deaths        int         `json:"deaths"`
	Active        int         `json:"active"`
	Admin2        interface{} `json:"admin2"`
	Fips          interface{} `json:"fips"`
	CombinedKey   string      `json:"combinedKey"`
	IncidentRate  interface{} `json:"incidentRate"`
	PeopleTested  interface{} `json:"peopleTested"`
	Iso2          string      `json:"iso2,omitempty"`
	Iso3          string      `json:"iso3,omitempty"`
}

// Read questions from file
func getQuestions() []DiagQuestion {
	questions := &[]DiagQuestion{}
	file, _ := ioutil.ReadFile("data/data.json")
	_ = json.Unmarshal([]byte(file), &questions)
	return *questions
}

// Get population
func getPopulation() map[string]int {

	var m map[string]int
	m = make(map[string]int)

	type Population struct {
		Country    string `json:"country"`
		Population string `json:"population"`
	}

	population := &[]Population{}
	file, _ := ioutil.ReadFile("data/population.json")
	_ = json.Unmarshal([]byte(file), &population)

	for _, p := range *population {
		i, err := strconv.Atoi(p.Population)
		fmt.Println(i)
		if err == nil {
			m[p.Country] = i
		}
	}

	return m
}

// Get diagnose
func getDiagnoses() []Diagnose {
	diagnose := &[]Diagnose{}
	file, _ := ioutil.ReadFile("data/eval.json")
	_ = json.Unmarshal([]byte(file), &diagnose)
	return *diagnose
}

// Analyse result
func analyseResult() string {
	message := ""
	diagnoses := getDiagnoses()
	for _, d := range diagnoses {
		if questions[d.QuestionID].Result == d.Result {
			message = d.Message
		}
	}
	return message
}

// Get next question ID
func getNextQuestionID(answerID int) int {
	answers := questions[currentQuestionID].Answers
	for _, a := range answers {
		if a.Value == answerID {
			return a.Next
		}
	}
	return 0
}

// Get CoVid Estonia data
func getCovCountry() (CovCountry, error) {
	var cov CovCountry
	url := "https://covid19.mathdro.id/api/countries/Estonia"
	res, err := http.Get(url)
	if err != nil {
		return cov, err
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return cov, err
	}

	json.Unmarshal(body, &cov)
	return cov, nil
}

// Get all data for map use
func getCoVidMap() ([]CoVidMap, error) {
	var cov []CoVidMap
	url := "https://covid19.mathdro.id/api/confirmed"
	res, err := http.Get(url)
	if err != nil {
		return cov, err
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return cov, err
	}

	json.Unmarshal(body, &cov)
	return cov, nil
}

// Proccess answerx
func proccessAnswer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	answerID, _ := strconv.Atoi(params["id"])

	questions[currentQuestionID].Result = answerID
	currentQuestionID = getNextQuestionID(answerID)

	if currentQuestionID == -1 {
		currentQuestionID = 0
		http.Redirect(w, r, "/done", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/questionnaire"+fmt.Sprintf("/%v", currentQuestionID), http.StatusSeeOther)
}

// Startig page
func getDefault(w http.ResponseWriter, r *http.Request) {

	// Show dashboard in case user is authenticated
	if page.IsAuthenticated {
		dashboard(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	page.Title = "Sisene"

	t, _ := template.ParseFiles("templates/index.html", "templates/login.html")
	t.ExecuteTemplate(w, "layout", page)

}

// Logout
func logout(w http.ResponseWriter, r *http.Request) {

	page.IsAuthenticated = false

	session, _ := store.Get(r, "kati-session")
	session.Values["user"] = ""
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Diary
func diary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Päevik"
	t, _ := template.ParseFiles("templates/index.html", "templates/diary.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Contact
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Kontaktid"
	t, _ := template.ParseFiles("templates/index.html", "templates/contact.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Support
func support(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Tugi"
	t, _ := template.ParseFiles("templates/index.html", "templates/support.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Maps
func maps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Kaardirakendus"
	t, _ := template.ParseFiles("templates/index.html", "templates/maps.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Api
func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Kati API"
	t, _ := template.ParseFiles("templates/index.html", "templates/api.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Privacy
func privacy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Privaatsustingimused"
	t, _ := template.ParseFiles("templates/index.html", "templates/privacy.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Faq
func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Korduvad küsimused"
	t, _ := template.ParseFiles("templates/index.html", "templates/faq.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Dashboard
func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	page.Title = "Dashboard"
	t, _ := template.ParseFiles("templates/index.html", "templates/dashboard.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Datafeed
func dataFeed(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	type FeedPage struct {
		Title           string
		IsAuthenticated bool
		Country         CovCountry
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

	var feed FeedPage
	feed.Title = "Andmestik"
	feed.IsAuthenticated = page.IsAuthenticated
	feed.Country = cov
	feed.World = m

	t, _ := template.ParseFiles("templates/index.html", "templates/datafeed.html")
	t.ExecuteTemplate(w, "layout", feed)
}

// Questionnaire
func questionnaire(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	var err error

	page.Title = "Küsimustik"

	params := mux.Vars(r)
	currentQuestionID, err = strconv.Atoi(params["id"])
	if err != nil {
		currentQuestionID = 0
	}

	page.CurrentQuestion = questions[currentQuestionID]
	questionType := questions[currentQuestionID].Type

	t, _ := template.ParseFiles("templates/index.html", "templates/form_"+questionType+".html")
	t.ExecuteTemplate(w, "layout", page)
}

// Last page after questions
func lastPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// TODO: Store result to database

	page.Title = "Test tehtud"
	page.DiagnoseHTML = template.HTML(analyseResult())

	questions = getQuestions()

	t, _ := template.ParseFiles("templates/index.html", "templates/done.html")
	t.ExecuteTemplate(w, "layout", page)
}

// Login
func loginPost(w http.ResponseWriter, r *http.Request) {

	//TODO: TARA login

	questions = getQuestions()

	session, _ := store.Get(r, "kati-session")

	session.Values["user"] = "37701130004"
	session.Values["authenticated"] = true

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.IsAuthenticated = true
	http.Redirect(w, r, "./questionnaire", http.StatusMovedPermanently)
}

// Session handler
var sessionHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "kati-session")

		authRequired := []string{"/questionnaire", "/answer", "/done", "/diary"}
		requestPath := r.URL.Path

		for _, value := range authRequired {
			if value == requestPath {
				if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})

}

func main() {

	page = PageData{}

	router := mux.NewRouter()

	serveStatic(router)

	router.HandleFunc("/", getDefault).Methods("GET")
	router.HandleFunc("/datafeed", dataFeed).Methods("GET")
	router.HandleFunc("/done", lastPage).Methods("GET")
	router.HandleFunc("/api", api).Methods("GET")
	router.HandleFunc("/maps", maps).Methods("GET")
	router.HandleFunc("/privacy", privacy).Methods("GET")
	router.HandleFunc("/faq", faq).Methods("GET")
	router.HandleFunc("/diary", diary).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	router.HandleFunc("/support", support).Methods("GET")
	router.HandleFunc("/questionnaire", questionnaire).Methods("GET")
	router.HandleFunc("/questionnaire/{id:[0-9]+}", questionnaire).Methods("GET")
	router.HandleFunc("/answer/{id:[0-9]+}", proccessAnswer).Methods("GET")
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

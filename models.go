package main

import (
	"encoding/json"
	"html/template"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// User holds a users account information
type User struct {
	Username      string
	Authenticated bool
	Questions     []Question
}

// PageData holds a generic page data
type PageData struct {
	Title                string
	IsAuthenticated      bool
	CurrentQuestionIndex int
	CurrentQuestion      Question
	DiagnoseHTML         template.HTML
	Diary                []DiaryEntry
}

// Question holds questions to users
type Question struct {
	ID          string   `json:"id"`
	Title       string   `json:"question"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Answers     []Answer `json:"answers"`
	Result      string
}

// Answer holds the question possible answers
type Answer struct {
	Caption string `json:"caption"`
	Value   string `json:"value"`
	Next    int    `json:"next"`
}

// Diagnose hold quetions result messages
type Diagnose struct {
	QuestionID int    `json:"question"`
	Result     string `json:"result"`
	Message    string `json:"diagnose"`
	Status     string `json:"status"`
}

// DiaryEntry holds single user diary entry
type DiaryEntry struct {
	ID                   uint `gorm:"primary_key"`
	CreatedAt            time.Time
	UserName             string
	QuestionnaireVersion string
	Answers              postgres.Jsonb
	AnswersMap           map[string]interface{} `gorm:"-" json:"-"`
	Result               string
}

// CoVidCountry holds CoVID-19 information from country
type CoVidCountry struct {
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

// CoVidMap holds COVID-19 information
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

func (u *User) getAnswersMap() map[string]interface{} {
	m := make(map[string]interface{})

	for _, q := range u.Questions {
		m[q.ID] = q.Result
	}

	return m
}

func (d *DiaryEntry) GetAnswerFor(id string) interface{} {

	if d.AnswersMap == nil {
		d.AnswersMap = make(map[string]interface{})
		json.Unmarshal(d.Answers.RawMessage, &d.AnswersMap)
	}

	if v, ok := d.AnswersMap[id]; ok {
		if id == "2" {
			switch v {
			case "1":
				return "Yes"
			case "0":
				return "No"
			}
		}
	}

	return "-"
}

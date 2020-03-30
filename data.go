package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Read questions from file
func getQuestions() []Question {
	questions := &[]Question{}
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

// Get CoVid Estonia data
func getCovCountry() (CoVidCountry, error) {
	var cov CoVidCountry
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

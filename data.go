package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	file, _ := ioutil.ReadFile("data/Estonia.json")
	_ = json.Unmarshal([]byte(file), &cov)
	return cov, nil
}

// Get all data for map use
func getCoVidMap() ([]CoVidMap, error) {
	var cov []CoVidMap
	file, _ := ioutil.ReadFile("data/confirmed_world.json")
	_ = json.Unmarshal([]byte(file), &cov)
	return cov, nil
}

// Get data of tested persons in Estonia
func getEstoniaCovidTested() ([]TestedPerson, error) {
	var testedPersons []TestedPerson
	file, _ := ioutil.ReadFile("data/opendata_covid19_test_results.json")
	_ = json.Unmarshal([]byte(file), &testedPersons)
	return testedPersons, nil
}

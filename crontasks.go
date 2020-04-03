package main

import "fmt"

// Download open data from Terviseamet
func downloadFromTerviseamet() {
	fmt.Println("Downloading from opendata.digilugu.ee")
	url := "https://opendata.digilugu.ee/opendata_covid19_test_results.json"
	if err := DownloadFile("data/opendata_covid19_test_results.json", url); err != nil {
		fmt.Println(err)
	}
}

// Download data from covid19.mathdro.id
// will retrive two files where one is country based json
func downloadCovidData() {
	fmt.Println("Downloading from covid19.mathdro.id")
	url := "https://covid19.mathdro.id/api/countries/Estonia"
	if err := DownloadFile("data/Estonia.json", url); err != nil {
		fmt.Println(err)
	}
	url = "https://covid19.mathdro.id/api/confirmed"
	if err := DownloadFile("data/confirmed_world.json", url); err != nil {
		fmt.Println(err)
	}
}

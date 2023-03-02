//This file holds refactored functionalities to send GET requests
//It has the general GET request function so other functions only need to format the value: "GetRequest"
//And more specific requests like "GetCountry" to find the Country
package coronastats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//GetCountry use this RESTservice: "https://restcountries.eu"
//to find the Official Name or Alpha3 code, this is decided by the @form parameter
//@return - the first Country object identified by the search
func GetCountry(w http.ResponseWriter, r *http.Request) SearchCountry {

	request := "https://restcountries.eu/rest/v2/name/" + GetKeyword(w, r, "Name") //Formating request

	byteValue := GetRequest(w, request)     //calls GetRequest to "GET" the data in byte from restcountries
	var countryList []SearchCountry         //Declare an array of SearchCountries, refer to structs.go for type structure
	json.Unmarshal(byteValue, &countryList) //Format the byte slice to Country objects stored in countryList
	return (countryList[0])                 //Return the first country, and hope restcountries wasn't being funky
}

//GetRequest is a function that does the HTTP get request and returns the output in []byte
//And then the function that wanted the data can format it into its own structure
func GetRequest(w http.ResponseWriter, httpRequest string) []byte {
	response, err := http.Get(httpRequest)

	if err != nil {
		errorReport := fmt.Errorf("error encountered in external API request: %v", err.Error())
		http.Error(w, errorReport.Error(), http.StatusNotFound)
	}

	byteValue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorReport := fmt.Errorf("error when reading the response: %v", err.Error())
		http.Error(w, errorReport.Error(), http.StatusBadRequest)
	}
	return byteValue
}

func GetRequestWebhook(httpRequest string) []byte {
	response, err := http.Get(httpRequest)

	if err != nil {
		errorReport := fmt.Errorf("error encountered in external API request: %v", err.Error())
		log.Fatal(errorReport.Error())
	}

	byteValue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorReport := fmt.Errorf("error when reading the response: %v", err.Error())
		log.Fatal(errorReport.Error())
	}
	return byteValue
}

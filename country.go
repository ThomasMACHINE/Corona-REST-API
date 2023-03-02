//File for the first endpoint

package coronastats

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

func GetCPC(w http.ResponseWriter, r *http.Request) {
	var country CountryCaseDisplay
	if r.FormValue("scope") == "" {
		country = GetCPC_No_Scope(w, r)
	} else {
		country = GetCPC_With_Scope(w, r)
	}

	data, err := json.Marshal(country)
	if err != nil {
		http.Error(w, "Something went wrong during marshalling data: "+err.Error(), http.StatusBadGateway)
	}
	//Setting content type to JSON
	http.Header.Add(w.Header(), "content-type", "application/json")
	fmt.Fprintln(w, string(data))
}
func GetCPC_No_Scope(w http.ResponseWriter, r *http.Request) CountryCaseDisplay {
	countryName := GetCountry(w, r) //Gets a country object from the GetCountry function in getRequestAPI.go
	//Uses GetRequest from networkingAPI, that returns the data in []Byte type
	confirmedByte := GetRequest(w, "https://covid-api.mmediagroup.fr/v1/history?country="+countryName.Name+"&status=Confirmed") //Country object holds a field, Name which is official name, ex: "France"
	//Do it again for recovered because I could not figure out how to send multiple Queries with status
	recoveredByte := GetRequest(w, "https://covid-api.mmediagroup.fr/v1/history?country="+countryName.Name+"&status=Recovered")
	//Initialize objects to format data into
	var confirmedData CovidCaseData
	var recoveredData CovidCaseData
	//Formating data into struct objects
	json.Unmarshal(confirmedByte, &confirmedData)
	json.Unmarshal(recoveredByte, &recoveredData)

	var displayData CountryCaseDisplay //Initialize displayObjectt
	displayData.Init(confirmedData)    //Constructor, copies name, continent and population

	startDate := "2020-01-22"  //Hardcoded startDate as that is the one I've seen from all reviewed country data
	endDate := GetDay(0, 0, 0) //Gets todays date, check Inputtformatting.go
	//Checking if recovered is zero, which for MOST of the cases means there is not reported data, however there are some countries with 0 cases
	//But in terms of data, it is still accurate as their reported data still is 0 in total
	if recoveredData.Key.Dates[endDate]-recoveredData.Key.Dates[startDate] <= 0 {
		endDate = GetDay(0, 0, -1) //Gets yesterdays date, check Inputtformatting.go
	}
	//Setting confirmed, recovered and scope field
	displayData.Confirmed = confirmedData.Key.Dates[endDate] - confirmedData.Key.Dates[startDate] //Change = endAmount - startAmount
	fmt.Println(confirmedData.Key.Dates[endDate], confirmedData.Key.Dates[startDate])
	displayData.Recovered = recoveredData.Key.Dates[endDate] - recoveredData.Key.Dates[startDate] //Change = endAmount - startAmount
	displayData.Scope = "total"

	//calculate population percentage
	var pop_perc float64 = float64(displayData.Confirmed) / float64(displayData.Population) //percentage = affected / population
	pop_perc = math.Round(pop_perc*100) / 100                                               //x = 0.5032 -> round(50.32) = 50, divide by 100 => 0.50
	displayData.Population_percentage = pop_perc

	return displayData //Return CountryCaseDisplay object
}
func GetCPC_With_Scope(w http.ResponseWriter, r *http.Request) CountryCaseDisplay {
	countryName := GetCountry(w, r) //Gets a country object from the GetCountry function in getRequestAPI.go
	//Uses GetRequest from networkingAPI, that returns the data in []Byte type
	confirmedByte := GetRequest(w, "https://covid-api.mmediagroup.fr/v1/history?country="+countryName.Name+"&status=Confirmed") //Country object holds a field, Name which is official name, ex: "France"
	//Do it again for recovered because I could not figure out how to send multiple Queries with status
	recoveredByte := GetRequest(w, "https://covid-api.mmediagroup.fr/v1/history?country="+countryName.Name+"&status=Recovered")
	//Initialize objects to format data into
	var confirmedData CovidCaseData
	var recoveredData CovidCaseData
	//Formating data into struct objects
	json.Unmarshal(confirmedByte, &confirmedData)
	json.Unmarshal(recoveredByte, &recoveredData)

	var displayData CountryCaseDisplay //Initialize displayObjectt
	displayData.Init(confirmedData)    //Constructor, copies name, continent and population

	startDate, endDate := GetKeyword(w, r, "StartDate"), GetKeyword(w, r, "EndDate") //Finds the scope using func - GetKeyword from inputFormatting.go

	//Checking if recovered for endDate is zero, if it is that could mean it has not been published yet
	//If not, it doesnt matter as all previous dates also are zero
	if recoveredData.Key.Dates[endDate] == 0 {
		endDate = GetDay(0, 0, -1) //Gets yesterdays date, check Inputtformatting.go
	}

	//Finds the scope using func: GetKeyword from inputFormatting.go
	//setting the Confirmed, Recovered, Scope
	displayData.Confirmed = confirmedData.Key.Dates[endDate] - confirmedData.Key.Dates[startDate] //Change = endAmount - startAmount
	displayData.Recovered = recoveredData.Key.Dates[endDate] - recoveredData.Key.Dates[startDate] //Change = endAmount - startAmount
	displayData.Scope = startDate + "-" + endDate
	//calculate population percentage and set the field
	var pop_perc float64 = float64(displayData.Confirmed) / float64(displayData.Population) //percentage = affected / population
	pop_perc = math.Round(pop_perc*100) / 100                                               //x = 0.5032 -> round(50.32) = 50, divide by 100 => 0.50
	displayData.Population_percentage = pop_perc

	return displayData //Return CountryCaseDisplay object
}

// ======= Webhook COMPATIBLE WOOOO

func GetCPC_Webhook(webhook WebhookRegistration) CountryCaseDisplay {
	//Uses GetRequest from networkingAPI, that returns the data in []Byte type
	confirmedByte := GetRequestWebhook("https://covid-api.mmediagroup.fr/v1/history?country=" + webhook.Country + "&status=Confirmed") //Country object holds a field, Name which is official name, ex: "France"
	//Do it again for recovered because I could not figure out how to send multiple Queries with status
	recoveredByte := GetRequestWebhook("https://covid-api.mmediagroup.fr/v1/history?country=" + webhook.Country + "&status=Recovered")
	//Initialize objects to format data into
	var confirmedData CovidCaseData
	var recoveredData CovidCaseData
	//Formating data into struct objects
	json.Unmarshal(confirmedByte, &confirmedData)
	json.Unmarshal(recoveredByte, &recoveredData)

	var displayData CountryCaseDisplay //Initialize displayObjectt
	displayData.Init(confirmedData)    //Constructor, copies name, continent and population

	startDate := "2020-01-22"
	endDate := GetDay(0, 0, 0) //Gets todays date, check Inputtformatting.go
	//Checking if recovered is zero, which for MOST of the cases means there is not reported data, however there are some countries with 0 cases
	//But in terms of data, it is still accurate as their reported data still is 0 in total
	if recoveredData.Key.Dates[endDate]-recoveredData.Key.Dates[startDate] <= 0 {
		endDate = GetDay(0, 0, -1) //Gets yesterdays date, check Inputtformatting.go
	}
	displayData.Confirmed = confirmedData.Key.Dates[endDate] - confirmedData.Key.Dates[startDate] //Change = endAmount - startAmount
	fmt.Println(confirmedData.Key.Dates[endDate], confirmedData.Key.Dates[startDate])
	displayData.Recovered = recoveredData.Key.Dates[endDate] - recoveredData.Key.Dates[startDate] //Change = endAmount - startAmount
	displayData.Scope = endDate

	//Finds the scope using func - GetKeyword from inputFormatting.go
	//setting the Confirmed, Recovered, Scope

	//calculate population percentage
	var pop_perc float64 = float64(displayData.Confirmed) / float64(displayData.Population) //percentage = affected / population
	pop_perc = math.Round(pop_perc*100) / 100                                               //x = 0.5032 -> round(50.32) = 50, divide by 100 => 0.50
	displayData.Population_percentage = pop_perc

	return displayData //Return CountryCaseDisplay object
}

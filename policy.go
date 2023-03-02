package coronastats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetPolicy(w http.ResponseWriter, r *http.Request) {
	var displayData StringencyDisplay
	//Test wether the url has scope key or not and call the correct function
	if r.FormValue("scope") == "" {
		displayData = GetPolicyNoScope(w, r)
	} else {
		displayData = GetPolicyWithScope(w, r)
	}
	//Decode data
	policyData, err := json.Marshal(displayData)
	if err != nil {
		http.Error(w, "Data was sent to server, but there was a problem while encoding the data: "+err.Error(), http.StatusBadRequest)
		return
	}
	//Set content type to JSON
	http.Header.Add(w.Header(), "content-type", "application/json")
	fmt.Fprintln(w, string(policyData))
}

func GetPolicyNoScope(w http.ResponseWriter, r *http.Request) StringencyDisplay {
	countryName := GetCountry(w, r) //get a SearchCountry object from GetCountry function, read GetKeyword from inputFormatting.go
	today := GetDay(0, 0, 0)        //Get todays date - inputFormatting.go
	byteData := GetRequest(w, "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/"+countryName.AlphaCode+"/"+today)

	var data StringencyData
	err := json.Unmarshal(byteData, &data)
	if err != nil {
		http.Error(w, "Could not format the Stringency Data", http.StatusBadRequest)
	}

	//Testing the data
	//If the data has not been posted, the Service provides a messsage field, where as if it has data it doesnt
	if data.Key.Message != "" { //Check if the string is not nil value
		twoDaysAgo := GetDay(0, 0, -2)
		byteData := GetRequest(w, "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/"+countryName.AlphaCode+"/"+twoDaysAgo)
		data.Key.Message = "" //Reset the Messsage field as correct data wont override it back to "" when unmarshalling
		json.Unmarshal(byteData, &data)
		if data.Key.Message != "" {
			http.Error(w, "Database has not updated for last 3 days, try this country later", http.StatusBadGateway)
		}
	}
	//Initializing StringencyDisplay object, read structs.go
	var displayData StringencyDisplay
	//Filling fields
	displayData.Country = countryName.Name //Using the official name
	displayData.Scope = "total"
	displayData.Stringency = data.Key.Stringency
	displayData.Trend = 0
	//Format data
	return displayData
}

func GetPolicyWithScope(w http.ResponseWriter, r *http.Request) StringencyDisplay {
	countryName := GetCountry(w, r) //get the name from the url, read GetKeyword from inputFormatting.go
	//get dates from the url, read GetKeyword from inputFormatting.go
	startDate := GetKeyword(w, r, "StartDate")
	endDate := GetKeyword(w, r, "EndDate")
	//Calling GetRequest in getRequestAPI to do the heavy work and get the data in []Byte type
	startByte := GetRequest(w, "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/"+countryName.AlphaCode+"/"+startDate)
	endByte := GetRequest(w, "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/"+countryName.AlphaCode+"/"+endDate)
	//Declaring variables to hold StringencyData for start and end
	var startData, endData StringencyData

	json.Unmarshal(startByte, &startData)
	json.Unmarshal(endByte, &endData)

	if endData.Key.Message != "" { //Don't have to check start as undeclared value for int is 0
		http.Error(w, "Requested  date is not in database yet, try yesterday", http.StatusBadRequest)
	}
	//Initialize object to hold the values we want to display stringency values - structs.go
	var displayData StringencyDisplay
	//Setting fields
	displayData.Country = countryName.Name
	displayData.Scope = startDate + "-" + endDate
	displayData.Stringency = endData.Key.Stringency
	displayData.Trend = endData.Key.Stringency - startData.Key.Stringency

	return displayData
}

// ====== used by Webhook
func GetPolicyWebhook(webhook WebhookRegistration) StringencyDisplay {
	date_used := GetDay(0, 0, 0) //Get todays date - inputFormatting.go
	byteData := GetRequestWebhook("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + webhook.Country + "/" + date_used)

	var data StringencyData
	err := json.Unmarshal(byteData, &data)
	if err != nil {
		log.Fatal("Couldn't unmarsal stringency data")
	}

	//Testing the data
	//If the data has not been posted, the Service provides a messsage field, where as if it has data it doesnt
	if data.Key.Message != "" { //Check if the string is not nil value
		date_used = GetDay(0, 0, -2)
		byteData := GetRequestWebhook("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + webhook.Country + "/" + date_used)
		data.Key.Message = "" //Reset the Messsage field as correct data wont override it back to "" when unmarshalling
		json.Unmarshal(byteData, &data)
	}
	//Initializing StringencyDisplay object, read structs.go
	var displayData StringencyDisplay
	//Filling fields
	displayData.Country = webhook.Country //Using the official name
	displayData.Scope = date_used
	displayData.Stringency = data.Key.Stringency
	displayData.Trend = 0
	//Format data
	return displayData
}

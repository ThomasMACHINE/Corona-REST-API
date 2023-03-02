package coronastats

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//Declaring global array of webhook objects
var webhookDB []WebhookRegistration

// ======= notifications endpoint ===========

func UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		// Expects incoming body in terms of WebhookRegistration struct
		var webhook WebhookRegistration
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
			return
		}
		webhook.Id = rand.Intn(100000000000) //Not proud of this one
		webhook.Active = true

		if webhook.Timeout < 10 {
			webhook.Timeout = 10
			fmt.Fprint(w, "Webhook timeout was changed to 10 seconds")
		}
		//Doing tests on the object to check that it has valid fields - check structs.go
		//Does 3 tests at once, could be improved by returning the tests that failed
		//As it could fail multiple and then give more informative error messages
		if !webhook.Test() {
			http.Error(w, "Bad data given into webhook post request, check the body of your request", http.StatusBadRequest)
			return
		} else { //The webhook was succesfully initialized
			webhookDB = append(webhookDB, webhook) //Add webhook to the webhook array
			fmt.Fprint(w, webhook.Id)              //Print the id to user
			if webhook.Trigger == "ON_TIMEOUT" {   //Check Trigger field to find correct routine
				runOnTimeOut(webhook)
			} else {
				runOnUpdate(webhook, "")
			}
		}
	case http.MethodGet:
		// For now just return all webhooks, don't respond to specific resource requests
		err := json.NewEncoder(w).Encode(webhookDB)
		if err != nil {
			http.Error(w, "Something went wrong wrong while encoding: "+err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Method not supported: "+r.Method, http.StatusBadRequest)
	}
}

//repeats until .Active field in webhhook is set to false (Delete method in notifications/id endpoint)
func runOnTimeOut(webhook WebhookRegistration) {
	//Set up a switch to check which data it wants
	switch webhook.Field {
	case "country":
		data := GetCPC_Webhook(webhook) //Nowhere to write responses to sounds very good to me
		content, err := json.Marshal(data)
		if err != nil {
			log.Println("Error during marshalling in invoctation")
		}
		CallUrl(webhook.Url, string(content))
	case "stringency":
		data := GetPolicyWebhook(webhook)
		content, err := json.Marshal(data)
		if err != nil {
			log.Println("Error during marshalling in invoctation")
		}
		CallUrl(webhook.Url, string(content))
	}
	//Pause execution
	time.Sleep(time.Duration(webhook.Timeout) * time.Second)
	//Check if webhook is still active
	if !webhook.Active {
		return
	} else { //Active is set to true run again!
		go runOnTimeOut(webhook)
	}
}
func runOnUpdate(webhook WebhookRegistration, date string) {
	var updated_date string
	switch webhook.Field {
	case "country":
		data := GetCPC_Webhook(webhook) //Nowhere to write responses to sounds very good to me
		if data.Scope == date {         //If date is same then it has not updated
		} else {
			updated_date = data.Scope //New data! Set updated_date to date from new data
			content, err := json.Marshal(data)
			if err != nil {
				log.Println("Error during marshalling in invoctation")
			}
			CallUrl(webhook.Url, string(content))
		}
	case "stringency":
		data := GetPolicyWebhook(webhook)
		if data.Scope == date {
		} else {
			updated_date = data.Scope
			content, err := json.Marshal(data)
			if err != nil {
				log.Println("Error during marshalling in invoctation")
			}
			fmt.Println(webhook, string(content))
			CallUrl(webhook.Url, string(content))
		}

	}
	//Pause execution
	time.Sleep(time.Duration(webhook.Timeout) * time.Second)
	//Check if webhook is still active
	if !webhook.Active {
		return
	} else { //Active is set to true run again!
		go runOnUpdate(webhook, updated_date)
	}
}

//Most of this function is inspiried from: https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/webhooksDemo/handlers.go
func CallUrl(url string, content string) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(content)))
	if err != nil {
		log.Printf("Error during request creation: %v", err.Error())
		return
	}
	// Hash content
	var byteKey = []byte("aafjapasdasfnefweg")
	mac := hmac.New(sha256.New, byteKey)
	_, err = mac.Write([]byte(content))
	if err != nil {
		log.Printf("%v", "Error during content hashing.")
		return
	}
	// Convert to string & add to header
	req.Header.Add("X-SIGNATURE", hex.EncodeToString(mac.Sum(nil)))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request: " + err.Error())
		return
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response: " + err.Error())
		return
	}
	fmt.Println("Invoked webhook, statuscode:" + strconv.Itoa(res.StatusCode) +
		", body: " + string(response))
}

// ====== notifications/id endpoint =======
func GetWebhook(w http.ResponseWriter, r *http.Request) {
	id_String := GetKeyword(w, r, "ID")
	id, err := strconv.Atoi(id_String)
	if err != nil {
		http.Error(w, "Non-existent ID", http.StatusBadRequest) //Best to not be too descriptive
		return
	}
	var webhook WebhookRegistration
	//Iteratting through the webhookDB to see if it is there
	//Could just change the structure of the storing of webhooks to index by ID to not have to iterate through every webhook
	for _, w := range webhookDB {
		if w.Id == id {
			webhook = w
			break
		}
	}
	//Fun fact, if i did this check by ID == 0 there would be a 1 in 100 billion chance the users webhook is unusable
	if webhook.Url == "" { //Check if we found a webhook object
		http.Error(w, "Non-existent ID", http.StatusBadRequest)
		return
	}
	//Valid webhook has been found
	switch r.Method {

	case http.MethodDelete: //For Delete requests
		webhook.Active = false //Set active to false, which will make the Go routine break once it is done sleeping
		fmt.Fprint(w, "Your webhook has now been deactivated")

	case http.MethodGet:
		webhookData, err := json.Marshal(webhook)
		if err != nil {
			http.Error(w, "I have no clue how this error could occur", http.StatusBadGateway)
			return
		}
		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Fprint(w, string(webhookData))
	}
}

//Returns the list of webhooks
func GetWebhookDB() []WebhookRegistration {
	return webhookDB
}

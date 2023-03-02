package coronastats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var startTime = time.Now()

func GetDiag(w http.ResponseWriter, r *http.Request) {
	var diag DiagInfo
	//Setting status, first we do a request to all the urls for the API's
	mmStatus, _ := http.Get("https://blog.mmediagroup.fr/post/m-media-launches-covid-19-api/")
	ctStatus, _ := http.Get("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-03-05/2021-03-19") //Doing a Get request on their API site gives a 404, so I am just checking with a dummy request
	rcStatus, _ := http.Get("https://restcountries.eu")
	//Then we assign the statuscode of the requests to their respective fields
	diag.Mmediagroupapi = mmStatus.StatusCode
	diag.Covidtrackerapi = ctStatus.StatusCode
	diag.RestCountriesAPI = rcStatus.StatusCode
	//Set registered
	webhookList := GetWebhookDB()
	diag.Registered = len(webhookList)
	//Set version
	diag.Version = 1
	//Set uptime
	diag.Uptime = int(time.Since(startTime).Seconds())
	//Encode the object
	diag_data, err := json.Marshal(diag)
	//Checking for error
	if err != nil {
		http.Error(w, "Something went wrong while encoding data: "+err.Error(), http.StatusBadGateway)
		return
	}
	//setting content type to json
	http.Header.Add(w.Header(), "content-type", "application/json")
	//Send out data
	fmt.Fprint(w, string(diag_data))
}

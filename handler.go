//This file holds the http handlers and is reesponsible
//For calling on the correct handler function for every request.
package coronastats

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

//HandleRequest sets up port and calls on functions for their respective Url call
func HandleRequest() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/corona/v1/country/", GetCPC)             //Cases per country endpoint -> country.go
	http.HandleFunc("/corona/v1/policy/", GetPolicy)           //Policy endpoint 			  -> policy.go
	http.HandleFunc("/corona/v1/notifications", UpdateWebhook) //Notification endpoint	  -> weebhooks.go
	http.HandleFunc("/corona/v1/notifications/", GetWebhook)   //Notification endpoint	  -> weebhooks.go
	http.HandleFunc("/corona/v1/diag/", GetDiag)
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

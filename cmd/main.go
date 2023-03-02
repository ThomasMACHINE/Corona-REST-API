package main

import (
	"coronastats"
	"math/rand"
	"time"
)

//Starts the program with calling main function and setting up the request handler
func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Set random seed so we get random number in rand calls
	coronastats.HandleRequest()            //sets up the http listener, found in the networking.go file in root directory
}

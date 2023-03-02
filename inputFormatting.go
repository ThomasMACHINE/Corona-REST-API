//This file contains all refactored functionalities to format and handle data from the request (URL)
package coronastats

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

//GetKeyword gets a requested Keyword from other function and then decodes
//the url to find the data requested
//@return - returns one of the Keywords requested (CountryName, dates, limit)
func GetKeyword(w http.ResponseWriter, r *http.Request, requestedData string) string {
	path := strings.Split(r.URL.Path, "/")

	switch requestedData {

	case "Name":
		//expected look of slice: /corona/v1/country/country_name
		Name := strings.Join(path[4:], "") // country_name
		return Name

	case "StartDate":
		dateString := r.FormValue("scope")         // "yyyy1-mm1-dd1-yyyy2-mm2-dd2"
		dateList := strings.Split(dateString, "-") // [yyyy_1 mm_1 dd_1 yyyy_2 mm_2 dd_2]
		//Check if the scope input has the proper formating + proper amount of elements
		if len(dateList) != 6 {
			err := fmt.Errorf("scope is written incorrectly, please write it in yyyy1-mm1-dd1-yyyy2-mm2-dd2 format"+"\n"+"Wrong input was: : %v", dateList) //Create error msg
			fmt.Fprint(w, err.Error())
			return "Crash"
		}
		startDate := strings.Join(dateList[:3], "-") //"yyyy_1-mm_1-dd_1"

		return startDate

	case "EndDate":
		dateString := r.FormValue("scope")         // "yyyy1-mm1-dd1-yy2-mm2-dd2"
		dateList := strings.Split(dateString, "-") // [yyyy_1 mm_1 dd_1 yy_2 mm_2 dd_2]

		//Check if the scope input has the proper formating + proper amount of elements
		if len(dateList) != 6 {
			err := fmt.Errorf("scope is written incorrectly, please write it in yyyy1-mm1-dd1-yyyy2-mm2-dd2 format"+"\n"+"Wrong input was: : %v", dateList) //Create error msg
			fmt.Fprint(w, err.Error())
			return "Crash"
		}
		endDate := strings.Join(dateList[3:], "-") //"yyyy_2-mm_2-dd_2"
		return endDate

	case "ID":
		{
			//Path = [corona v1 notifications id ""]
			if len(path) != 5 {
				http.Error(w, "Malformed url: "+r.URL.String(), http.StatusBadRequest)
			}
			id := strings.Join(path[4:5], "") // id
			return id
		}
	default:
		err := fmt.Errorf("wrong keyword used in function call to GetKeyword: %v", requestedData) //Create error
		fmt.Println(err.Error())                                                                  //Print error to terminal as this is an error only developer can cause
		return "Crash"                                                                            //remaking return value to interface sounds like bad practice
	}
}

//GetDay returns the date into the format we use for our GET requests
//The parameters: year, month, day are number of years, months, day we want to add/subtract from todays date
func GetDay(year int, month int, day int) string {
	dateObject := time.Now().AddDate(year, month, day) //Gets the desired date as a time object
	dateString := dateObject.String()                  // "2006-01-02 15:04:05.999999999 -0700 MST"
	dateSlice := strings.Split(dateString, " ")        //[2006-01-02 15:04:05.999999999 -0700 MST]
	date := strings.Join(dateSlice[0:1], "")
	return date
}

//If you spent several hours trying to format time objects before finding time.String()  (What a painful experience)
//Here is my tribute to you: https://www.youtube.com/watch?v=BG6EtT-mReM

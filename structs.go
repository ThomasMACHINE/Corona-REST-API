//this file holds all the structs, its "constructors" and documentation
package coronastats

//SearchCountry acts as a placeholder to unmarshal the JSON values from the GetCountry
//As we currently are
type SearchCountry struct {
	Name      string `json:"name"`       //Country Name
	AlphaCode string `json:"alpha3Code"` //Alpha 3 code is an official 3 character unique ID for every country
}

/*
Population int    `json:"population"` //No countries have more than 2^31 people, so we wont overflow
	//cool fact: unlike java and cpp, int size is platform dependent and is 2^63 on 64 bit machines, read: https://golangbyexample.com/go-size-range-int-uint/
*/

//=========== FIRST ENDPOINT STRUCTS ============

//mmediagroup's covid api returns in a
type CovidCaseData struct {
	Key CountryCase `json:"All"`
}

//CountryCase is the structure used to format and hold data it sends out for the CasesPerCountry endpoint
type CountryCase struct {
	Country    string         `json:"country"`
	Continent  string         `json:"continent"`
	Population int            `json:"population"`
	Dates      map[string]int `json:"dates"`
}

//CountryCaseDisplay holds all the fields we want to export for Cases Per Country endpoint,
//could have made it inherit all the fields by plaacing it in countryCase
//Struct, but i feel like it would just look very messy when reading
type CountryCaseDisplay struct {
	Country               string  `json:"country"`
	Continent             string  `json:"continent"`
	Population            int     `json:"population"`
	Scope                 string  `json:"scope"`
	Recovered             int     `json:"recovered"`
	Confirmed             int     `json:"confirmed"`
	Population_percentage float64 `json:"population_percentage"`
}

//Constructor for displayData struct
func (c *CountryCaseDisplay) Init(data CovidCaseData) {
	c2 := data.Key //Key holds all CountryCase object
	c.Country = c2.Country
	c.Continent = c2.Continent
	c.Population = c2.Population
}

// ========== SECOND ENDPOINT STRUCTS ====================
type StringencyData struct {
	Key StringencyValues `json:"stringencyData"`
}

type StringencyValues struct {
	Date              string   `json:"date_value"`
	Stringency_Actual *float32 `json:"stringency_actual"`
	Stringency        float32  `json:"stringency"`
	Message           string   `json:"msg"`
}

type StringencyDisplay struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float32 `json:"stringency"`
	Trend      float32 `json:"trend"`
}

// ======= WEBHOOK STRUCTS ==========

type WebhookRegistration struct {
	Id      int    `json:"id"`
	Url     string `json:"url"`     //https://localhost:8080/client/
	Timeout int    `json:"timeout"` //3600
	Field   string `json:"field"`   //stringency
	Country string `json:"country"` //France
	Trigger string `json:"trigger"` //ON_CHANGE
	Active  bool   `json:"active"`
}

func (wr *WebhookRegistration) Test() bool {
	if wr.MissingFields() {
		return false
	} else if !wr.ValidTrigger() {
		return false
	} else if !wr.ValidField() {
		return false
	} else {
		return true
	}
}

//Method for webhookregistration to see if any fiels are empty (default value)
//Could improve it by doing an if statement for every field and send back as slice, which would greatly improve error message
func (wr *WebhookRegistration) MissingFields() bool {
	if wr.Url == "" || wr.Field == "" || wr.Country == "" || wr.Trigger == "" {
		return true
	}
	if wr.Timeout == 0 {
		return true
	}
	return false
}

//Webhook method to check if trigger is correct
func (wr *WebhookRegistration) ValidTrigger() bool {
	//Check if Trigger is one of the accepted methods
	if wr.Trigger == "ON_CHANGE" || wr.Trigger == "ON_TIMEOUT" { //Could improve with a slice of accepted triggers but i only have 2 so it ok for now
		return true
	}
	return false
}

func (wr *WebhookRegistration) ValidField() bool {
	//Check if Trigger is one of the accepted methods
	if wr.Field == "stringency" || wr.Field == "country" { //Could improve with a slice of accepted triggers but i only have 2 so it ok for now
		return true
	}
	return false
}

//======= DIAG STRUCT ==========

type DiagInfo struct {
	Mmediagroupapi   int // <http status code for mmediagroupapi API>",
	Covidtrackerapi  int // <http status code for covidtrackerapi API>"
	RestCountriesAPI int // <http status code for restcountries API>"
	Registered       int //  <number of registered webhooks>,
	Version          int // "version of applicatiton",
	Uptime           int // <time in seconds from the last service restart>

}

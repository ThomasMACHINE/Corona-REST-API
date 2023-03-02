# Assignment 2

to build and run:
```
go build cmd/main.go
./main
```
# Purpose
This is a RESTservice made to provide information about corona cases occuring in different countries and associated government responses. It makes use of three REST web services to provide you with this information:

[Blog.mmediagroup](https://blog.mmediagroup.fr/post/m-media-launches-covid-19-api/) gives information about covid cases,

[Covidtracker](https://covidtracker.bsg.ox.ac.uk/about-api) gives information about stringency policies

[RestCountries](https://restcountries.eu) is used to find associated name for country used in endpoints

There are four endpoints for this service. And this syntax when reading the valuess of the path:

* ```{:value}``` indicates mandatory input parameters specified by the user (i.e., the one using *your* service).
* ```{value}``` indicates optional input specified by the user (i.e., the one using *your* service), where `value' can itself contain further optional input. The same notation applies for HTTP parameter specifications (e.g., ```{?param}```).

The url provided for forming requests is: http://10.212.136.111:8080

# Country endpoint
Provides data for recovered and confirmed covid19 victims from Covidtracker
```Path: /corona/v1/country/{:country_name}{?scope=begin_date-end_date}```
If scope is not set, latest available data will be provided. When scope is set it should be written in a "yyyy-mm-dd-yyyy-mm-dd" form.
```Example Request: http://10.212.136.111:8080/corona/v1/country/France```

``` Example Response:
{
    "country": "France",
    "continent": "Europe",
    "population": 64979548,
    "scope": "total",
    "recovered": 260389,
    "confirmed": 4513685,
    "population_percentage": 0.07
}
```
* Content type: `application/json`

# Stringency endpoint

Provides data for stringency in given country and the change of stringency: "trend" collected by Blog.mmediagroup
```Path: /corona/v1/policy/{:country_name}{?scope=begin_date-end_date}```
If scope is not set, latest available data will be provided and trend is set to 0. When scope is set it should be written in a "yyyy-mm-dd-yyyy-mm-dd" form.
```Example Request: http://10.212.136.111:8080/corona/v1/policy/France```

```
Example response:
{
    "country": "France",
    "scope": "total",
    "stringency": 68.52,
    "trend": 0
} 
```
* Content type: `application/json`

For both policy and country endpoints the country_name value is flexible in the sense that you can write the name in many different forms. Examples are that you can input "Norge" for norway and you can also make small typos. However there is a chance that you might not get any country if your value is too obscure. Restcountries.eu is the webservice used to make this possible.

Please take note that this functionality is not implemented in the webhook registration.

# Diag endpoint
Sends the diagnostics data of the RESTservice
### - Request

```
Path: /corona/v1/country/diag
```
```
Example Response:
{
    "Mmediagroupapi": 200, //Status code for web service
    "Covidtrackerapi": 200, //Status code for web service
    "RestCountriesAPI": 200, //Status code for web service
    "Registered": 0, //Registered webhooks
    "Version": 0, // Version of the service
    "Uptime": 5788 //runtime meassured in seconds
}
```

# Webhook endpoint

### Registration of Webhook

### - Request
This application also has webhooks, but invocation does not currently work.
```
Method: POST
Path: /corona/v1/notifications/
```


# RequiredBody:
```
{
   "url": "https://localhost:8080/client/", //Url that you want to be invoked
   "timeout": 3600,                         //Time in seconds for intervals inbetween invocations
   "field": "country",                      //which endpoint you want data from
   "country": "France",                     //Country that you want data on
   "trigger": "ON_CHANGE"                   //When the webhook should invoke the url
}
```
Note, 
Accepted triggers: ON_CHANGE and ON_TIMEOUT. 
Accepted fields: "country" and "stringency"

For country you have to use the Official name if field is country(case sensitive), and the ALPHA3-code for the country if field is stringency


### Disabling Webhook
You can also disable your webhook
### - Request

```
Method: DELETE
Path: /corona/v1/notifications/{id}
```

* {id] is the ID for the webhook registration

### - Response

```
Your webhook has now been deactivated
```
* Content type: plain text

### View registered webhook
You can get all the stored fields from the object by doing a GET request on the same url.
### - Request

```
Method: GET
Path: /corona/v1/notifications/{id}
```

* {id] is the ID for the webhook registration

### - Response

Body (Example):
```
{
   "id": "OIdksUDwveiwe",
   "url": "http://localhost:8080/client/",
   "timeout": 3600,
   "field": "stringency",
   "country": "France",
   "trigger": "ON_CHANGE"
}
```
* Content type: `application/json`

### View all registered webhooks
You can also see all the registered webhooks with their respective field
### - Request

```
Method: GET
Path: /corona/v1/notifications/
```

### - Response
The response is a collection of registered webhooks as specified in the POST body, alongside the server-defined ID.

Body (Example):
```
[{
   "id": "OIdksUDwveiwe",
   "url": "https://localhost:8080/client/",
   "timeout": 3600,
   "information": "stringency",
   "country": "France",
   "trigger": "ON_CHANGE"
},
...
]
```
* Content type: `application/json`


### Webhook Invocation

```
Method: POST
Path: <url specified in the corresponding webhook registration>
```

Any invocation of the registered webhook has the format of the output corresponding to the `information` that has been registered. For example, if `stringency` was specified as information during webhook registration, the structure of the body would follow the policy endpoint output; conversely, if registering `confirmed` as information value for the webhook, the latest (no date ranges) information about confirmed Covid-19 cases for the given country are returned.

The `scope`field will read the date used in acquiring the information as the external webservices used update infrequently.

### Descrepencies

[covidtracker]https://covidtracker.bsg.ox.ac.uk/about-api has some issues where it can take 26+ hours to update the latest day, this is fixed by calling back two days if todays data has not been published.

## Descrepencies in the RESTservice: https://blog.mmediagroup.fr/post/m-media-launches-covid-19-api/

## Mayor Descrepencies:

Some countries only update values  monthly:
 Norway has only been irregulary updated on: 
```"2020-11-16", "2020-10-07", "2020-09-26", and once every month on a random day since "2020-01-22"```
Inbetween these slots the number of recovered/cases does not change, which makes our reported data very inaccurate

# Minor descrepencies:

for France: https://covid-api.mmediagroup.fr/v1/history?country=France&status=Recovered there are no recovered or confirmed cases between the two dates: "2021-03-21": 251238, "2021-03-20": 251238, 

Recovered not updated but confirmed cases was updated:
```"2020-10-18": 85008, "2020-10-17": 85008```

```"2020-08-30": 73279, "2020-08-29": 73279```

A drop in recovered, followed by a day of no update
```"2020-08-31": 73537, "2020-08-30": 73279, "2020-08-29": 73279, "2020-08-28": 73501```


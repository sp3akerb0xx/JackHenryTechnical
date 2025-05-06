Nitin Chennam's Weather API

## Startup
- In order to run this project, you will need to be on `go v1.24.x`. 
- Once you have confirmed that you are running the correct version of go, then you can run the startup script provided by running `./run.sh` in your terminal

## Available Endpoints
- This API has 3 exposed endpoints
    - `/sample/`
    - `/weather/simple/{latitude},{longitude}`
    - `/weather/detailed/{latitude},{longitude}`

### /sample
- This endpoint exists as a form of "health check" for the API, taking no arguments and returning back a sample JSON body 

### /weather/simple/{latitude},{longitude}
- Example Request: `curl --location 'http://localhost:8080/weather/simple/43.680031,-70.310425'`
- This endpoint accepts longitude and latitude coordinates and returns back an simple response body of the following format:

```
{
	"shortForecast": "Partly Cloudy",
	"tldr": "moderate"
}
```
- This was the original output requested by the terms of this technical interview
    - Obviously, there are additional features in this API, but this endpoint will do exactly what was requested
- "tldr" stands for "Too Long, Didn't Read", which how the young folk say "To summarize". This field provides my characterization of the weather in one of three ways: "cold", "moderate", or "hot"
    - I used a very simple mapping to determine this, based entirely on my preference. 
        - If the temperature is below 45 degrees Fahrenheit, it is labelled as "cold"
        - If the temperature is above 75 degrees Fahrenheit, it is labelled as "hot"
        - Every other temperature was classified as "moderate"
- This endpoint utilizes the NWS `/gridpoints/{wfo}/{x},{y}/forecast` endpoint to pull the necessary data
- Given more time, I would have come up with a way to integrate wind chill and heat index into the classification as well, but I already spent too much time working on the rest of this. So instead, I created the last endpoint that this API exposes to allow for additional context should the end user desire the information

### /weather/detailed/{latitude},{longitude}
- Example Request: `curl --location 'http://localhost:8080/weather/detailed/43.680031,-70.310425'`
- This endpoint accepts longitude and latitude coordinates and returns back a detailed response body of the following format:

```
{
	"temperature": 53.96,
	"windChill": "Temperature too warm for wind chill",
	"heatIndex": "Temperature too cold for heat index",
	"relativeHumidity": "61%"
}
```
- I used a different endpoint to retrieve this data: `/zones/forecast/{zoneId}/observations`, as the endpoint used in the simple endpoint above did not contain the necessary data points
- This endpoint is more of an ongoing experiment, as I want to implement some sort of logic to provide a "felt temperature" by taking into account the relative humidity, the heat index, and the wind chill, should those factors come into play
    - This is clearly beyond the scope of the assignment, but I thought including something like that would be the obvious next iteration of this API in a world where this was a project someone was working on
    - So, in the spirit of iterative development, I implemented the first part of that new feature: hitting a different NWS endpoint to pull the relative humidity, wind chill, and heat index numbers to make them available to use in any future logic that I might work on
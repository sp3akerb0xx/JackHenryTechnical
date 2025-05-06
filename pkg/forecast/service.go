package forecast

import (
	"encoding/json"
	"fmt"
	"goAPI/v2/models"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Handles the bulk of the logic for the simple request endpoint. 
func SimpleRequest(latitude float64, longitude float64, client *http.Client, w http.ResponseWriter) (http.ResponseWriter, error) {
	var forecastResponse models.ForecastResponse

	pointResponse, err := initialRequest(latitude, longitude, client)
	if err != nil {
		return nil, err
	}
	newUrl := pointResponse.Properties.Forecast

	log.Debug().Msgf("Making GET request at the following address: %s", newUrl)
	log.Info().Msg("Making GET request to NWS gridpoints endpoint")
	res, err := client.Get(newUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&forecastResponse)
	if err != nil {
		return nil, err
	}

	weatherReport := models.SimpleReport{
		ShortForecast: forecastResponse.Properties.Periods[0].ShortForecast,
		TLDR: applyLogic(forecastResponse),
	}

	w = returnHandler(w, weatherReport)
	return w, nil
}

// Handles the logic for the Detailed Report endpoint
func DetailedRequest(latitude float64, longitude float64, client *http.Client, w http.ResponseWriter) (http.ResponseWriter, error){
	var observationResponse models.ObservationResponse
	pointResponse, err := initialRequest(latitude, longitude, client)
	if err != nil {
		return nil, err
	}
	newUrl := fmt.Sprintf("%s/observations", pointResponse.Properties.ForecastZone)
	log.Debug().Msgf("Making GET request at the following address: %s", newUrl)
	log.Info().Msg("Making GET request to NWS /zones/forecast endpoint")

	res, err := client.Get(newUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&observationResponse)
	if err != nil {
		return nil, err
	}

	detailedReport := models.DetailedReport{
		Temperature:		observationResponse.Features[0].Properties.Temperature.Value * 1.80 + 32.00,
		RelativeHumidity: 	fmt.Sprintf("%v%%", int(observationResponse.Features[0].Properties.RelativeHumidity.Value)),
	}

	if observationResponse.Features[0].Properties.WindChill.Value == nil {
		detailedReport.WindChill = "Temperature too warm for wind chill"
	} else {
		detailedReport.WindChill = observationResponse.Features[0].Properties.WindChill.Value
	}
	
	if observationResponse.Features[0].Properties.HeatIndex.Value == nil {
		detailedReport.HeatIndex = "Temperature too cold for heat index"
	} else {
		detailedReport.HeatIndex =  observationResponse.Features[0].Properties.HeatIndex.Value
	}

	w = returnHandler(w, detailedReport)
	return w, nil
}

// Since both of the endpoints require this initial API call, I decided to isolate it into a seperate function
// This call is necessary as it returns the subsequent GET request we need to make to get our data
func initialRequest(latitude float64, longitude float64, client *http.Client) (models.PointResponse, error){
	var pointResponse models.PointResponse
	url := fmt.Sprintf("https://api.weather.gov/points/%f,%f", latitude, longitude)
	log.Debug().Msgf("Making GET request at the following address: %s", url)
	log.Info().Msg("Making initial request to NWS")
	resp, err := client.Get(url)
	if err != nil {
		return pointResponse, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&pointResponse)
	return pointResponse, nil
}

// handles encoding final return struct into JSON while doing some simple formatting
func returnHandler(w http.ResponseWriter, x interface{}) http.ResponseWriter{
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(x)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return w
}

// Isolates the mapping of the temperatures and the TLDR
func applyLogic(forecast models.ForecastResponse) (string) {
	temp := forecast.Properties.Periods[0].Temperature
	log.Debug().Msgf("Temperature: %v", temp)
	if temp < 45 {
		return "cold"
	} else if temp > 75 {
		return "hot"
	} else {
		return "moderate"
	}
}

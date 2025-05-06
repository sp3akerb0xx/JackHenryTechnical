package api

import (
	
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"goAPI/v2/models"
	"goAPI/v2/pkg/forecast"
	"net/http"
	"encoding/json"
	"time"
	"strconv"

)

var sample models.SimpleReport
var client *http.Client


func StartApi(){

	log.Info().Msg("Creating Router...")
	router := mux.NewRouter()

	log.Info().Msg("Creating http client...")
	client = &http.Client{
		Timeout: time.Duration(1) * time.Second,
	}


	//Define the endpoints for operations
	log.Info().Msg("Creating handlers for endpoints...")
	router.HandleFunc("/weather/simple/{latitude},{longitude}", getSimpleForecast).Methods("GET")
	router.HandleFunc("/weather/detailed/{latitude},{longitude}", getDetailedForecast).Methods("GET")
	router.HandleFunc("/sample", getSample).Methods("GET")

	//Start the HTTP server
	log.Info().Msg("Starting up server on port 8080...")
	log.Fatal().Err(http.ListenAndServe(":8080", router))
}

// SimpleForecast returns the response that was requested as a part of this interview 
func getSimpleForecast(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Retrieving Simple Forecast")
	vars := mux.Vars(r)

	latitude, err := strconv.ParseFloat(vars["latitude"], 32)
	if err != nil{
		log.Fatal().Msg(err.Error())
	}
	longitude, err := strconv.ParseFloat(vars["longitude"], 32)
	if err != nil{
		log.Fatal().Msg(err.Error())
	}
	log.Debug().Msgf("Coordinates Parsed. Verify for accuracy: \t %f, %f", latitude, longitude)

	w, err = forecast.SimpleRequest(latitude, longitude, client, w)
	if err != nil {
		log.Err(err)
	}
}

// DetailedForecast returns some more technical weather metrics from a different NWS endpoint
// You can ask me why I chose to do this if you want. There's a bit of a story behind it...
func getDetailedForecast(w http.ResponseWriter, r *http.Request){
	log.Info().Msg("Retrieving Detailed Forecast")
	vars := mux.Vars(r)

	latitude, err := strconv.ParseFloat(vars["latitude"], 32)
	if err != nil{
		log.Fatal().Msg(err.Error())
	}
	longitude, err := strconv.ParseFloat(vars["longitude"], 32)
	if err != nil{
		log.Fatal().Msg(err.Error())
	}
	log.Debug().Msgf("Coordinates Parsed. Verify for accuracy: \t %f, %f", latitude, longitude)

	w, err = forecast.DetailedRequest(latitude, longitude, client, w)
	if err != nil {
		log.Err(err)
	}
}

// A sanity check function to make sure that the server is actually working
func getSample(w http.ResponseWriter, r *http.Request){
	//retrieving sample response
	log.Info().Msg("Retrieving sample response")
	localSample, err := models.ReturnSample()
	if err != nil{
		log.Fatal().Msg(err.Error())
	}

	// some simple formatting to clean up the response on the user side
	sample = localSample
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(sample)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}	
package models

type SimpleReport struct {
	ShortForecast		string		`json:"shortForecast"`
	TLDR				string		`json:"tldr"`
}

type DetailedReport struct {
	Temperature			float64		`json:"temperature"`
	WindChill			any			`json:"windChill"`
	HeatIndex			any			`json:"heatIndex"`
	RelativeHumidity	any			`json:"relativeHumidity"`
}

type Info struct {
	Temperature	int `json:"temperature"`
	WindSpeed	string	`json:"windSpeed"`
	RealFeel	int `json:"realFeel"`
}

type PointResponse struct {
	Properties struct {
		Forecast            string `json:"forecast"`
		ForecastZone    string `json:"forecastZone"`
	} `json:"properties"`
}

type ForecastResponse struct {
	Properties struct {
		Periods []struct {
			Temperature                int    `json:"temperature"`
			TemperatureUnit            string `json:"temperatureUnit"`
			ShortForecast    string `json:"shortForecast"`
			DetailedForecast string `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

type ObservationResponse struct {
	Features []struct {
		Properties struct {
			Temperature     struct {
				Value          float64	`json:"value"`
			} `json:"temperature"`
			WindChill struct {
				Value          any    `json:"value"`
			} `json:"windChill"`
			HeatIndex struct {
				Value          any    `json:"value"`
			} `json:"heatIndex"`
			RelativeHumidity struct {
				Value          float64 `json:"value"`
			} `json:"relativeHumidity"`
		} `json:"properties"`
	} `json:"features"`
}

func ReturnSample() (SimpleReport, error) {
	sample := SimpleReport{
		ShortForecast: 	"Cloudy with a chance of meatballs",
		TLDR: 			"moderate",
	}
	return sample, nil
}
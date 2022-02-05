package accuweather

import (
	"encoding/json"
	"fmt"
	"github.com/blakehartshorn/go-nest-temp-monitor/configuration"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// Weather - parse weather details
type Weather struct {
	Timestamp   string `json:"LocalObservationDateTime"`
	Temperature struct {
		Metric struct {
			Value float32 `json:"Value"`
		} `json:"Metric"`
	} `json:"Temperature"`
	Humidity int `json:"RelativeHumidity"`
	Pressure struct {
		Metric struct {
			Value float32 `json:"Value"`
		} `json:"Metric"`
	} `json:"Pressure"`
}

// WriteWeather - write weather metrics to InfluxDB
func WriteWeather(
	config configuration.AccuWeatherConfig,
	influx api.WriteAPI) {
	url := fmt.Sprintf(
		"https://dataservice.accuweather.com/currentconditions/v1/%s?apikey=%s&details=true",
		config.Location,
		config.APIKey,
	)
	for {
		httpClient := &http.Client{Timeout: time.Second * 10}
		res, err := httpClient.Get(url)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		var response []Weather
		jsonErr := json.Unmarshal(body, &response)
		if jsonErr != nil {
			log.Print(string(body))
			log.Println("ERROR: Could not unmarshal json!")
			log.Print(jsonErr)
		} else {
			weather := response[0]

			p := influxdb2.NewPoint("accuweather",
				map[string]string{
					"locationKey": config.Location,
				},
				map[string]interface{}{
					"temperature": weather.Temperature.Metric.Value,
					"humidity":    weather.Humidity,
					"pressure":    weather.Pressure.Metric.Value,
				},
				time.Now())
			influx.WritePoint(p)

			log.Printf("Wrote weather metrics from AccuWeather. Sleeping for %d minute(s).\n", config.Interval)
		}
		time.Sleep(time.Minute * time.Duration(config.Interval))
	}
}

package openweathermap

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

// Response - Parse openweathermap json output
type Response struct {
	Timestamp int64  `json:"dt"`
	City      string `json:"name"`
	Sys       struct {
		Country string `json:"country"`
	} `json:"sys"`
	Weather struct {
		Humidity    int     `json:"humidity"`
		Pressure    int     `json:"pressure"`
		Temperature float32 `json:"temp"`
	} `json:"main"`
}

// WriteWeather - write weather metrics to InfluxDB
func WriteWeather(
	config configuration.OpenWeatherMapConfig,
	influx api.WriteAPI) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?id=%s&appid=%s&units=metric",
		config.CityID,
		config.AppID,
	)
	for {
		httpClient := &http.Client{Timeout: time.Second * 10}
		res, err := httpClient.Get(url)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		var weather Response
		jsonErr := json.Unmarshal(body, &weather)
		if jsonErr != nil {
			log.Println("ERROR: Could not unmarshal json!")
			log.Print(jsonErr)
		} else {

			p := influxdb2.NewPoint("openweathermap",
				map[string]string{
					"city": weather.City,
					"country": weather.Sys.Country,
				},
				map[string]interface{}{
					"temperature": weather.Weather.Temperature,
					"humidity":    weather.Weather.Humidity,
					"pressure":    weather.Weather.Pressure,
				},
				time.Unix(weather.Timestamp, 0))
			influx.WritePoint(p)

			log.Printf("Wrote weather metrics from openweathermap. Sleeping for %d minute(s).\n", config.Interval)
		}
		time.Sleep(time.Minute * time.Duration(config.Interval))
	}
}

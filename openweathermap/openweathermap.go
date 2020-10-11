package openweathermap

import (
	"encoding/json"
	"fmt"
	"go-nest-temp-monitor/configuration"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	influxClient "github.com/influxdata/influxdb1-client"
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
	influx *influxClient.Client,
	database string) {
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

			pts := make([]influxClient.Point, 1)
			pts[0] = influxClient.Point{
				Measurement: "openweathermap",
				Fields: map[string]interface{}{
					"temperature": weather.Weather.Temperature,
					"humidity":    weather.Weather.Humidity,
					"pressure":    weather.Weather.Pressure,
				},
				Tags: map[string]string{
					"city":    weather.City,
					"country": weather.Sys.Country,
				},
				Time:      time.Unix(weather.Timestamp, 0),
				Precision: "rfc3339",
			}
			bps := influxClient.BatchPoints{
				Points:   pts,
				Database: database,
			}
			_, err = influx.Write(bps)
			if err != nil {
				log.Println("ERROR: Could not write data point!")
				log.Print(bps)
				log.Print(err)
			} else {
				log.Printf("Wrote weather metrics from openweathermap. Sleeping for %d minute(s).\n", config.Interval)
			}
		}
		time.Sleep(time.Minute * time.Duration(config.Interval))
	}
}

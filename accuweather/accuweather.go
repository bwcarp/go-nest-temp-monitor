package accuweather

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
	influx *influxClient.Client,
	database string) {
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
			timestamp, _ := time.Parse(time.RFC3339, weather.Timestamp)
			pts := make([]influxClient.Point, 1)
			pts[0] = influxClient.Point{
				Measurement: "accuweather",
				Fields: map[string]interface{}{
					"temperature": weather.Temperature.Metric.Value,
					"humidity":    weather.Humidity,
					"pressure":    weather.Pressure.Metric.Value,
				},
				Tags: map[string]string{
					"locationKey": config.Location,
				},
				Time:      timestamp,
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
				log.Printf("Wrote weather metrics from AccuWeather. Sleeping for %d minute(s).\n", config.Interval)
			}
		}
		time.Sleep(time.Minute * time.Duration(config.Interval))
	}
}

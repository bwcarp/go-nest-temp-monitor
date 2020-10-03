package weathergov

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

// Response - root of JSON object returned by weather.gov
type Response struct {
	Properties struct {
		Timestamp   string  `json:"timestamp"`
		Temperature Weather `json:"temperature"`
		Humidity    Weather `json:"relativeHumidity"`
		Pressure    Weather `json:"barometricPressure"`
	} `json:"properties"`
}

// Weather - parse json data for value and unit
type Weather struct {
	Value          float64 `json:"value"`
	UnitCode       string  `json:"unitCode"`
	QualityControl string  `json:"qualityControl"`
}

// WriteWeather - write weather metrics to InfluxDB
func WriteWeather(
	config configuration.WeatherGovConfig,
	influx *influxClient.Client,
	database string) {
	url := fmt.Sprintf(
		"https://api.weather.gov/stations/%s/observations/latest",
		config.Station,
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

			timestamp, _ := time.Parse(time.RFC3339, weather.Properties.Timestamp)
			temperature := weather.Properties.Temperature.Value
			humidity := weather.Properties.Humidity.Value
			// Convert Pa to hPa for consistency with other apps
			pressure := weather.Properties.Pressure.Value * 0.01

			pts := make([]influxClient.Point, 1)
			pts[0] = influxClient.Point{
				Measurement: "weathergov",
				Fields: map[string]interface{}{
					"temperature": temperature,
					"humidity":    humidity,
					"pressure":    pressure,
				},
				Tags: map[string]string{
					"station": config.Station,
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
			}
		}
		log.Printf("Wrote weather metrics from weather.gov. Sleeping for %d minute(s).\n", config.Interval)
		time.Sleep(time.Minute * time.Duration(config.Interval))
	}
}

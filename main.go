package main

import (
	"flag"
	"fmt"
	"go-nest-temp-monitor/accuweather"
	"go-nest-temp-monitor/configuration"
	"go-nest-temp-monitor/nest"
	"go-nest-temp-monitor/openweathermap"
	"go-nest-temp-monitor/weathergov"
	"log"
	"net/url"
	"time"

	influxClient "github.com/influxdata/influxdb1-client"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "c", "./config.json", "Specify path to config.json")
	config := configuration.GetConfig(configFile)
	database := config.InfluxConfig.Database
	influxURL, _ := url.Parse(
		fmt.Sprintf(
			"%s://%s:%d",
			config.InfluxConfig.Protocol,
			config.InfluxConfig.Hostname,
			config.InfluxConfig.Port,
		),
	)
	influxConf := influxClient.Config{
		URL:      *influxURL,
		Username: config.InfluxConfig.Username,
		Password: config.InfluxConfig.Password,
	}

	influxCon, err := influxClient.NewClient(influxConf)
	if err != nil {
		log.Print("Could not connect to InfluxDB")
		log.Fatal(err)
	}

	if config.NestConfig.Enabled == true {
		nest.Initialize(config.NestConfig)
		go nest.RefreshLogin()
		time.Sleep(time.Second * 10)
		go nest.WriteNest(influxCon, database)
	}

	if config.AccuWeatherConfig.Enabled == true {
		go accuweather.WriteWeather(config.AccuWeatherConfig, influxCon, database)
	}

	if config.OpenWeatherMapConfig.Enabled == true {
		go openweathermap.WriteWeather(config.OpenWeatherMapConfig, influxCon, database)
	}

	if config.WeatherGovConfig.Enabled == true {
		go weathergov.WriteWeather(config.WeatherGovConfig, influxCon, database)
	}

	for {
		time.Sleep(time.Hour)
		log.Print("Keep alive")
	}

}

package main

import (
	"flag"
	"github.com/blakehartshorn/go-nest-temp-monitor/accuweather"
	"github.com/blakehartshorn/go-nest-temp-monitor/configuration"
	"github.com/blakehartshorn/go-nest-temp-monitor/nest"
	"github.com/blakehartshorn/go-nest-temp-monitor/openweathermap"
	"github.com/blakehartshorn/go-nest-temp-monitor/weathergov"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {

	var configFile = flag.String("c", "./config.json", "Specify path to config.json")
	flag.Parse()

	config := configuration.GetConfig(*configFile)

	influxClient := influxdb2.NewClientWithOptions(
		config.InfluxConfig.URL, config.InfluxConfig.Token,
		influxdb2.DefaultOptions().
			SetBatchSize(config.InfluxConfig.BatchSize).
			SetFlushInterval(config.InfluxConfig.FlushInt*1000).
			SetPrecision(time.Second))
	defer influxClient.Close()

	influxWriter := influxClient.WriteAPI(
		config.InfluxConfig.Org, config.InfluxConfig.Bucket)
	defer influxWriter.Flush()

	if config.NestConfig.Enabled == true {
		nest.Initialize(config.NestConfig)
		go nest.RefreshLogin()
		time.Sleep(time.Second * 10)
		go nest.WriteNest(influxWriter)
	}

	if config.AccuWeatherConfig.Enabled == true {
		go accuweather.WriteWeather(config.AccuWeatherConfig, influxWriter)
	}

	if config.OpenWeatherMapConfig.Enabled == true {
		go openweathermap.WriteWeather(config.OpenWeatherMapConfig, influxWriter)
	}

	if config.WeatherGovConfig.Enabled == true {
		go weathergov.WriteWeather(config.WeatherGovConfig, influxWriter)
	}

	for {
		time.Sleep(time.Hour)
		log.Println("Keep alive")
	}
}

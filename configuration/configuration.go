package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// ConfigRoot - load json data from config file
type ConfigRoot struct {
	InfluxConfig         InfluxConfig         `json:"influx"`
	NestConfig           NestConfig           `json:"nest"`
	AccuWeatherConfig    AccuWeatherConfig    `json:"accuweather"`
	OpenWeatherMapConfig OpenWeatherMapConfig `json:"openweathermap"`
	WeatherGovConfig     WeatherGovConfig     `json:"weather.gov"`
}

// InfluxConfig - InfluxDB configuration
type InfluxConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NestConfig - Google Nest configuration
type NestConfig struct {
	Enabled      bool   `json:"enable"`
	Interval     int    `json:"interval"`
	ProjectID    string `json:"project_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

// AccuWeatherConfig - AccuWeather configuration
type AccuWeatherConfig struct {
	Enabled  bool   `json:"enable"`
	Interval int    `json:"interval"`
	APIKey   string `json:"apikey"`
	Location string `json:"locationkey"`
}

// OpenWeatherMapConfig - OpenWeatherMap configuration
type OpenWeatherMapConfig struct {
	Enabled  bool   `json:"enable"`
	Interval int    `json:"interval"`
	AppID    string `json:"appid"`
	CityID   string `json:"cityid"`
}

// WeatherGovConfig - weather.gov configuration
type WeatherGovConfig struct {
	Enabled  bool   `json:"enable"`
	Interval int    `json:"interval"`
	Station  string `json:"station"`
}

// GetConfig - read config file
func GetConfig(filePath string) ConfigRoot {
	configFile, _ := ioutil.ReadFile(filePath)
	configRoot := ConfigRoot{}
	err := json.Unmarshal(configFile, &configRoot)
	if err != nil {
		log.Print(err)
		log.Fatal("Could not parse config file.")
	}
	return configRoot
}

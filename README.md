# go-nest-temp-monitor
This is a small daemon for gathering Nest thermostat data and comparing it to different weather sites. It writes to InfluxDB using the influxdb1 client in Go.

#### Celsius vs Fahrenheit
All metrics for all components of this app are collected in celsius because it is standard across all services. You can get fahrenheit in your influx queries and this works in Grafana as well:
```
> select (temperature * 1.8 + 32) from nest where time >= now() - 1h
name: nest
time                           temperature
----                           -----------
2020-10-04T13:06:56.40389252Z  63.4639928
```

## Google Nest Device Access Sandbox
In order to make use of the Nest monitoring features, you need to sign up for Google's developer program and jump through a number of hoops. The documentation is here: https://developers.google.com/nest/device-access/registration

You'll need to get to the point where you have a client id, client secret and refresh token. The access of this token should be limited to this app. This will update the access_token every 45 minutes. 

Stats written to Influx include temperature, humidity, mode, heat/cool settings, device/parent relationship, and whether HVAC is currently running.

## Weather sites
3 weather sites are currently available to monitor. You can monitor one or all of them. Stats gather include temperature (celsius), humidity, and pressure.

### OpenWeatherMap
This weather source may be preferable as it's free, international, updates frequently, and supports frequent API calls. There is no way you'll go over the API limit using just this app, but I recommend not setting it to run too frequently out of courtesy. You can get the `cityid` value for the config by searching for your city and grabbing it from the URI. You can sign up and get an API key from here: https://home.openweathermap.org/users/sign_up

### AccuWeather
This works well, but the free tier is limited to 50 API calls a day, so you wont want to set this lower than 30 minutes for your check interval. You can run more frequent queries starting at $25/mo. It also assumes you're registering a commercial app even if you're a hobbyist. https://developer.accuweather.com/packages

You can get the `locationid` value by search for your city on the website and grabbing it from the URI.

### Weather.gov
This is a good free option for Americans, but it doesn't update very frequently. No API key or account are required. Search for your city on weather.gov and you should find the station code, e.g. KBOS for Boston. Place this in the config and this will just work. Data is written using the timestamp on the output and influxdb deduplicates this, so don't be surprised if you only get new data every hour or so.

## Installation
No build scripts are provided, but you can do the following:
```
mkdir -p ~/go/src
cd ~/go/src
git clone https://github.com/blakehartshorn/go-nest-temp-monitor.git
cd go-nest-temp-monitor
go build go-nest-temp-monitor.go
```
Copy the binary to where you would prefer to run it from.

This is presently tested on Debian Bullseye on both ARM64 and AMD64.

## Config
You'll need to create a database in influxdb for this app in advance. If unsecured, username and password can be left blank. Mark the modules as true or false in config.json and add the appropriate API keys and station identifiers. Low intervals aren't that useful because the data provided isn't always new. The defaults shown in `config.json.example` are probably fine for most use cases.

## Running
```
Usage of ./go-nest-temp-monitor:
  -c string
        Specify path to config.json (default "./config.json")
```
The only argument is the placement of config.json. An example systemd script:
```
[Unit]
Description=Nest Thermostat Monitor
After=network.target influxdb.service
Requires=network.target

[Service]
User=username
Group=username
ExecStart=/home/username/go-nest-temp-monitor/go-nest-temp-monitor -c /home/username/go-nest-temp-monitor/config.json
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=nest

[Install]
WantedBy=multi-user.target
```

## Graphing
There are two popular solutions for graphing metrics with InfluxDB and both are easy to setup. Default configs aren't provided here, but get creative and lay the data out how you want it.
* Grafana: https://github.com/grafana/grafana
* Chronograf: https://github.com/influxdata/chronograf

# WeatherByZip
### Description
This is a simple tool written in Go that can tell you basic info about the weather depending on the zipcode you input
### Setup
#### Prerequisites:
 - Go version 1.15
Clone this repo using `git clone https://github.com/JamilDroubi/WeatherByZip.git`

Move into the folder `cd WeatherByZip`

Edit WeatherByZip.go `vim WeatherByZip.go` to replace the variable apikey with your own apikey from [openweathermap.org](https://openweathermap.org) (line 74) and remove the secrets import (delete line 12)

Move into src/aqi `cd src/aqi`

Edit aqi.go `vim aqi.go` to replace the apikey variable with your apikey from [aqicn.org](https://aqicn.org) and remove the secrets import (delete line 10)

If you replaced the api keys correctly it should look like this:
```go
apikey := "fmtjaswejsdoje"
```
Go back into the WeatherByZip folder `cd ../..`

Build the project `go build WeatherByZip.go`

### Usage
```
$ ./weatherbyzip <zipcode> 
```
If you happen to know the latitude and longitude of a place then you can use them to get weather data for a location by using
```
$ Lat_and_Long/Lat_and_Long <latitude> <longitude>
```
You can get weather data in those coordinates for the next week by using
```
$ Lat_and_Long/Lat_and_Long <latitude> <longitude> -w
```

package main

import (
    "fmt"
    "os"
    "encoding/json"
    "net/http"
    "net/url"
    "log"
    "time"
    s    "secrets"
    a    "aqi"
)

type OWM struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		Dt         int     `json:"dt"`
		Sunrise    int     `json:"sunrise"`
		Sunset     int     `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        float64 `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		Weather    []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Daily []struct {
		Dt      int `json:"dt"`
		Sunrise int `json:"sunrise"`
		Sunset  int `json:"sunset"`
		Temp    struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float64 `json:"day"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		DewPoint  float64 `json:"dew_point"`
		WindSpeed float64 `json:"wind_speed"`
		WindDeg   int     `json:"wind_deg"`
		Weather   []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds int     `json:"clouds"`
		Pop    float64     `json:"pop"`
		Uvi    float64 `json:"uvi"`
	} `json:"daily"`
}

func epochToHumanReadable(epoch int64) time.Time {
    return time.Unix(epoch, 0)
}

func main() {

    argc := len(os.Args)

    var weekly bool
    var lat string
    var lon string

    if argc < 3 || argc > 4 {
        fmt.Println("usage: weather <latitude> <longitude>")
        os.Exit(3)
    }

    if argc == 4 {
        for i := 1; i < argc; i++ {
            switch os.Args[i] {
                case "--week" :
                        weekly = true
                        break
                case "-w"     :
                        weekly = true
                        break
                case "-W"     :
                        weekly = true
                        break
                default       :
                        break
            }

            if weekly {
                switch i {
                    case  1 :
                        lat = os.Args[i + 1]
                        lon = os.Args[i + 2]
                        break
                    case  2 :
                        lat = os.Args[i - 1]
                        lon = os.Args[i + 1]
                        break
                    case  3 :
                        lat = os.Args[i - 2]
                        lon = os.Args[i - 1]
                        break
                    default :
                        break
                }
                break
            }
        }
    } else {
        lat = os.Args[1]
        lon = os.Args[2]
    }

    apikey := s.GetKey() // Get your own key from openweathermap.org

    // QueryEscape lets us put it into a url
    safelat := url.QueryEscape(lat)
    safelon := url.QueryEscape(lon)
    safekey := url.QueryEscape(apikey)
    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&exclude=minutely,hourly&units=imperial&appid=%s", safelat, safelon, safekey)

    // Build the request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal("Do: ", err)
        return
    }

    // Create an http client
    client := &http.Client{}

    // send the request
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal("Do: ", err)
        return
    }

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

    // Fill the record with our data
    var record OWM

    // Read the json data
    if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
        log.Println(err)
    }

    aqival, city, mainpol := a.GetAqi(3, lat, lon)

    fmt.Printf("\tCity:            %s\n", city)
    fmt.Printf("\tLatitude:     %g\n", record.Lat)
    fmt.Printf("\tLongitude:    %g\n", record.Lon)
    fmt.Println("Today: ")
    fmt.Printf("\tWeather:      %s\n", record.Current.Weather[0].Main)
    fmt.Printf("\tTemperature:  %g\n", record.Current.Temp)
    fmt.Printf("\tAir Quality:     %d\n", aqival)
    fmt.Printf("\tMain Pollutant:  %s\n", mainpol)
    fmt.Printf("\tSunrise:      %s\n", epochToHumanReadable(int64(record.Current.Sunrise)).Format("15:04:05"))
    fmt.Printf("\tSunset:       %s\n", epochToHumanReadable(int64(record.Current.Sunset)).Format("15:04:05"))
    fmt.Printf("\tFeels Like:   %g\n", record.Current.FeelsLike)
    fmt.Printf("\tHumidity:     %d%%\n", record.Current.Humidity)
    fmt.Printf("\tCurrent Time: %s\n", time.Now().Format("15:04:05"))

    if weekly {
        for i := 1; i < len(record.Daily); i++ {
            fmt.Printf("%d days from now:\n", i)
            fmt.Printf("\tWeather:      %s\n", record.Daily[i].Weather[0].Main)
            fmt.Printf("\tTemperature:  %g\n", record.Daily[i].Temp.Day)
            fmt.Printf("\tSunrise:      %s\n", epochToHumanReadable(int64(record.Daily[i].Sunrise)).Format("15:04:05"))
            fmt.Printf("\tSunset:       %s\n", epochToHumanReadable(int64(record.Daily[i].Sunset)).Format("15:04:05"))
            fmt.Printf("\tFeels Like:   %g\n", record.Daily[i].FeelsLike.Day)
            fmt.Printf("\tHumidity:     %d%%\n", record.Daily[i].Humidity)
        }
    }

}


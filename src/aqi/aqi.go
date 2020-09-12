package aqi

import (
    "fmt"
    "encoding/json"
    "net/http"
    "net/url"
    "log"
    "time"
    s    "secrets"
)

type WAQI struct {
	Status string `json:"status"`
	Data   struct {
		Aqi          int `json:"aqi"`
		Idx          int `json:"idx"`
		Attributions []struct {
			URL  string `json:"url"`
			Name string `json:"name"`
			Logo string `json:"logo,omitempty"`
		} `json:"attributions"`
		City struct {
			Geo  []float64 `json:"geo"`
			Name string    `json:"name"`
			URL  string    `json:"url"`
		} `json:"city"`
		Dominentpol string `json:"dominentpol"`
		Iaqi        struct {
			H struct {
				V float64 `json:"v"`
			} `json:"h"`
			P struct {
				V float64 `json:"v"`
			} `json:"p"`
			Pm25 struct {
				V int `json:"v"`
			} `json:"pm25"`
			T struct {
				V float64 `json:"v"`
			} `json:"t"`
			W struct {
				V float64 `json:"v"`
			} `json:"w"`
			Wg struct {
				V float64 `json:"v"`
			} `json:"wg"`
		} `json:"iaqi"`
		Time struct {
			S   string `json:"s"`
			Tz  string `json:"tz"`
			V   int    `json:"v"`
			Iso string `json:"iso"`
		} `json:"time"`
		Forecast struct {
			Daily struct {
				O3 []struct {
					Avg int    `json:"avg"`
					Day string `json:"day"`
					Max int    `json:"max"`
					Min int    `json:"min"`
				} `json:"o3"`
				Pm10 []struct {
					Avg int    `json:"avg"`
					Day string `json:"day"`
					Max int    `json:"max"`
					Min int    `json:"min"`
				} `json:"pm10"`
				Pm25 []struct {
					Avg int    `json:"avg"`
					Day string `json:"day"`
					Max int    `json:"max"`
					Min int    `json:"min"`
				} `json:"pm25"`
				Uvi []struct {
					Avg int    `json:"avg"`
					Day string `json:"day"`
					Max int    `json:"max"`
					Min int    `json:"min"`
				} `json:"uvi"`
			} `json:"daily"`
		} `json:"forecast"`
		Debug struct {
			Sync time.Time `json:"sync"`
		} `json:"debug"`
	} `json:"data"`
}

func epochToHumanReadable(epoch int64) time.Time {
    return time.Unix(epoch, 0)
}

func GetAqi(argc int, lat string , lon string) (aqival int, city string, mainpol string){

    if argc != 3 {
        fmt.Println("usage: AQI <latitude> <longitude>")
    }

    apikey := s.AqiKey() // Get your own key from aqicn.org

    // QueryEscape lets us put it into a url
    safelat := url.QueryEscape(lat)
    safelon := url.QueryEscape(lon)
    safekey := url.QueryEscape(apikey)
    url := fmt.Sprintf("http://api.waqi.info/feed/geo:%s;%s/?token=%s", safelat, safelon, safekey)

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
    var record WAQI

    // Read the json data
    if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
        log.Println(err)
    }

    // fmt.Printf("\tCity:            %s\n", record.Data.City.Name)
    // fmt.Println("Today: ")
    // fmt.Printf("\tAir Quality:     %d\n", record.Data.Aqi)
    // fmt.Printf("\tMain Pollutant:  %s\n", record.Data.Dominentpol)
    // fmt.Printf("\tCurrent Time: %s\n", time.Now().Format("15:04:05"))

    city    = record.Data.City.Name
    aqival  = record.Data.Aqi
    mainpol = record.Data.Dominentpol

    return
}


package main

import (
    "fmt"
    "os"
    "encoding/json"
    "net/http"
    "net/url"
    "log"
    "time"
)

type OWM struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64    `json:"sunrise"`
		Sunset  int64    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func epochToHumanReadable(epoch int64) time.Time {
    return time.Unix(epoch, 0)
}

func main() {
    zip := os.Args[1]
    argc := len(os.Args)
    if argc != 2 {
        fmt.Println("usage: weather <zip>")
        os.Exit(3)
    }
    // fmt.Println(zip)
    if zip[0] < '0' || zip[0] > '9' || len(zip) != 5 {
        fmt.Println("usage: weather <zip>")
        os.Exit(3)
    }
    var apikey string = "Don't take my Secrets!"
    // QueryEscape lets us put it into a url
    safezip := url.QueryEscape(zip)
    safekey := url.QueryEscape(apikey)
    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s&units=imperial&appid=%s", safezip, safekey)

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

    fmt.Printf("City: %s\n", record.Name)
    fmt.Printf("Temperature:  %g\n", record.Main.Temp)
    fmt.Printf("Feels Like:   %g\n", record.Main.FeelsLike)
    fmt.Printf("Humidity:     %d%%\n", record.Main.Humidity)
    fmt.Printf("Current Time: %s\n", time.Now().Format("15:04:05"))
    fmt.Printf("Sunrise:      %s\n", epochToHumanReadable(record.Sys.Sunrise).Format("15:04:05"))
    fmt.Printf("Sunset:       %s\n", epochToHumanReadable(record.Sys.Sunset).Format("15:04:05"))

}


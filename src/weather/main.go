package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var icons = map[string]string{
	"01d": "\uF185", // day clear
	"01n": "\uF186", // night clear
	"02d": "\uF0C2", // day few clouds
	"02n": "\uF0C2", // night few clouds
	"03d": "\uF0C2", // day scattered clouds
	"03n": "\uF0C2", // night scattered clouds
	"04d": "\uF0C2", // day broken clouds
	"04n": "\uF0C2", // night broken clouds
	"09d": "\uF0E9", // day showers
	"09n": "\uF73C", // night showers
	"10d": "\uF0E9", // day rain
	"10n": "\uF0E9", // night rain
	"11d": "\uF0E7", // day thunderstorm
	"11n": "\uF0E7", // night thunderstorm
	"13d": "\uF2DC", // day snow
	"13n": "\uF2DC", // night snow
	"50d": "\uF773", // day mist
	"50n": "\uF773", // night mist
}

type weather struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		FeelsLike float32 `json:"feels_like"`
	} `json:"main"`
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: weather [location]")
		os.Exit(1)
	}
	appId := getAppId()
	location := os.Args[1]
	degrees := '\u2103'
	weatherUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?units=Metric&appid=%s&q=%s", appId, location)
	resp, err := http.Get(weatherUrl)
	if err != nil || resp.StatusCode != http.StatusOK {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	var w weather
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		panic(err)
	}
	data, _ := json.Marshal(struct {
		Icon        string `json:"icon"`
		Description string `json:"description"`
	}{
		Icon:        icons[w.Weather[0].Icon],
		Description: fmt.Sprintf("%s, %.0f%c ", strings.Title(w.Weather[0].Description), w.Main.FeelsLike, degrees),
	})
	fmt.Printf("%s", string(data))
}
func getAppId() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	weatherFile := filepath.Join(homeDir, ".weatherkey")
	weatherKey, err := os.ReadFile(weatherFile)
	if err != nil {
		panic(fmt.Sprintf("Could not find weatherfile; %s", weatherFile))
	}

	return strings.TrimSpace(string(weatherKey))
}

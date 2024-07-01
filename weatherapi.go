package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const apiEndpoint = "https://api.weatherapi.com/v1"

type Error struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type WeatherAPIClient struct {
	key string
	c   *http.Client
}

func NewWeatherAPIClient(key string) (*WeatherAPIClient, error) {
	c := &WeatherAPIClient{}

	if len(strings.TrimSpace(key)) == 0 {
		return c, errors.New("please provide a valid api key")
	}

	c.key = key
	c.c = http.DefaultClient

	return c, nil
}

type IPLookupResponse struct {
	City      string  `json:"city"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

// GetLocationFromIP
func (w *WeatherAPIClient) GetLocationFromIP(ip string) (*IPLookupResponse, error) {
	p := new(IPLookupResponse)

	r, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("/%s?key=%s&q=%s", "ip.json", w.key, ip), nil)
	if err != nil {
		return p, err
	}

	resp, err := w.c.Do(r)
	if err != nil {
		return p, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > http.StatusCreated {
		var e = new(Error)

		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return nil, err
		}

		return nil, errors.New(e.Error.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(p); err != nil {
		return p, err
	}

	return p, nil
}

type WeatherLookupResponse struct {
	Current struct {
		TempC float32 `json:"temp_c"`
	} `json:"current"`
}

// GetCurrentWeather
func (w *WeatherAPIClient) GetCurrentWeather(lat, lon float32) (*WeatherLookupResponse, error) {
	p := new(WeatherLookupResponse)

	r, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("/%s?key=%s&q=%f,%f", "current.json", w.key, lat, lon), nil)
	if err != nil {
		return p, err
	}

	resp, err := w.c.Do(r)
	if err != nil {
		return p, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > http.StatusCreated {
		var e = new(Error)

		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return nil, err
		}

		return nil, errors.New(e.Error.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(p); err != nil {
		return p, err
	}

	return p, nil
}

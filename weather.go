package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// This should not be hard coded, but still fine
const endpoint = "https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m"

type WeatherResponse struct {
	Current struct {
		Temperature float32 `json:"temperature_2m"`
	} `json:"current"`
}

func GetWeather(latitude, longitude float32) (*WeatherResponse, error) {
	p := new(WeatherResponse)

	r, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(endpoint, latitude, longitude), nil)
	if err != nil {
		return p, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return p, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > http.StatusCreated {
		var s struct {
			Error struct {
				Message string `json:"error_message"`
			} `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}

		return nil, errors.New(s.Error.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(p); err != nil {
		return p, err
	}

	return p, nil
}

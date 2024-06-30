package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const baseEndpoint = "https://api.ip2location.io"

type Ip2LocationResponse struct {
	IP        string  `json:"ip_address"`
	City      string  `json:"city_name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func GetLocationFromIP(key, ip string) (*Ip2LocationResponse, error) {
	p := new(Ip2LocationResponse)

	r, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/?key=%s&ip=%s", baseEndpoint, key, ip), nil)
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

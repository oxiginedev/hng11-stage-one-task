package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Response struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	http.HandleFunc("GET /api/hello", HandleIncomingRequest)

	log.Println("Server started and listening on", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func HandleIncomingRequest(w http.ResponseWriter, r *http.Request) {
	res := &Response{}

	// get the visitor name parameter from url
	name := r.URL.Query().Get("visitor_name")

	ip, err := GetIpAddress(r)
	if err != nil {
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	res.ClientIP = ip

	// Get api key from environment
	ip2key := os.Getenv("IP2LOCATION_KEY")

	ip2l, err := GetLocationFromIP(ip2key, ip)
	if err != nil {
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	res.Location = ip2l.City

	weather, err := GetWeather(ip2l.Latitude, ip2l.Longitude)
	if err != nil {
		respond(w, http.StatusInternalServerError, "Could not fetch weather details")
		return
	}

	res.Greeting = fmt.Sprintf("Hello, %s!, the temperature is %.2f degrees Celsius in %s",
		name, weather.Current.Temperature, ip2l.City)

	_ = respond(w, http.StatusOK, res)
}

func respond(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}

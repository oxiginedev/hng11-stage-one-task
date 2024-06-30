package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Response struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

func main() {
	port := os.Getenv("HNG_PORT")
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
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Something went wrong"))
	}

	res.ClientIP = ip

	// Get api key from environment
	ip2key := os.Getenv("HNG_IP2LOCATION_KEY")

	ip2l, err := GetLocationFromIP(ip2key, ip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Something went wrong"))
	}

	res.Location = ip2l.City
	res.Greeting = fmt.Sprintf("Hello, %s!, the temperature is 11 degrees Celsius in %s", name, ip2l.City)

	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Something went wrong"))
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type M map[string]interface{}

type Response struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9000"
	}

	h := http.NewServeMux()

	h.HandleFunc("GET /api/hello", HandleIncomingRequest)

	s := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           h,
		ReadHeaderTimeout: time.Second * 2,
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 15,
	}

	log.Printf("Listening on :%s...\n", port)

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("server is attempting shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("server failed to shutdown")
	}

	log.Println("server exited!")
}

func HandleIncomingRequest(w http.ResponseWriter, r *http.Request) {
	res := &Response{}

	// get the visitor name parameter from url
	name := r.URL.Query().Get("visitor_name")

	ip, err := GetIpAddress(r)
	if err != nil {
		respond(w, http.StatusInternalServerError, M{"message": err.Error()})
		return
	}

	res.ClientIP = ip

	// Get api key from environment
	wk, ok := os.LookupEnv("WEATHERAPI_KEY")
	if !ok {
		log.Println("WEATHERAPI_KEY must be set")
		respond(w, http.StatusInternalServerError, M{"message": "Weather API key not set"})
		return
	}

	wc, err := NewWeatherAPIClient(wk)
	if err != nil {
		respond(w, http.StatusInternalServerError, M{"message": err.Error()})
		return
	}

	l, err := wc.GetLocationFromIP(ip)
	if err != nil {
		respond(w, http.StatusInternalServerError, M{"message": err.Error()})
		return
	}

	res.Location = l.City

	cw, err := wc.GetCurrentWeather(l.Latitude, l.Longitude)
	if err != nil {
		respond(w, http.StatusInternalServerError, M{"message": err.Error()})
		return
	}

	res.Greeting = fmt.Sprintf("Hello, %s!, the temperature is %.2f degrees Celsius in %s",
		name, cw.Current.TempC, l.City)

	_ = respond(w, http.StatusOK, res)
}

func respond(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}

package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// TemperatureResponse matches what smart_home TemperatureService expects.
type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func locationFromSensorID(sensorID string) string {
	switch sensorID {
	case "1":
		return "Living Room"
	case "2":
		return "Bedroom"
	case "3":
		return "Kitchen"
	default:
		return "Unknown"
	}
}

func sensorIDFromLocation(location string) string {
	switch location {
	case "Living Room":
		return "1"
	case "Bedroom":
		return "2"
	case "Kitchen":
		return "3"
	default:
		return "0"
	}
}

func randomTemperature() float64 {
	// Random value between 15.0 and 30.0, one decimal place.
	return float64(150+rand.Intn(151)) / 10.0
}

func buildResponse(location, sensorID string) TemperatureResponse {
	if location == "" {
		location = locationFromSensorID(sensorID)
	}
	if sensorID == "" {
		sensorID = sensorIDFromLocation(location)
	}

	return TemperatureResponse{
		Value:       randomTemperature(),
		Unit:        "°C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Simulated remote temperature sensor",
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func handleTemperature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/temperature")
	path = strings.Trim(path, "/")

	location := r.URL.Query().Get("location")
	sensorID := ""

	if path != "" {
		sensorID = path
	}

	if location == "" && sensorID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "location query or sensor id path is required",
		})
		return
	}

	writeJSON(w, http.StatusOK, buildResponse(location, sensorID))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8081"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/temperature", handleTemperature)
	mux.HandleFunc("/temperature/", handleTemperature)

	log.Printf("temperature-api listening on %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

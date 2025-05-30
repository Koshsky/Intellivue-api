package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	intellivue "github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

// ApiServer handles HTTP requests and returns data from ComputerClient
type ApiServer struct {
	port           string
	computerClient *intellivue.ComputerClient
}

// NewApiServer creates a new instance of ApiServer
func NewApiServer(port string, client *intellivue.ComputerClient) *ApiServer {
	log.Printf("Initializing API Server on port %s", port)
	return &ApiServer{
		port:           port,
		computerClient: client,
	}
}

// Run starts the HTTP server
func (s *ApiServer) Run() error {
	// Register routes
	http.HandleFunc("/api/vitals", s.handleVitals)
	http.HandleFunc("/api/waveforms", s.handleWaveforms)
	http.HandleFunc("/api/alarms", s.handleAlarms)

	log.Printf("API Server starting on port %s", s.port)
	log.Println("Registered endpoints:")
	log.Println("  - GET /api/vitals")
	log.Println("  - GET /api/waveforms")
	log.Println("  - GET /api/alarms")

	// Start server
	addr := fmt.Sprintf(":%s", s.port)
	return http.ListenAndServe(addr, nil)
}

// handleVitals returns the latest vital signs
func (s *ApiServer) handleVitals(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for /api/vitals from %s", r.Method, r.RemoteAddr)

	if r.Method != http.MethodGet {
		log.Printf("Method %s not allowed for /api/vitals", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, exists := s.computerClient.GetLatestData("vitals")
	if !exists {
		log.Println("No vital signs data available")
		http.Error(w, "No vital signs data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data.Value); err != nil {
		log.Printf("Error encoding vital signs response: %v", err)
	} else {
		log.Println("Successfully sent vital signs data")
	}
}

// handleWaveforms returns the latest waveform data
func (s *ApiServer) handleWaveforms(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for /api/waveforms from %s", r.Method, r.RemoteAddr)

	if r.Method != http.MethodGet {
		log.Printf("Method %s not allowed for /api/waveforms", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, exists := s.computerClient.GetLatestData("waveforms")
	if !exists {
		log.Println("No waveform data available")
		http.Error(w, "No waveform data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data.Value); err != nil {
		log.Printf("Error encoding waveforms response: %v", err)
	} else {
		log.Println("Successfully sent waveforms data")
	}
}

// handleAlarms returns the latest alarm states
func (s *ApiServer) handleAlarms(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for /api/alarms from %s", r.Method, r.RemoteAddr)

	if r.Method != http.MethodGet {
		log.Printf("Method %s not allowed for /api/alarms", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, exists := s.computerClient.GetLatestData("alarms")
	if !exists {
		log.Println("No alarm data available")
		http.Error(w, "No alarm data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data.Value); err != nil {
		log.Printf("Error encoding alarms response: %v", err)
	} else {
		log.Println("Successfully sent alarms data")
	}
}

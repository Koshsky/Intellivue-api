package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/api"
	intellivue "github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func init() {
	// Configure logging with timestamp
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[Intellivue] ")
}

func main() {
	log.Println("Starting Intellivue API application...")

	// Configuration
	host := "host1"       // Replace with actual host
	clientPort := "port1" // Replace with actual port
	apiPort := "8989"     // Replace with actual port

	log.Printf("Configuration loaded: host=%s, clientPort=%s, apiPort=%s", host, clientPort, apiPort)

	// Create ComputerClient
	log.Println("Initializing ComputerClient...")
	client := intellivue.NewComputerClient(host, clientPort)

	// Create API Server
	log.Println("Initializing API Server...")
	server := api.NewApiServer(apiPort, client)

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start ComputerClient
	log.Println("Starting ComputerClient...")
	client.Connect(ctx)
	log.Println("ComputerClient started successfully")

	// Start API server in a goroutine
	log.Println("Starting API Server...")
	go func() {
		if err := server.Run(); err != nil {
			log.Printf("API server error: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Application is ready! Press Ctrl+C to shutdown")

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("Received signal %v, initiating shutdown...", sig)

	// Cancel context to stop all goroutines
	cancel()

	// Give some time for graceful shutdown
	log.Println("Waiting for all goroutines to finish...")
	time.Sleep(2 * time.Second)

	log.Println("Application shutdown complete")
}

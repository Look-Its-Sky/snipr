package main

import (
	"os"
	"log"
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

func auth() (*ethclient.Client) {
	wsNodeURL := os.Getenv("NODE_URL_WSS")
	if wsNodeURL == "" {
		log.Fatalln("NODE_URL environment variable is not set. This should be your WebSocket endpoint (e.g., wss://...). ")
	}

	httpNodeURL := os.Getenv("NODE_URL_HTTP")
	if httpNodeURL == "" {
		log.Println("Warning: NODE_URL_HTTP is not set. Falling back to NODE_URL for RPC calls. This may not work with all node providers.")
		httpNodeURL = wsNodeURL
	}

	ethClient, err := ethclient.Dial(wsNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket endpoint (%s): %v", wsNodeURL, err)
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}
	log.Printf("Chain ID: %s", chainID.String())

	// Failure
	if err != nil {
		panic("Error connecting to either / both clients! Check logs")
	}

	// Success
	log.Println("Connected to WebSocket endpoint: " + wsNodeURL)

	return ethClient 
}
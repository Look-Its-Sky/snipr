package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
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

	rpcClient, err := rpc.DialContext(context.Background(), httpNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to HTTP endpoint (%s): %v", httpNodeURL, err)
	}
	log.Println("Connected to WebSocket endpoint: " + wsNodeURL)
	log.Println("Connected to HTTP RPC endpoint: " + httpNodeURL)

	headers := make(chan *types.Header)
	sub, err := ethClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new heads, please ensure your NODE_URL (%s) is a WebSocket endpoint: %v", wsNodeURL, err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)

		case header := <-headers:
			block, err := ethClient.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("Error getting block by hash: %v", err)
				continue
			}

			// New block
			log.Println("----------------------------------------")
			log.Printf("📦 New Block!\n")
			log.Printf("- Number: %s\n", header.Number.String())
			log.Printf("- Hash: %s\n", header.Hash().Hex())
			log.Printf("- Miner: %s\n", header.Coinbase.Hex())

			// Scrape block for each transaction
			for _, tx := range block.Transactions() {
				scrape(tx, ethClient, rpcClient)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

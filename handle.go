package main 

import (
	"log"
	"sync"
	"time"
	"context"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"snipr/schemas"
)

func listenForPools(exchange *schemas.Exchange, wg *sync.WaitGroup, client *ethclient.Client) {
	defer wg.Done()
	log.Printf("Connecting to exchange at %s via WebSocket (%s)", exchange.Address, exchange.WssURL)

	contractAddress := common.HexToAddress(exchange.Address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Printf("Failed to subscribe to logs for exchange %s: %v", exchange.Address, err)
		return
	}
	defer sub.Unsubscribe()

	contractAbi, err := abi.JSON(strings.NewReader(exchange.ABI))
	if err != nil {
		log.Printf("Failed to parse ABI for exchange %s: %v", exchange.Address, err)
		return
	}

	// TODO: find a more robost way to do this
	var eventName string
	if _, ok := contractAbi.Events["PoolCreated"]; ok {
		eventName = "PoolCreated"
	} else if _, ok := contractAbi.Events["PairCreated"]; ok {
		eventName = "PairCreated"
	} else {
		log.Printf("No 'PoolCreated' or 'PairCreated' event found in ABI for %s", exchange.Address)
		return
	}

	eventID := contractAbi.Events[eventName].ID
	log.Printf("Listening for %s events on %s", eventName, exchange.Address)

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription error for exchange %s: %v", exchange.Address, err)
			return

		case vLog := <-logs:
			if len(vLog.Topics) > 0 && vLog.Topics[0] == eventID {
				exchange.Process(vLog, contractAbi, eventName)
			}
			time.Sleep(5000 * time.Millisecond) // Avoid rate limiting
		}
	}
}

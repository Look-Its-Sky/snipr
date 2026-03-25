package main 

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"snipr/schemas"
)

func listenForPools(exchange *schemas.Exchange, wg *sync.WaitGroup, client *ethclient.Client) {
	defer wg.Done()

	contractAddress := common.HexToAddress(exchange.Address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	contractAbi, err := abi.JSON(strings.NewReader(exchange.ABI))
	if err != nil {
		log.Printf("Failed to parse ABI for exchange %s: %v", exchange.Address, err)
		return
	}

	var eventName string
	if _, ok := contractAbi.Events["PairCreated"]; ok {
		eventName = "PairCreated" // Uniswap V2
	} else if _, ok := contractAbi.Events["PoolCreated"]; ok {
		eventName = "PoolCreated" // Uniswap V3
	} else if _, ok := contractAbi.Events["Initialize"]; ok {
		eventName = "Initialize" // Uniswap V4
	} else {
		log.Printf("No 'PoolCreated' or 'PairCreated' event found in ABI for %s", exchange.Address)
		return
	}

	eventID := contractAbi.Events[eventName].ID
	log.Printf("Listening for %s events on contract: %s", eventName, exchange.Address)

	// wss reconnection loop
	for {
		logs := make(chan types.Log)
		sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			log.Printf("Failed to subscribe to logs for exchange %s: %v. Retrying in 5s...", exchange.Address, err)
			time.Sleep(5 * time.Second)
			continue // resubscribe
		}

		// process events
		func() {
			defer sub.Unsubscribe()
			for {
				select {
				case err := <-sub.Err():
					log.Printf("Subscription dropped for exchange %s: %v. Reconnecting...", exchange.Address, err)
					return 

				case vLog := <-logs:
					if len(vLog.Topics) == 0 || vLog.Topics[0] != eventID {
						continue 
					}

					contract, err := exchange.Process(vLog, contractAbi, eventName)
					if err != nil {
						log.Printf("Error processing log for exchange %s: %v", exchange.Address, err)
						continue 
					}

					if !*disableDB { go pushNewContract(contract) }
				}
			}
		}()

		// avoid spamming node on reconnect
		time.Sleep(2 * time.Second)
	}
}

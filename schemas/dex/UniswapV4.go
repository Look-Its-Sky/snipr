package dex

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"snipr/schemas"

// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// )

// // Listen for 'PoolCreated'
// func UniswapV4() schemas.Exchange {
// 	process := func(vLog types.Log, contractAbi abi.ABI, eventName string) { 
// 		var eventData struct {
// 			TickSpacing int32
// 			Hooks       common.Address
// 			Currency0   common.Address
// 			Currency1   common.Address
// 			Fee         uint64
// 		}
		
// 		err := contractAbi.UnpackIntoInterface(&eventData, eventName, vLog.Data)
// 		if err != nil {
// 			log.Printf("Failed to unpack V4 %s\nevent data: %v", eventName, err)
// 		}

// 		// fee := vLog.Topics[3].Big().Uint64()

// 		fmt.Printf("V4 Pool Details: [Token0: %s, Token1: %s, Fee: %d, TickSpacing: %d, Hooks: %s]\n",
// 			eventData.Currency0.Hex(), eventData.Currency1.Hex(), eventData.Fee, eventData.TickSpacing, eventData.Hooks.Hex())

// 		// Database logic would go here
// 	}

// 	return schemas.Exchange{
// 		Address: "0x000000000004444c5dc75cB358380D2e3dE08A90",
// 		ABI:     `[{"inputs":[{"internalType":"address","name":"initialOwner","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"currency0","type":"address"},{"indexed":true,"internalType":"address","name":"currency1","type":"address"},{"indexed":true,"internalType":"uint24","name":"fee","type":"uint24"},{"indexed":false,"internalType":"int24","name":"tickSpacing","type":"int24"},{"indexed":false,"internalType":"address","name":"hooks","type":"address"}],"name":"PoolCreated","type":"event"}]`,
// 		WssURL:  os.Getenv("NODE_URL_WSS"),
// 		HttpURL: os.Getenv("NODE_URL_HTTP"),
// 		Process: process,
// 	}
// }

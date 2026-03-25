package dex

// TODO: check this shit 

import (
	"log"

	"snipr/schemas"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Listen for 'PairCreated' on PancakeSwap V2
func PancakeSwapV2(disableDB *bool) *schemas.Exchange {
	process := func(vLog types.Log, contractAbi abi.ABI, eventName string) (*schemas.Contract, error) {
		created_coin := common.HexToAddress(vLog.Topics[1].Hex())
		backing_coin := common.HexToAddress(vLog.Topics[2].Hex())

		log.Printf("Token created on PancakeSwap V2 -\nCreated Coin: %s\nBacking Coin: %s\n",
			created_coin.Hex(),
			backing_coin.Hex(),
		)

		c := schemas.Contract{
			Address:            created_coin.Hex(),
			BackingCoinAddress: backing_coin.Hex(),
			Exchange:						"PancakeSwapV2",
			BlockNumber:				vLog.BlockNumber,
		}

		return &c, nil
	}

	return &schemas.Exchange{
		Name: "PancakeSwapV2",
		// Official PancakeSwap V2 Factory on BSC
		Address: "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73", 
		
		// Lightweight ABI containing ONLY the 'PairCreated' event
		ABI:     `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"token0","type":"address"},{"indexed":true,"internalType":"address","name":"token1","type":"address"},{"indexed":false,"internalType":"address","name":"pair","type":"address"},{"indexed":false,"internalType":"uint256","name":"","type":"uint256"}],"name":"PairCreated","type":"event"}]`,
		Process: process,
	}
}

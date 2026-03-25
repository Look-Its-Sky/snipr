package dex

// TODO: check this shit

import (
	"log"

	"snipr/schemas"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Listen for 'PoolCreated' on PancakeSwap V3
func PancakeSwapV3(disableDB *bool) *schemas.Exchange {
	process := func(vLog types.Log, contractAbi abi.ABI, eventName string) (*schemas.Contract, error) {
		created_coin := common.HexToAddress(vLog.Topics[1].Hex())
		backing_coin := common.HexToAddress(vLog.Topics[2].Hex())

		log.Printf("Token created on PancakeSwap V3 -\nCreated Coin: %s\nBacking Coin: %s\n",
			created_coin.Hex(),
			backing_coin.Hex(),
		)

		c := schemas.Contract{
			Address:            created_coin.Hex(),
			BackingCoinAddress: backing_coin.Hex(),
			Exchange:						"PancakeSwapV3",
			BlockNumber:				vLog.BlockNumber,
		}

		return &c, nil
	}

	return &schemas.Exchange{
		Name: "PancakeSwapV3",

		// Official PancakeSwap V3 Factory on BSC
		Address: "0x0BFbCF9fa4f9C56B0F40a671Ad40E0805A091865", 
		
		// Lightweight ABI containing ONLY the 'PoolCreated' event
		ABI:     `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"token0","type":"address"},{"indexed":true,"internalType":"address","name":"token1","type":"address"},{"indexed":true,"internalType":"uint24","name":"fee","type":"uint24"},{"indexed":false,"internalType":"int24","name":"tickSpacing","type":"int24"},{"indexed":false,"internalType":"address","name":"pool","type":"address"}],"name":"PoolCreated","type":"event"}]`,
		Process: process,
	}
}

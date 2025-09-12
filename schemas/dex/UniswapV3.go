package dex

import (
	"log"
	"math/big"

	"snipr/schemas"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Listen for 'PoolCreated'
func UniswapV3(disableDB *bool) *schemas.Exchange {
	process := func(vLog types.Log, contractAbi abi.ABI, eventName string) (*schemas.Contract, error) {
		created_coin := common.HexToAddress(vLog.Topics[1].Hex())
		backing_coin := common.HexToAddress(vLog.Topics[2].Hex())

		// Used for validation
		var pairCreated struct {
			Pair           common.Address
			AllPairsLength *big.Int
		}

		err := contractAbi.UnpackIntoInterface(&pairCreated, eventName, vLog.Data)
		if err != nil {
			log.Printf("Uniswap V3: Failed to unpack PairCreated event data: %v", err)
			return nil, err
		}

		log.Printf("Token created on Uniswap V3 -\nCreated Coin: %s\nBacking Coin: %s\n",
			created_coin.Hex(),
			backing_coin.Hex(),
		)

		c := schemas.Contract{
			Address:            created_coin.Hex(),
			BackingCoinAddress: backing_coin.Hex(),
		}

		return &c, nil
	}

	return &schemas.Exchange{
		Address: "0x1F98431c8aD98523631AE4a59f267346ea31F984", // Uniswap V3 factory address
		ABI:     `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint24","name":"fee","type":"uint24"},{"indexed":true,"internalType":"int24","name":"tickSpacing","type":"int24"}],"name":"FeeAmountEnabled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"oldOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnerChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"token0","type":"address"},{"indexed":true,"internalType":"address","name":"token1","type":"address"},{"indexed":true,"internalType":"uint24","name":"fee","type":"uint24"},{"indexed":false,"internalType":"int24","name":"tickSpacing","type":"int24"},{"indexed":false,"internalType":"address","name":"pool","type":"address"}],"name":"PoolCreated","type":"event"},{"inputs":[{"internalType":"address","name":"tokenA","type":"address"},{"internalType":"address","name":"tokenB","type":"address"},{"internalType":"uint24","name":"fee","type":"uint24"}],"name":"createPool","outputs":[{"internalType":"address","name":"pool","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint24","name":"fee","type":"uint24"},{"internalType":"int24","name":"tickSpacing","type":"int24"}],"name":"enableFeeAmount","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint24","name":"","type":"uint24"}],"name":"feeAmountTickSpacing","outputs":[{"internalType":"int24","name":"","type":"int24"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"},{"internalType":"uint24","name":"","type":"uint24"}],"name":"getPool","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"parameters","outputs":[{"internalType":"address","name":"factory","type":"address"},{"internalType":"address","name":"token0","type":"address"},{"internalType":"address","name":"token1","type":"address"},{"internalType":"uint24","name":"fee","type":"uint24"},{"internalType":"int24","name":"tickSpacing","type":"int24"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_owner","type":"address"}],"name":"setOwner","outputs":[],"stateMutability":"nonpayable","type":"function"}]`,
		Process: process,
	}
}

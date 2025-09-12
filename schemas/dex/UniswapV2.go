package dex

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"snipr/schemas"
)

// Listens for 'PairCreated'
func UniswapV2(disableDB *bool) *schemas.Exchange {
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
			log.Printf("UniswapV2: Failed to unpack PairCreated event data: %v", err)
			return nil, err
		}

		log.Printf("Token created on Uniswap V2 -\nCreated Coin: %s\nBacking Coin: %s\n",
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
		Address: "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f", // Uniswap V2 Factory Address
		ABI:     `[{"inputs":[{"internalType":"address","name":"_feeToSetter","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"token0","type":"address"},{"indexed":true,"internalType":"address","name":"token1","type":"address"},{"indexed":false,"internalType":"address","name":"pair","type":"address"},{"indexed":false,"internalType":"uint256","name":"allPairsLength","type":"uint256"}],"name":"PairCreated","type":"event"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"allPairs","outputs":[{"internalType":"address","name":"pair","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"allPairsLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"tokenA","type":"address"},{"internalType":"address","name":"tokenB","type":"address"}],"name":"createPair","outputs":[{"internalType":"address","name":"pair","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeTo","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeToSetter","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"getPair","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_feeTo","type":"address"}],"name":"setFeeTo","stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_feeToSetter","type":"address"}],"name":"setFeeToSetter","stateMutability":"nonpayable","type":"function"}]`,
		Process: process,
	}
}

package dex

import (
	"log"
	_ "math/big"

	"snipr/schemas"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Listen for 'Initialize' 
func UniswapV4(disableDB *bool) *schemas.Exchange {
	process := func(vLog types.Log, contractAbi abi.ABI, eventName string) (*schemas.Contract, error) {
		// V4 Paradigm Shift: Topics[1] is the Pool ID (bytes32), NOT a contract address!
		poolId := vLog.Topics[1].Hex()
		
		// Topics[2] and Topics[3] are the two tokens making up the pool
		currency0 := common.HexToAddress(vLog.Topics[2].Hex())
		currency1 := common.HexToAddress(vLog.Topics[3].Hex())

		// ========================================================================
		// HIGHLY RECOMMENDED: Unpack the data to check for malicious Hooks!
		// If 'Hooks' is not the zero address (0x000...000), the pool has custom 
		// logic attached that could restrict selling or honeypot your bot.
		// ========================================================================
		//
		// var initData struct {
		// 	Fee          uint32
		// 	TickSpacing  int32
		// 	Hooks        common.Address
		// 	SqrtPriceX96 *big.Int
		// 	Tick         int32
		// }
		// err := contractAbi.UnpackIntoInterface(&initData, eventName, vLog.Data)
		// if err != nil {
		// 	log.Printf("Uniswap V4: Failed to unpack Initialize event data: %v", err)
		// 	return nil, err
		// }

		log.Printf("Pool initialized on Uniswap V4 -\nPool ID: %s\nCurrency0: %s\nCurrency1: %s\n",
			poolId,
			currency0.Hex(),
			currency1.Hex(),
		)

		c := schemas.Contract{
			Address:            currency0.Hex(), // new token
			BackingCoinAddress: currency1.Hex(), // backing coin
			Exchange:						"UniswapV4",
			BlockNumber:				vLog.BlockNumber,
		}

		return &c, nil
	}

	return &schemas.Exchange{
		Name: "UniswapV4",
		Address: "0x28e2ea090877bf75740558f6bfb36a5ffee9e9df", 
		
		// Lightweight ABI containing ONLY the 'Initialize' event
		ABI:     `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"PoolId","name":"id","type":"bytes32"},{"indexed":true,"internalType":"Currency","name":"currency0","type":"address"},{"indexed":true,"internalType":"Currency","name":"currency1","type":"address"},{"indexed":false,"internalType":"uint24","name":"fee","type":"uint24"},{"indexed":false,"internalType":"int24","name":"tickSpacing","type":"int24"},{"indexed":false,"internalType":"contract IHooks","name":"hooks","type":"address"},{"indexed":false,"internalType":"uint160","name":"sqrtPriceX96","type":"uint160"},{"indexed":false,"internalType":"int24","name":"tick","type":"int24"}],"name":"Initialize","type":"event"}]`,
		Process: process,
	}
}

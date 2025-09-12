package schemas

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Exchange struct {
	Address 		string 
	ABI  			string
	WssURL  		string
	HttpURL 		string
	Process			func(vLog types.Log, contractAbi abi.ABI, eventName string) (*Contract, error)
}

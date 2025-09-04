package main

import (
	"gorm.io/gorm"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// For stores
type Contract struct {
	gorm.Model 				// ID, CreatedAt, UpdatedAt, DeletedAt
	ContractAddress 	string `gorm:"uniqueIndex;not null"`
	CreatorAddress 		string `gorm:"index"`
	TransactionHash 	string `gorm:"uniqueIndex;not null"`
	BlockNumber 			uint64 
	Symbol						string
	ByteCode          []byte
	TotalSupply			  uint64	
	Decimals					uint8
	Blacklisted       bool
}

// For validation
type ERC20Token struct {
	Address     common.Address
	Name        string
	Symbol      string
	TotalSupply *big.Int
	Decimals    uint8
}

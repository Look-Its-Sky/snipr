package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// ERC20_ABI is the minimal Application Binary Interface needed to get token details.
const ERC20_ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"}]`

// essentialSelectors are the 4-byte function selectors for key ERC-20 functions.
var essentialSelectors = map[string]string{
	"totalSupply()":      "18160ddd",
	"balanceOf(address)": "70a08231",
	"transfer(address,uint256)": "a9059cbb",
}

// This is the **reliable** way to find contracts, as it includes internal transactions (contract-deploys-contract).
func findCreatedContracts(rpcClient *rpc.Client, txHash common.Hash) ([]common.Address, error) {
	var trace map[string]interface{}
	var err error

	//  provider requires 'callTracer' instead of the default 'structTracer'.
	tracerConfig := map[string]string{"tracer": "callTracer"}

	// Retry logic for intermittent node errors
	for i := 0; i < 3; i++ {
		err = rpcClient.CallContext(context.Background(), &trace, "debug_traceTransaction", txHash.Hex(), tracerConfig)
		if err == nil {
			break
		}
		log.Printf("Error calling debug_traceTransaction (attempt %d): %v. Retrying in %d seconds...", i+1, err, (i+1)*2)
		time.Sleep(time.Duration((i+1)*2) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("error calling debug_traceTransaction with callTracer after multiple retries: %w", err)
	}

	var createdContracts []common.Address
	findCreatesInCallTrace(&createdContracts, trace)
	return createdContracts, nil
}

// findCreatesInCallTrace recursively walks a 'callTracer' result to find CREATE/CREATE2 calls.
func findCreatesInCallTrace(contracts *[]common.Address, callFrame map[string]interface{}) {
	// Check if the current frame is a contract creation
	if callType, ok := callFrame["type"].(string); ok {
		if callType == "CREATE" || callType == "CREATE2" {
			// For CREATE/CREATE2, the 'to' field holds the new contract address
			if to, ok := callFrame["to"].(string); ok {
				*contracts = append(*contracts, common.HexToAddress(to))
			}
		}
	}

	// Recursively check nested calls
	if calls, ok := callFrame["calls"].([]interface{}); ok {
		for _, call := range calls {
			if nextFrame, ok := call.(map[string]interface{}); ok {
				findCreatesInCallTrace(contracts, nextFrame)
			}
		}
	}
}

// isERC20Token checks if a contract's bytecode contains the essential ERC-20 function selectors.
func isERC20Token(client *ethclient.Client, address common.Address) (bool, error) {
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil for latest block
	if err != nil {
		return false, fmt.Errorf("failed to get bytecode for address %s: %w", address.Hex(), err)
	}
	hexBytecode := hex.EncodeToString(bytecode)

	// Check for the presence of essential function selectors.
	for _, selector := range essentialSelectors {
		if !strings.Contains(hexBytecode, selector) {
			// If even one essential function is missing, it's likely not a standard token.
			return false, nil
		}
	}

	return true, nil
}

//getTokenDetails retrieves metadata (name, symbol, etc.) from a token contract.
func getTokenDetails(client *ethclient.Client, tokenAddress common.Address) (*ERC20Token, error) {
	parsedABI, err := abi.JSON(strings.NewReader(ERC20_ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	token := &ERC20Token{Address: tokenAddress}

	// Call the `name()` function
	nameData, err := parsedABI.Pack("name")
	if err != nil {
		return nil, err
	}
	nameResult, err := client.CallContract(context.Background(), newCallMsg(tokenAddress, nameData), nil)
	if err != nil {
		// Some tokens might not have a name function 
		token.Name = "N/A"
	} else {
		err = parsedABI.UnpackIntoInterface(&token.Name, "name", nameResult)
		if err != nil {
			token.Name = "Unparsable"
		}
	}

	// Call `symbol()` 
	symbolData, err := parsedABI.Pack("symbol")
	if err != nil {
		return nil, err
	}
	symbolResult, err := client.CallContract(context.Background(), newCallMsg(tokenAddress, symbolData), nil)
	if err != nil {
		token.Symbol = "N/A"
	} else {
		err = parsedABI.UnpackIntoInterface(&token.Symbol, "symbol", symbolResult)
		if err != nil {
			token.Symbol = "Unparsable"
		}
	}

	// Call the `totalSupply()` function
	totalSupplyData, err := parsedABI.Pack("totalSupply")
	if err != nil {
		return nil, err
	}
	totalSupplyResult, err := client.CallContract(context.Background(), newCallMsg(tokenAddress, totalSupplyData), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call totalSupply: %w", err)
	}
	token.TotalSupply = new(big.Int)
	err = parsedABI.UnpackIntoInterface(&token.TotalSupply, "totalSupply", totalSupplyResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack totalSupply: %w", err)
	}

	// Call `decimals()`
	decimalsData, err := parsedABI.Pack("decimals")
	if err != nil {
		return nil, err
	}
	decimalsResult, err := client.CallContract(context.Background(), newCallMsg(tokenAddress, decimalsData), nil)
	if err != nil {
		token.Decimals = 18 // Default to 18
	} else {
		err = parsedABI.UnpackIntoInterface(&token.Decimals, "decimals", decimalsResult)
		if err != nil {
			token.Decimals = 18 // Default to 18 if unparsable
		}
	}

	return token, nil
}

func newCallMsg(to common.Address, data []byte) ethereum.CallMsg {
	return ethereum.CallMsg{
		To:   &to,
		Data: data,
	}
}

func scrape(tx *types.Transaction, ethClient *ethclient.Client, rpcClient *rpc.Client) {
	log.Println("----------------------------------------")
	log.Printf("🔍 Analyzing transaction: %s\n", tx.Hash())
	log.Println("----------------------------------------")

	txHash := tx.Hash()

	contracts, err := findCreatedContracts(rpcClient, txHash)
	if err != nil {
		log.Printf("Failed to find created contracts!\n%v", err)
		return
	}

	if len(contracts) == 0 {
		log.Println("No new contracts were created in this transaction.")
		return
	}

	for _, addr := range contracts {
		log.Printf("Analyzing %s...\n", addr.Hex())

		isToken, err := isERC20Token(ethClient, addr)
		if err != nil {
			log.Printf("Could not check if %s is a token!\n%v", addr.Hex(), err)
			continue
		}

		// Not token... skip
		if !isToken {
			log.Printf("%s does not appear to be an ERC20 token\n", addr.Hex())
			continue
		}

		// Found token
		log.Printf("📝 It's an ERC20 token! Fetching details...\n")
		tokenDetails, err := getTokenDetails(ethClient, addr)
		if err != nil {
			log.Printf("    Could not get token details: %v", err)
			continue
		}

		log.Printf("    - Name: %s\n", tokenDetails.Name)
		log.Printf("    - Symbol: %s\n", tokenDetails.Symbol)
		log.Printf("    - Total Supply: %s\n", tokenDetails.TotalSupply.String())
		log.Printf("    - Decimals: %d\n", tokenDetails.Decimals)

		// Get additional details for the Contract struct
		receipt, err := ethClient.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			log.Printf("Failed to get transaction receipt: %v", err)
			continue
		}

		from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			log.Printf("Failed to get sender from transaction: %v", err)
			continue
		}

		bytecode, err := ethClient.CodeAt(context.Background(), addr, nil)
		if err != nil {
			log.Printf("Failed to get bytecode: %v", err)
			continue
		}

		// Dump token to database
		c := Contract{
			ContractAddress: addr.Hex(),
			CreatorAddress:  from.Hex(),
			TransactionHash: txHash.Hex(),
			BlockNumber:     receipt.BlockNumber.Uint64(),
			Symbol:          tokenDetails.Symbol,
			ByteCode:        bytecode,
			TotalSupply:     tokenDetails.TotalSupply.Uint64(),
			Decimals:        tokenDetails.Decimals,
			Blacklisted:     false,
		}
		pushNew(c)
	}
}

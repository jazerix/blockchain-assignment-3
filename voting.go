package main

import (
	"log"
	"voting/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	partyChainCode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := partyChainCode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}

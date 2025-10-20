package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SettlementEngineContract struct {
	contractapi.Contract
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SettlementEngineContract{})
	if err != nil {
		panic(err)
	}
	if err := chaincode.Start(); err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type EnergyTokenContract struct {
	contractapi.Contract
}

func main() {
	chaincode, err := contractapi.NewChaincode(&EnergyTokenContract{})
	if err != nil {
		panic(err)
	}
	if err := chaincode.Start(); err != nil {
		panic(err)
	}
}

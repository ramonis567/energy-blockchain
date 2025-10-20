package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CreditMarketContract struct {
	contractapi.Contract
}

func main() {
	chaincode, err := contractapi.NewChaincode(&CreditMarketContract{})
	if err != nil {
		panic(err)
	}
	if err := chaincode.Start(); err != nil {
		panic(err)
	}
}

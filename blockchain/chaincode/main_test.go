package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)

func TestCloseContract_ReturnEscrow(t *testing.T) {
	// Initialize Chaincode
	cc, err := contractapi.NewChaincode(&CombinedEnergyContract{})
	assert.NoError(t, err)

	stub := shimtest.NewMockStub("energy_cc", cc)

	// Register Agents
	res := stub.MockInvoke("1", [][]byte{[]byte("RegisterAgent"), []byte("seller"), []byte("producer"), []byte("Seller"), []byte("Address1")})
	assert.Equal(t, int32(200), res.Status, "Register seller failed: %s", res.Message)

	res = stub.MockInvoke("2", [][]byte{[]byte("RegisterAgent"), []byte("buyer"), []byte("consumer"), []byte("Buyer"), []byte("Address2")})
	assert.Equal(t, int32(200), res.Status, "Register buyer failed: %s", res.Message)

	// Mint Tokens
	res = stub.MockInvoke("3", [][]byte{[]byte("Mint"), []byte("seller"), []byte("ECR"), []byte("1000")})
	assert.Equal(t, int32(200), res.Status, "Mint ECR failed: %s", res.Message)

	res = stub.MockInvoke("4", [][]byte{[]byte("Mint"), []byte("buyer"), []byte("ENGT"), []byte("1000")})
	assert.Equal(t, int32(200), res.Status, "Mint ENGT failed: %s", res.Message)

	// Create Supply Contract
	// CreateSupplyContract(id, sellerID, buyerID, energyTotal, pricePerKWh, startDate, endDate, settlementFreq, sellerCollateralECR, buyerCollateralENGT)
	// energyTotal=100, price=1, sellerCol=0, buyerCol=50
	res = stub.MockInvoke("5", [][]byte{
		[]byte("CreateSupplyContract"),
		[]byte("contract1"),
		[]byte("seller"),
		[]byte("buyer"),
		[]byte("100"), // energyTotal
		[]byte("1"),   // price
		[]byte("2023-01-01"),
		[]byte("2023-12-31"),
		[]byte("MONTHLY"),
		[]byte("0"),  // sellerCollateral
		[]byte("50"), // buyerCollateral
	})
	assert.Equal(t, int32(200), res.Status, "CreateSupplyContract failed: %s", res.Message)

	// Verify Balance Decreased
	balBytes, err := stub.GetState("energytoken:balance:buyer")
	assert.NoError(t, err)
	var bal TokenBalance
	err = json.Unmarshal(balBytes, &bal)
	assert.NoError(t, err)
	assert.Equal(t, 950.0, bal.ENGT)

	// Verify Agent Record Balance (should be 950)
	agentBytes, err := stub.GetState("agentreg:agent:buyer")
	assert.NoError(t, err)
	var agent Agent
	err = json.Unmarshal(agentBytes, &agent)
	assert.NoError(t, err)
	assert.Equal(t, 950.0, agent.ENGTBalance)

	// Close Contract
	// CloseContract(id)
	res = stub.MockInvoke("6", [][]byte{[]byte("CloseContract"), []byte("contract1")})
	assert.Equal(t, int32(200), res.Status, "CloseContract failed: %s", res.Message)

	// Verify Escrow Returned
	balBytes, err = stub.GetState("energytoken:balance:buyer")
	assert.NoError(t, err)
	err = json.Unmarshal(balBytes, &bal)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, bal.ENGT, "Buyer balance should be restored to 1000")

	// Verify Agent Record Balance Updated
	agentBytes, err = stub.GetState("agentreg:agent:buyer")
	assert.NoError(t, err)
	err = json.Unmarshal(agentBytes, &agent)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, agent.ENGTBalance, "Buyer agent record should show 1000 ENGT")
}

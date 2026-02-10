package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func TestGetAllAgentsSynchronizesBalance(t *testing.T) {
	// Setup
	contract := new(CombinedEnergyContract)
	cc, err := contractapi.NewChaincode(contract)
	if err != nil {
		t.Fatalf("Error creating chaincode: %s", err)
	}
	stub := shimtest.NewMockStub("energy_contract", cc)

	// 1. Register Agent
	// RegisterAgent(ctx, id, type, name, address)
	res := stub.MockInvoke("1", [][]byte{[]byte("RegisterAgent"), []byte("producer1"), []byte("producer"), []byte("Producer One"), []byte("Address 1")})
	if res.Status != shim.OK {
		t.Fatalf("RegisterAgent failed: %s", res.Message)
	}

	// 2. Mint Tokens (Updates Balance and Agent)
	// Mint(ctx, agentID, tokenType, amountStr)
	res = stub.MockInvoke("2", [][]byte{[]byte("Mint"), []byte("producer1"), []byte("ECR"), []byte("100.0")})
	if res.Status != shim.OK {
		t.Fatalf("Mint failed: %s", res.Message)
	}

	// Verify Balance is 100 via GetAgent (which also syncs, but let's trust Mint for now)
	// We want to verify the state in the ledger is correct before we mess with it.
	// Balance Key: energytoken:balance:producer1
	balKey := "energytoken:balance:producer1"
	balBytes, err := stub.GetState(balKey)
	if err != nil {
		t.Fatalf("Failed to get balance state: %s", err)
	}
	var bal TokenBalance
	json.Unmarshal(balBytes, &bal)
	if bal.ECR != 100.0 {
		t.Fatalf("Ledger balance incorrect. Expected 100.0, got %f", bal.ECR)
	}

	// 3. Manually Corrupt Agent State (Stale Data)
	// We simulate that the Agent record was somehow not updated or is stale.
	agentKey := "agentreg:agent:producer1"
	// Create stale agent
	staleAgent := Agent{
		ID: "producer1",
		Type: "producer",
		Name: "Producer One",
		Address: "Address 1",
		ECRBalance: 50.0, // Stale!
		ENGTBalance: 0.0,
		RegisteredAt: "timestamp",
	}
	staleBytes, _ := json.Marshal(staleAgent)

	// MockStub.PutState needs to be in a transaction context if checking from outside?
	// Actually MockStub.PutState simply writes to state. But to emulate a "past" write, we just write it.
	// We wrap in MockTransactionStart/End just to be safe and emulate a commit.
	stub.MockTransactionStart("tx_corrupt")
	err = stub.PutState(agentKey, staleBytes)
	stub.MockTransactionEnd("tx_corrupt")
	if err != nil {
		t.Fatalf("Failed to corrupt state: %s", err)
	}

	// Verify it is corrupted
	corruptBytes, _ := stub.GetState(agentKey)
	var checkAgent Agent
	json.Unmarshal(corruptBytes, &checkAgent)
	if checkAgent.ECRBalance != 50.0 {
		t.Fatalf("Failed to corrupt agent state for test. It is %f", checkAgent.ECRBalance)
	}

	// 4. Call GetAllAgents
	res = stub.MockInvoke("3", [][]byte{[]byte("GetAllAgents")})
	if res.Status != shim.OK {
		t.Fatalf("GetAllAgents failed: %s", res.Message)
	}

	// 5. Verify the returned agent has the correct balance (100, not 50)
	var agents []*Agent
	err = json.Unmarshal(res.Payload, &agents)
	if err != nil {
		t.Fatalf("Failed to unmarshal agents: %s", err)
	}

	if len(agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(agents))
	}

	if agents[0].ECRBalance != 100.0 {
		t.Errorf("Expected ECRBalance 100.0, got %f. The balance was not synchronized!", agents[0].ECRBalance)
	} else {
		fmt.Println("Success: Agent balance was synchronized.")
	}
}

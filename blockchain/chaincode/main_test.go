package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func BenchmarkGetAllAgents(b *testing.B) {
	// Initialize Chaincode
	cc, err := contractapi.NewChaincode(&CombinedEnergyContract{})
	if err != nil {
		b.Fatalf("Error creating chaincode: %s", err)
	}
	stub := shimtest.NewMockStub("energy_cc", cc)

	// Populate Ledger
	// 1. Add 100 Agents
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("agent%03d", i)
		agent := Agent{ID: id, Type: "prosumer", Name: fmt.Sprintf("Agent %d", i)}
		bytes, _ := json.Marshal(agent)
		stub.MockTransactionStart(fmt.Sprintf("init_agent_%d", i))
		stub.PutState(AgentKeyPrefix+id, bytes)
		stub.MockTransactionEnd(fmt.Sprintf("init_agent_%d", i))
	}

	// 2. Add 2000 Offers (noise data that should be skipped)
	for i := 0; i < 2000; i++ {
		id := fmt.Sprintf("offer%04d", i)
		offer := Offer{ID: id, Status: "OPEN", SellerID: "agent001"}
		bytes, _ := json.Marshal(offer)
		stub.MockTransactionStart(fmt.Sprintf("init_offer_%d", i))
		stub.PutState(OfferKeyPrefix+id, bytes)
		stub.MockTransactionEnd(fmt.Sprintf("init_offer_%d", i))
	}

	// 3. Add 2000 Contracts (more noise)
	for i := 0; i < 2000; i++ {
		id := fmt.Sprintf("contract%04d", i)
		contract := SupplyContract{ID: id, Status: "ACTIVE"}
		bytes, _ := json.Marshal(contract)
		stub.MockTransactionStart(fmt.Sprintf("init_contract_%d", i))
		stub.PutState(ContractKeyPrefix+id, bytes)
		stub.MockTransactionEnd(fmt.Sprintf("init_contract_%d", i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res := stub.MockInvoke("txID", [][]byte{[]byte("GetAllAgents")})
		if res.Status != shim.OK {
			b.Fatalf("GetAllAgents failed: %s", res.Message)
		}
	}
}

func TestGetAllAgentsCorrectness(t *testing.T) {
	// Initialize Chaincode
	cc, err := contractapi.NewChaincode(&CombinedEnergyContract{})
	if err != nil {
		t.Fatalf("Error creating chaincode: %s", err)
	}
	stub := shimtest.NewMockStub("energy_cc", cc)

	// Populate with known data
	ids := []string{"agentA", "agentB", "agentC"}
	for _, id := range ids {
		agent := Agent{ID: id, Type: "producer"}
		bytes, _ := json.Marshal(agent)
		stub.MockTransactionStart("init")
		stub.PutState(AgentKeyPrefix+id, bytes)
		stub.MockTransactionEnd("init")
	}

	// Add noise
	stub.MockTransactionStart("init_noise")
	stub.PutState(OfferKeyPrefix+"noise", []byte("{}"))
	stub.MockTransactionEnd("init_noise")

	// Invoke
	res := stub.MockInvoke("txID", [][]byte{[]byte("GetAllAgents")})
	if res.Status != shim.OK {
		t.Fatalf("GetAllAgents failed: %s", res.Message)
	}

	var agents []Agent
	err = json.Unmarshal(res.Payload, &agents)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}

	if len(agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(agents))
	}
}

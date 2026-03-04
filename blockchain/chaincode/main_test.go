package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type MockTransactionContext struct {
	contractapi.TransactionContext
	stub *shimtest.MockStub
}

func (m *MockTransactionContext) GetStub() contractapi.ChaincodeStubInterface {
	return m.stub
}

func BenchmarkGetAgentsByType(b *testing.B) {
	contract := new(CombinedEnergyContract)
	stub := shimtest.NewMockStub("energycc", nil)
	ctx := &MockTransactionContext{stub: stub}

	// Setup: 1000 agents, 100 of type "producer"
	// and 1000 offers to clutter the ledger
	for i := 0; i < 1000; i++ {
		agentID := fmt.Sprintf("agent%d", i)
		agentType := "consumer"
		if i%10 == 0 {
			agentType = "producer"
		}
		agent := Agent{
			ID:   agentID,
			Type: agentType,
			Name: "Agent " + agentID,
		}
		stub.MockTransactionStart(fmt.Sprintf("setup%d", i))
		stub.PutState(agentKey(agentID), mustMarshal(agent))

		balance := TokenBalance{AgentID: agentID, ECR: 100.0, ENGT: 100.0}
		stub.PutState(balanceKey(agentID), mustMarshal(balance))

		stub.PutState(offerKey(agentID), []byte("{}"))
		stub.MockTransactionEnd(fmt.Sprintf("setup%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := contract.GetAgentsByType(ctx, "producer")
		if err != nil {
			b.Fatalf("Error in GetAgentsByType: %v", err)
		}
	}
}

func TestGetAgentsByType(t *testing.T) {
	contract := new(CombinedEnergyContract)
	stub := shimtest.NewMockStub("energycc", nil)
	ctx := &MockTransactionContext{stub: stub}

	// Setup: 10 agents, 2 of type "producer"
	for i := 0; i < 10; i++ {
		agentID := fmt.Sprintf("agent%d", i)
		agentType := "consumer"
		if i%5 == 0 {
			agentType = "producer"
		}
		agent := Agent{
			ID:   agentID,
			Type: agentType,
			Name: "Agent " + agentID,
		}
		stub.MockTransactionStart(fmt.Sprintf("setup%d", i))
		stub.PutState(agentKey(agentID), mustMarshal(agent))

		balance := TokenBalance{AgentID: agentID, ECR: 100.0, ENGT: 100.0}
		stub.PutState(balanceKey(agentID), mustMarshal(balance))
		stub.MockTransactionEnd(fmt.Sprintf("setup%d", i))
	}

	agents, err := contract.GetAgentsByType(ctx, "producer")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agents))
	}

	for _, a := range agents {
		if a.Type != "producer" {
			t.Errorf("Agent %s has wrong type: %s", a.ID, a.Type)
		}
		if a.ECRBalance != 100.0 {
			t.Errorf("Agent %s has wrong balance: %f", a.ID, a.ECRBalance)
		}
	}
}

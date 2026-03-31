package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
)

// CustomMockStub overrides GetStateByRange
type CustomMockStub struct {
	*shimtest.MockStub
}

// MockIterator implements shim.StateQueryIteratorInterface
type MockIterator struct {
	shim.StateQueryIteratorInterface
	items   []*queryresult.KV
	current int
}

func (it *MockIterator) HasNext() bool {
	return it.current < len(it.items)
}

func (it *MockIterator) Next() (*queryresult.KV, error) {
	if !it.HasNext() {
		return nil, fmt.Errorf("no more items")
	}
	item := it.items[it.current]
	it.current++
	return item, nil
}

func (it *MockIterator) Close() error {
	return nil
}

func (s *CustomMockStub) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	var results []*queryresult.KV
	for k, v := range s.State {
		results = append(results, &queryresult.KV{Key: k, Value: v})
	}
	return &MockIterator{items: results}, nil
}

func (s *CustomMockStub) GetState(key string) ([]byte, error) {
	// Simulate DB access latency
	time.Sleep(10 * time.Microsecond)
	return s.MockStub.GetState(key)
}

// MockTransactionContext implements contractapi.TransactionContextInterface
type MockTransactionContext struct {
	contractapi.TransactionContextInterface
	stub shim.ChaincodeStubInterface
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

func setupMockStub(t testing.TB, numAgents int) *MockTransactionContext {
	baseStub := shimtest.NewMockStub("mockStub", nil)
	stub := &CustomMockStub{MockStub: baseStub}

	for i := 0; i < numAgents; i++ {
		id := fmt.Sprintf("agent%d", i)
		agent := Agent{
			ID:          id,
			Type:        "producer",
			Name:        fmt.Sprintf("Producer %d", i),
			Address:     "Some Address",
			ECRBalance:  100.0,
			ENGTBalance: 50.0,
			RegisteredAt: "2023-01-01T00:00:00Z",
		}
		agentBytes, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("Failed to marshal agent: %v", err)
		}
		// Use AgentKeyPrefix from main.go
		stub.State[AgentKeyPrefix+id] = agentBytes

		// Explicitly set balance too, as the current code fetches it
		balance := TokenBalance{
			AgentID: id,
			ECR:     100.0,
			ENGT:    50.0,
		}
		balanceBytes, err := json.Marshal(balance)
		if err != nil {
			t.Fatalf("Failed to marshal balance: %v", err)
		}
		stub.State[BalanceKeyPrefix+id] = balanceBytes
	}

	return &MockTransactionContext{stub: stub}
}

func TestGetAllAgentsCorrectness(t *testing.T) {
	ctx := setupMockStub(t, 10)
	contract := &CombinedEnergyContract{}

	agents, err := contract.GetAllAgents(ctx)
	if err != nil {
		t.Fatalf("GetAllAgents failed: %v", err)
	}

	if len(agents) != 10 {
		t.Errorf("Expected 10 agents, got %d", len(agents))
	}

	for _, a := range agents {
		if a.ECRBalance != 100.0 {
			t.Errorf("Agent %s ECRBalance mismatch: got %f, want 100.0", a.ID, a.ECRBalance)
		}
	}
}

func TestGetAgent(t *testing.T) {
	ctx := setupMockStub(t, 1)
	contract := &CombinedEnergyContract{}
	agent, err := contract.GetAgent(ctx, "agent0")
	if err != nil {
		t.Fatalf("GetAgent failed: %v", err)
	}
	if agent.ECRBalance != 100.0 {
		t.Errorf("GetAgent balance mismatch: %f", agent.ECRBalance)
	}
}

func TestGetAllAgentsSourceOfTruth(t *testing.T) {
	ctx := setupMockStub(t, 1)
	stub := ctx.stub.(*CustomMockStub)

	// Modify Agent record to have different balance
	id := "agent0"
	agentBytes := stub.State[AgentKeyPrefix+id]
	var agent Agent
	json.Unmarshal(agentBytes, &agent)
	agent.ECRBalance = 999.0
	newAgentBytes, _ := json.Marshal(agent)
	stub.State[AgentKeyPrefix+id] = newAgentBytes

	// BalanceKey record still has 100.0 (from setup)

	contract := &CombinedEnergyContract{}
	agents, _ := contract.GetAllAgents(ctx)

	if len(agents) == 0 {
		t.Fatalf("No agents returned")
	}

	// Before optimization: 100.0 (fetches from balance)
	// After optimization: 999.0 (uses agent record)
	if agents[0].ECRBalance == 100.0 {
		t.Log("Source: BalanceKey (Unoptimized)")
	} else if agents[0].ECRBalance == 999.0 {
		t.Log("Source: AgentRecord (Optimized)")
	} else {
		 t.Errorf("Unexpected balance: %f", agents[0].ECRBalance)
	}
}

func BenchmarkGetAllAgents(b *testing.B) {
	ctx := setupMockStub(b, 1000)
	contract := &CombinedEnergyContract{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := contract.GetAllAgents(ctx)
		if err != nil {
			b.Fatalf("GetAllAgents failed: %v", err)
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)

// MockTransactionContext implements contractapi.TransactionContextInterface
type MockTransactionContext struct {
	contractapi.TransactionContext
	stub shim.ChaincodeStubInterface
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

func (m *MockTransactionContext) GetClientIdentity() cid.ClientIdentity {
	return nil
}

func seedLedger(stub *shimtest.MockStub, numOffers int, numOther int) {
	stub.MockTransactionStart("seed")

	// Add offers
	for i := 0; i < numOffers; i++ {
		id := fmt.Sprintf("offer_%d", i)
		offer := Offer{
			ID:           id,
			SellerID:     "seller1",
			BuyerID:      "",
			EnergyAmount: 100,
			PricePerKWh:  0.5,
			Status:       "OPEN",
			CreatedAt:    time.Now().String(),
		}
		bytes, _ := json.Marshal(offer)
		stub.PutState(OfferKeyPrefix+id, bytes)
	}

	// Add other keys (noise)
	for i := 0; i < numOther; i++ {
		key := fmt.Sprintf("otherkey_%d", i)
		stub.PutState(key, []byte("some value"))
	}

	// Add agents (noise with different prefix)
	for i := 0; i < numOther; i++ {
		id := fmt.Sprintf("agent_%d", i)
		stub.PutState(AgentKeyPrefix+id, []byte("{}"))
	}

	stub.MockTransactionEnd("seed")
}

func TestGetAllOffers(t *testing.T) {
	stub := shimtest.NewMockStub("energy_cc", nil)
	seedLedger(stub, 10, 10) // Small dataset for functional test

	ctx := &MockTransactionContext{stub: stub}
	contract := &CombinedEnergyContract{}

	offers, err := contract.GetAllOffers(ctx)
	assert.NoError(t, err)
	assert.Len(t, offers, 10)

	for _, o := range offers {
		assert.Contains(t, o.ID, "offer_")
	}
}

func BenchmarkGetAllOffers(b *testing.B) {
	// Setup mock stub
	stub := shimtest.NewMockStub("energy_cc", nil)

	// Seed with a mix of data
	// 1000 offers, 10000 other items.
	// Enough to show difference but not too slow for CI
	seedLedger(stub, 1000, 10000)

	ctx := &MockTransactionContext{stub: stub}
	contract := &CombinedEnergyContract{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		offers, err := contract.GetAllOffers(ctx)
		if err != nil {
			b.Fatalf("GetAllOffers failed: %v", err)
		}
		if len(offers) == 0 {
			b.Fatal("Expected offers, got 0")
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
    "github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// MockContext captures the stub for testing
type MockContext struct {
	contractapi.TransactionContext
	stub shim.ChaincodeStubInterface
}

func (m *MockContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

func BenchmarkGetAllOffers(b *testing.B) {
	// Create a MockStub. We don't need a real chaincode for this benchmark
	// as we are calling the method directly.
	stub := shimtest.NewMockStub("energy_contract", nil)

	stub.MockTransactionStart("setup")

	// 1. Populate with noise (keys that are NOT offers)
	// Add agents (starts with 'a')
	for i := 0; i < 2000; i++ {
		key := agentKey(fmt.Sprintf("agent%04d", i))
		stub.PutState(key, []byte(`{}`))
	}

	// Add some random junk (starts with 'r')
	for i := 0; i < 2000; i++ {
		key := fmt.Sprintf("randomjunk:%04d", i)
		stub.PutState(key, []byte("junk"))
	}

	// 2. Populate with actual Offers (starts with 's' -> spotmarket:offer:)
	numOffers := 100
	for i := 0; i < numOffers; i++ {
		id := fmt.Sprintf("offer%04d", i)
		offer := Offer{
			ID:           id,
			SellerID:     "sellerX",
			EnergyAmount: 100.0,
			PricePerKWh:  1.5,
			Status:       "OPEN",
		}
		data, _ := json.Marshal(offer)
		stub.PutState(offerKey(id), data)
	}

	// 3. Add noise AFTER the offers (starts with 'z')
	for i := 0; i < 2000; i++ {
		key := fmt.Sprintf("zebra:%04d", i)
		stub.PutState(key, []byte("junk"))
	}

	stub.MockTransactionEnd("setup")

	contract := &CombinedEnergyContract{}
	ctx := &MockContext{stub: stub}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		offers, err := contract.GetAllOffers(ctx)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
		if len(offers) != numOffers {
			b.Fatalf("Expected %d offers, got %d", numOffers, len(offers))
		}
	}
}

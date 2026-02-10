package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionContext implements contractapi.TransactionContextInterface
type MockTransactionContext struct {
	mock.Mock
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(shim.ChaincodeStubInterface)
}

func (m *MockTransactionContext) GetClientIdentity() cid.ClientIdentity {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(cid.ClientIdentity)
}

func TestSettleContractPeriod_UpdatesAgentRecords(t *testing.T) {
	cc := new(CombinedEnergyContract)
	stub := shimtest.NewMockStub("energy_contract", nil)

	ctx := new(MockTransactionContext)
	ctx.On("GetStub").Return(stub)

	// 1. Register Agents
	stub.MockTransactionStart("tx1")
	err := cc.RegisterAgent(ctx, "seller1", "producer", "Seller 1", "Addr 1")
	stub.MockTransactionEnd("tx1")
	assert.NoError(t, err)

	stub.MockTransactionStart("tx2")
	err = cc.RegisterAgent(ctx, "buyer1", "consumer", "Buyer 1", "Addr 2")
	stub.MockTransactionEnd("tx2")
	assert.NoError(t, err)

	// 2. Mint Tokens
	stub.MockTransactionStart("tx3")
	err = cc.Mint(ctx, "seller1", "ECR", "1000")
	stub.MockTransactionEnd("tx3")
	assert.NoError(t, err)

	stub.MockTransactionStart("tx4")
	err = cc.Mint(ctx, "buyer1", "ENGT", "1000")
	stub.MockTransactionEnd("tx4")
	assert.NoError(t, err)

	// Verify initial balances in agent record
	sellerJSON := stub.State["agentreg:agent:seller1"]
	var sellerAgent Agent
	err = json.Unmarshal(sellerJSON, &sellerAgent)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, sellerAgent.ECRBalance)

	// 3. Create Contract
	stub.MockTransactionStart("tx5")
	err = cc.CreateSupplyContract(ctx, "contract1", "seller1", "buyer1", 100, 2.0, "2023-01-01", "2023-12-31", "MONTHLY", 0, 0)
	stub.MockTransactionEnd("tx5")
	assert.NoError(t, err)

	// 4. Report Delivery
	stub.MockTransactionStart("tx6")
	err = cc.ReportDelivery(ctx, "contract1", 50)
	stub.MockTransactionEnd("tx6")
	assert.NoError(t, err)

	// 5. Settle
	stub.MockTransactionStart("tx7")
	err = cc.SettleContractPeriod(ctx, "contract1", 50)
	stub.MockTransactionEnd("tx7")
	assert.NoError(t, err)

	// 6. Verify Balances
	// We can inspect balances directly from state or using GetBalance helper (which reads state)
	// But GetBalance uses stub, so we need a transaction context or just call stub.GetState inside it.
	// Since GetBalance calls GetStub().GetState(), we need MockTransactionStart if GetState requires it.
	// Actually MockStub.GetState usually works without transaction?
	// The error "cannot PutState without a transactions" implies PutState needs it. GetState usually doesn't enforce it strictly in MockStub but let's be safe.

	// Let's just inspect stub.State directly for verification to avoid complexity

	sellerBalKey := "energytoken:balance:seller1"
	buyerBalKey := "energytoken:balance:buyer1"

	var sellerBal TokenBalance
	json.Unmarshal(stub.State[sellerBalKey], &sellerBal)
	assert.Equal(t, 950.0, sellerBal.ECR)
	assert.Equal(t, 100.0, sellerBal.ENGT)

	var buyerBal TokenBalance
	json.Unmarshal(stub.State[buyerBalKey], &buyerBal)
	assert.Equal(t, 50.0, buyerBal.ECR)
	assert.Equal(t, 900.0, buyerBal.ENGT)

	// 7. Verify Agent Records (The actual task)

	sellerJSON = stub.State["agentreg:agent:seller1"]
	err = json.Unmarshal(sellerJSON, &sellerAgent)
	assert.NoError(t, err)
	assert.Equal(t, 950.0, sellerAgent.ECRBalance, "Seller Agent ECRBalance should be updated in record")
	assert.Equal(t, 100.0, sellerAgent.ENGTBalance, "Seller Agent ENGTBalance should be updated in record")

	buyerJSON := stub.State["agentreg:agent:buyer1"]
	var buyerAgent Agent
	err = json.Unmarshal(buyerJSON, &buyerAgent)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, buyerAgent.ECRBalance, "Buyer Agent ECRBalance should be updated in record")
	assert.Equal(t, 900.0, buyerAgent.ENGTBalance, "Buyer Agent ENGTBalance should be updated in record")
}

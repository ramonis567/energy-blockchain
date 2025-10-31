package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SpotMarketContract struct {
	contractapi.Contract
}

type Offer struct {
	ID           string  `json:"id"`
	SellerID     string  `json:"seller_id"`
	BuyerID      string  `json:"buyer_id,omitempty"`
	EnergyAmount float64 `json:"energy_amount"` // em kWh
	PricePerKWh  float64 `json:"price_per_kwh"` // em ENGT/kWh
	TotalPrice   float64 `json:"total_price"`
	Status       string  `json:"status"` // "OPEN", "ACCEPTED", "SETTLED"
	CreatedAt    string  `json:"created_at"`
	AcceptedAt   string  `json:"accepted_at,omitempty"`
	SettledAt    string  `json:"settled_at,omitempty"`
}

func (s *SpotMarketContract) CreateOffer(ctx contractapi.TransactionContextInterface,
	id string, sellerID string, energyAmount float64, pricePerKWh float64) error {

	if energyAmount <= 0 || pricePerKWh <= 0 {
		return fmt.Errorf("valores inválidos: energia e preço devem ser positivos")
	}

	exists, err := s.OfferExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("oferta já existe: %s", id)
	}

	total := energyAmount * pricePerKWh
	offer := Offer{
		ID:           id,
		SellerID:     sellerID,
		EnergyAmount: energyAmount,
		PricePerKWh:  pricePerKWh,
		TotalPrice:   total,
		Status:       "OPEN",
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	data, _ := json.Marshal(offer)
	return ctx.GetStub().PutState(id, data)
}

// AcceptOffer aceita uma oferta existente
func (s *SpotMarketContract) AcceptOffer(ctx contractapi.TransactionContextInterface, id string, buyerID string) error {
	data, err := ctx.GetStub().GetState(id)
	if err != nil || data == nil {
		return fmt.Errorf("oferta não encontrada: %s", id)
	}

	var offer Offer
	json.Unmarshal(data, &offer)

	if offer.Status != "OPEN" {
		return fmt.Errorf("a oferta %s já foi aceita ou liquidada", id)
	}

	// Verifica saldo do comprador (ENGT)
	check := ctx.GetStub().InvokeChaincode("energytoken", [][]byte{
		[]byte("GetBalance"), []byte(buyerID),
	}, "mychannel")
	if check.Status != 200 {
		return fmt.Errorf("erro ao verificar saldo do comprador")
	}

	var buyerBalance map[string]interface{}
	json.Unmarshal(check.Payload, &buyerBalance)
	if buyerBalance["engt"].(float64) < offer.TotalPrice {
		return fmt.Errorf("saldo insuficiente de ENGT")
	}

	// Transferência financeira (ENGT)
	ctx.GetStub().InvokeChaincode("energytoken", [][]byte{
		[]byte("Transfer"), []byte(buyerID), []byte(offer.SellerID),
		[]byte("ENGT"), []byte(fmt.Sprintf("%f", offer.TotalPrice)),
	}, "mychannel")

	// Transferência energética (ECR)
	ctx.GetStub().InvokeChaincode("energytoken", [][]byte{
		[]byte("Transfer"), []byte(offer.SellerID), []byte(buyerID),
		[]byte("ECR"), []byte(fmt.Sprintf("%f", offer.EnergyAmount)),
	}, "mychannel")

	offer.BuyerID = buyerID
	offer.Status = "SETTLED"
	offer.AcceptedAt = time.Now().UTC().Format(time.RFC3339)
	offer.SettledAt = offer.AcceptedAt

	newData, _ := json.Marshal(offer)
	ctx.GetStub().PutState(id, newData)
	return nil
}

// GetAllOffers lista todas as ofertas
func (s *SpotMarketContract) GetAllOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, error) {
	it, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer it.Close()

	var offers []*Offer
	for it.HasNext() {
		res, _ := it.Next()
		var offer Offer
		if json.Unmarshal(res.Value, &offer) == nil && offer.ID != "" {
			offers = append(offers, &offer)
		}
	}
	return offers, nil
}

func (s *SpotMarketContract) OfferExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func main() {
	cc, err := contractapi.NewChaincode(&SpotMarketContract{})
	if err != nil {
		panic(err)
	}
	if err := cc.Start(); err != nil {
		panic(err)
	}
}

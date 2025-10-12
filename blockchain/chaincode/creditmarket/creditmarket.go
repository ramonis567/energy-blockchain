package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct{ contractapi.Contract }

// Prefixo apenas para ofertas (mantém separado)
const OfferPrefix = "OFFER_"

func offerKey(id string) string { return OfferPrefix + id }

// ----- Estruturas -----
type Agent struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

type Offer struct {
	ID       string  `json:"id"`
	SellerID string  `json:"sellerId"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"` // OPEN / ACCEPTED / CANCELLED
}

// ----- Agentes -----
// Registrar novo agente
func (s *SmartContract) RegisterAgent(ctx contractapi.TransactionContextInterface, id, agentType string) error {
	exists, err := s.AgentExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("agente já existe: %s", id)
	}

	agent := Agent{
		ID:      id,
		Type:    agentType,
		Balance: 0,
	}

	bz, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	// Grava com o ID puro
	return ctx.GetStub().PutState(id, bz)
}

func (s *SmartContract) AgentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	bz, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return bz != nil, nil
}

// Listar todos os agentes
func (s *SmartContract) GetAllAgents(ctx contractapi.TransactionContextInterface) ([]*Agent, error) {
	it, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("erro ao obter range: %v", err)
	}
	defer it.Close()

	var agents []*Agent
	for it.HasNext() {
		qr, err := it.Next()
		if err != nil {
			return nil, err
		}

		// Filtra apenas o que não for uma oferta
		if strings.HasPrefix(qr.Key, OfferPrefix) {
			continue
		}

		var a Agent
		if err := json.Unmarshal(qr.Value, &a); err != nil {
			continue
		}
		agents = append(agents, &a)
	}
	return agents, nil
}

// Consultar saldo
func (s *SmartContract) GetAgentBalance(ctx contractapi.TransactionContextInterface, id string) (float64, error) {
	bz, err := ctx.GetStub().GetState(id)
	if err != nil || bz == nil {
		return 0, fmt.Errorf("agente não encontrado: %s", id)
	}
	var a Agent
	_ = json.Unmarshal(bz, &a)
	return a.Balance, nil
}

// Atualizar saldo (Interconnection)
func (s *SmartContract) UpdateBalance(ctx contractapi.TransactionContextInterface, id string, delta float64) error {
	bz, err := ctx.GetStub().GetState(id)
	if err != nil || bz == nil {
		return fmt.Errorf("agente não encontrado: %s", id)
	}
	var a Agent
	_ = json.Unmarshal(bz, &a)
	a.Balance += delta
	return ctx.GetStub().PutState(id, mustJSON(a))
}

// ----- Ofertas -----
// Criar oferta (aplicação)
func (s *SmartContract) CreateOffer(ctx contractapi.TransactionContextInterface, id, seller string, qty, price float64) error {
	// Valida existência do vendedor
	exists, err := s.AgentExists(ctx, seller)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("vendedor não existe: %s", seller)
	}

	offer := Offer{ID: id, SellerID: seller, Quantity: qty, Price: price, Status: "OPEN"}
	return ctx.GetStub().PutState(offerKey(id), mustJSON(offer))
}

// Aceitar oferta
func (s *SmartContract) AcceptOffer(ctx contractapi.TransactionContextInterface, offerId, buyer string) error {
	offerBZ, err := ctx.GetStub().GetState(offerKey(offerId))
	if err != nil || offerBZ == nil {
		return fmt.Errorf("oferta não encontrada: %s", offerId)
	}
	var offer Offer
	_ = json.Unmarshal(offerBZ, &offer)
	if offer.Status != "OPEN" {
		return fmt.Errorf("oferta não disponível")
	}

	// Carregar agentes
	sellerBZ, _ := ctx.GetStub().GetState(offer.SellerID)
	buyerBZ, _ := ctx.GetStub().GetState(buyer)
	if sellerBZ == nil {
		return fmt.Errorf("vendedor não encontrado: %s", offer.SellerID)
	}
	if buyerBZ == nil {
		return fmt.Errorf("comprador não encontrado: %s", buyer)
	}

	var seller, buyerAgent Agent
	_ = json.Unmarshal(sellerBZ, &seller)
	_ = json.Unmarshal(buyerBZ, &buyerAgent)

	// Checar saldo do vendedor
	if seller.Balance < offer.Quantity {
		return fmt.Errorf("saldo insuficiente do vendedor")
	}

	// Transferir energia (créditos)
	seller.Balance -= offer.Quantity
	buyerAgent.Balance += offer.Quantity
	offer.Status = "ACCEPTED"

	// Persistir no ledger
	if err := ctx.GetStub().PutState(seller.ID, mustJSON(seller)); err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(buyerAgent.ID, mustJSON(buyerAgent)); err != nil {
		return err
	}
	return ctx.GetStub().PutState(offerKey(offer.ID), mustJSON(offer))
}

// ----- util -----
func mustJSON(v any) []byte { b, _ := json.Marshal(v); return b }

func main() {
	cc, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(err)
	}
	if err := cc.Start(); err != nil {
		panic(err)
	}
}

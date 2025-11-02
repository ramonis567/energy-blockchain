// SPDX-License-Identifier: Apache-2.0
// ========================================================
// Merged Chaincode: AgentRegistry + EnergyToken + SpotMarket
// ========================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// -------------------------
// Enhanced Key Prefixes
// -------------------------

const (
	AgentKeyPrefix   = "agentreg:agent:"
	BalanceKeyPrefix = "energytoken:balance:"
	OfferKeyPrefix   = "spotmarket:offer:"
)

// -------------------------
// Utility Functions
// -------------------------

func agentKey(id string) string {
	return AgentKeyPrefix + id
}

func balanceKey(id string) string {
	return BalanceKeyPrefix + id
}

func offerKey(id string) string {
	return OfferKeyPrefix + id
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// IMPROVED: Better floating point comparison with higher tolerance
func almostEqual(a, b float64) bool {
	const eps = 1e-6 // Increased tolerance for floating point arithmetic
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < eps
}

// -------------------------
// Data Structures
// -------------------------

var validAgentTypes = map[string]bool{
	"producer":    true,
	"consumer":    true,
	"prosumer":    true,
	"distributor": true,
	"battery":     true,
}

type Agent struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Address string `json:"address"`

	ECRBalance   float64 `json:"ecr_balance"`
	ENGTBalance  float64 `json:"engt_balance"`
	RegisteredAt string  `json:"registered_at"`
}

type AgentRegistrationEvent struct {
	AgentID    string `json:"agentId"`
	Type       string `json:"type"`
	Timestamp  string `json:"timestamp"`
	EmitTokens bool   `json:"emitTokens"`
}

type TokenBalance struct {
	AgentID string  `json:"agentId"`
	ECR     float64 `json:"ecr"`
	ENGT    float64 `json:"engt"`
}

type Offer struct {
	ID           string  `json:"id"`
	SellerID     string  `json:"seller_id"`
	BuyerID      string  `json:"buyer_id"`
	EnergyAmount float64 `json:"energy_amount"`
	PricePerKWh  float64 `json:"price_per_kwh"`
	TotalPrice   float64 `json:"total_price"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	AcceptedAt   string  `json:"accepted_at"`
	SettledAt    string  `json:"settled_at"`
}

// -------------------------
// Combined Contract
// -------------------------

type CombinedEnergyContract struct {
	contractapi.Contract
}

// ========================================================
// Agent Registry Methods
// ========================================================

func (s *CombinedEnergyContract) RegisterAgent(ctx contractapi.TransactionContextInterface,
	id string, agentType string, name string, address string) error {

	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("ID do agente não pode estar vazio")
	}
	if strings.TrimSpace(agentType) == "" {
		return fmt.Errorf("tipo do agente não pode estar vazio")
	}
	if !validAgentTypes[agentType] {
		return fmt.Errorf("tipo do agente inválido: %s", agentType)
	}

	exists, err := s.AgentExists(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência do agente: %v", err)
	}
	if exists {
		return fmt.Errorf("agente existente, crie outro: %s", id)
	}

	agent := Agent{
		ID:           id,
		Type:         agentType,
		Name:         name,
		Address:      address,
		ECRBalance:   0.0,
		ENGTBalance:  0.0,
		RegisteredAt: time.Now().UTC().Format(time.RFC3339),
	}

	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("erro ao serializar agente: %v", err)
	}

	if err := ctx.GetStub().PutState(agentKey(id), agentJSON); err != nil {
		return fmt.Errorf("erro ao salvar agente no ledger: %v", err)
	}

	event := AgentRegistrationEvent{
		AgentID:    id,
		Type:       agentType,
		Timestamp:  agent.RegisteredAt,
		EmitTokens: true,
	}
	eventJSON, _ := json.Marshal(event)
	if err := ctx.GetStub().SetEvent("AgentRegistry:AgentRegistered", eventJSON); err != nil {
		log.Printf("Aviso: não foi possível emitir evento AgentRegistered: %v", err)
	}

	return nil
}

func (s *CombinedEnergyContract) AgentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	agentJSON, err := ctx.GetStub().GetState(agentKey(id))
	if err != nil {
		return false, fmt.Errorf("erro ao acessar ledger: %v", err)
	}
	return agentJSON != nil, nil
}

func (s *CombinedEnergyContract) GetAgent(ctx contractapi.TransactionContextInterface, id string) (*Agent, error) {
	agentJSON, err := ctx.GetStub().GetState(agentKey(id))
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar ledger: %v", err)
	}
	if agentJSON == nil {
		return nil, fmt.Errorf("agente não encontrado: %s", id)
	}

	var agent Agent
	if err := json.Unmarshal(agentJSON, &agent); err != nil {
		return nil, fmt.Errorf("erro ao desserializar agente: %v", err)
	}
	return &agent, nil
}

func (s *CombinedEnergyContract) GetAllAgents(ctx contractapi.TransactionContextInterface) ([]*Agent, error) {
	it, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre o ledger: %v", err)
	}
	defer it.Close()

	var agents []*Agent
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, fmt.Errorf("erro ao ler próximo item: %v", err)
		}

		if !strings.HasPrefix(kv.Key, AgentKeyPrefix) {
			continue
		}

		var a Agent
		if err := json.Unmarshal(kv.Value, &a); err != nil {
			continue
		}
		if a.ID != "" {
			agents = append(agents, &a)
		}
	}
	return agents, nil
}

func (s *CombinedEnergyContract) GetAgentsByType(ctx contractapi.TransactionContextInterface, agentType string) ([]*Agent, error) {
	all, err := s.GetAllAgents(ctx)
	if err != nil {
		return nil, err
	}

	var out []*Agent
	for _, a := range all {
		if a.Type == agentType {
			out = append(out, a)
		}
	}
	return out, nil
}

func (s *CombinedEnergyContract) UpdateAgent(ctx contractapi.TransactionContextInterface,
	id string, name string, address string) error {

	agent, err := s.GetAgent(ctx, id)
	if err != nil {
		return err
	}

	agent.Name = name
	agent.Address = address

	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("erro ao serializar agente: %v", err)
	}
	return ctx.GetStub().PutState(agentKey(id), agentJSON)
}

func (s *CombinedEnergyContract) GetAgentCount(ctx contractapi.TransactionContextInterface) (int, error) {
	agents, err := s.GetAllAgents(ctx)
	if err != nil {
		return 0, err
	}
	return len(agents), nil
}

func (s *CombinedEnergyContract) GetAgentFullInfo(ctx contractapi.TransactionContextInterface, id string) (map[string]interface{}, error) {
	agent, err := s.GetAgent(ctx, id)
	if err != nil {
		return nil, err
	}

	balance, err := s.GetBalance(ctx, id)
	if err != nil {
		return map[string]interface{}{
			"agent":   agent,
			"balance": map[string]interface{}{"agentId": id, "ecr": 0.0, "engt": 0.0, "status": "unavailable"},
		}, nil
	}

	return map[string]interface{}{
		"agent":   agent,
		"balance": balance,
	}, nil
}

// ========================================================
// Energy Token Methods
// ========================================================

func (c *CombinedEnergyContract) Mint(ctx contractapi.TransactionContextInterface, agentID string, tokenType string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("quantidade inválida: %.2f", amount)
	}
	key := balanceKey(agentID)

	balance := TokenBalance{AgentID: agentID}
	existing, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("erro ao acessar ledger: %v", err)
	}
	if existing != nil {
		_ = json.Unmarshal(existing, &balance)
	}

	switch tokenType {
	case "ECR":
		balance.ECR += amount
	case "ENGT":
		balance.ENGT += amount
	default:
		return fmt.Errorf("tokenType inválido: %s", tokenType)
	}

	data, _ := json.Marshal(balance)
	if err := ctx.GetStub().PutState(key, data); err != nil {
		return fmt.Errorf("erro ao salvar saldo: %v", err)
	}

	event := map[string]interface{}{
		"agentId": agentID,
		"token":   tokenType,
		"amount":  amount,
	}
	eventBytes, _ := json.Marshal(event)
	if err := ctx.GetStub().SetEvent("EnergyToken:TokenMinted", eventBytes); err != nil {
		log.Printf("Aviso: falha ao emitir evento TokenMinted: %v", err)
	}

	return nil
}

func (c *CombinedEnergyContract) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, tokenType string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("quantidade inválida")
	}
	if from == to {
		return fmt.Errorf("transação inválida: remetente e destinatário iguais")
	}

	fromKey := balanceKey(from)
	toKey := balanceKey(to)

	var fromBal, toBal TokenBalance

	fromData, _ := ctx.GetStub().GetState(fromKey)
	toData, _ := ctx.GetStub().GetState(toKey)

	if fromData != nil {
		_ = json.Unmarshal(fromData, &fromBal)
	} else {
		fromBal = TokenBalance{AgentID: from, ECR: 0, ENGT: 0}
	}
	if toData != nil {
		_ = json.Unmarshal(toData, &toBal)
	} else {
		toBal = TokenBalance{AgentID: to, ECR: 0, ENGT: 0}
	}

	switch tokenType {
	case "ECR":
		if fromBal.ECR < amount {
			return fmt.Errorf("saldo insuficiente de ECR (%.2f < %.2f)", fromBal.ECR, amount)
		}
		fromBal.ECR -= amount
		toBal.ECR += amount
	case "ENGT":
		if fromBal.ENGT < amount {
			return fmt.Errorf("saldo insuficiente de ENGT (%.2f < %.2f)", fromBal.ENGT, amount)
		}
		fromBal.ENGT -= amount
		toBal.ENGT += amount
	default:
		return fmt.Errorf("tokenType inválido: %s", tokenType)
	}

	ctx.GetStub().PutState(fromKey, mustMarshal(fromBal))
	ctx.GetStub().PutState(toKey, mustMarshal(toBal))

	event := map[string]interface{}{
		"from":    from,
		"to":      to,
		"token":   tokenType,
		"amount":  amount,
		"fromBal": fromBal,
		"toBal":   toBal,
	}
	eventBytes, _ := json.Marshal(event)
	if err := ctx.GetStub().SetEvent("EnergyToken:TokenTransferred", eventBytes); err != nil {
		log.Printf("Aviso: falha ao emitir evento TokenTransferred: %v", err)
	}

	return nil
}

func (c *CombinedEnergyContract) GetBalance(ctx contractapi.TransactionContextInterface, agentID string) (*TokenBalance, error) {
	key := balanceKey(agentID)
	data, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar ledger: %v", err)
	}

	if data == nil {
		return &TokenBalance{
			AgentID: agentID,
			ECR:     0.0,
			ENGT:    0.0,
		}, nil
	}

	var balance TokenBalance
	if err := json.Unmarshal(data, &balance); err != nil {
		return nil, fmt.Errorf("erro ao desserializar saldo: %v", err)
	}
	return &balance, nil
}

// ========================================================
// Spot Market Methods - FIXED VERSION
// ========================================================

func (s *CombinedEnergyContract) CreateOffer(
	ctx contractapi.TransactionContextInterface,
	id string, sellerID string, energyAmount float64, pricePerKWh float64,
) error {

	if strings.TrimSpace(id) == "" || strings.TrimSpace(sellerID) == "" {
		return fmt.Errorf("id e sellerID são obrigatórios")
	}
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

	offer := Offer{
		ID:           id,
		SellerID:     sellerID,
		BuyerID:      "",
		EnergyAmount: energyAmount,
		PricePerKWh:  pricePerKWh,
		TotalPrice:   energyAmount * pricePerKWh,
		Status:       "OPEN",
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		AcceptedAt:   "",
		SettledAt:    "",
	}
	data, _ := json.Marshal(offer)
	if err := ctx.GetStub().PutState(offerKey(id), data); err != nil {
		return fmt.Errorf("erro ao salvar oferta: %v", err)
	}

	evt := map[string]interface{}{
		"id":            offer.ID,
		"seller_id":     offer.SellerID,
		"energy_amount": offer.EnergyAmount,
		"price_per_kwh": offer.PricePerKWh,
		"total_price":   offer.TotalPrice,
		"created_at":    offer.CreatedAt,
	}
	if b, _ := json.Marshal(evt); ctx.GetStub().SetEvent("SpotMarket:OfferCreated", b) != nil {
		log.Printf("Aviso: falha ao emitir OfferCreated")
	}
	return nil
}

// FIXED: Completely rewritten AcceptOffer function
func (s *CombinedEnergyContract) AcceptOffer(
	ctx contractapi.TransactionContextInterface,
	id string, buyerID string,
) error {
	// Load offer
	raw, err := ctx.GetStub().GetState(offerKey(id))
	if err != nil || raw == nil {
		return fmt.Errorf("oferta não encontrada: %s", id)
	}
	var offer Offer
	if err := json.Unmarshal(raw, &offer); err != nil {
		return fmt.Errorf("erro ao desserializar oferta: %v", err)
	}

	// Validate offer status
	if offer.Status != "OPEN" {
		return fmt.Errorf("a oferta %s já foi aceita ou liquidada", id)
	}
	if offer.SellerID == buyerID {
		return fmt.Errorf("um agente não pode comprar a própria oferta")
	}

	// Check balances BEFORE transfers
	sellerBalance, err := s.GetBalance(ctx, offer.SellerID)
	if err != nil {
		return fmt.Errorf("falha ao ler saldo do vendedor: %v", err)
	}
	buyerBalance, err := s.GetBalance(ctx, buyerID)
	if err != nil {
		return fmt.Errorf("falha ao ler saldo do comprador: %v", err)
	}

	// Verify sufficient balances
	if buyerBalance.ENGT < offer.TotalPrice {
		return fmt.Errorf("saldo insuficiente de ENGT (%.2f < %.2f)", buyerBalance.ENGT, offer.TotalPrice)
	}
	if sellerBalance.ECR < offer.EnergyAmount {
		return fmt.Errorf("saldo insuficiente de ECR (%.2f < %.2f)", sellerBalance.ECR, offer.EnergyAmount)
	}

	// Calculate expected final balances
	expectedBuyerECR := buyerBalance.ECR + offer.EnergyAmount
	expectedBuyerENGT := buyerBalance.ENGT - offer.TotalPrice
	expectedSellerECR := sellerBalance.ECR - offer.EnergyAmount
	expectedSellerENGT := sellerBalance.ENGT + offer.TotalPrice

	// Execute ENGT transfer (buyer pays seller)
	err = s.Transfer(ctx, buyerID, offer.SellerID, "ENGT", offer.TotalPrice)
	if err != nil {
		return fmt.Errorf("falha ao transferir ENGT: %v", err)
	}

	// Execute ECR transfer (seller delivers energy to buyer)
	err = s.Transfer(ctx, offer.SellerID, buyerID, "ECR", offer.EnergyAmount)
	if err != nil {
		// If ECR transfer fails, we should revert the ENGT transfer
		// For simplicity, we rely on the atomic nature of blockchain transactions
		// In production, you might want more sophisticated compensation logic
		return fmt.Errorf("falha ao transferir ECR: %v", err)
	}

	// Verify final balances (with detailed logging for debugging)
	finalBuyerBalance, err := s.GetBalance(ctx, buyerID)
	if err != nil {
		return fmt.Errorf("falha ao verificar saldo final do comprador: %v", err)
	}
	finalSellerBalance, err := s.GetBalance(ctx, offer.SellerID)
	if err != nil {
		return fmt.Errorf("falha ao verificar saldo final do vendedor: %v", err)
	}

	// Debug logging
	log.Printf("Balance verification:")
	log.Printf("Buyer - Expected: ECR=%.6f, ENGT=%.6f", expectedBuyerECR, expectedBuyerENGT)
	log.Printf("Buyer - Actual:   ECR=%.6f, ENGT=%.6f", finalBuyerBalance.ECR, finalBuyerBalance.ENGT)
	log.Printf("Seller - Expected: ECR=%.6f, ENGT=%.6f", expectedSellerECR, expectedSellerENGT)
	log.Printf("Seller - Actual:   ECR=%.6f, ENGT=%.6f", finalSellerBalance.ECR, finalSellerBalance.ENGT)

	// Verify with tolerance for floating point arithmetic
	if !almostEqual(finalBuyerBalance.ECR, expectedBuyerECR) ||
		!almostEqual(finalBuyerBalance.ENGT, expectedBuyerENGT) ||
		!almostEqual(finalSellerBalance.ECR, expectedSellerECR) ||
		!almostEqual(finalSellerBalance.ENGT, expectedSellerENGT) {

		log.Printf("Balance mismatch detected!")
		return fmt.Errorf("inconsistência na liquidação: deltas inesperados de saldo")
	}

	// Update offer status to SETTLED
	now := time.Now().UTC().Format(time.RFC3339)
	offer.BuyerID = buyerID
	offer.Status = "SETTLED"
	offer.AcceptedAt = now
	offer.SettledAt = now

	offerJSON, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("erro ao serializar oferta atualizada: %v", err)
	}

	if err := ctx.GetStub().PutState(offerKey(id), offerJSON); err != nil {
		return fmt.Errorf("erro ao atualizar oferta liquidada: %v", err)
	}

	// Emit settlement event
	evt := map[string]interface{}{
		"id":          offer.ID,
		"seller_id":   offer.SellerID,
		"buyer_id":    offer.BuyerID,
		"energy_kwh":  offer.EnergyAmount,
		"price_total": offer.TotalPrice,
		"accepted_at": offer.AcceptedAt,
		"buyer_ecr":   finalBuyerBalance.ECR,
		"buyer_engt":  finalBuyerBalance.ENGT,
		"seller_ecr":  finalSellerBalance.ECR,
		"seller_engt": finalSellerBalance.ENGT,
	}
	if b, err := json.Marshal(evt); err == nil {
		if err := ctx.GetStub().SetEvent("SpotMarket:OfferSettled", b); err != nil {
			log.Printf("Aviso: falha ao emitir OfferSettled: %v", err)
		}
	}

	return nil
}

func (s *CombinedEnergyContract) GetAllOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, error) {
	it, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre o ledger: %v", err)
	}
	defer it.Close()

	var offers []*Offer
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, fmt.Errorf("erro ao ler próximo item: %v", err)
		}
		if !strings.HasPrefix(kv.Key, OfferKeyPrefix) {
			continue
		}
		var o Offer
		if json.Unmarshal(kv.Value, &o) == nil && o.ID != "" {
			offers = append(offers, &o)
		}
	}
	return offers, nil
}

func (s *CombinedEnergyContract) OfferExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(offerKey(id))
	if err != nil {
		return false, fmt.Errorf("erro ao verificar oferta: %v", err)
	}
	return data != nil, nil
}

// -------------------------
// Bootstrap
// -------------------------

func main() {
	cc, err := contractapi.NewChaincode(&CombinedEnergyContract{})
	if err != nil {
		log.Panicf("Erro ao criar chaincode: %v", err)
	}
	if err := cc.Start(); err != nil {
		log.Panicf("Erro ao iniciar chaincode: %v", err)
	}
}

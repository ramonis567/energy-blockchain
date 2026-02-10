// SPDX-License-Identifier: Apache-2.0
// ========================================================
// Merged Chaincode: AgentRegistry + EnergyToken + SpotMarket + ContractMarket - COMPLETELY FIXED VERSION
// ========================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// -------------------------
// Enhanced Key Prefixes
// -------------------------

const (
	AgentKeyPrefix    = "agentreg:agent:"
	BalanceKeyPrefix  = "energytoken:balance:"
	OfferKeyPrefix    = "spotmarket:offer:"
	ContractKeyPrefix = "supply:contract:"
	EscrowKeyPrefix   = "supply:escrow:"
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

func contractKey(id string) string {
	return ContractKeyPrefix + id
}

func escrowKey(contractID, party, token string) string {
	return EscrowKeyPrefix + contractID + ":" + party + ":" + token
}

func getRangeForPrefix(prefix string) (string, string) {
	if len(prefix) == 0 {
		return "", ""
	}
	prefixBytes := []byte(prefix)
	endKeyBytes := make([]byte, len(prefixBytes))
	copy(endKeyBytes, prefixBytes)
	endKeyBytes[len(endKeyBytes)-1]++
	return prefix, string(endKeyBytes)
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
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

type SupplyContract struct {
	ID               string  `json:"id"`
	SellerID         string  `json:"seller_id"`
	BuyerID          string  `json:"buyer_id"`
	EnergyTotal      float64 `json:"energy_total"`      // kWh contratados
	DeliveredTotal   float64 `json:"delivered_total"`   // kWh já entregues (relatados)
	UnsettledEnergy  float64 `json:"unsettled_energy"`  // kWh entregues e ainda não liquidados
	PricePerKWh      float64 `json:"price_per_kwh"`     // ENGT/kWh
	TotalValue       float64 `json:"total_value"`       // EnergyTotal * PricePerKWh
	SellerCollateral float64 `json:"seller_collateral"` // ECR bloqueado
	BuyerCollateral  float64 `json:"buyer_collateral"`  // ENGT bloqueado
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
	SettlementFreq   string  `json:"settlement_freq"` // DAILY, WEEKLY, MONTHLY
	Status           string  `json:"status"`          // "ACTIVE", "CLOSED", "CANCELLED"
	CreatedAt        string  `json:"created_at"`
	LastSettleAt     string  `json:"last_settled_at"`
}

type Escrow struct {
	ContractID string  `json:"contract_id"`
	OwnerID    string  `json:"owner_id"` // "seller" ou "buyer"
	Token      string  `json:"token"`
	Amount     float64 `json:"amount"`
	CreatedAt  string  `json:"created_at"`
}

// -------------------------
// Combined Contract
// -------------------------

type CombinedEnergyContract struct {
	contractapi.Contract
}

// ========================================================
// Agent Registry Methods - FIXED VERSION
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

	// FIX: Update agent balances with current token balances
	balance, err := s.GetBalance(ctx, id)
	if err == nil {
		agent.ECRBalance = balance.ECR
		agent.ENGTBalance = balance.ENGT
	}

	return &agent, nil
}

// FIXED: GetAllAgents now returns agents with updated balances
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
			// FIX: Update agent balances with current token balances
			balance, err := s.GetBalance(ctx, a.ID)
			if err == nil {
				a.ECRBalance = balance.ECR
				a.ENGTBalance = balance.ENGT
			}
			agents = append(agents, &a)
		}
	}
	return agents, nil
}

// FIXED: GetAgentsByType now returns agents with updated balances
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
	agent, err := s.GetAgent(ctx, id) // This now returns agent with updated balances
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
// Energy Token Methods - FIXED VERSION
// ========================================================

// FIXED: Mint with agent balance synchronization
func (c *CombinedEnergyContract) Mint(ctx contractapi.TransactionContextInterface, agentID string, tokenType string, amountStr string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return fmt.Errorf("valor inválido para amount: %s", amountStr)
	}

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
		if err := json.Unmarshal(existing, &balance); err != nil {
			return fmt.Errorf("erro ao desserializar saldo: %v", err)
		}
	}

	switch tokenType {
	case "ECR":
		balance.ECR += amount
	case "ENGT":
		balance.ENGT += amount
	default:
		return fmt.Errorf("tokenType inválido: %s", tokenType)
	}

	data, err := json.Marshal(balance)
	if err != nil {
		return fmt.Errorf("erro ao serializar saldo: %v", err)
	}

	if err := ctx.GetStub().PutState(key, data); err != nil {
		return fmt.Errorf("erro ao salvar saldo: %v", err)
	}

	// FIX: Also update the agent record to keep balances synchronized
	agent, err := c.GetAgent(ctx, agentID)
	if err == nil {
		switch tokenType {
		case "ECR":
			agent.ECRBalance = balance.ECR
		case "ENGT":
			agent.ENGTBalance = balance.ENGT
		}
		agentJSON, err := json.Marshal(agent)
		if err == nil {
			if err := ctx.GetStub().PutState(agentKey(agentID), agentJSON); err != nil {
				log.Printf("Aviso: não foi possível atualizar saldo do agente: %v", err)
			}
		}
	}

	event := map[string]interface{}{
		"agentId": agentID,
		"token":   tokenType,
		"amount":  amount,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Aviso: erro ao serializar evento: %v", err)
	} else {
		if err := ctx.GetStub().SetEvent("EnergyToken:TokenMinted", eventBytes); err != nil {
			log.Printf("Aviso: falha ao emitir evento TokenMinted: %v", err)
		}
	}

	return nil
}

// FIXED: Transfer with agent balance synchronization
func (c *CombinedEnergyContract) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, tokenType string, amountStr string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return fmt.Errorf("valor inválido para amount: %s", amountStr)
	}

	if amount <= 0 {
		return fmt.Errorf("quantidade inválida")
	}
	if from == to {
		return fmt.Errorf("transação inválida: remetente e destinatário iguais")
	}

	fromKey := balanceKey(from)
	toKey := balanceKey(to)

	var fromBal, toBal TokenBalance

	fromData, err := ctx.GetStub().GetState(fromKey)
	if err != nil {
		return fmt.Errorf("erro ao ler saldo do remetente: %v", err)
	}
	toData, err := ctx.GetStub().GetState(toKey)
	if err != nil {
		return fmt.Errorf("erro ao ler saldo do destinatário: %v", err)
	}

	if fromData != nil {
		if err := json.Unmarshal(fromData, &fromBal); err != nil {
			return fmt.Errorf("erro ao desserializar saldo do remetente: %v", err)
		}
	} else {
		fromBal = TokenBalance{AgentID: from, ECR: 0, ENGT: 0}
	}

	if toData != nil {
		if err := json.Unmarshal(toData, &toBal); err != nil {
			return fmt.Errorf("erro ao desserializar saldo do destinatário: %v", err)
		}
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

	fromDataUpdated, err := json.Marshal(fromBal)
	if err != nil {
		return fmt.Errorf("erro ao serializar saldo do remetente: %v", err)
	}
	toDataUpdated, err := json.Marshal(toBal)
	if err != nil {
		return fmt.Errorf("erro ao serializar saldo do destinatário: %v", err)
	}

	if err := ctx.GetStub().PutState(fromKey, fromDataUpdated); err != nil {
		return fmt.Errorf("erro ao salvar saldo do remetente: %v", err)
	}
	if err := ctx.GetStub().PutState(toKey, toDataUpdated); err != nil {
		return fmt.Errorf("erro ao salvar saldo do destinatário: %v", err)
	}

	// FIX: Also update agent records to keep balances synchronized
	fromAgent, err := c.GetAgent(ctx, from)
	if err == nil {
		switch tokenType {
		case "ECR":
			fromAgent.ECRBalance = fromBal.ECR
		case "ENGT":
			fromAgent.ENGTBalance = fromBal.ENGT
		}
		fromAgentJSON, err := json.Marshal(fromAgent)
		if err == nil {
			if err := ctx.GetStub().PutState(agentKey(from), fromAgentJSON); err != nil {
				log.Printf("Aviso: não foi possível atualizar saldo do agente remetente: %v", err)
			}
		}
	}

	toAgent, err := c.GetAgent(ctx, to)
	if err == nil {
		switch tokenType {
		case "ECR":
			toAgent.ECRBalance = toBal.ECR
		case "ENGT":
			toAgent.ENGTBalance = toBal.ENGT
		}
		toAgentJSON, err := json.Marshal(toAgent)
		if err == nil {
			if err := ctx.GetStub().PutState(agentKey(to), toAgentJSON); err != nil {
				log.Printf("Aviso: não foi possível atualizar saldo do agente destinatário: %v", err)
			}
		}
	}

	event := map[string]interface{}{
		"from":    from,
		"to":      to,
		"token":   tokenType,
		"amount":  amount,
		"fromBal": fromBal,
		"toBal":   toBal,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Aviso: erro ao serializar evento de transferência: %v", err)
	} else {
		if err := ctx.GetStub().SetEvent("EnergyToken:TokenTransferred", eventBytes); err != nil {
			log.Printf("Aviso: falha ao emitir evento TokenTransferred: %v", err)
		}
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
		return fmt.Errorf("erro ao verificar existência da oferta: %v", err)
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

	data, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("erro ao serializar oferta: %v", err)
	}

	if err := ctx.GetStub().PutState(offerKey(id), data); err != nil {
		return fmt.Errorf("erro ao salvar oferta no ledger: %v", err)
	}

	log.Printf("DEBUG: Offer %s successfully created and saved", id)

	evt := map[string]interface{}{
		"id":            offer.ID,
		"seller_id":     offer.SellerID,
		"energy_amount": offer.EnergyAmount,
		"price_per_kwh": offer.PricePerKWh,
		"total_price":   offer.TotalPrice,
		"created_at":    offer.CreatedAt,
	}

	eventBytes, err := json.Marshal(evt)
	if err != nil {
		log.Printf("Aviso: erro ao serializar evento: %v", err)
		return nil
	}

	if err := ctx.GetStub().SetEvent("SpotMarket:OfferCreated", eventBytes); err != nil {
		log.Printf("Aviso: falha ao emitir evento OfferCreated: %v", err)
	}

	return nil
}

// FIXED: AcceptOffer with agent balance synchronization
func (s *CombinedEnergyContract) AcceptOffer(
	ctx contractapi.TransactionContextInterface,
	id string, buyerID string,
) error {
	log.Printf("=== ACCEPT OFFER START ===")
	log.Printf("DEBUG: Aceitando oferta %s para comprador %s", id, buyerID)

	raw, err := ctx.GetStub().GetState(offerKey(id))
	if err != nil {
		return fmt.Errorf("erro ao carregar oferta: %v", err)
	}
	if raw == nil {
		return fmt.Errorf("oferta não encontrada: %s", id)
	}

	var offer Offer
	if err := json.Unmarshal(raw, &offer); err != nil {
		return fmt.Errorf("erro ao desserializar oferta: %v", err)
	}

	if offer.Status != "OPEN" {
		return fmt.Errorf("a oferta %s já foi aceita ou liquidada", id)
	}
	if offer.SellerID == buyerID {
		return fmt.Errorf("um agente não pode comprar a própria oferta")
	}

	sellerBalance, err := s.GetBalance(ctx, offer.SellerID)
	if err != nil {
		return fmt.Errorf("falha ao ler saldo do vendedor: %v", err)
	}
	buyerBalance, err := s.GetBalance(ctx, buyerID)
	if err != nil {
		return fmt.Errorf("falha ao ler saldo do comprador: %v", err)
	}

	if sellerBalance.ECR < offer.EnergyAmount {
		return fmt.Errorf("saldo insuficiente de ECR do vendedor (%.2f < %.2f)", sellerBalance.ECR, offer.EnergyAmount)
	}
	if buyerBalance.ENGT < offer.TotalPrice {
		return fmt.Errorf("saldo insuficiente de ENGT do comprador (%.2f < %.2f)", buyerBalance.ENGT, offer.TotalPrice)
	}

	log.Printf("DEBUG: Executando liquidação - %s vende %.2f ECR por %.2f ENGT para %s",
		offer.SellerID, offer.EnergyAmount, offer.TotalPrice, buyerID)

	sellerBalance.ECR -= offer.EnergyAmount
	sellerBalance.ENGT += offer.TotalPrice

	buyerBalance.ECR += offer.EnergyAmount
	buyerBalance.ENGT -= offer.TotalPrice

	if err := ctx.GetStub().PutState(balanceKey(offer.SellerID), mustMarshal(sellerBalance)); err != nil {
		return fmt.Errorf("erro ao salvar saldo do vendedor: %v", err)
	}
	if err := ctx.GetStub().PutState(balanceKey(buyerID), mustMarshal(buyerBalance)); err != nil {
		sellerBalance.ECR += offer.EnergyAmount
		sellerBalance.ENGT -= offer.TotalPrice
		_ = ctx.GetStub().PutState(balanceKey(offer.SellerID), mustMarshal(sellerBalance))
		return fmt.Errorf("falha ao salvar saldo do comprador: %v", err)
	}

	// FIX: Also update agent records after spot market transaction
	sellerAgent, err := s.GetAgent(ctx, offer.SellerID)
	if err == nil {
		sellerAgent.ECRBalance = sellerBalance.ECR
		sellerAgent.ENGTBalance = sellerBalance.ENGT
		sellerAgentJSON, err := json.Marshal(sellerAgent)
		if err == nil {
			if err := ctx.GetStub().PutState(agentKey(offer.SellerID), sellerAgentJSON); err != nil {
				log.Printf("Aviso: não foi possível atualizar saldo do agente vendedor: %v", err)
			}
		}
	}

	buyerAgent, err := s.GetAgent(ctx, buyerID)
	if err == nil {
		buyerAgent.ECRBalance = buyerBalance.ECR
		buyerAgent.ENGTBalance = buyerBalance.ENGT
		buyerAgentJSON, err := json.Marshal(buyerAgent)
		if err == nil {
			if err := ctx.GetStub().PutState(agentKey(buyerID), buyerAgentJSON); err != nil {
				log.Printf("Aviso: não foi possível atualizar saldo do agente comprador: %v", err)
			}
		}
	}

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

	event := map[string]interface{}{
		"offer_id":    offer.ID,
		"seller_id":   offer.SellerID,
		"buyer_id":    offer.BuyerID,
		"energy_kwh":  offer.EnergyAmount,
		"price_total": offer.TotalPrice,
		"settled_at":  offer.SettledAt,
		"seller_ecr":  sellerBalance.ECR,
		"buyer_ecr":   buyerBalance.ECR,
		"seller_engt": sellerBalance.ENGT,
		"buyer_engt":  buyerBalance.ENGT,
		"status":      "SETTLED",
	}
	eventBytes, _ := json.Marshal(event)
	if err := ctx.GetStub().SetEvent("SpotMarket:OfferSettled", eventBytes); err != nil {
		log.Printf("Aviso: falha ao emitir evento OfferSettled: %v", err)
	}

	log.Printf("DEBUG: Liquidação concluída com sucesso para oferta %s", id)
	log.Printf("=== ACCEPT OFFER END ===")
	return nil
}

func (s *CombinedEnergyContract) GetAllOffers(ctx contractapi.TransactionContextInterface) ([]*Offer, error) {
	startKey, endKey := getRangeForPrefix(OfferKeyPrefix)
	it, err := ctx.GetStub().GetStateByRange(startKey, endKey)
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
		// Prefix check is now implicit due to range query
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

// ========================================================
// Contract Market Methods - FIXED VERSION
// ========================================================

// FIXED: CreateSupplyContract with agent balance synchronization
func (s *CombinedEnergyContract) CreateSupplyContract(
	ctx contractapi.TransactionContextInterface,
	id, sellerID, buyerID string,
	energyTotal, pricePerKWh float64,
	startDate, endDate, settlementFreq string,
	sellerCollateralECR, buyerCollateralENGT float64,
) error {
	// validações básicas
	if strings.TrimSpace(id) == "" || strings.TrimSpace(sellerID) == "" || strings.TrimSpace(buyerID) == "" {
		return fmt.Errorf("id, sellerID e buyerID são obrigatórios")
	}
	if energyTotal <= 0 || pricePerKWh <= 0 {
		return fmt.Errorf("energyTotal e pricePerKWh devem ser positivos")
	}

	// existir
	exists, err := s.ContractExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("contrato já existe: %s", id)
	}
	if ok, _ := s.AgentExists(ctx, sellerID); !ok {
		return fmt.Errorf("sellerID não encontrado: %s", sellerID)
	}
	if ok, _ := s.AgentExists(ctx, buyerID); !ok {
		return fmt.Errorf("buyerID não encontrado: %s", buyerID)
	}

	// ler saldos para validar colateral
	sellerBal, err := s.GetBalance(ctx, sellerID)
	if err != nil {
		return fmt.Errorf("erro ao ler saldo do vendedor: %v", err)
	}
	buyerBal, err := s.GetBalance(ctx, buyerID)
	if err != nil {
		return fmt.Errorf("erro ao ler saldo do comprador: %v", err)
	}

	if sellerCollateralECR < 0 || buyerCollateralENGT < 0 {
		return fmt.Errorf("colaterais não podem ser negativos")
	}
	if sellerBal.ECR < sellerCollateralECR {
		return fmt.Errorf("vendedor sem ECR suficiente para colateral (%.2f < %.2f)", sellerBal.ECR, sellerCollateralECR)
	}
	if buyerBal.ENGT < buyerCollateralENGT {
		return fmt.Errorf("comprador sem ENGT suficiente para colateral (%.2f < %.2f)", buyerBal.ENGT, buyerCollateralENGT)
	}

	// debitar colaterais (se > 0) e gravar escrows
	now := time.Now().UTC().Format(time.RFC3339)

	if sellerCollateralECR > 0 {
		sellerBal.ECR -= sellerCollateralECR
		if err := ctx.GetStub().PutState(balanceKey(sellerID), mustMarshal(sellerBal)); err != nil {
			return fmt.Errorf("erro ao salvar saldo vendedor: %v", err)
		}

		// FIX: Update agent record
		sellerAgent, err := s.GetAgent(ctx, sellerID)
		if err == nil {
			sellerAgent.ECRBalance = sellerBal.ECR
			sellerAgentJSON, err := json.Marshal(sellerAgent)
			if err == nil {
				if err := ctx.GetStub().PutState(agentKey(sellerID), sellerAgentJSON); err != nil {
					log.Printf("Aviso: não foi possível atualizar saldo do agente vendedor: %v", err)
				}
			}
		}

		sellerEscrow := Escrow{
			ContractID: id, OwnerID: sellerID, Token: "ECR", Amount: sellerCollateralECR, CreatedAt: now,
		}
		if err := ctx.GetStub().PutState(escrowKey(id, "seller", "ECR"), mustMarshal(sellerEscrow)); err != nil {
			return fmt.Errorf("erro ao gravar escrow vendedor: %v", err)
		}
	}

	if buyerCollateralENGT > 0 {
		buyerBal.ENGT -= buyerCollateralENGT
		if err := ctx.GetStub().PutState(balanceKey(buyerID), mustMarshal(buyerBal)); err != nil {
			return fmt.Errorf("erro ao salvar saldo comprador: %v", err)
		}

		// FIX: Update agent record
		buyerAgent, err := s.GetAgent(ctx, buyerID)
		if err == nil {
			buyerAgent.ENGTBalance = buyerBal.ENGT
			buyerAgentJSON, err := json.Marshal(buyerAgent)
			if err == nil {
				if err := ctx.GetStub().PutState(agentKey(buyerID), buyerAgentJSON); err != nil {
					log.Printf("Aviso: não foi possível atualizar saldo do agente comprador: %v", err)
				}
			}
		}

		buyerEscrow := Escrow{
			ContractID: id, OwnerID: buyerID, Token: "ENGT", Amount: buyerCollateralENGT, CreatedAt: now,
		}
		if err := ctx.GetStub().PutState(escrowKey(id, "buyer", "ENGT"), mustMarshal(buyerEscrow)); err != nil {
			return fmt.Errorf("erro ao gravar escrow comprador: %v", err)
		}
	}

	// construir e salvar contrato
	contract := SupplyContract{
		ID:               id,
		SellerID:         sellerID,
		BuyerID:          buyerID,
		EnergyTotal:      energyTotal,
		DeliveredTotal:   0,
		UnsettledEnergy:  0,
		PricePerKWh:      pricePerKWh,
		TotalValue:       energyTotal * pricePerKWh,
		SellerCollateral: sellerCollateralECR,
		BuyerCollateral:  buyerCollateralENGT,
		StartDate:        startDate,
		EndDate:          endDate,
		SettlementFreq:   settlementFreq,
		Status:           "ACTIVE",
		CreatedAt:        now,
		LastSettleAt:     "",
	}
	if err := ctx.GetStub().PutState(contractKey(id), mustMarshal(contract)); err != nil {
		return fmt.Errorf("erro ao salvar contrato: %v", err)
	}

	// evento
	ev := map[string]interface{}{
		"id": id, "seller": sellerID, "buyer": buyerID,
		"energy_total": energyTotal, "price": pricePerKWh,
		"seller_collateral_ecr": sellerCollateralECR,
		"buyer_collateral_engt": buyerCollateralENGT,
		"created_at":            now,
	}
	_ = ctx.GetStub().SetEvent("Supply:ContractCreated", mustMarshal(ev))
	return nil
}

func (s *CombinedEnergyContract) ContractExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(contractKey(id))
	if err != nil {
		return false, fmt.Errorf("erro ao verificar contrato: %v", err)
	}
	return data != nil, nil
}

func (s *CombinedEnergyContract) GetContract(ctx contractapi.TransactionContextInterface, id string) (*SupplyContract, error) {
	data, err := ctx.GetStub().GetState(contractKey(id))
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar contrato: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("contrato não encontrado: %s", id)
	}
	var c SupplyContract
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("erro ao desserializar contrato: %v", err)
	}
	return &c, nil
}

func (s *CombinedEnergyContract) GetAllContracts(ctx contractapi.TransactionContextInterface) ([]*SupplyContract, error) {
	it, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("erro ao iterar ledger: %v", err)
	}
	defer it.Close()

	var out []*SupplyContract
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, fmt.Errorf("erro ao ler item: %v", err)
		}
		if !strings.HasPrefix(kv.Key, ContractKeyPrefix) {
			continue
		}
		var c SupplyContract
		if json.Unmarshal(kv.Value, &c) == nil && c.ID != "" {
			out = append(out, &c)
		}
	}
	return out, nil
}

// ReportDelivery adiciona kWh entregues ao contrato.
func (s *CombinedEnergyContract) ReportDelivery(ctx contractapi.TransactionContextInterface, id string, deliveredKWh float64) error {
	if deliveredKWh <= 0 {
		return fmt.Errorf("deliveredKWh deve ser positivo")
	}
	contract, err := s.GetContract(ctx, id)
	if err != nil {
		return err
	}
	if contract.Status != "ACTIVE" {
		return fmt.Errorf("contrato não está ativo")
	}

	// limitar para não exceder o total contratado (opcional)
	maxAdd := contract.EnergyTotal - contract.DeliveredTotal
	if maxAdd < 0 {
		maxAdd = 0
	}
	add := deliveredKWh
	if add > maxAdd {
		add = maxAdd
	}

	contract.DeliveredTotal += add
	contract.UnsettledEnergy += add

	if err := ctx.GetStub().PutState(contractKey(id), mustMarshal(contract)); err != nil {
		return fmt.Errorf("erro ao atualizar contrato: %v", err)
	}

	ev := map[string]interface{}{"contract_id": id, "delivered_add": add, "delivered_total": contract.DeliveredTotal, "unsettled": contract.UnsettledEnergy}
	_ = ctx.GetStub().SetEvent("Supply:DeliveryReported", mustMarshal(ev))
	return nil
}

// FIXED: SettleContractPeriod with agent balance synchronization
func (s *CombinedEnergyContract) SettleContractPeriod(ctx contractapi.TransactionContextInterface, id string, kwhToSettle float64) error {
	if kwhToSettle <= 0 {
		return fmt.Errorf("kwhToSettle deve ser positivo")
	}
	contract, err := s.GetContract(ctx, id)
	if err != nil {
		return err
	}
	if contract.Status != "ACTIVE" {
		return fmt.Errorf("contrato não está ativo")
	}
	if contract.UnsettledEnergy <= 0 {
		return fmt.Errorf("não há energia pendente para liquidar")
	}

	// limitar ao que está pendente
	amount := kwhToSettle
	if amount > contract.UnsettledEnergy {
		amount = contract.UnsettledEnergy
	}

	payTotal := amount * contract.PricePerKWh

	// saldos atuais
	sellerBal, err := s.GetBalance(ctx, contract.SellerID)
	if err != nil {
		return fmt.Errorf("erro ao ler saldo vendedor: %v", err)
	}
	buyerBal, err := s.GetBalance(ctx, contract.BuyerID)
	if err != nil {
		return fmt.Errorf("erro ao ler saldo comprador: %v", err)
	}

	// ler escrows
	sellerEscData, _ := ctx.GetStub().GetState(escrowKey(contract.ID, "seller", "ECR"))
	buyerEscData, _ := ctx.GetStub().GetState(escrowKey(contract.ID, "buyer", "ENGT"))
	var sellerEsc Escrow
	var buyerEsc Escrow
	if sellerEscData != nil {
		_ = json.Unmarshal(sellerEscData, &sellerEsc)
	}
	if buyerEscData != nil {
		_ = json.Unmarshal(buyerEscData, &buyerEsc)
	}

	// 1) Entrega de energia: prioriza saldo livre do vendedor; faltando, usa escrow do vendedor
	needECR := amount
	useFreeECR := minF(sellerBal.ECR, needECR)
	sellerBal.ECR -= useFreeECR
	needECR -= useFreeECR

	useEscECR := 0.0
	if needECR > 0 {
		useEscECR = minF(sellerEsc.Amount, needECR)
		sellerEsc.Amount -= useEscECR
		needECR -= useEscECR
	}
	if needECR > 0 {
		return fmt.Errorf("ECR insuficiente (livre+escrow) para liquidar %.2f kWh", amount)
	}

	// Creditar ECR no comprador
	buyerBal.ECR += amount

	// 2) Pagamento: prioriza saldo livre do comprador; faltando, usa escrow do comprador
	needENGT := payTotal
	useFreeENGT := minF(buyerBal.ENGT, needENGT)
	buyerBal.ENGT -= useFreeENGT
	needENGT -= useFreeENGT

	useEscENGT := 0.0
	if needENGT > 0 {
		useEscENGT = minF(buyerEsc.Amount, needENGT)
		buyerEsc.Amount -= useEscENGT
		needENGT -= useEscENGT
	}
	if needENGT > 0 {
		// rollback simples do crédito ECR no comprador e débitos feitos
		buyerBal.ECR -= amount
		sellerBal.ECR += useFreeECR
		sellerEsc.Amount += useEscECR
		return fmt.Errorf("ENGT insuficiente (livre+escrow) para pagar %.2f", payTotal)
	}

	// Creditar ENGT no vendedor
	sellerBal.ENGT += payTotal

	// Persistir saldos atualizados
	if err := ctx.GetStub().PutState(balanceKey(contract.SellerID), mustMarshal(sellerBal)); err != nil {
		return fmt.Errorf("erro ao salvar saldo vendedor: %v", err)
	}
	if err := ctx.GetStub().PutState(balanceKey(contract.BuyerID), mustMarshal(buyerBal)); err != nil {
		// rollback vendedor se falhar
		sellerBal.ENGT -= payTotal
		sellerBal.ECR += useFreeECR
		sellerEsc.Amount += useEscECR
		_ = ctx.GetStub().PutState(balanceKey(contract.SellerID), mustMarshal(sellerBal))
		return fmt.Errorf("erro ao salvar saldo comprador: %v", err)
	}

	// FIX: Update agent records after successful settlement
	sellerAgent, err := s.GetAgent(ctx, contract.SellerID)
	if err == nil {
		sellerAgent.ECRBalance = sellerBal.ECR
		sellerAgent.ENGTBalance = sellerBal.ENGT
		sellerAgentJSON, err := json.Marshal(sellerAgent)
		if err == nil {
			if err := ctx.GetStub().PutState(agentKey(contract.SellerID), sellerAgentJSON); err != nil {
				log.Printf("Aviso: não foi possível atualizar saldo do agente vendedor: %v", err)
			}
		}
	}

	buyerAgent, err := s.GetAgent(ctx, contract.BuyerID)
	if err == nil {
		buyerAgent.ECRBalance = buyerBal.ECR
		buyerAgent.ENGTBalance = buyerBal.ENGT
		buyerAgentJSON, err := json.Marshal(buyerAgent)
		if err == nil {
			if err := ctx.GetStub().PutState(agentKey(contract.BuyerID), buyerAgentJSON); err != nil {
				log.Printf("Aviso: não foi possível atualizar saldo do agente comprador: %v", err)
			}
		}
	}

	// persistir escrows (se existirem)
	if sellerEscData != nil {
		if err := ctx.GetStub().PutState(escrowKey(contract.ID, "seller", "ECR"), mustMarshal(sellerEsc)); err != nil {
			return fmt.Errorf("erro ao salvar escrow vendedor: %v", err)
		}
	}
	if buyerEscData != nil {
		if err := ctx.GetStub().PutState(escrowKey(contract.ID, "buyer", "ENGT"), mustMarshal(buyerEsc)); err != nil {
			return fmt.Errorf("erro ao salvar escrow comprador: %v", err)
		}
	}

	// atualizar contrato
	contract.UnsettledEnergy -= amount
	contract.LastSettleAt = time.Now().UTC().Format(time.RFC3339)
	if err := ctx.GetStub().PutState(contractKey(contract.ID), mustMarshal(contract)); err != nil {
		return fmt.Errorf("erro ao atualizar contrato: %v", err)
	}

	// evento
	ev := map[string]interface{}{
		"contract_id": contract.ID, "kwh_settled": amount,
		"price_total": payTotal, "last_settle_at": contract.LastSettleAt,
	}
	_ = ctx.GetStub().SetEvent("Supply:ContractSettled", mustMarshal(ev))
	return nil
}

// FIXED: CloseContract with agent balance synchronization
func (s *CombinedEnergyContract) CloseContract(ctx contractapi.TransactionContextInterface, id string) error {
	contract, err := s.GetContract(ctx, id)
	if err != nil {
		return err
	}
	if contract.Status != "ACTIVE" {
		return fmt.Errorf("contrato não está ativo")
	}

	// carregar saldos e escrows
	sellerBal, err := s.GetBalance(ctx, contract.SellerID)
	if err != nil {
		return fmt.Errorf("erro saldo vendedor: %v", err)
	}
	buyerBal, err := s.GetBalance(ctx, contract.BuyerID)
	if err != nil {
		return fmt.Errorf("erro saldo comprador: %v", err)
	}

	// escrows
	sellerEscData, _ := ctx.GetStub().GetState(escrowKey(contract.ID, "seller", "ECR"))
	buyerEscData, _ := ctx.GetStub().GetState(escrowKey(contract.ID, "buyer", "ENGT"))
	var sellerEsc Escrow
	var buyerEsc Escrow
	if sellerEscData != nil {
		_ = json.Unmarshal(sellerEscData, &sellerEsc)
	}
	if buyerEscData != nil {
		_ = json.Unmarshal(buyerEscData, &buyerEsc)
	}

	undelivered := contract.EnergyTotal - contract.DeliveredTotal
	if undelivered < 0 {
		undelivered = 0
	}

	// Caso A: tudo entregue e liquidado (sem pendências)
	if undelivered == 0 && contract.UnsettledEnergy == 0 {
		// devolver colaterais
		if sellerEsc.Amount > 0 {
			sellerBal.ECR += sellerEsc.Amount
			sellerEsc.Amount = 0
			if err := ctx.GetStub().PutState(balanceKey(contract.SellerID), mustMarshal(sellerBal)); err != nil {
				return fmt.Errorf("erro ao devolver escrow do vendedor: %v", err)
			}

			// FIX: Update agent record
			sellerAgent, err := s.GetAgent(ctx, contract.SellerID)
			if err == nil {
				sellerAgent.ECRBalance = sellerBal.ECR
				sellerAgentJSON, err := json.Marshal(sellerAgent)
				if err == nil {
					if err := ctx.GetStub().PutState(agentKey(contract.SellerID), sellerAgentJSON); err != nil {
						log.Printf("Aviso: não foi possível atualizar saldo do agente vendedor: %v", err)
					}
				}
			}

			_ = ctx.GetStub().DelState(escrowKey(contract.ID, "seller", "ECR"))
		}
		if buyerEsc.Amount > 0 {
			buyerBal.ENGT += buyerEsc.Amount
			buyerEsc.Amount = 0
			if err := ctx.GetStub().PutState(balanceKey(contract.BuyerID), mustMarshal(buyerBal)); err != nil {
				return fmt.Errorf("erro ao devolver escrow do comprador: %v", err)
			}

			// FIX: Update agent record
			buyerAgent, err := s.GetAgent(ctx, contract.BuyerID)
			if err == nil {
				buyerAgent.ENGTBalance = buyerBal.ENGT
				buyerAgentJSON, err := json.Marshal(buyerAgent)
				if err == nil {
					if err := ctx.GetStub().PutState(agentKey(contract.BuyerID), buyerAgentJSON); err != nil {
						log.Printf("Aviso: não foi possível atualizar saldo do agente comprador: %v", err)
					}
				}
			}

			_ = ctx.GetStub().DelState(escrowKey(contract.ID, "buyer", "ENGT"))
		}
	} else {
		// Caso B: há entrega pendente → penaliza vendedor com seu escrow ECR
		penaltyECR := sellerEsc.Amount
		if penaltyECR > 0 {
			// transfere ECR do escrow do vendedor para o comprador
			buyerBal.ECR += penaltyECR
			sellerEsc.Amount = 0

			if err := ctx.GetStub().PutState(balanceKey(contract.BuyerID), mustMarshal(buyerBal)); err != nil {
				return fmt.Errorf("erro ao aplicar penalidade (creditar comprador): %v", err)
			}

			// FIX: Update both agent records
			sellerAgent, err := s.GetAgent(ctx, contract.SellerID)
			if err == nil {
				sellerAgent.ECRBalance = sellerBal.ECR
				sellerAgentJSON, err := json.Marshal(sellerAgent)
				if err == nil {
					if err := ctx.GetStub().PutState(agentKey(contract.SellerID), sellerAgentJSON); err != nil {
						log.Printf("Aviso: não foi possível atualizar saldo do agente vendedor: %v", err)
					}
				}
			}

			buyerAgent, err := s.GetAgent(ctx, contract.BuyerID)
			if err == nil {
				buyerAgent.ECRBalance = buyerBal.ECR
				buyerAgentJSON, err := json.Marshal(buyerAgent)
				if err == nil {
					if err := ctx.GetStub().PutState(agentKey(contract.BuyerID), buyerAgentJSON); err != nil {
						log.Printf("Aviso: não foi possível atualizar saldo do agente comprador: %v", err)
					}
				}
			}

			_ = ctx.GetStub().DelState(escrowKey(contract.ID, "seller", "ECR"))
		}
		// Buyer escrow retorna ao comprador (não há débito a aplicar aqui)
		if buyerEsc.Amount > 0 {
			buyerBal.ENGT += buyerEsc.Amount
			buyerEsc.Amount = 0
			if err := ctx.GetStub().PutState(balanceKey(contract.BuyerID), mustMarshal(buyerBal)); err != nil {
				return fmt.Errorf("erro ao devolver escrow do comprador: %v", err)
			}

			// FIX: Update agent record
			buyerAgent, err := s.GetAgent(ctx, contract.BuyerID)
			if err == nil {
				buyerAgent.ENGTBalance = buyerBal.ENGT
				buyerAgentJSON, err := json.Marshal(buyerAgent)
				if err == nil {
					if err := ctx.GetStub().PutState(agentKey(contract.BuyerID), buyerAgentJSON); err != nil {
						log.Printf("Aviso: não foi possível atualizar saldo do agente comprador: %v", err)
					}
				}
			}

			_ = ctx.GetStub().DelState(escrowKey(contract.ID, "buyer", "ENGT"))
		}
	}

	contract.Status = "CLOSED"
	if err := ctx.GetStub().PutState(contractKey(contract.ID), mustMarshal(contract)); err != nil {
		return fmt.Errorf("erro ao fechar contrato: %v", err)
	}

	ev := map[string]interface{}{
		"contract_id": id,
		"undelivered": undelivered,
		"status":      "CLOSED",
	}
	_ = ctx.GetStub().SetEvent("Supply:ContractClosed", mustMarshal(ev))
	return nil
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
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

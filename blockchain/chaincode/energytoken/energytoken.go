package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type EnergyTokenContract struct {
	contractapi.Contract
}

type TokenBalance struct {
	AgentID string  `json:"agentId"`
	ECR     float64 `json:"ecr"`
	ENGT    float64 `json:"engt"`
}

// Mint tokens (ECR or ENGT)
func (c *EnergyTokenContract) Mint(ctx contractapi.TransactionContextInterface, agentID string, tokenType string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("quantidade inválida: %.2f", amount)
	}
	key := fmt.Sprintf("balance:%s", agentID)

	// Verifica saldo atual
	balance := TokenBalance{AgentID: agentID}
	existing, _ := ctx.GetStub().GetState(key)
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
	return ctx.GetStub().PutState(key, data)
}

// Transfer tokens
func (c *EnergyTokenContract) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, tokenType string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("quantidade inválida")
	}

	fromKey := fmt.Sprintf("balance:%s", from)
	toKey := fmt.Sprintf("balance:%s", to)

	fromData, _ := ctx.GetStub().GetState(fromKey)
	toData, _ := ctx.GetStub().GetState(toKey)

	if fromData == nil || toData == nil {
		return fmt.Errorf("um dos agentes não existe")
	}

	var fromBal, toBal TokenBalance
	_ = json.Unmarshal(fromData, &fromBal)
	_ = json.Unmarshal(toData, &toBal)

	switch tokenType {
	case "ECR":
		if fromBal.ECR < amount {
			return fmt.Errorf("saldo insuficiente de ECR")
		}
		fromBal.ECR -= amount
		toBal.ECR += amount
	case "ENGT":
		if fromBal.ENGT < amount {
			return fmt.Errorf("saldo insuficiente de ENGT")
		}
		fromBal.ENGT -= amount
		toBal.ENGT += amount
	default:
		return fmt.Errorf("tokenType inválido")
	}

	// Salva novamente
	ctx.GetStub().PutState(fromKey, mustMarshal(fromBal))
	ctx.GetStub().PutState(toKey, mustMarshal(toBal))
	return nil
}

// Consultar saldo
func (c *EnergyTokenContract) GetBalance(ctx contractapi.TransactionContextInterface, agentID string) (*TokenBalance, error) {
	key := fmt.Sprintf("balance:%s", agentID)
	data, err := ctx.GetStub().GetState(key)
	if err != nil || data == nil {
		return nil, fmt.Errorf("saldo não encontrado para o agente %s", agentID)
	}
	var balance TokenBalance
	_ = json.Unmarshal(data, &balance)
	return &balance, nil
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func main() {
	cc, err := contractapi.NewChaincode(&EnergyTokenContract{})
	if err != nil {
		panic(err)
	}
	if err := cc.Start(); err != nil {
		panic(err)
	}
}

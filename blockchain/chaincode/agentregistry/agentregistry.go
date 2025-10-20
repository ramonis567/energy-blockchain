package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract estrutura principal
type AgentRegistryContract struct {
	contractapi.Contract
}

// Agent representa um participante do sistema
type Agent struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"` // "producer", "consumer", "distributor"
	Balance float64 `json:"balance"`
	Name    string  `json:"name"`    // Novo campo para evolução futura
	Address string  `json:"address"` // Novo campo para evolução futura
}

// AgentRegistrationEvent evento para notificar registro de agente
type AgentRegistrationEvent struct {
	AgentID string `json:"agentId"`
	Type    string `json:"type"`
}

// RegisterAgent registra um novo agente no ledger
func (s *AgentRegistryContract) RegisterAgent(ctx contractapi.TransactionContextInterface,
	id string, agentType string, name string, address string) error {

	// Validações básicas
	if id == "" {
		return fmt.Errorf("ID do agente não pode estar vazio")
	}
	if agentType == "" {
		return fmt.Errorf("tipo do agente não pode estar vazio")
	}

	// Verifica se o agente já existe
	exists, err := s.AgentExists(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência do agente: %v", err)
	}
	if exists {
		return fmt.Errorf("agente já registrado: %s", id)
	}

	// Cria novo agente
	agent := Agent{
		ID:      id,
		Type:    agentType,
		Balance: 0.0,
		Name:    name,
		Address: address,
	}

	// Serializa para JSON
	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("erro ao serializar agente: %v", err)
	}

	// Salva no ledger
	err = ctx.GetStub().PutState(id, agentJSON)
	if err != nil {
		return fmt.Errorf("erro ao salvar agente no ledger: %v", err)
	}

	// Emite evento para notificar o registro
	event := AgentRegistrationEvent{
		AgentID: id,
		Type:    agentType,
	}
	eventJSON, _ := json.Marshal(event)
	err = ctx.GetStub().SetEvent("AgentRegistered", eventJSON)
	if err != nil {
		// Não falha a transação se o evento não for emitido
		log.Printf("Aviso: não foi possível emitir evento: %v", err)
	}

	return nil
}

// AgentExists verifica se um agente existe
func (s *AgentRegistryContract) AgentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	agentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("erro ao acessar ledger: %v", err)
	}
	return agentJSON != nil, nil
}

// GetAgent retorna um agente específico
func (s *AgentRegistryContract) GetAgent(ctx contractapi.TransactionContextInterface, id string) (*Agent, error) {
	agentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar ledger: %v", err)
	}
	if agentJSON == nil {
		return nil, fmt.Errorf("agente não encontrado: %s", id)
	}

	var agent Agent
	err = json.Unmarshal(agentJSON, &agent)
	if err != nil {
		return nil, fmt.Errorf("erro ao desserializar agente: %v", err)
	}

	return &agent, nil
}

// GetAllAgents retorna todos os agentes registrados
func (s *AgentRegistryContract) GetAllAgents(ctx contractapi.TransactionContextInterface) ([]*Agent, error) {
	// Obtém todos os estados do ledger
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre o ledger: %v", err)
	}
	defer resultsIterator.Close()

	var agents []*Agent

	// Itera sobre todos os resultados
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("erro ao ler próximo item: %v", err)
		}

		// Tenta desserializar como Agent
		var agent Agent
		err = json.Unmarshal(queryResponse.Value, &agent)
		if err != nil {
			// Se não for um Agent válido, ignora (pode ser outro tipo de dado)
			continue
		}

		// Verifica se tem ID (campo obrigatório)
		if agent.ID != "" {
			agents = append(agents, &agent)
		}
	}

	return agents, nil
}

// GetAgentsByType retorna agentes filtrados por tipo (para evolução futura)
func (s *AgentRegistryContract) GetAgentsByType(ctx contractapi.TransactionContextInterface, agentType string) ([]*Agent, error) {
	allAgents, err := s.GetAllAgents(ctx)
	if err != nil {
		return nil, err
	}

	var filteredAgents []*Agent
	for _, agent := range allAgents {
		if agent.Type == agentType {
			filteredAgents = append(filteredAgents, agent)
		}
	}

	return filteredAgents, nil
}

// UpdateAgent atualiza informações do agente (para evolução futura)
func (s *AgentRegistryContract) UpdateAgent(ctx contractapi.TransactionContextInterface,
	id string, name string, address string) error {

	agent, err := s.GetAgent(ctx, id)
	if err != nil {
		return err
	}

	// Atualiza campos permitidos
	agent.Name = name
	agent.Address = address

	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("erro ao serializar agente: %v", err)
	}

	return ctx.GetStub().PutState(id, agentJSON)
}

// GetAgentCount retorna o número total de agentes (para evolução futura)
func (s *AgentRegistryContract) GetAgentCount(ctx contractapi.TransactionContextInterface) (int, error) {
	agents, err := s.GetAllAgents(ctx)
	if err != nil {
		return 0, err
	}
	return len(agents), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&AgentRegistryContract{})
	if err != nil {
		log.Panicf("Erro ao criar chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Erro ao iniciar chaincode: %v", err)
	}
}

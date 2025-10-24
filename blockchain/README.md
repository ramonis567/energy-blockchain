# Camada Blockchain — Sistema de Gerenciamento de Transações Energéticas com foco em Microgrids (MMGD)

Este diretório contém toda a **camada blockchain** do sistema de gerenciamento de transações energéticas, responsável por registrar, liquidar e auditar as trocas de energia entre agentes (consumidores, produtores, prosumers e baterias) em uma camada blockchain permissionada **Hyperledger Fabric**.

A camada é composta por múltiplos **smart contracts (chaincodes)** interconectados, que implementam os mecanismos de mercado **Spot**, **Créditos** e **Contratos Bilaterais**, além de contratos de base e de liquidação.

---

## 🧩 Estrutura da Camada Blockchain

blockchain/
├── chaincode/
│ ├── agentregistry/ # Cadastro e gestão de agentes
│ │ ├── agentregistry.go
│ │ ├── go.mod
│ │ └── README.md
│ │
│ ├── energytoken/ # Token energético (ECR) e financeiro (ENGT)
│ │ ├── energytoken.go
│ │ ├── go.mod
│ │ └── README.md
│ │
│ ├── spotmarket/ # Mercado à vista de energia (transações instantâneas)
│ │ ├── spotmarket.go
│ │ ├── go.mod
│ │ └── README.md
│ │
│ ├── creditmarket/ # Mercado de créditos energéticos (compensação temporal)
│ │ ├── creditmarket.go
│ │ ├── go.mod
│ │ └── README.md
│ │
│ ├── supplycontract/ # Contratos bilaterais e PPAs locais
│ │ ├── supplycontract.go
│ │ ├── go.mod
│ │ └── README.md
│ │
│ ├── settlementengine/ # Liquidação e reconciliação de saldos
│ │ ├── settlementengine.go
│ │ ├── go.mod
│ │ └── README.md
│ │
│ ├── utils/ # Funções auxiliares reutilizáveis
│ │ ├── ledger_utils.go
│ │ └── token_types.go
│ │
│ └── shared/ # Estruturas e tipos compartilhados entre contratos
│ ├── structs.go
│ ├── events.go
│ └── README.md
│
├── scripts/
│ ├── deploy-chaincode.sh # Automação de deploy para qualquer contrato
│ ├── test-agentregistry.sh
│ ├── test-spotmarket.sh
│ ├── test-supplycontract.sh
│ └── test-creditmarket.sh
│ └── start.sh # Inicia a rede de teste e faz deploy do primeiro chaincode
│ └── stop.sh # Finaliza a rede e remove serviços auxiliares
│ └── update-go.sh # Atualiza e configura o Go 1.21+
│ └── README.md  # Anotações gerais sobre a simulação  (arquivo temporário)
└── README.md # (este arquivo)

## 🧠 Arquitetura dos Smart Contracts

| **Smart Contract** | **Função Principal** | **Tipo** | **Depende de** |
|--------------------|----------------------|----------|----------------|
| `AgentRegistry`    | Cadastro e gestão de agentes | Base | — |
| `EnergyToken`      | Emissão e transferência de tokens ECR/ENGT | Base | AgentRegistry |
| `SpotMarket`       | Transações instantâneas de energia | Mercado | AgentRegistry, EnergyToken |
| `CreditMarket`     | Créditos energéticos e compensação temporal | Mercado | AgentRegistry, EnergyToken |
| `SupplyContract`   | Contratos bilaterais de longo prazo | Mercado | AgentRegistry, EnergyToken |
| `SettlementEngine` | Liquidação e reconciliação geral | Core | Todos os anteriores |

---

## 💰 Tokens Utilizados

- **ECR (Energy Credit Token):** Representa créditos energéticos em kWh, não conversíveis em moeda.  
- **ENGT (Energy Trade Token):** Token financeiro fungível usado nas transações spot e contratuais.

Cada agente possui saldos em ambos os tokens, registrados no ledger de forma imutável.

---

## 🚀 Execução e Deploy

### Atualizar o Go (apenas uma vez)
```bash
cd blockchain
chmod +x update-go.sh
./update-go.sh
```

### Iniciar a rede Fabric e fazer deploy do primeiro chaincode
```bash
cd scripts
chmod +x start.sh
./start.sh
```

Isso inicia:
Rede de teste do Fabric (mychannel)
Chaincode inicial (por padrão, agentregistry)
Serviços auxiliares via Docker (Mosquitto, InfluxDB, Grafana)

### Deploy de novos contratos
```bash
./scripts/deploy-chaincode.sh <nome>
# Exemplo:
./scripts/deploy-chaincode.sh energytoken
```

### Testes via CLI
```bash



```

## Testes e Validação  (REVISAR)

Após o deploy de cada contrato:
Use os scripts em scripts/test-*.sh para validar as funções básicas.
Verifique se os eventos (AgentRegistered, TokenMinted, OfferCreated, etc.) aparecem nos logs do peer.
Valide integridade dos saldos com o SettlementEngine.

## Integração com Backend 

A comunicação entre a camada blockchain e a aplicação é feita via Fabric SDK for Node.js, no diretório /backend.
As rotas REST previstas são:

Rota	Método	Função
/api/agents	POST / GET / GET/:id	Registrar e consultar agentes
/api/tokens/mint	POST	Emitir tokens ECR/ENGT
/api/tokens/transfer	POST	Transferir tokens entre agentes
/api/spot/offers	POST / GET	Ofertas no mercado spot
/api/credits	POST / GET	Gerenciar créditos energéticos
/api/contracts	POST / GET	Criar e executar contratos bilaterais
/api/market/summary	GET	Consultar relatórios consolidados


Go: 1.21.5+
Hyperledger Fabric: 2.5+
Docker e Docker Compose
WSL2 + Ubuntu 22.04 (recomendado para Windows)
Fabric Samples: Instalado em ~/go/src/github.com/fabric-samples/test-network


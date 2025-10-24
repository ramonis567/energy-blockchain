Verificar caminho dos chaincodes

echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

Iniciar dependências e gerar .sum
go mod init agentregistry
go mod tidy
go mod vendor
go build

cd ~/go/src/github.com/fabric-samples/test-network

# Set the PATH to include Fabric binaries
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
# Set TLS enabled
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

Teste que funcionaram:
# Registrar agentes
# Produtor
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n agentregistry \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"RegisterAgent","Args":["producer1", "producer", "Solar Farm North", "Sun Street 100"]}'

# Consumidor
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n agentregistry \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"RegisterAgent","Args":["consumer1", "consumer", "Shopping Center", "Main Avenue 500"]}'

# Distribuidor
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n agentregistry \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"RegisterAgent","Args":["distrib1", "distributor", "Energy Distributor", "Power Street 200"]}'


# Consultar todos os agentes
# Listar todos os agentes
peer chaincode query \
  -C mychannel \
  -n agentregistry \
  -c '{"function":"GetAllAgents","Args":[]}'

# Consultar agente específico
peer chaincode query \
  -C mychannel \
  -n agentregistry \
  -c '{"function":"GetAgent","Args":["producer1"]}'

# Filtrar por tipo
peer chaincode query \
  -C mychannel \
  -n agentregistry \
  -c '{"function":"GetAgentsByType","Args":["producer"]}'

# Contar agentes
peer chaincode query \
  -C mychannel \
  -n agentregistry \
  -c '{"function":"GetAgentCount","Args":[]}'

# Atualizar agente específico
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n agentregistry \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"UpdateAgent","Args":["producer1", "Solar Farm North EXPANDED", "New Sun Street 150"]}'

# Listar funções disponíveis (se suportado)
peer chaincode query \
  -C mychannel \
  -n agentregistry \
  -c '{"function":"org.hyperledger.fabric:GetMetadata","Args":[]}'
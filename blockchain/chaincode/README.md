Iniciar dependências e gerar .sum
go mod init agentregistry
go mod tidy
go mod vendor
go build

cd ~/go/src/github.com/fabric-samples/test-network

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

peer lifecycle chaincode queryinstalled
peer lifecycle chaincode querycommitted -C mychannel

# Register agents
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"RegisterAgent","Args":["producer1","producer","Solar Plant North","Sun Street 10"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"RegisterAgent","Args":["consumer1","consumer","Shopping Center","Main Avenue 500"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"RegisterAgent","Args":["prosumer1","prosumer","Smart Building","Tech Park 123"]}'


# Mint ECR (energia) para o produtor
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"Mint","Args":["producer1","ECR","200"]}'

# Mint ENGT (moeda) para o consumidor
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"Mint","Args":["consumer1","ENGT","1000"]}'

# Mint ECR e ENGT para prosumer
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"Mint","Args":["prosumer1","ECR","100"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"Mint","Args":["prosumer1","ENGT","500"]}'

# Checking balances
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetBalance","Args":["producer1"]}'
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetBalance","Args":["consumer1"]}'
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetBalance","Args":["prosumer1"]}'

# Criando oferta de energia do produtor (50 kWh a 4 ENGT/kWh)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"CreateOffer","Args":["offer1","producer1","50","4"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"CreateOffer","Args":["offer2","prosumer1","25","3.5"]}'

# Listing open offers
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetAllOffers","Args":[]}'


# ACEITAR OFERTA DO PRODUTOR PELO CONSUMIDOR
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"AcceptOffer","Args":["offer1","consumer1"]}'


# ACEITAR OFERTA DO PROSUMER PELO CONSUMIDOR
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"AcceptOffer","Args":["offer2","consumer1"]}'

# Listing all offers (should be settled)
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetAllOffers","Args":[]}'

# Show agent full info
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetAgentFullInfo","Args":["producer1"]}'


# All registered agents:
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetAllAgents","Args":[]}'

# Producers only:
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetAgentsByType","Args":["producer"]}'

# Agent count:
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetAgentCount","Args":[]}'


# Trying to accept already settled offer (should fail):
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"AcceptOffer","Args":["offer1","consumer1"]}'

# Trying to accept non-existent offer (should fail):
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"AcceptOffer","Args":["offer999","consumer1"]}'

# Trying to register duplicate agent (should fail):
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"RegisterAgent","Args":["producer1","producer","Duplicate","Address"]}' 

# Trying to mint negative amount (should fail):
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"Mint","Args":["producer1","ECR","-100"]}'


# Transferring ECR from producer to consumer:
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"Transfer","Args":["producer1","consumer1","ECR","10"]}'

# Transferring ENGT from consumer to prosumer:
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile $ORDERER_CA \
-C mychannel -n main_cc \
--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
-c '{"function":"Transfer","Args":["consumer1","prosumer1","ENGT","25"]}'

# Final balances after direct transfers:
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetBalance","Args":["producer1"]}'
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetBalance","Args":["consumer1"]}'
peer chaincode query -C mychannel -n main_cc -c '{"function":"GetBalance","Args":["prosumer1"]}'
https://hyperledger-fabric.readthedocs.io/en/release-2.5/getting_started.html

Prerequisites: 
Docker latest version
WSL2 (also download and install a linux distro like Ubuntu 22.04)

Install using this guide: https://hyperledger-fabric.readthedocs.io/en/release-2.5/install.html

To run test network:
Go to fabric-samples folder, located in ~home/go/src/github.com/fabric-samples/test-network$ open docker and run:

./network.sh up createChannel -c mychannel -ca

To see docker containers that are already running:
docker ps

To see peer logs:
docker logs peer0.org1.example.com



TO CHANGE:

Para rodar com chaincode: ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go/ -ccl go

Executar transações: export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

Adicionar ativos: peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "$PEER0_ORG1_CA" -c '{"function":"CreateAsset","Args":["asset1", "blue", "5", "Tom", "100"]}'


Consultar ativo: peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","asset1"]}'

Rodar meu proprio chaincode: ./network.sh deployCC -ccn energycc -ccp ./chaincode/energy-chaincode/ -ccl go

Transações: peer chaincode invoke -C mychannel -n energycc -c '{"function":"RegistrarMedicao","Args":["med1", "bus01", "1.5", "geracao"]}'

Criar .env para armazenar caminhos das credenciais
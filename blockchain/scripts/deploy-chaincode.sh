#!/bin/bash
# Uso: ./scripts/deploy-chaincode.sh <name> <path>
NAME=$1

cd ~/go/src/github.com/fabric-samples/test-network

./network.sh deployCC \
    -ccn $NAME \
    -ccp ~/energy-blockchain/blockchain/chaincode/$NAME \
    -ccl go \
    -ccv 1.0

#!/bin/bash
# setenv.sh - Fabric environment setup
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
# Org1 environment
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "✅ Chaincode $NAME implantado a partir de $NAME"

#!/bin/bash
set -e

echo "🚀 Starting Hyperledger Fabric test network..."
# Path to test-network
cd ~/go/src/github.com/fabric-samples/test-network
#  Start newtowk and create a channel 'mychannel'

# ./network.sh up createChannel -c mychannel -ca             -< test
./network.sh up createChannel
# Deploy chaincode
# ./network.sh deployCC -ccn creditmarket -ccp /home/ramon/energy-blockchain/blockchain/chaincode/agentregistry -ccl go -ccv 1.2
cd ~/energy-blockchain/blockchain/scripts
./deploy-chaincode.sh

cd ~/go/src/github.com/fabric-samples/test-network

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem


echo "✅ Blockchain ok!"
echo "Channel: mychannel"
echo "Chaincode: energycc"

echo "🚀 Starting Mosquitto, InfluxDB and Grafana..."
cd ~/energy-blockchain
docker compose up -d

echo "✅ Services ok!"
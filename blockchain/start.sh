#!/bin/bash
set -e

echo "🚀 Starting Hyperledger Fabric test network..."
# Path to test-network
cd ~/go/src/github.com/fabric-samples/test-network
#  Start newtowk and create a channel 'mychannel'
./network.sh up createChannel -c mychannel -ca
# Deploy chaincode
./network.sh deployCC -ccn creditmarket -ccp /home/ramon/energy-blockchain/blockchain/chaincode/creditmarket -ccl go -ccv 1.5

echo "✅ Blockchain ok!"
echo "Channel: mychannel"
echo "Chaincode: energycc"

echo "🚀 Starting Mosquitto, InfluxDB and Grafana..."
cd ~/energy-blockchain
docker compose up -d

echo "✅ Services ok!"
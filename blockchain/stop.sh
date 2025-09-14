#!/bin/bash
set -e

echo "🛑 Stoping Hyperledger Fabric test network..."
# Path to test-network
cd ~/go/src/github.com/fabric-samples/test-network
./network.sh down


echo "🛑 Stoping services..."
cd ~/energy-blockchain
docker compose down

echo "✅ Enviroment down"

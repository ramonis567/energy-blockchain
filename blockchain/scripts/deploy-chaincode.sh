#!/bin/bash
# Uso: ./scripts/deploy-chaincode.sh <name> <path>
NAME=$1

cd ~/go/src/github.com/fabric-samples/test-network

./network.sh deployCC \
    -ccn main_cc \
    -ccp ~/energy-blockchain/blockchain/chaincode \
    -ccl go \
    -ccv 1.1

echo "✅ Chaincode main implantado a partir de chaincode"

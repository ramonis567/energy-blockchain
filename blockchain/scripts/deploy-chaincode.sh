#!/bin/bash
# Uso: ./scripts/deploy-chaincode.sh <name> <path>
NAME=$1

cd ~/go/src/github.com/fabric-samples/test-network

./network.sh deployCC \
    -ccn $NAME \
    -ccp ~/energy-blockchain/blockchain/chaincode/$NAME \
    -ccl go \
    -ccv 1.0

echo "✅ Chaincode $NAME implantado a partir de $NAME"

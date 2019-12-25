package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/blclibs/impl"

)

var logger = shim.NewLogger("faceid.chaincode")

func main() {
	cc := impl.NewFabricChaincode()
	cc.Register(NewFaceIDService())

	err := shim.Start(cc)
	if err != nil {
		logger.Errorf("Error starting chaincode - %s", err)
	}
}

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/snlansky/apk/rpc"
)

var log = shim.NewLogger("servicejournal.chaincode")

func main() {
	helper := rpc.NewHelper(InitFunction)

	helper.RegisterByRule(DefaultTicketService, rpc.DefaultRule)
	helper.RegisterByRule(DefaultNodeService, rpc.DefaultRule)

	err := shim.Start(helper)
	if err != nil {
		log.Errorf("Error starting chaincode - %s", err)
	}
}

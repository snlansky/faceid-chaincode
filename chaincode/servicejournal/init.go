package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snlansky/apk/utils"
)

const SUPER = "SuperAccount"

func InitFunction(stub shim.ChaincodeStubInterface) pb.Response {

	addr, err := utils.GetAddress(stub)
	if err != nil {
		logger.Errorf("get super account cert error：%s", err)
		panic(err)
	}
	err = stub.PutState(SUPER, addr)
	if err != nil {
		logger.Errorf("save super account error：%s", err)
		panic(err)
	}

	logger.Infof("init by super account:%s", string(addr))
	return shim.Success([]byte("SUCCESS"))
}

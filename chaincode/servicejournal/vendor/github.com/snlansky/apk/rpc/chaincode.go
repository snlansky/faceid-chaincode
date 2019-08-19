package rpc

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"errors"
)

type ChaincodeClient struct {
	Name string
}

func NewChaincodeClient(name string) *ChaincodeClient {
	return &ChaincodeClient{Name: name}
}

func (c *ChaincodeClient) Invoke(stub shim.ChaincodeStubInterface, method string, params ... interface{}) ([]byte, error) {
	buf, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	ret := stub.InvokeChaincode(c.Name, [][]byte{[]byte(method), buf}, stub.GetChannelID())
	if ret.GetStatus() != int32(shim.OK) {
		return nil, errors.New(ret.GetMessage())
	}
	return ret.GetPayload(), nil
}

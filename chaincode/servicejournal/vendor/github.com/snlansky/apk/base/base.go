package base

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/rpc"
	"github.com/snlansky/apk/utils"
	"strings"
)

const Separator = ":"

func GetAddress(stub shim.ChaincodeStubInterface) (string, error) {
	addr, err := utils.GetAddress(stub)
	if err != nil {
		return "", err
	}
	return string(addr), nil
}

func MustGetAddress(stub shim.ChaincodeStubInterface) string {
	addr, err := utils.GetAddress(stub)
	rpc.Check(err, rpc.ERR_INVALID_CERT)
	return string(addr)
}

func GetTimestamp(stub shim.ChaincodeStubInterface) (int64, error) {
	ts, err := stub.GetTxTimestamp()
	if err != nil {
		return 0, err
	}
	return ts.GetSeconds(), nil
}

func MustGetTimestamp(stub shim.ChaincodeStubInterface) int64 {
	ts, err := stub.GetTxTimestamp()
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
	return ts.GetSeconds()
}

func CreateEvent(stub shim.ChaincodeStubInterface, appName, eventName, value string) error {
	return stub.SetEvent(strings.Join([]string{appName, eventName}, Separator), []byte(value))
}

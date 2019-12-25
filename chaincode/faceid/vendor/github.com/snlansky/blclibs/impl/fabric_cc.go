package impl

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snlansky/blclibs"
	"github.com/snlansky/glibs/rpc"
	"time"
)

var logger = shim.NewLogger("chaincode")

type FabricChaincode struct {
	rpc rpc.Rpc
}

func NewFabricChaincode() *FabricChaincode {
	return &FabricChaincode{rpc: rpc.New()}
}

func (cc *FabricChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("SUCCESS"))
}

func (cc *FabricChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	stb := NewFabricContractStub(stub)
	args := stb.GetArgs()
	if len(args) <= 0 || len(args) > 2 {
		return shim.Error(rpc.ERR_PARAM_INVALID)
	}

	method := string(args[0])
	var param []interface{}

	if len(args) == 2 {
		err := json.Unmarshal(args[1], &param)
		if err != nil {
			logger.Errorf("json.Unmarshal error:%s, date:%s", err.Error(), string(args[1]))
			return shim.Error(rpc.ERR_JSON_UNMARSHAL)
		}
	}

	addr, err := stb.GetAddress()
	if err != nil {
		logger.Errorf("auth user failed, error:%s", err.Error())
		return shim.Error("ERR_INVALID_CERT")
	}

	logger.Infof(">>> address:%s, method:%s, params:%v", addr, method, param)

	req := &rpc.Request{
		Method: method,
		Params: param,
	}

	return cc.handler(stb, req)
}

func (cc *FabricChaincode) handler(stub blclibs.IContractStub, req *rpc.Request) pb.Response {
	var (
		ret interface{}
		err error
	)

	startTime := time.Now()

	ret, err = cc.rpc.Handler(req, stub)
	if err != nil {
		logger.Errorf("response error:%s", err.Error())
		return shim.Error(err.Error())
	}

	if ret == nil {
		logger.Infof("process takes %v, response success:null", time.Since(startTime))
		return shim.Success(nil)
	}

	buf, err := json.Marshal(ret)
	if err != nil {
		logger.Errorf("response error:%s", err.Error())
		return shim.Error(rpc.ERR_JSON_MARSHAL)
	}
	logger.Infof("process takes %v, response success:%s", time.Since(startTime), string(buf))
	return shim.Success(buf)
}

func (cc *FabricChaincode) Register(i interface{}) {
	cc.rpc.Register(i, rpc.DefaultRule)
}

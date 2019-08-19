package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snlansky/apk/utils"
	"reflect"
	"time"
)

type Filter func(stub shim.ChaincodeStubInterface, request *Request) error

type helper struct {
	init    func(stub shim.ChaincodeStubInterface) pb.Response
	rpc     *Rpc
	filters []Filter
	batch   bool
}

func NewHelper(init func(stub shim.ChaincodeStubInterface) pb.Response, filters ...Filter) *helper {
	return &helper{init: init, rpc: NewRpcHelper(), filters: filters, batch: false}
}

func (h *helper) Init(stub shim.ChaincodeStubInterface) pb.Response {
	if h.init == nil {
		return shim.Error(ERR_NOT_FIND_INIT_FUNCTION)
	}
	return h.init(stub)
}

func (h *helper) EnableBatchInvoke() {
	h.batch = true
}

func (h *helper) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	ctx := newStubImpl(stub)
	args := ctx.GetArgs()
	if len(args) <= 0 || len(args) > 2 {
		return shim.Error(ERR_PARAM_INVALID)
	}

	method := string(args[0])
	var param []interface{}

	if len(args) == 2 {
		err := json.Unmarshal(args[1], &param)
		if err != nil {
			logger.Errorf("json.Unmarshal error:%s, date:%s", err.Error(), string(args[1]))
			return shim.Error(ERR_JSON_UNMARSHAL)
		}
	}

	addr, err := utils.GetAddress(stub)
	if err != nil {
		logger.Errorf("auth user failed, error:%s", err.Error())
		return shim.Error(ERR_INVALID_CERT)
	}

	logger.Infof(">>> address:%s, method:%s, params:%v", addr, method, param)

	req := &Request{
		Method: method,
		Params: param,
	}

	for _, f := range h.filters {
		err := f(ctx, req)
		if err != nil {
			logger.Errorf("filter error:%v, request:%v", err.Error(), req)
			return shim.Error(err.Error())
		}
	}
	return h.handler(ctx, req)
}

func (h *helper) handler(stub shim.ChaincodeStubInterface, req *Request) pb.Response {
	var (
		ret interface{}
		err error
	)

	startTime := time.Now()

	if req.Method == SERVICE_BATCH_INVOKE && h.batch {
		ret, err = h.batchInvoke(stub, req)
	} else {
		ret, err = h.simpleInvoke(stub, req)
	}

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
		return shim.Error(ERR_JSON_MARSHAL)
	}
	logger.Infof("process takes %v, response success:%s", time.Since(startTime), string(buf))
	return shim.Success(buf)
}

func (h *helper) batchInvoke(stub shim.ChaincodeStubInterface, req *Request) ([]*Response, error) {
	typeOfRequest := reflect.TypeOf(req)
	rets := make([]*Response, len(req.Params))
	for i, param := range req.Params {
		value, ok := convertParam(param, typeOfRequest)
		if !ok {
			return nil, fmt.Errorf("[RPC]: invoke index [%d], convert param failed: %v is not <Request> type", i, param)
		}
		inv, ok := value.Interface().(*Request)
		if !ok {
			return nil, fmt.Errorf("[RPC]: invoke index [%d], param %v is not <Request> type", i, param)
		}

		ret, err := h.rpc.Handler(inv, stub)
		if err != nil {
			return nil, fmt.Errorf("[RPC]: invoke index [%d], param %+v, error:%s", i, inv, err.Error())
		}
		rets[i] = &Response{inv, ret}
	}
	return rets, nil
}

func (h *helper) simpleInvoke(stub shim.ChaincodeStubInterface, req *Request) (interface{}, error) {
	return h.rpc.Handler(req, stub)
}

func (h *helper) Register(i interface{}) {
	h.rpc.RegisterMethod(i)
}

func (h *helper) RegisterByRule(i interface{}, r Rule) {
	h.rpc.RegisterMethodByRule(i, r)
}

func DefaultInitFunction(_ shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("SUCCESS"))
}

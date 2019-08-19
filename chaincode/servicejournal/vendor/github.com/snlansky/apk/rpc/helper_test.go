package rpc

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"testing"
	"strconv"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"fmt"
	"errors"
	"reflect"
)

type Msg struct {
	Info    string  `json:"info"`
	Value   float32 `json:"value"`
	Forward bool    `json:"forward"`
}

type MyChaincode struct {
}

func (m *MyChaincode) init(ctx Context) pb.Response {
	return shim.Success([]byte("SUCCESS"))
}

func (m *MyChaincode) Store(stub shim.ChaincodeStubInterface, key string, v int) error {
	return stub.PutState(key, []byte(strconv.Itoa(v)))
}

func (m *MyChaincode) Query(stub shim.ChaincodeStubInterface, key string) []byte {
	buf, err := stub.GetState(key)
	check(err)
	return buf
}

func (m *MyChaincode) Clean(stub shim.ChaincodeStubInterface) {
}

func (m *MyChaincode) Throw(stub shim.ChaincodeStubInterface, t bool) error {
	if t {
		return NewInternalError(errors.New("something wrong"), "内部错误")
	}
	return nil
}

func (m *MyChaincode) SetMsg(stub shim.ChaincodeStubInterface, msg *Msg) *Msg {
	return &Msg{
		Info:msg.Info,
		Value:msg.Value,
		Forward:msg.Forward,
	}
}

func (m *MyChaincode) SetValue(stub shim.ChaincodeStubInterface, msg Msg, info string) *Msg {
	msg.Info = info
	return &msg
}

func TestNewHelper(t *testing.T) {
	mcc := new(MyChaincode)

	helper := NewHelper(mcc.init)

	helper.Register(mcc)

	stub := shim.NewMockStub("MyChaincode", helper)
	stub.MockTransactionStart("init")

	var parems []interface{}
	var resp pb.Response

	resp = stub.MockInit("1", [][]byte{[]byte("init")})
	assert.Equal(t, resp.Payload, []byte("SUCCESS"))

	parems = []interface{}{"age", 23}
	buf, err := json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("Store"), buf})
	fmt.Println(string(resp.Message))
	assert.Equal(t, string(resp.Payload), "")

	parems = []interface{}{"age"}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("Query"), buf})
	fmt.Println(string(resp.Message))
	assert.Equal(t, resp.Payload, getJson(t, []byte(strconv.Itoa(23))))

	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{})
	assert.Equal(t, resp.Status, int32(shim.ERROR))
	assert.Equal(t, resp.Message, ERR_PARAM_INVALID)

	parems = []interface{}{"age"}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("Query"), buf, []byte("other")})
	assert.Equal(t, resp.Status, int32(shim.ERROR))
	assert.Equal(t, resp.Message, ERR_PARAM_INVALID)

	parems = []interface{}{"age", "other"}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("Query"), buf})
	assert.Equal(t, resp.Status, int32(shim.ERROR))
	assert.Equal(t, resp.Message, "[RPC]: params not matched. got 2, need 1.")

	resp = stub.MockInvoke("2", [][]byte{[]byte("Clean")})
	assert.Equal(t, resp.Status, int32(shim.OK))
	assert.Equal(t, string(resp.Payload), "")

	parems = []interface{}{false}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("Throw"), buf})
	assert.Equal(t, resp.Status, int32(shim.OK))
	assert.Equal(t, string(resp.Payload), "")

	parems = []interface{}{true}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("Throw"), buf})
	assert.Equal(t, resp.Status, int32(shim.ERROR))
	assert.Equal(t, resp.Message, "内部错误")

	parems = []interface{}{}
	parems = append(parems, &Msg{Info:"this is msg", Value:34.678, Forward: true})
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte("SetMsg"), buf})
	assert.Equal(t, resp.Status, int32(shim.OK))
	assert.NotEqual(t, resp.Payload, "")
	var m Msg
	err = json.Unmarshal(resp.Payload, &m)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(m, Msg{Info:"this is msg", Value:34.678, Forward: true}))
}

func getJson(t *testing.T, i interface{}) []byte {
	buf, err := json.Marshal(i)
	assert.NoError(t, err)
	return buf
}

func TestHelper_BatchInvoke(t *testing.T) {
	mcc := new(MyChaincode)

	helper := NewHelper(mcc.init)

	helper.Register(mcc)
	helper.EnableBatchInvoke()

	stub := shim.NewMockStub("MyChaincode", helper)
	stub.MockTransactionStart("init")


	var resp pb.Response

	resp = stub.MockInit("1", [][]byte{[]byte("init")})
	assert.Equal(t, resp.Payload, []byte("SUCCESS"))



	var parems []*Request
	parems = append(parems, &Request{"Store", []interface{}{"lili", 23}})
	parems = append(parems, &Request{"Store", []interface{}{"lucy", 23}})
	parems = append(parems, &Request{"Throw", []interface{}{false}})
	parems = append(parems, &Request{"SetMsg", []interface{}{&Msg{Info:"this is msg", Value:34.678, Forward: true}}})
	parems = append(parems, &Request{"SetValue", []interface{}{&Msg{Info:"this is msg", Value:34.678, Forward: true}, "new msg"}})
	buf, err := json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("2", [][]byte{[]byte(SERVICE_BATCH_INVOKE), buf})
	fmt.Println(string(resp.Message))
}
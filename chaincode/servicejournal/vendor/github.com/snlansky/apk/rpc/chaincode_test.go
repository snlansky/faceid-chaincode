package rpc

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"errors"
	"fmt"
)

type Subcc struct {
}

func (m *Subcc) NoParam(stub shim.ChaincodeStubInterface) {
}

func (m *Subcc) Store(stub shim.ChaincodeStubInterface, key string) string {
	return key
}

func (m *Subcc) GetError(stub shim.ChaincodeStubInterface, key string) error {
	return errors.New(key)
}

type Pricc struct {
}

func _init(_ Context) pb.Response {
	return shim.Success([]byte("SUCCESS"))
}

type Subcc1 struct {
}

func (s *Subcc1) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("init Subcc1"))
}

func (s *Subcc1) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	fmt.Println(f, args)
	return shim.Success([]byte("invoke Subcc1"))
}

func TestNewChaincodeClient(t *testing.T) {
	phelper := NewHelper(_init)
	shelper := NewHelper(_init)

	phelper.Register(new(Pricc))
	shelper.Register(new(Subcc))

	stub := shim.NewMockStub("Pricc", phelper)
	stub.MockTransactionStart("init")
	stub.ChannelID = "mychan"

	stub2 := shim.NewMockStub("Subcc", shelper)
	stub2.MockTransactionStart("init")
	stub2.ChannelID = "mychan"

	stub3 := shim.NewMockStub("Subcc1", new(Subcc1))
	stub3.MockTransactionStart("init")
	stub3.ChannelID = "mychan"

	stub.MockPeerChaincode("othercc/mychan", stub2)
	stub.MockPeerChaincode("othercc1/mychan", stub3)

	cli := NewChaincodeClient("othercc")
	cli1 := NewChaincodeClient("othercc1")

	buf, err := cli.Invoke(stub, "NoParam")
	assert.NoError(t, err)
	assert.Equal(t, len(buf), 0)

	buf, err = cli.Invoke(stub, "Store", "test")
	assert.NoError(t, err)
	var s string
	err = json.Unmarshal(buf, &s)
	assert.NoError(t, err)
	assert.Equal(t, s, "test")

	buf, err = cli.Invoke(stub, "GetError", "error info")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error info")

	buf1, err := cli1.Invoke(stub, "functionName", "123", "abc", 45, true, 12.456, nil, map[string]interface{}{"a": "b", "a1": 12}, []interface{}{12, "fsf", false})
	assert.NoError(t, err)
	assert.Equal(t, string(buf1), "invoke Subcc1")
}

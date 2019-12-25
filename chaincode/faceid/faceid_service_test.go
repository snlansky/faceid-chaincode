package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snlansky/blclibs/impl"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewFaceIDService(t *testing.T) {
	cc := impl.NewFabricChaincode()
	cc.Register(NewFaceIDService())

	stub := NewMockStub("FaceIDCC", cc)

	stub.MockTransactionStart("init")

	var parems []interface{}
	var resp pb.Response

	resp = stub.MockInit("tx1", [][]byte{[]byte("init")})
	assert.Equal(t, resp.Payload, []byte("SUCCESS"))

	parems = []interface{}{&FaceID{
		ID:         "",
		SourceType: "proto",
		SourceHash: "has-xxxx",
		Algorithm:  "hash256",
		Labels:     []string{"test"},
		Metadata:   nil,
		Timestamp:  0,
	}}
	buf, err := json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("tx2", [][]byte{[]byte("FaceIDService.RegisterFaceID"), buf})
	assert.Equal(t, resp.Payload, []byte(nil))

	parems = []interface{}{}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("tx3", [][]byte{[]byte("FaceIDService.GetFaceID"), buf})
	assert.NotEmpty(t, resp.Payload)
	fmt.Println(string(resp.Payload))
}

func TestFaceIDService_Record(t *testing.T) {
	cc := impl.NewFabricChaincode()
	cc.Register(NewFaceIDService())

	stub := NewMockStub("FaceIDCC", cc)

	stub.MockTransactionStart("init")

	var parems []interface{}
	var resp pb.Response

	resp = stub.MockInit("tx1", [][]byte{[]byte("init")})
	assert.Equal(t, resp.Payload, []byte("SUCCESS"))

	parems = []interface{}{&FaceID{
		ID:         "",
		SourceType: "proto",
		SourceHash: "has-xxxx",
		Algorithm:  "hash256",
		Labels:     []string{"test"},
		Metadata:   nil,
		Timestamp:  0,
	}}
	buf, err := json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("tx2", [][]byte{[]byte("FaceIDService.RegisterFaceID"), buf})
	assert.Equal(t, resp.Payload, []byte(nil))

	parems = []interface{} {
		&FaceID{
			ID:         "",
			SourceType: "proto",
			SourceHash: "has-yyy",
			Algorithm:  "hash256",
			Labels:     []string{"test"},
			Metadata:   nil,
			Timestamp:  0,
		}}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("tx3", [][]byte{[]byte("FaceIDService.Record"), buf})
	assert.Equal(t, resp.Payload, []byte(nil))

	parems = []interface{} {
		&RequestFaceIDHistory{
			StartTime: 0,
			EndTime:   time.Now().Unix() + 10,
			Labels:    nil,
		}}
	buf, err = json.Marshal(parems)
	assert.NoError(t, err)
	resp = stub.MockInvoke("tx4", [][]byte{[]byte("FaceIDService.HistoryFaceIDs"), buf})
	var faces []*FaceID
	err = json.Unmarshal(resp.Payload, &faces)
	assert.NoError(t, err)
	for _, id := range faces {
		fmt.Println(id)
	}
}

package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/lib"
)

const NodeTableKey = "TicketTable"

type NodeTable struct {
	ck *lib.CompositeKey
}

func NewNodeTable() *NodeTable {
	return &NodeTable{ck: lib.NewCompositeKey(NodeTableKey, "TicketId", "NodeId")}
}

func (t *NodeTable) save(stub shim.ChaincodeStubInterface, node *NodeCommon) error {
	buf, err := json.Marshal(node)
	if err != nil {
		return err
	}
	return t.ck.Insert(stub, []string{node.TicketID, node.Id}, buf)
}

func (t *NodeTable) find(stub shim.ChaincodeStubInterface, ticketId, nodeId string) (*NodeCommon, error) {
	bytes, err := t.ck.GetValue(stub, []string{ticketId, nodeId})
	if err != nil {
		return nil, err
	}

	if bytes == nil || len(bytes) == 0 {
		return nil, nil
	}

	var node NodeCommon
	err = json.Unmarshal(bytes, &node)
	if err != nil {
		return nil, err
	}

	return &node, err
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/base"
	"github.com/snlansky/apk/rpc"
)

var DefaultNodeService = &NodeService{nodeTable: NewNodeTable(), ticketTable: NewTicketTable()}

type NodeService struct {
	nodeTable   *NodeTable
	ticketTable *TicketTable
}

func (s *NodeService) Create(stub shim.ChaincodeStubInterface, nodeJson string) *NodeCommon {
	var node NodeCommon
	err := json.Unmarshal([]byte(nodeJson), &node)
	rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	mustValidate(&node)

	ticket, err := s.ticketTable.find(stub, node.TicketID)
	rpc.Check(err, InternalError)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}

	old, err := s.nodeTable.find(stub, ticket.Id, node.Id)
	rpc.Check(err, InternalError)
	if old != nil {
		rpc.Throw(ReduplicateCreate)
	}

	err = s.nodeTable.save(stub, &node)
	rpc.Check(err, InternalError)

	err = base.CreateEvent(stub, AppName, CreateNodeEvent, fmt.Sprintf("%s:%s", node.TicketID, node.Id))
	rpc.Check(err, InternalError)
	return &node
}

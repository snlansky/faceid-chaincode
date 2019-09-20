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

func (s *NodeService) Create(stub shim.ChaincodeStubInterface, nodeJson string) *NodeResponse {
	var req NodeRequest
	err := json.Unmarshal([]byte(nodeJson), &req)
	rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	mustValidate(&req)

	ticket, err := s.ticketTable.find(stub, req.TicketID)
	rpc.Check(err, InternalError)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}

	old, err := s.nodeTable.find(stub, ticket.Id, req.Id)
	rpc.Check(err, InternalError)
	if old != nil {
		rpc.Throw(ReduplicateCreate)
	}

	addr := Address(base.MustGetAddress(stub))
	timestamp, err := stub.GetTxTimestamp()
	rpc.Check(err, InternalError)

	node := &Node{
		NodeCommon: NodeCommon{
			Id:         stub.GetTxID(),
			TicketID:   req.TicketID,
			HandlerId:  req.HandlerId,
			Status:     req.Status,
			CreateTime: req.CreateTime,
			UploadTime: timestamp.Seconds,
			Source:     req.Source,
		},
		Metadata: map[Address]string{},
	}

	node.Metadata[addr] = req.Extension

	err = s.nodeTable.save(stub, node)
	rpc.Check(err, InternalError)

	ticket.NodeList = append(ticket.NodeList, node.Id)
	err = s.ticketTable.update(stub, ticket)
	rpc.Check(err, rpc.ERR_PARAM_INVALID)

	err = base.CreateEvent(stub, AppName, CreateNodeEvent, fmt.Sprintf("%s:%s", req.TicketID, req.Id))
	rpc.Check(err, InternalError)
	return s.getResponse(addr, node)
}

func (s *NodeService) Update(stub shim.ChaincodeStubInterface, nodeJson string) *NodeResponse {
	var req NodeRequest
	err := json.Unmarshal([]byte(nodeJson), &req)
	rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	mustValidate(&req)

	ticket, err := s.ticketTable.find(stub, req.TicketID)
	rpc.Check(err, InternalError)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}

	node, err := s.nodeTable.find(stub, ticket.Id, req.Id)
	rpc.Check(err, InternalError)
	if node == nil {
		rpc.Throw(NotExistResource)
	}

	addr := Address(base.MustGetAddress(stub))
	node.Metadata[addr] = req.Extension

	err = s.ticketTable.update(stub, ticket)
	rpc.Check(err, InternalError)

	err = base.CreateEvent(stub, AppName, UpdateNodeEvent, fmt.Sprintf("%s:%s", node.TicketID, node.Id))
	rpc.Check(err, InternalError)
	return s.getResponse(addr, node)
}

func (s *NodeService) getResponse(addr Address, node *Node) *NodeResponse {
	return &NodeResponse{
		NodeCommon: node.NodeCommon,
		Extension:  node.Metadata[addr],
	}
}

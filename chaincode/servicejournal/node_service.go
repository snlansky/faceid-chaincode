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
	rpc.Check(err, FindResourceFailed)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}

	old, err := s.nodeTable.find(stub, ticket.Id, req.Id)
	rpc.Check(err, FindResourceFailed)
	if old != nil {
		rpc.Throw(ReduplicateCreate)
	}

	details := map[string]string{}
	if req.Details != nil {
		for k, v := range req.Details {
			details[k] = v
		}
	}

	addr := Address(base.MustGetAddress(stub))

	node := &Node{
		NodeCommon: NodeCommon{
			Id:          stub.GetTxID(),
			TicketID:    req.TicketID,
			HandlerId:   req.HandlerId,
			Status:      req.Status,
			CreateTime:  req.CreateTime,
			UpdateTime:  req.CreateTime,
			UploadTime:  base.MustGetTimestamp(stub),
			Description: req.Description,
			System:      req.System,
			Details:     details,
		},
		SourceList: []Address{},
		Metadata:   map[Address]string{},
	}

	node.SourceList = append(node.SourceList, addr)
	node.Metadata[addr] = req.Extension

	err = s.nodeTable.save(stub, node)
	rpc.Check(err, SaveResourcFailed)

	ticket.NodeList = append(ticket.NodeList, node.Id)
	err = s.ticketTable.update(stub, ticket)
	rpc.Check(err, "ERR_UPDATE_FAILED")

	err = base.CreateEvent(stub, AppName, CreateNodeEvent, fmt.Sprintf("%s:%s", node.TicketID, node.Id))
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
	rpc.Check(err, FindResourceFailed)
	if node == nil {
		rpc.Throw(NotExistResource)
	}

	node.UpdateTime = base.MustGetTimestamp(stub)

	addr := Address(base.MustGetAddress(stub))
	if node.SourceList[len(node.SourceList)-1] != addr {
		node.SourceList = append(node.SourceList, addr)
	}
	node.Metadata[addr] = req.Extension

	if req.Details != nil {
		for k, v := range req.Details {
			node.Details[k] = v
		}
	}

	err = s.ticketTable.update(stub, ticket)
	rpc.Check(err, "ERR_UPDATE_FAILED")

	err = base.CreateEvent(stub, AppName, UpdateNodeEvent, fmt.Sprintf("%s:%s", node.TicketID, node.Id))
	rpc.Check(err, InternalError)
	return s.getResponse(addr, node)
}

func (s *NodeService) GetById(stub shim.ChaincodeStubInterface, ticketId, nodeId string) *NodeResponse {
	node, err := s.nodeTable.find(stub, ticketId, nodeId)
	rpc.Check(err, FindResourceFailed)
	if node == nil {
		rpc.Throw(NotExistResource)
	}
	addr := Address(base.MustGetAddress(stub))
	return s.getResponse(addr, node)
}

func (s *NodeService) GetNodeById(stub shim.ChaincodeStubInterface, ticketId, nodeId string) *Node {
	node, err := s.nodeTable.find(stub, ticketId, nodeId)
	rpc.Check(err, FindResourceFailed)
	if node == nil {
		rpc.Throw(NotExistResource)
	}
	return node
}

func (s *NodeService) getResponse(addr Address, node *Node) *NodeResponse {
	return &NodeResponse{
		NodeCommon: node.NodeCommon,
		Extension:  node.Metadata[addr],
	}
}

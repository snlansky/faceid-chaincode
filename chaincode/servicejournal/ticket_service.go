package main

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/base"
	"github.com/snlansky/apk/rpc"
)

var DefaultTicketService = &TicketService{ticketTable: NewTicketTable()}

type TicketService struct {
	ticketTable *TicketTable
}

func (s *TicketService) Create(stub shim.ChaincodeStubInterface, ticketRequestJson string) *TicketResponse {
	var req TicketRequest
	err := json.Unmarshal([]byte(ticketRequestJson), &req)
	rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	mustValidate(&req)

	old, err := s.ticketTable.find(stub, req.Id)
	rpc.Check(err, InternalError)
	if old != nil {
		rpc.Throw(ReduplicateCreate)
	}

	timestamp, err := stub.GetTxTimestamp()
	rpc.Check(err, InternalError)

	address, err := base.GetAddress(stub)
	rpc.Check(err, InternalError)

	id := stub.GetTxID()
	ticket := &Ticket{
		TicketCommon: TicketCommon{
			Id:          id,
			Title:       req.Title,
			Description: req.Description,
			Products:    req.Products,
			System:      req.System,
			Status:      req.Status,
			OwnerId:     req.OwnerId,
			SubmitterId: req.SubmitterId,
			HandlerId:   req.HandlerId,
			CreateTime:  req.CreateTime,
			UpdateTime:  req.UpdateTime,
			UploadTime:  timestamp.Seconds,
		},
		SourceList: []Address{},
		NodeList:   []string{},
		Metadata:   map[Address]string{},
	}

	ticket.SourceList = append(ticket.SourceList, Address(address))
	ticket.Metadata[Address(address)] = req.Extension

	err = s.ticketTable.save(stub, ticket)
	rpc.Check(err, InternalError)

	rpc.Check(base.CreateEvent(stub, AppName, CreateTicketEvent, ticket.Id), InternalError)

	return s.getTicketResponse(Address(address), ticket)
}

func (s *TicketService) getTicketResponse(addr Address, ticket *Ticket) *TicketResponse {
	return &TicketResponse{
		TicketCommon: ticket.TicketCommon,
		SourceList:   ticket.SourceList,
		NodeList:     ticket.NodeList,
		Extension:    ticket.Metadata[addr],
	}
}

func (s *TicketService) Update(stub shim.ChaincodeStubInterface, ticketRequestJson string) *TicketResponse {
	var req TicketRequest
	err := json.Unmarshal([]byte(ticketRequestJson), &req)
	rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	mustValidate(&req)

	ticket, err := s.ticketTable.find(stub, req.Id)
	rpc.Check(err, InternalError)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}

	ticket.Status = req.Status
	ticket.UpdateTime = req.UpdateTime

	addr := Address(base.MustGetAddress(stub))
	findSource := false
	for _, s := range ticket.SourceList {
		if s == req.Source {
			findSource = true
		}
	}
	if !findSource {
		ticket.SourceList = append(ticket.SourceList, req.Source)
	}

	ticket.Metadata[addr] = req.Extension

	err = s.ticketTable.update(stub, ticket)
	rpc.Check(err, rpc.ERR_PARAM_INVALID)

	rpc.Check(base.CreateEvent(stub, AppName, UpdateTicketEvent, ticket.Id), InternalError)
	return s.getTicketResponse(addr, ticket)
}

func (s *TicketService) FindByID(stub shim.ChaincodeStubInterface, id string) *TicketResponse {
	ticket, err := s.ticketTable.find(stub, id)
	rpc.Check(err, InternalError)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}
	addr := Address(base.MustGetAddress(stub))
	return s.getTicketResponse(addr, ticket)
}

func mustValidate(obj interface{}) {
	validate, err := govalidator.ValidateStruct(obj)
	rpc.Check(err, rpc.ERR_PARAM_INVALID)
	if !validate {
		rpc.Throw(rpc.ERR_PARAM_INVALID)
	}
}

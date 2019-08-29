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

func (s *TicketService) Create(stub shim.ChaincodeStubInterface, ticket *TicketCommon) *TicketCommon {
	//var ticket TicketCommon
	//err := json.Unmarshal([]byte(ticketJson), &ticket)
	//rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	mustValidate(&ticket)

	creator, err := base.GetAddress(stub)
	rpc.Check(err, InternalError)

	if creator != ticket.SubmitterId {
		rpc.Throw(rpc.ERR_PARAM_INVALID)
	}

	old, err := s.ticketTable.find(stub, ticket.Id)
	rpc.Check(err, InternalError)
	if old != nil {
		rpc.Throw(ReduplicateCreate)
	}

	err = s.ticketTable.save(stub, ticket)
	rpc.Check(err, InternalError)

	rpc.Check(base.CreateEvent(stub, AppName, CreateTicketEvent, ticket.Id), InternalError)

	return ticket
}

func (s *TicketService) Update(stub shim.ChaincodeStubInterface, ticketJson string) *TicketCommon {
	var ticket TicketCommon
	err := json.Unmarshal([]byte(ticketJson), &ticket)
	rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	mustValidate(&ticket)

	old, err := s.ticketTable.find(stub, ticket.Id)
	rpc.Check(err, InternalError)
	if old == nil {
		rpc.Throw(NotExistResource)
	}

	err = s.ticketTable.update(stub, &ticket)
	rpc.Check(err, rpc.ERR_PARAM_INVALID)

	rpc.Check(base.CreateEvent(stub, AppName, UpdateTicketEvent, ticket.Id), InternalError)
	return &ticket
}

func (s *TicketService) FindByID(stub shim.ChaincodeStubInterface, id string) *TicketCommon {
	ticket, err := s.ticketTable.find(stub, id)
	rpc.Check(err, InternalError)
	return ticket
}

func mustValidate(obj interface{}) {
	validate, err := govalidator.ValidateStruct(obj)
	rpc.Check(err, rpc.ERR_PARAM_INVALID)
	if !validate {
		rpc.Throw(rpc.ERR_PARAM_INVALID)
	}
}

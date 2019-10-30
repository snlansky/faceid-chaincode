package main

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/base"
	"github.com/snlansky/apk/rpc"
)

var DefaultTicketService = &TicketService{
	ticketTable:  NewTicketTable(),
	mappingTable: NewMappingTable(),
}

type TicketService struct {
	ticketTable  *TicketTable
	mappingTable *MappingTable
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

	address, err := base.GetAddress(stub)
	rpc.Check(err, InternalError)

	details := map[string]string{}
	if req.Details != nil {
		for k, v := range req.Details {
			details[k] = v
		}
	}

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
			UploadTime:  base.MustGetTimestamp(stub),
			Details:     details,
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

// 构建各个系统的Id映射
func (s *TicketService) BuildIdsMapping(stub shim.ChaincodeStubInterface, bindJson string) *IdMapping {
	var req IdMapping
	err := json.Unmarshal([]byte(bindJson), &req)
	rpc.Check(err, rpc.ERR_JSON_UNMARSHAL)

	if req.TicketId == "" || req.InternalIds == nil || len(req.InternalIds) == 0 {
		rpc.Throw("internal ids or ticket id is null")
	}

	ticket, err := s.ticketTable.find(stub, req.TicketId)
	rpc.Check(err, InternalError)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}

	addr, err := base.GetAddress(stub)
	rpc.Check(err, InternalError)
	err = s.mappingTable.save(stub, addr, req.TicketId, req.InternalIds)
	rpc.Check(err, InternalError)
	return &req
}

// 通过个别系统的Id查询各个系统的映射表
func (s *TicketService) FindIdMapping(stub shim.ChaincodeStubInterface, id string) *IdMapping {
	if id == "" {
		rpc.Throw("id is null")
	}

	addr, err := base.GetAddress(stub)
	rpc.Check(err, InternalError)

	ticketId, ids, err := s.mappingTable.find(stub, addr, id)
	rpc.Check(err, InternalError)


	// 这里不判读空
	return &IdMapping{
		TicketId:    ticketId,
		InternalIds: ids,
	}
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

	ticket.Title = req.Title
	ticket.Status = req.Status
	ticket.UpdateTime = req.UpdateTime

	if req.Details != nil {
		for k, v := range req.Details {
			ticket.Details[k] = v
		}
	}

	addr := Address(base.MustGetAddress(stub))

	// ticket.SourceList length must > 0
	// 相同就不追加了
	if ticket.SourceList[len(ticket.SourceList)-1] != addr {
		ticket.SourceList = append(ticket.SourceList, addr)
	}
	ticket.HandlerId = req.HandlerId
	ticket.Metadata[addr] = req.Extension

	err = s.ticketTable.update(stub, ticket)
	rpc.Check(err, "ERR_UPDATE_FAILED")

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

func (s *TicketService) FindTicketByID(stub shim.ChaincodeStubInterface, id string) *Ticket {
	ticket, err := s.ticketTable.find(stub, id)
	rpc.Check(err, InternalError)
	if ticket == nil {
		rpc.Throw(NotExistResource)
	}
	return ticket
}

func mustValidate(obj interface{}) {
	validate, err := govalidator.ValidateStruct(obj)
	rpc.Check(err, rpc.ERR_PARAM_INVALID)
	if !validate {
		rpc.Throw("ERR_STRUCT_VALIDATE" )
	}
}

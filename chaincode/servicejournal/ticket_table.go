package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/lib"
)

const TicketTableKey = "TicketTable"

type TicketTable struct {
	ck *lib.CompositeKey
}

func NewTicketTable() *TicketTable {
	return &TicketTable{ck: lib.NewCompositeKey(TicketTableKey, "Id")}
}

func (t *TicketTable) save(stub shim.ChaincodeStubInterface, ticket *TicketCommon) error {
	bytes, err := json.Marshal(ticket)
	if err != nil {
		return err
	}
	return t.ck.Insert(stub, []string{ticket.Id}, bytes)
}

func (t *TicketTable) find(stub shim.ChaincodeStubInterface, id string) (*TicketCommon, error) {
	bytes, err := t.ck.GetValue(stub, []string{id})
	if err != nil {
		return nil, err
	}
	if bytes == nil || len(bytes) == 0 {
		return nil, nil
	}
	var ticket TicketCommon
	err = json.Unmarshal(bytes, &ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (t *TicketTable) update(stub shim.ChaincodeStubInterface, ticket *TicketCommon) error {
	bytes, err := json.Marshal(ticket)
	if err != nil {
		return err
	}
	return t.ck.Update(stub, []string{ticket.Id}, bytes)
}

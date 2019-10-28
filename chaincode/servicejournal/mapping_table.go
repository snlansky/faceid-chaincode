package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/lib"
)

const MappingTableKey = "MappingTable"

const (
	TicketIdType       = "TicketIdId"     // 值类型为ticketId
	InternalIdListType = "InternalIdList" // 值类型为InternalIdList
)

type Value struct {
	Type    string
	Content []byte
}

func (v *Value) serialize() ([]byte, error) {
	return json.Marshal(v)
}

func unserialize(data []byte) (Value, error) {
	var v Value
	err := json.Unmarshal(data, &v)
	return v, err
}

type MappingTable struct {
	ck *lib.CompositeKey
}

// MappingTableKey:Address:ID -> Value
func NewMappingTable() *MappingTable {
	return &MappingTable{ck: lib.NewCompositeKey(MappingTableKey, "Address", "Id")}
}

func (t *MappingTable) save(stub shim.ChaincodeStubInterface, address string, ticketId string, internalIds []string) error {
	// mapping ticketId -> internalIdList
	err := t.saveTicketIdMapping(stub, address, ticketId, internalIds)
	if err != nil {
		return err
	}

	for i := range internalIds {
		err = t.saveInternalIdMapping(stub, address, ticketId, internalIds[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *MappingTable) saveTicketIdMapping(stub shim.ChaincodeStubInterface, address string, ticketId string, internalIds []string) error {
	idsBytes, err := json.Marshal(internalIds)
	if err != nil {
		return err
	}

	idsValue := Value{
		Type:    InternalIdListType,
		Content: idsBytes,
	}
	idsVBytes, err := idsValue.serialize()
	if err != nil {
		return err
	}

	return t.ck.Insert(stub, []string{address, ticketId}, idsVBytes)
}

func (t *MappingTable) saveInternalIdMapping(stub shim.ChaincodeStubInterface, address string, ticketId string, internalId string) error {
	tidValue := Value{
		Type:    TicketIdType,
		Content: []byte(ticketId),
	}
	tidVBytes, err := tidValue.serialize()
	if err != nil {
		return err
	}

	return t.ck.Insert(stub, []string{address, internalId}, tidVBytes)
}

func (t *MappingTable) find(stub shim.ChaincodeStubInterface, address string, id string) (string, []string, error) {
	bytes, err := t.ck.GetValue(stub, []string{address, id})
	if err != nil {
		return "", nil, err
	}

	if bytes == nil || len(bytes) == 0 {
		return "", nil, nil
	}

	value, err := unserialize(bytes)
	if err != nil {
		return "", nil, nil
	}

	switch value.Type {
	case TicketIdType:
		ticketId := string(value.Content)
		internalIdList, err := t.findInternalListByTicketId(stub, address, ticketId)
		if err != nil {
			return "", nil, err
		}
		return ticketId, internalIdList, err

	case InternalIdListType:
		var list []string
		err := json.Unmarshal(value.Content, &list)
		if err != nil {
			return "", nil, err
		}
		return id, list, nil
	default:
		return "", nil, fmt.Errorf("unreachable")
	}
}

func (t *MappingTable) findInternalListByTicketId(stub shim.ChaincodeStubInterface, address string, ticketId string) ([]string, error) {
	bytes, err := t.ck.GetValue(stub, []string{address, ticketId})
	if err != nil {
		return nil, err
	}

	value, err := unserialize(bytes)
	if err != nil {
		return nil, err
	}

	var list []string
	err = json.Unmarshal(value.Content, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

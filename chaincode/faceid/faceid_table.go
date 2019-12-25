package main

import (
	"encoding/json"
	"github.com/snlansky/blclibs"
)

const FaceIDTableKey = "FaceIDTable"

type FaceIDTable struct {
	table *blclibs.Table
}

func NewFaceIDTable() *FaceIDTable {
	return &FaceIDTable{table: blclibs.NewTable(FaceIDTableKey, "address", "faceid")}
}

func (t *FaceIDTable) Save(stub blclibs.IContractStub, addr blclibs.Address, id *FaceID) error {
	bytes, err := json.Marshal(id)
	if err != nil {
		return err
	}
	return t.table.Insert(stub, []string{string(addr), id.ID}, bytes)
}

func (t *FaceIDTable) Get(stub blclibs.IContractStub, addr blclibs.Address, id string) (*FaceID, error) {
	value, err := t.table.GetValue(stub, []string{string(addr), id})
	if err != nil {
		return nil, err
	}
	var faceId FaceID
	err = json.Unmarshal(value, &faceId)
	if err != nil {
		return nil, err
	}
	return &faceId, nil
}

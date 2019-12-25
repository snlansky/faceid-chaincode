package main

import (
	"encoding/json"
	"github.com/snlansky/blclibs"
)

const FaceIDUserTableKey = "FaceIDUserTable"

type FaceIDUserTable struct {
	table *blclibs.Table
}

func NewFaceIDUserTable() *FaceIDUserTable {
	return &FaceIDUserTable{table: blclibs.NewTable(FaceIDUserTableKey, "address")}
}

func (userTable *FaceIDUserTable) Save(stub blclibs.IContractStub, address blclibs.Address, user *User) error {
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return userTable.table.Insert(stub, []string{string(address)}, bytes)
}

func (userTable *FaceIDUserTable) Get(stub blclibs.IContractStub, address blclibs.Address) (*User, error) {
	value, err := userTable.table.GetValue(stub, []string{string(address)})
	if err != nil {
		return nil, err
	}

	if value == nil || len(value) == 0 {
		return nil, nil
	}

	var u User
	err = json.Unmarshal(value, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

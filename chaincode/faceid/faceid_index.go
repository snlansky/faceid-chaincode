package main

import (
	"encoding/json"
	"github.com/snlansky/blclibs"
)

type FaceIDIndex struct {
	index *blclibs.Index
}

func NewFaceIDIndex() *FaceIDIndex {
	return &FaceIDIndex{index: blclibs.NewIndex(AppName, "FaceIDTable", "ID")}
}

func (index *FaceIDIndex) Save(stub blclibs.IContractStub, addr blclibs.Address, id *FaceID) error {
	idx := &TimeIndex{
		FaceID:    id.ID,
		Timestamp: id.Timestamp,
	}
	bytes, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	_, err = index.index.Save(stub, addr, bytes)
	return err
}

func (index *FaceIDIndex) Get(stub blclibs.IContractStub, addr blclibs.Address, history *RequestFaceIDHistory) ([]string, error) {
	var list []string
	err := index.index.Filter(stub, addr, func(value []byte) (bool, error) {
		var idx TimeIndex
		err := json.Unmarshal(value, &idx)
		if err != nil {
			return false, err
		}
		if idx.Timestamp >= history.StartTime && idx.Timestamp <= history.EndTime {
			list = append(list, idx.FaceID)
		}
		return true, nil
	})
	return list, err
}

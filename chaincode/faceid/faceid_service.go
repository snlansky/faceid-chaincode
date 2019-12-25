package main

import (
	"github.com/snlansky/blclibs"
	"github.com/snlansky/glibs/rpc"
)

type FaceIDService struct {
	faceIDUserTable *FaceIDUserTable
	faceIDTable     *FaceIDTable
	faceIDIndex     *FaceIDIndex
}

func NewFaceIDService() *FaceIDService {
	return &FaceIDService{
		faceIDUserTable: NewFaceIDUserTable(),
		faceIDTable:     NewFaceIDTable(),
		faceIDIndex:     NewFaceIDIndex(),
	}
}

func (svc *FaceIDService) RegisterFaceID(stub blclibs.IContractStub, id *FaceID) {
	address, err := stub.GetAddress()
	rpc.Check(err)
	addr := blclibs.Address(address)

	mustValidate(id)
	id.ID = stub.GetTxID()
	ts, err := stub.GetTxTimestamp()
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
	id.Timestamp = ts.Unix()

	err = svc.faceIDTable.Save(stub, addr, id)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	err = svc.faceIDUserTable.Save(stub, addr, &User{RegisterFaceID: id.ID})
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	err = svc.faceIDIndex.Save(stub, addr, id)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
}

func (svc *FaceIDService) GetFaceID(stub blclibs.IContractStub) *FaceID {
	address, err := stub.GetAddress()
	rpc.Check(err)
	addr := blclibs.Address(address)
	user, err := svc.faceIDUserTable.Get(stub, addr)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	if user == nil {
		return &FaceID{}
	}

	faceID, err := svc.faceIDTable.Get(stub, addr, user.RegisterFaceID)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
	return faceID
}

func (svc *FaceIDService) Record(stub blclibs.IContractStub, id *FaceID) {
	address, err := stub.GetAddress()
	rpc.Check(err)
	addr := blclibs.Address(address)

	mustValidate(id)
	id.ID = stub.GetTxID()
	ts, err := stub.GetTxTimestamp()
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
	id.Timestamp = ts.Unix()

	err = svc.faceIDTable.Save(stub, addr, id)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	err = svc.faceIDIndex.Save(stub, addr, id)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
}

func (svc *FaceIDService) HistoryFaceIDs(stub blclibs.IContractStub, req *RequestFaceIDHistory) []*FaceID {
	address, err := stub.GetAddress()
	rpc.Check(err)
	addr := blclibs.Address(address)

	mustValidate(req)

	ids, err := svc.faceIDIndex.Get(stub, addr, req)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	var faces []*FaceID
	for _, id := range ids {
		faceID, err := svc.faceIDTable.Get(stub, addr, id)
		rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
		faces = append(faces, faceID)
	}
	return faces
}

func mustValidate(obj Validator) {
	if !obj.Validate() {
		rpc.Throw("ERR_STRUCT_VALIDATE")
	}
}

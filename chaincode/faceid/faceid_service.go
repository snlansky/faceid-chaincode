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

	user := svc.makeUser(stub, addr)
	user.RegisterFaceID = id

	err = svc.faceIDTable.Save(stub, addr, id)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	err = svc.faceIDIndex.Save(stub, addr, id)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	err = svc.faceIDUserTable.Save(stub, addr, user)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
}

func (svc *FaceIDService) RegisterCertificate(stub blclibs.IContractStub, id *FaceID) {
	address, err := stub.GetAddress()
	rpc.Check(err)
	addr := blclibs.Address(address)

	mustValidate(id)
	id.ID = stub.GetTxID()
	ts, err := stub.GetTxTimestamp()
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
	id.Timestamp = ts.Unix()

	user := svc.makeUser(stub, addr)
	user.RegisterCertificate = id

	err = svc.faceIDUserTable.Save(stub, addr, user)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
}

func (svc *FaceIDService) GetUser(stub blclibs.IContractStub) *User {
	address, err := stub.GetAddress()
	rpc.Check(err)
	addr := blclibs.Address(address)
	return svc.makeUser(stub, addr)
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

	// registration check
	check, err := svc.registrationCheck(stub, addr)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)
	if !check {
		rpc.Throw("not registration")
	}

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

func (svc *FaceIDService) makeUser(stub blclibs.IContractStub, addr blclibs.Address) *User {
	user, err := svc.faceIDUserTable.Get(stub, addr)
	rpc.Check(err, rpc.ERR_INTERNAL_INVALID)

	if user == nil {
		user = &User{
			RegisterFaceID:      nil,
			RegisterCertificate: nil,
		}
	}
	return user
}

func (svc *FaceIDService) registrationCheck(stub blclibs.IContractStub, addr blclibs.Address) (bool, error) {
	user, err := svc.faceIDUserTable.Get(stub, addr)
	return user != nil, err
}

func mustValidate(obj Validator) {
	if !obj.Validate() {
		rpc.Throw("ERR_STRUCT_VALIDATE")
	}
}

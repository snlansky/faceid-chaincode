package impl

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/snlansky/blclibs"
	"time"
)

type FabricContractStub struct {
	stub    shim.ChaincodeStubInterface
	creator func() []byte
}

func NewFabricContractStub(stub shim.ChaincodeStubInterface) blclibs.IContractStub {
	return &FabricContractStub{stub: stub}
}

func (f *FabricContractStub) setCreatorFactory(creator func() []byte) {
	f.creator = creator
}

func (f *FabricContractStub) GetArgs() [][]byte {
	return f.stub.GetArgs()
}

func (f *FabricContractStub) GetTxID() string {
	return f.stub.GetTxID()
}

func (f *FabricContractStub) GetChannelID() string {
	return f.stub.GetChannelID()
}

func (f *FabricContractStub) GetAddress() ([]byte, error) {
	creatorByte, err := f.stub.GetCreator()
	if err != nil {
		return nil, err
	}

	certStart := bytes.Index(creatorByte, []byte("-----BEGIN"))
	if certStart == -1 {
		return nil, errors.New("No creator certificate found")
	}
	certText := creatorByte[certStart:]

	bl, _ := pem.Decode(certText)
	if bl == nil {
		return nil, errors.New("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return nil, errors.New("Parse Certificate failed")
	}

	if pub, ok := cert.PublicKey.(*ecdsa.PublicKey); ok {
		pubKey := append(pub.X.Bytes(), pub.Y.Bytes()...)
		publicSHA256 := sha256.Sum256(pubKey)
		address := blclibs.Base58Encode(publicSHA256[:])

		return address, nil
	}

	return nil, errors.New("Only support ECDSA")
}

func (f *FabricContractStub) GetState(key string) ([]byte, error) {
	return f.stub.GetState(key)
}

func (f *FabricContractStub) PutState(key string, value []byte) error {
	return f.stub.PutState(key, value)
}

func (f *FabricContractStub) DelState(key string) ([]byte, error) {
	buf, err := f.stub.GetState(key)
	if err != nil {
		return nil, err
	}
	err = f.stub.DelState(key)
	return buf, err
}

func (f *FabricContractStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return f.stub.CreateCompositeKey(objectType, attributes)
}

func (f *FabricContractStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return f.stub.SplitCompositeKey(compositeKey)
}

func (f *FabricContractStub) GetTxTimestamp() (time.Time, error) {
	ts, err := f.stub.GetTxTimestamp()
	if err != nil {
		return time.Time{}, err
	}
	if ts == nil {
		return time.Time{}, errors.New("timestamp: nil Timestamp")
	}
	if ts.Seconds < -62135596800 {
		return time.Time{}, fmt.Errorf("timestamp: %v before 0001-01-01", ts)
	}
	if ts.Seconds >= 253402300800 {
		return time.Time{}, fmt.Errorf("timestamp: %v after 10000-01-01", ts)
	}
	if ts.Nanos < 0 || ts.Nanos >= 1e9 {
		return time.Time{}, fmt.Errorf("timestamp: %v: nanos not in range [0, 1e9)", ts)
	}

	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC(), nil
}

func (f *FabricContractStub) SetEvent(name string, payload []byte) error {
	return f.stub.SetEvent(name, payload)
}

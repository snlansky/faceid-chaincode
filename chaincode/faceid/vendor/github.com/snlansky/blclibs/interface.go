package blclibs

import (
	"time"
)

type IContractStub interface {
	GetArgs() [][]byte
	GetTxID() string
	GetChannelID() string
	GetAddress() ([]byte, error)
	GetState(key string) ([]byte, error)
	PutState(key string, value []byte) error
	DelState(key string) ([]byte, error)
	CreateCompositeKey(objectType string, attributes []string) (string, error)
	SplitCompositeKey(compositeKey string) (string, []string, error)
	GetTxTimestamp() (time.Time, error)
	SetEvent(name string, payload []byte) error
}

type Address string

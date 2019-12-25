package impl

import (
	"github.com/snlansky/blclibs"
	"time"
)

type MemImpl struct {
	states map[string][]byte
	events map[string][]byte
}

func (m *MemImpl) GetArgs() [][]byte {
	panic("implement me")
}

func (m *MemImpl) GetTxID() string {
	panic("implement me")
}

func (m *MemImpl) GetChannelID() string {
	return "mem-channel"
}

func (m *MemImpl) GetAddress() ([]byte, error) {
	return []byte("my-address-1"), nil
}

func (m *MemImpl) GetState(key string) ([]byte, error) {
	v := m.states[key]
	return v, nil
}

func (m *MemImpl) PutState(key string, value []byte) error {
	m.states[key] = value
	return nil
}

func (m *MemImpl) DelState(key string) ([]byte, error) {
	v, ok := m.states[key]
	if ok {
		delete(m.states, key)
		return v, nil
	}
	return nil, nil
}

func (m *MemImpl) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return blclibs.CreateCompositeKey(objectType, attributes)
}

func (m *MemImpl) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return blclibs.SplitCompositeKey(compositeKey)
}

func (m *MemImpl) GetTxTimestamp() (time.Time, error) {
	return time.Now(), nil
}

func (m *MemImpl) SetEvent(name string, payload []byte) error {
	m.events[name] = payload
	return nil
}

func NewMemImpl() blclibs.IContractStub {
	return &MemImpl{states: map[string][]byte{}, events: map[string][]byte{}}
}

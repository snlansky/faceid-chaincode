package rpc

import "github.com/hyperledger/fabric/core/chaincode/shim"

type StubImpl struct {
	shim.ChaincodeStubInterface
	cache map[string][]byte
}

func newStubImpl(stub shim.ChaincodeStubInterface) *StubImpl {
	return &StubImpl{ChaincodeStubInterface: stub, cache: make(map[string][]byte)}
}

func (c *StubImpl) GetState(key string) ([]byte, error) {
	if value, found := c.cache[key]; found {
		return value, nil
	}
	value, err := c.ChaincodeStubInterface.GetState(key)
	if err != nil {
		return nil, err
	}
	c.cache[key] = value
	return value, nil
}

func (c *StubImpl) PutState(key string, value []byte) error {
	err := c.ChaincodeStubInterface.PutState(key, value)
	if err != nil {
		return err
	}
	c.cache[key] = value
	return nil
}

func (c *StubImpl) DelState(key string) error {
	err := c.ChaincodeStubInterface.DelState(key)
	if err != nil {
		return err
	}
	delete(c.cache, key)
	return nil
}

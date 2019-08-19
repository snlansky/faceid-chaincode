package lib

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
	"fmt"
)

const sep = ":"

type KV struct {
	Keys  []string
	Value []byte
}

type CompositeKey struct {
	object string
	fields []string
}

func NewCompositeKey(obj string, fields ...string) *CompositeKey {
	return &CompositeKey{object: obj, fields: fields}
}

func (s *CompositeKey) GetType() string {
	return strings.Join(append([]string{s.object}, s.fields...), sep)
}

func (s *CompositeKey) createCompositeKey(stub shim.ChaincodeStubInterface, keys []string) (string, error) {
	if err := s.check(keys); err != nil {
		return "", err
	}
	return stub.CreateCompositeKey(s.GetType(), keys)
}

func (s *CompositeKey) Insert(stub shim.ChaincodeStubInterface, keys []string, value []byte) error {
	key, err := s.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	if value == nil {
		value = []byte{0x00}
	}
	return stub.PutState(key, value)
}

func (s *CompositeKey) Delete(stub shim.ChaincodeStubInterface, keys []string) error {
	key, err := s.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	return stub.DelState(key)
}

func (s *CompositeKey) Update(stub shim.ChaincodeStubInterface, keys []string, value []byte) error {
	key, err := s.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	return stub.PutState(key, value)
}

func (s *CompositeKey) UpdateFunc(stub shim.ChaincodeStubInterface, keys []string, f func(oldValue []byte) ([]byte, error)) error {
	key, err := s.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	buf, err := stub.GetState(key)
	if err != nil {
		return err
	}
	value, err := f(buf)
	if err != nil {
		return err
	}
	return stub.PutState(key, value)
}

func (s *CompositeKey) Select(stub shim.ChaincodeStubInterface, keys ...string) (list []*KV, err error) {
	iter, err := stub.GetStateByPartialCompositeKey(s.GetType(), keys)
	if err != nil {
		return
	}
	defer iter.Close()

	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}
		ks, err := s.SplitKey(stub, kv.GetKey())
		if err != nil {
			return nil, err
		}
		list = append(list, &KV{
			Keys:  ks,
			Value: kv.GetValue(),
		})
	}
	return
}

func (s *CompositeKey) SelectFunc(stub shim.ChaincodeStubInterface, f func(kv *KV) error, keys ...string) (err error) {
	iter, err := stub.GetStateByPartialCompositeKey(s.GetType(), keys)
	if err != nil {
		return
	}
	defer iter.Close()

	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return err
		}
		ks, err := s.SplitKey(stub, kv.GetKey())
		if err != nil {
			return err
		}
		err = f(&KV{
			Keys:  ks,
			Value: kv.GetValue(),
		})
		if err != nil {
			return err
		}
	}
	return
}

func (s *CompositeKey) SplitKey(stub shim.ChaincodeStubInterface, key string) ([]string, error) {
	_, ks, err := stub.SplitCompositeKey(key)
	if err != nil {
		return nil, err
	}
	err = s.check(ks)
	if err != nil {
		return nil, err
	}
	return ks, nil
}

func (s *CompositeKey) GetValue(stub shim.ChaincodeStubInterface, keys []string) ([]byte, error) {
	key, err := s.createCompositeKey(stub, keys)
	if err != nil {
		return nil, err
	}
	return stub.GetState(key)
}

func (s *CompositeKey) GetHistory(stub shim.ChaincodeStubInterface, keys []string, f func(buf []byte) (bool, error)) error {
	key, err := s.createCompositeKey(stub, keys)
	if err != nil {
		return err
	}
	iter, err := stub.GetHistoryForKey(key)
	if err != nil {
		return err
	}
	defer iter.Close()
	for iter.HasNext() {
		data, err := iter.Next()
		if err != nil {
			return err
		}
		ok, err := f(data.GetValue())
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

func (s *CompositeKey) check(keys []string) error {
	if len(keys) != len(s.fields) {
		return fmt.Errorf("keys count not matched. got %d, need %d", len(keys), len(s.fields))
	}
	return nil
}

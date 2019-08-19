package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snlansky/apk/rpc"
	"strings"
	"fmt"
)

func TestCompositeKey(t *testing.T) {
	ck := NewCompositeKey("UserTable", "Name", "Address")
	assert.Equal(t, ck.GetType(), strings.Join([]string{"UserTable", "Name", "Address"}, sep))
	helper := rpc.NewHelper(rpc.DefaultInitFunction)
	stub := shim.NewMockStub("Test", helper)
	stub.MockTransactionStart("init")

	err := ck.Insert(stub, []string{"lili", "xx-xxxx-xxx"}, []byte("v1"))
	assert.NoError(t, err)
	err = ck.Insert(stub, []string{"bill", "yy-yyyy-yyy"}, []byte("v2"))
	assert.NoError(t, err)

	buf, err := ck.GetValue(stub, []string{"lili", "xx-xxxx-xxx"})
	assert.NoError(t, err)
	assert.Equal(t, buf, []byte("v1"))

	err = ck.Delete(stub, []string{"lili", "xx-xxxx-xxx"})
	assert.NoError(t, err)
	buf, err = ck.GetValue(stub, []string{"lili", "xx-xxxx-xxx"})
	assert.NoError(t, err)
	assert.Equal(t, buf==nil, true)

	err = ck.Insert(stub, []string{"lili1", "xx-xxxx-xxx1"}, []byte("v1"))
	assert.NoError(t, err)
	err = ck.Insert(stub, []string{"lili2", "xx-xxxx-xxx2"}, []byte("v1"))
	assert.NoError(t, err)
	err = ck.Insert(stub, []string{"lili1", "xx-xxxx-xxx3"}, []byte("v1"))
	assert.NoError(t, err)

	err = ck.Update(stub, []string{"lili1", "xx-xxxx-xxx3"}, []byte("v3"))
	assert.NoError(t, err)

	buf, err = ck.GetValue(stub, []string{"lili1", "xx-xxxx-xxx3"})
	assert.NoError(t, err)
	assert.Equal(t, buf, []byte("v3"))


	kvList, err := ck.Select(stub, "lili1")
	assert.NoError(t, err)
	for _, kv := range kvList {
		fmt.Println(kv.Keys, string(kv.Value))
	}
	key, err:= ck.createCompositeKey(stub, []string{"lili1", "xx-xxxx-xxx3"})
	assert.NoError(t, err)
	fmt.Println("------------------------------", key, )

	err = ck.SelectFunc(stub, func(kv *KV) error {
		fmt.Println(kv.Keys, string(kv.Value))
		return nil
	}, "lili1")
	assert.NoError(t, err)
}

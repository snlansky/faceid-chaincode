package rpc

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strconv"
	"fmt"
	"errors"
	"runtime"
)

type mockChainCode struct {
}

func (m *mockChainCode) Add(a, b int) int {
	return a + b
}

func (m *mockChainCode) Addf(a float64, b int16) float64 {
	c := float64(b) + a
	return c
}

func (m *mockChainCode) SayHello() []byte {
	return []byte("hello")
}

func (m *mockChainCode) Print(s string, name string, age int) string {
	return fmt.Sprintf(s, name, age)
}

func (m *mockChainCode) internal(a int) []byte {
	return []byte(strconv.Itoa(a))
}

func (m *mockChainCode) GetPanic(s string) []byte {
	panic(s)
}

func (m *mockChainCode) GetPanicError() []byte {
	panic(errors.New("something is wrong"))
}

func (m *mockChainCode) GerPanicRuntime() []byte {
	err := &runtime.TypeAssertionError{}
	panic(err)
}

func (m *mockChainCode) GetNull() {
}

func (m *mockChainCode) GetError(t bool) error {
	if t {
		return errors.New("error")
	}
	return nil
}

func TestNewRpcHelper(t *testing.T) {
	rpc := NewRpcHelper()
	rpc.RegisterMethod(&mockChainCode{})

	ret, err := rpc.Handler(&Request{
		Method: "Add",
		Params: []interface{}{12, 14},
	})
	assert.NoError(t, err)
	assert.Equal(t, ret, 26)

	ret, err = rpc.Handler(&Request{
		Method: "Addf",
		Params: []interface{}{123.563, float64(56)},
	})
	assert.Equal(t, ret, 123.563+float64(56))

	ret, err = rpc.Handler(&Request{
		Method: "SayHello",
		Params: []interface{}{},
	})
	assert.Equal(t, ret, []byte("hello"))

	s := "hello name:%s, age:%d"
	name := "lilith"
	age := 12
	ret, err = rpc.Handler(&Request{
		Method: "Print",
		Params: []interface{}{s, name, age},
	})
	assert.Equal(t, ret, fmt.Sprintf(s, name, age))

	ret, err = rpc.Handler(&Request{
		Method: "Print",
		Params: []interface{}{"hello name:%s, age:%d", "lilith"},
	})
	assert.Error(t, err)

	ret, err = rpc.Handler(&Request{
		Method: "Print",
		Params: []interface{}{"hello name:%s, age:%d", "lilith", ""},
	})
	assert.Error(t, err)

	ret, err = rpc.Handler(&Request{
		Method: "Print",
		Params: []interface{}{"12", "lilith", 12},
	})
	assert.NoError(t, err)

	ret, err = rpc.Handler(&Request{
		Method: "internal",
		Params: []interface{}{12},
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "[RPC]: method not found.[method=internal]")

	ret, err = rpc.Handler(&Request{
		Method: "GetPanic",
		Params: []interface{}{"lilith"},
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "lilith")

	ret, err = rpc.Handler(&Request{
		Method: "GetPanicError",
		Params: []interface{}{"lilith"},
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "[RPC]: params not matched. got 1, need 0.")

	ret, err = rpc.Handler(&Request{
		Method: "GetPanicError",
		Params: []interface{}{},
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "something is wrong")

	ret, err = rpc.Handler(&Request{
		Method: "GerPanicRuntime",
		Params: []interface{}{},
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), ERR_RUNTIME)

	ret, err = rpc.Handler(&Request{
		Method: "GetNull",
		Params: []interface{}{},
	})
	assert.NoError(t, err)
	assert.Nil(t, ret)

	ret, err = rpc.Handler(&Request{
		Method: "GetError",
		Params: []interface{}{false},
	})
	assert.NoError(t, err)
	assert.Nil(t, ret)

	ret, err = rpc.Handler(&Request{
		Method: "GetError",
		Params: []interface{}{true},
	})
	assert.Error(t, err)
	assert.Equal(t, err, errors.New("error"))
}

type Super struct {
}

func (s *Super) Print() string {
	return "super print"
}

type Sub struct {
	s Super
}

func (s *Sub) Say() string {
	return "sub say"
}

func TestHelper_Register(t *testing.T) {
	rpc := NewRpcHelper()
	rpc.RegisterMethod(new(Sub))

	v, err := rpc.Handler(&Request{
		Method: "Print",
		Params: []interface{}{},
	})
	fmt.Println(v, err)
	assert.Error(t, err)
}
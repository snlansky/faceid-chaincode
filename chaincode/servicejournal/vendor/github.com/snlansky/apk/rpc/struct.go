package rpc

import (
	"reflect"
	"strings"
)

type Request struct {
	Method string        `json:"func_name"`
	Params []interface{} `json:"params"`
}

type Method struct {
	Method reflect.Method
	host   reflect.Value
	idx    int
}

type InternalError struct {
	err  error
	info []string
}

func NewInternalError(err error, info ... string) *InternalError {
	return &InternalError{err: err, info: info}
}

func (e *InternalError) Error() string {
	msg := e.info[:]
	if e.err != nil {
		msg = append(msg, e.err.Error())
	}
	return strings.Join(msg, ",")
}

func (e *InternalError) External() string {
	return strings.Join(e.info, ",")
}

func Check(err error, info ...string) {
	if err != nil {
		switch e := err.(type) {
		case *InternalError:
			panic(e)
		default:
			panic(&InternalError{err: err, info: info})
		}
	}
}

func Throw(info ...string) {
	panic(NewInternalError(nil, info...))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Response struct {
	Request *Request    `json:"request"`
	Result  interface{} `json:"result"`
}

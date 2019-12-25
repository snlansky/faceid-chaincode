package rpc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"unicode"
	"unicode/utf8"
)

// 方法名注册规则
type Rule func(typeName string, funcName string) string

// RPC interface
type Rpc interface {
	Register(v interface{}, rule ...Rule) // 将对象的方法按照一定的变换规则进行映射
	Handler(req *Request, baseParam ...interface{}) (interface{}, error)
}

type method struct {
	m    reflect.Method
	host reflect.Value
	idx  int
}

type rpcImpl struct {
	mths map[string]*method
}

func New() Rpc {
	return &rpcImpl{make(map[string]*method)}
}

//处理客户端消息
// baseParam 一般为基础参数如client 等
func (h *rpcImpl) Handler(req *Request, baseParam ...interface{}) (interface{}, error) {
	method, params, err := h.parse(req, baseParam...)
	if err != nil {
		return nil, err
	}

	return h.call(method, params)
}

// 解析客户端请求，
// default_params : 服务器端调用时自动带入的参数, 和客户端请求的参数共同组成method的参数。
func (h *rpcImpl) parse(req *Request, defaultParams ...interface{}) (*method, []reflect.Value, error) {
	method, ok := h.mths[req.Method]
	if !ok {
		return nil, nil, fmt.Errorf("[RPC]: method not found.[method=%s]", req.Method)
	}

	defaultParamsLen := len(defaultParams)
	//长度应减去method的receiver
	var lens int
	if req.Params == nil {
		lens = 0
	} else {
		lens = len(req.Params)
	}

	if lens != (method.m.Type.NumIn() - defaultParamsLen - 1) {
		return nil, nil, fmt.Errorf("[RPC]: params not matched. got %d, need %d", lens, method.m.Type.NumIn()-defaultParamsLen-1)
	}

	params := make([]reflect.Value, lens+defaultParamsLen)
	for idx, hdnParam := range defaultParams {
		params[idx] = reflect.ValueOf(hdnParam)
	}

	for i := 0; i < lens; i++ {
		targetType := method.m.Type.In(i + 1 + defaultParamsLen) //跳过receiver和default_params
		newParam, ok := convertParam(req.Params[i], targetType)
		if !ok {
			return nil, nil, fmt.Errorf("[RPC]: convert param faild. expect %s, found=%v value=%v", targetType, reflect.TypeOf(req.Params[i]), req.Params[i])
		}
		params[i+defaultParamsLen] = newParam
	}

	return method, params, nil
}

func (h *rpcImpl) call(mth *method, params []reflect.Value) (ret interface{}, err error) {

	defer func() {
		if re := recover(); re != nil {
			switch v := re.(type) {
			case InternalError:
				//logger.Error(v.Error())
				err = errors.New(v.External())
			case *InternalError:
				//logger.Error(v.Error())
				err = errors.New(v.External())
			case string:
				err = errors.New(v)
			case runtime.Error:
				fmt.Printf("runtime error:%v \nstack :%s\n", v, string(debug.Stack()))
				err = ErrRuntime
			case error:
				err = v
			default:
				err = fmt.Errorf("ohter error type:%v, value:%v", reflect.TypeOf(re), v)
			}
		}
	}()

	result := mth.host.Method(mth.idx).Call(params)
	if len(result) > 0 {
		ret = result[0].Interface()
		// if return error
		if ret != nil {
			if v, ok := ret.(error); ok {
				err = v
			}
		}
	}
	return
}

func (h *rpcImpl) Register(v interface{}, rule ...Rule) {
	reflectType := reflect.TypeOf(v)
	host := reflect.ValueOf(v)
	for i := 0; i < reflectType.NumMethod(); i++ {
		m := reflectType.Method(i)
		char, _ := utf8.DecodeRuneInString(m.Name)
		//非导出函数不注册
		if !unicode.IsUpper(char) {
			continue
		}

		fn := m.Name
		if len(rule) > 0 {
			fn = rule[0](h.splitTypeName(reflectType.String()), m.Name)
		}

		fmt.Printf("[RPC] register functon: %s\n", fn)
		h.mths[fn] = &method{m, host, m.Index}
	}
}

func (h *rpcImpl) splitTypeName(complete string) string {
	part := strings.Split(complete, ".")
	return part[len(part)-1]
}
func DefaultRule(typeName, funcName string) string {
	return fmt.Sprintf("%s.%s", typeName, funcName)
}

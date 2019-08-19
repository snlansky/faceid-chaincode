package rpc

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"unicode"
	"unicode/utf8"
)

var logger = shim.NewLogger("RPC")

// 方法名注册规则
type Rule func(typeName, funcName string) string

//----------------------------------------------------------------------------------------

type Rpc struct {
	Methods map[string]*Method
}

func NewRpcHelper() *Rpc {
	return &Rpc{make(map[string]*Method)}
}

//处理客户端消息
// baseParam 一般为基础参数如client 等
func (h *Rpc) Handler(req *Request, baseParam ...interface{}) (interface{}, error) {
	method, params, err := h.Parse(req, baseParam...)
	if err != nil {
		return nil, err
	}

	return h.Call(method, params)
}

// 解析客户端请求，
// default_params : 服务器端调用时自动带入的参数, 和客户端请求的参数共同组成method的参数。
func (h *Rpc) Parse(req *Request, defaultParams ...interface{}) (*Method, []reflect.Value, error) {
	method, ok := h.Methods[req.Method]
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

	if lens != (method.Method.Type.NumIn() - defaultParamsLen - 1) {
		return nil, nil, fmt.Errorf("[RPC]: params not matched. got %d, need %d.", lens, method.Method.Type.NumIn()-defaultParamsLen-1)
	}

	params := make([]reflect.Value, lens+defaultParamsLen)
	for idx, hdnParam := range defaultParams {
		params[idx] = reflect.ValueOf(hdnParam)
	}

	for i := 0; i < lens; i++ {
		targetType := method.Method.Type.In(i + 1 + defaultParamsLen) //跳过receiver和default_params
		newParam, ok := convertParam(req.Params[i], targetType)
		if !ok {
			return nil, nil, fmt.Errorf("[RPC]: convert param faild. expect %s, found=%v value=%v.", targetType, reflect.TypeOf(req.Params[i]), req.Params[i])
		}
		params[i+defaultParamsLen] = newParam
	}

	return method, params, nil
}

func (h *Rpc) Call(method *Method, params []reflect.Value) (ret interface{}, err error) {

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
				logger.Errorf("runtime error:%v \nstack :%s", v, string(debug.Stack()))
				err = errors.New(ERR_RUNTIME)
			case error:
				err = v
			default:
				err = fmt.Errorf("ohter error type:%v, value:%v", reflect.TypeOf(re), v)
			}
		}
	}()

	result := method.host.Method(method.idx).Call(params)
	if len(result) > 0 {
		ret = result[0].Interface()
		// if return error
		if ret != nil {
			if v, ok := ret.(error); ok {
				panic(v)
			}
		}
	}
	return
}

func (h *Rpc) RegisterMethod(v interface{}) {
	reflectType := reflect.TypeOf(v)
	host := reflect.ValueOf(v)
	for i := 0; i < reflectType.NumMethod(); i++ {
		m := reflectType.Method(i)
		char, _ := utf8.DecodeRuneInString(m.Name)
		if !unicode.IsUpper(char) {
			continue
		}
		h.Methods[m.Name] = &Method{m, host, m.Index}
	}
}

// 将对象的方法按照一定的变换规则进行映射
func (h *Rpc) RegisterMethodByRule(v interface{}, rule Rule) {
	reflectType := reflect.TypeOf(v)
	host := reflect.ValueOf(v)
	for i := 0; i < reflectType.NumMethod(); i++ {
		m := reflectType.Method(i)
		char, _ := utf8.DecodeRuneInString(m.Name)
		//非导出函数不注册
		if !unicode.IsUpper(char) {
			continue
		}
		h.Methods[rule(h.SplitTypeName(reflectType.String()), m.Name)] = &Method{m, host, m.Index}
	}
}

func (h *Rpc) SplitTypeName(complete string) string {
	part := strings.Split(complete, ".")
	return part[len(part)-1]
}

func DefaultRule(typeName, funcName string) string {
	name := fmt.Sprintf("%s.%s", typeName, funcName)
	logger.Infof("register : %s", name)
	return name
}

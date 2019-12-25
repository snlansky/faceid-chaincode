package rpc

type Request struct {
	Method string        `json:"func_name"`
	Params []interface{} `json:"params"`
}

type Response struct {
	Request *Request    `json:"request"`
	Result  interface{} `json:"result"`
}

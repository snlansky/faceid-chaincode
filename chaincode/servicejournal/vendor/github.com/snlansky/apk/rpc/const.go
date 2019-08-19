package rpc

const (
	ERR_NOT_FIND_INIT_FUNCTION = "ERR_NOT_FIND_INIT_FUNCTION" // 没有没找初始化函数
	ERR_JSON_UNMARSHAL         = "ERR_JSON_UNMARSHAL"         // 读取json数据错误
	ERR_PARSE_RPC_REQ          = "ERR_PARSE_RPC_REQ"          // 解析RPC请求错误
	ERR_METHOD_NOT_FOUND       = "ERR_METHOD_NOT_FOUND"       // 没有找到方法
	ERR_PARAM_COUNT_NOT_MATCH  = "ERR_PARAM_COUNT_NOT_MATCH"  // 调用参数数量不匹配
	ERR_PARAM_INVALID          = "ERR_PARAM_INVALID"          // 参数错误
	ERR_JSON_MARSHAL           = "ERR_JSON_MARSHAL"           // json数据错误
	ERR_RUNTIME                = "ERR_RUNTIME"                // 运行时错误
	ERR_INTERNAL_INVALID       = "ERR_INTERNAL_INVALID"       // 内部错误
	ERR_INVALID_CERT           = "ERR_INVALID_CERT"           // 错误的证书
)

const SERVICE_BATCH_INVOKE = "Service.BatchInvoke"

package main

var SUCCESS = []byte("SUCCESS")

const (
	CertificateValidationFailure = "证书验证失败"
	InternalError                = "内部错误"
	ReduplicateCreate            = "重复创建"
	NotExistResource             = "资源不存在"
	FindResourceFailed           = "查找资源失败"
	SaveResourcFailed            = "保存资源失败"
	NotValidationResource        = "资源验证失败"
	Unauthorized                 = "未授权"
	NotFoundApplyRecord          = "未找到申请记录"
)

const NotImplemented = "NotImplemented"

type OptionType = string

const (
	AppName = "faceid"
)

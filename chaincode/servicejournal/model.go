package main

type Address string

type TicketCommon struct {
	Id          string            `json:"id" valid:"optional"`           // 提单id
	Title       string            `json:"title" valid:"required"`        // 标题（需要能更新）
	Description string            `json:"description" valid:"optional"`  // 描述
	Products    []string          `json:"produces" valid:"optional"`     // 产品信息
	System      string            `json:"system" valid:"required"`       // 来源系统
	Status      string            `json:"status" valid:"required"`       // 状态
	OwnerId     Address           `json:"owner_id" valid:"required"`     // 所有者数字身份
	SubmitterId Address           `json:"submitter_id" valid:"required"` // 提单人数字身份
	HandlerId   Address           `json:"handler_id" valid:"optional"`   // 处理人数字身份（需要能更新）
	CreateTime  int64             `json:"create_time" valid:"required"`  // 创建时间
	UpdateTime  int64             `json:"update_time" valid:"required"`  // 更新时间
	UploadTime  int64             `json:"upload_time" valid:"optional"`  // 上链时间 只计算第一次
	Details     map[string]string `json:"details" valid:"optional"`      // 详情
}

type Ticket struct {
	TicketCommon
	SourceList []Address          `json:"source_list" valid:"optional"` // 数据来源:企业的数字身份列表
	NodeList   []string           `json:"node_list" valid:"optional"`   // 节点列表
	Metadata   map[Address]string `json:"metadata" valid:"optional"`    // 元数据只能链上可见 key:数字身份，value:详细信息json
}

type TicketRequest struct {
	TicketCommon
	//Source    Address `json:"source" valid:"optional"`    // 数据来源:企业的数字身份
	Extension string `json:"extension" valid:"optional"` // 扩展信息
}

type TicketResponse struct {
	TicketCommon
	SourceList []Address `json:"source_list" valid:"optional"` // 数据来源:企业的数字身份列表
	NodeList   []string  `json:"node_list" valid:"optional"`   // 节点列表
	Extension  string    `json:"extension" valid:"optional"`   // 扩展信息
}

// ---------------------------------------------------------------

type NodeCommon struct {
	Id          string  `json:"id" valid:"optional"`           // 节点id
	TicketID    string  `json:"ticket_id" valid:"required"`    // 提单id
	HandlerId   Address `json:"handler_id" valid:"optional"`   // 处理人数字身份
	Status      string  `json:"status" valid:"required"`       // 状态
	CreateTime  int64   `json:"create_time" valid:"required"`  // 创建时间
	UpdateTime  int64   `json:"update_time" valid:"optional"`  // 更新时间
	UploadTime  int64   `json:"upload_time" valid:"optional"`  // 上链时间
	Description string  `json:"description"  valid:"optional"` // 处理意见
	System      string  `json:"system" valid:"optional"`       // 来源系统
}

type Node struct {
	NodeCommon
	SourceList []Address          `json:"source_list" valid:"optional"` // 数据来源:企业的数字身份列表
	Metadata   map[Address]string `json:"metadata" valid:"optional"`    // 元数据只能链上可见 key:数字身份，value:详细信息json
}

type NodeRequest struct {
	NodeCommon
	Extension string `json:"extension" valid:"optional"` // 扩展信息
}

type NodeResponse struct {
	NodeCommon
	Extension string `json:"extension" valid:"optional"` // 扩展信息
}

// -----------------------------------------------------------------
type IdMapping struct {
	TicketId    string   `json:"ticket_id" valid:"required"`    // 提单id
	InternalIds []string `json:"internal_ids" valid:"required"` // 其他系统的内部ID
}

package main

type TicketCommon struct {
	Id          string `json:"id" valid:"required"`           // 提单id
	Title       string `json:"title" valid:"required"`        // 标题（需要能更新）
	ProductName string `json:"product_name" valid:"required"` // 产品名
	ServiceName string `json:"service_name" valid:"required"` // 服务名
	Status      string `json:"status" valid:"required"`       // 状态
	OwnerId     string `json:"owner_id" valid:"required"`     // 所有者id
	SubmitterId string `json:"submitter_id" valid:"required"` // 提单人id
	HandlerId   string `json:"handler_id" valid:"optional"`   // 处理人id（需要能更新）
	CreateTime  int64  `json:"create_time" valid:"required"`  // 创建时间
	UpdateTime  int64  `json:"update_time" valid:"required"`  // 更新时间
	Extension   string `json:"extension" valid:"optional"`    // 扩展信息
}

type NodeCommon struct {
	Id          string `json:"id" valid:"required"`          // 节点id
	TicketID    string `json:"ticket_id" valid:"required"`   // 提单id
	Operation   string `json:"operation" valid:"required"`   // 操作内容
	Description string `json:"description" valid:"required"` // 处理意见（不知道叫什么）
	HandlerId   string `json:"handler_id" valid:"required"`  // 处理人id
	CreateTime  string `json:"create_time" valid:"required"` // 创建时间
	Extension   string `json:"extension" valid:"optional"`   // 扩展信息
}

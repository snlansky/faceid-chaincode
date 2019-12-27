package main

type Validator interface {
	Validate() bool
}

type FaceID struct {
	ID         string                 `json:"id" valid:"optional"`          // ID
	SourceType string                 `json:"source_type" valid:"required"` // 资源类型
	SourceHash string                 `json:"source_hash" valid:"required"` // 资源hash
	Algorithm  string                 `json:"algorithm" valid:"required"`   // hash 算法
	Labels     []string               `json:"labels" valid:"optional"`      // 标签
	Metadata   map[string]interface{} `json:"metadata" valid:"optional"`    // 元数据
	Timestamp  int64                  `json:"timestamp" valid:"optional"`   // 时间戳(s)
}

func (f *FaceID) Validate() bool {
	if f.SourceType == "" || f.SourceHash == "" || f.Algorithm == "" {
		return false
	}
	return true
}

type RequestFaceIDHistory struct {
	StartTime int64    `json:"start_time" valid:"optional"` // 开始时间
	EndTime   int64    `json:"end_time" valid:"required"`   // 结束时间
	Labels    []string `json:"labels" valid:"optional"`     // 标签
}

func (r *RequestFaceIDHistory) Validate() bool {
	if r.EndTime <= 0 || r.StartTime < 0 {
		return false
	}
	return true
}

type TimeIndex struct {
	FaceID    string
	Timestamp int64
}

type User struct {
	RegisterFaceID      *FaceID
	RegisterCertificate *FaceID
}

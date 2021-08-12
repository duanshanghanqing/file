package util

// 定义接口返回的消息结构体
type responseMessage struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
}

const (
	success     = iota // 0: 成功
	fail               // 1：失败
	authInvalid        // 2：auth 过期
)

// 空消息
func NewResponseMessage() *responseMessage {
	return &responseMessage{}
}

// 成功
func NewResponseSuccessMessage() *responseMessage {
	rm := NewResponseMessage()
	rm.Code = success
	rm.Message = "请求成功"
	return rm
}

// 失败
func NewResponseFailMessage() *responseMessage {
	rm := NewResponseMessage()
	rm.Code = fail
	rm.Message = "请求失败"
	return rm
}

// 身份无效
func NewResponseAuthInvalidMessage() *responseMessage {
	rm := NewResponseMessage()
	rm.Code = authInvalid
	rm.Message = "auth 过期"
	return rm
}

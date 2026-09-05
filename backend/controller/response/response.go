package response

import (
	"encoding/json"

	"github.com/ydzydzydz/pmail_telegram_push/logger"
)

// Response 响应体
type Response struct {
	Code    int    `json:"code"`    // 状态码 0 成功 -1 失败
	Message string `json:"message"` // 提示信息
	Data    any    `json:"data"`    // 数据
}

// SuccessResponse 成功响应
func SuccessResponse(message string, data any) *Response {
	return &Response{
		Code:    0,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(message string) *Response {
	return &Response{
		Code:    -1,
		Message: message,
	}
}

// Json 序列化响应体
func (r *Response) Json() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("marshal response failed")
		return ""
	}
	return string(jsonBytes)
}

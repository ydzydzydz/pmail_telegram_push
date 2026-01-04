package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	"github.com/ydzydzydz/pmail_telegram_push/logger"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

var (
	//go:embed dist/index.html
	SettingHtml string // 设置页面
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
	json, err := json.Marshal(r)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("marshal response failed")
		return ""
	}
	return string(json)
}

// getSetting 获取Telegram Push设置
func (h *PmailTelegramPushHook) getSetting(userID int) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("获取Telegram Push设置")

	setting, err := h.settingService.GetSetting(userID)
	if err != nil {
		return ErrorResponse("获取Telegram Push设置失败").Json()
	}

	return SuccessResponse("获取Telegram Push设置成功", setting).Json()
}

func (h *PmailTelegramPushHook) testMessage(userID int, requestData string) string {
	const testMessageResource = "testMessage"
	logger.PluginLogger.Info().Int("user_id", userID).Msg("测试Telegram Push消息")
	var setting model.TelegramPushSetting
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("反序列化测试消息请求失败")
		return ErrorResponse("反序列化测试消息请求失败").Json()
	}
	if setting.ChatID == "" {
		logger.PluginLogger.Error().Msg("测试消息 Chat ID 不能为空")
		return ErrorResponse("测试消息 Chat ID 不能为空").Json()
	}

	_, err := h.sendTestMessage(&setting)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("发送测试消息失败")
		return ErrorResponse("发送测试消息失败").Json()
	}

	return SuccessResponse("测试Telegram Push消息成功", nil).Json()
}

// updateSetting 更新Telegram Push设置
func (h *PmailTelegramPushHook) updateSetting(userID int, requestData string) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("更新Telegram Push设置")

	var setting model.TelegramPushSetting
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("反序列化设置请求失败")
		return ErrorResponse("反序列化设置请求失败").Json()
	}

	setting.UserID = userID
	setting.ChatID = strings.TrimSpace(setting.ChatID)
	if err := h.settingService.UpdateSetting(userID, &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("更新Telegram Push设置失败")
		return ErrorResponse("更新Telegram Push设置失败").Json()
	}

	return SuccessResponse("更新Telegram Push设置成功", nil).Json()
}

// BotInfo 机器人信息
type BotInfo struct {
	Username string `json:"username"` // 机器人用户名
	BotLink  string `json:"bot_link"` // 机器人链接
}

// NewBotInfo 创建机器人信息
func NewBotInfo(username string) BotInfo {
	return BotInfo{
		Username: username,
		BotLink:  fmt.Sprintf("https://t.me/%s", username),
	}
}

// getBotInfo 获取Telegram Bot信息
func (h *PmailTelegramPushHook) getBotInfo() string {
	logger.PluginLogger.Info().Msg("获取Telegram Bot信息")

	me, err := h.bot.GetMe(context.Background())
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("获取Telegram Bot信息失败")
		return ErrorResponse("获取Telegram Bot信息失败").Json()
	}

	return SuccessResponse("获取Telegram Bot信息成功", NewBotInfo(me.Username)).Json()
}

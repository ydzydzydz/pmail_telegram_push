package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ydzydzydz/pmail_telegram_push/controller/response"
	"github.com/ydzydzydz/pmail_telegram_push/logger"
	"github.com/ydzydzydz/pmail_telegram_push/model"
	"github.com/ydzydzydz/pmail_telegram_push/sender"
)

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

// BotController 机器人控制器
type BotController struct {
	sender *sender.TelegramBotSender
}

// NewBotController 创建机器人控制器
func NewBotController(sender *sender.TelegramBotSender) *BotController {
	return &BotController{
		sender: sender,
	}
}

// getBotInfo 获取Telegram Bot信息
func (c *BotController) GetBotInfo() string {
	logger.PluginLogger.Info().Msg("获取Telegram Bot信息")

	me, err := c.sender.GetBot().GetMe(context.Background())
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("获取Telegram Bot信息失败")
		return response.ErrorResponse("获取Telegram Bot信息失败").Json()
	}

	return response.SuccessResponse("获取Telegram Bot信息成功", NewBotInfo(me.Username)).Json()
}

// SendTestMessage 发送测试消息
func (c *BotController) SendTestMessage(userID int, requestData string) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("测试Telegram Push消息")
	var setting model.PluginTelegramPushSettingModel
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("反序列化测试消息请求失败")
		return response.ErrorResponse("反序列化测试消息请求失败").Json()
	}
	if setting.ChatID == "" {
		logger.PluginLogger.Error().Msg("测试消息 Chat ID 不能为空")
		return response.ErrorResponse("测试消息 Chat ID 不能为空").Json()
	}

	err := c.sender.SendTestMessage(context.Background(), &setting)
	if err != nil {
		logger.PluginLogger.Error().Err(err).Msg("发送测试消息失败")
		return response.ErrorResponse("发送测试消息失败").Json()
	}

	return response.SuccessResponse("测试Telegram Push消息成功", nil).Json()
}

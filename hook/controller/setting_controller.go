package controller

import (
	"encoding/json"
	"strings"

	"github.com/ydzydzydz/pmail_telegram_push/hook/controller/response"
	"github.com/ydzydzydz/pmail_telegram_push/hook/logger"
	"github.com/ydzydzydz/pmail_telegram_push/hook/model"
	"github.com/ydzydzydz/pmail_telegram_push/hook/service"
)

type SettingController struct {
	settingService *service.SettingService
}

// NewSettingController 创建设置控制器实例
func NewSettingController(settingService *service.SettingService) *SettingController {
	return &SettingController{settingService: settingService}
}

// GetSetting 获取Telegram Push设置
func (c *SettingController) GetSetting(userID int) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("获取Telegram Push设置")

	setting, err := c.settingService.GetSetting(userID)
	if err != nil {
		return response.ErrorResponse("获取Telegram Push设置失败").Json()
	}

	return response.SuccessResponse("获取Telegram Push设置成功", setting).Json()
}

// UpdateSetting 更新Telegram Push设置
func (c *SettingController) UpdateSetting(userID int, requestData string) string {
	logger.PluginLogger.Info().Int("user_id", userID).Msg("更新Telegram Push设置")

	var setting model.PluginTelegramPushSettingModel
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("反序列化设置请求失败")
		return response.ErrorResponse("反序列化设置请求失败").Json()
	}

	setting.UserID = userID
	setting.ChatID = strings.TrimSpace(setting.ChatID)
	if err := c.settingService.UpdateSetting(userID, &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("更新Telegram Push设置失败")
		return response.ErrorResponse("更新Telegram Push设置失败").Json()
	}

	return response.SuccessResponse("更新Telegram Push设置成功", nil).Json()
}

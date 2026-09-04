package dao

import "github.com/ydzydzydz/pmail_telegram_push/hook/model"

// ISettingDao 是设置数据访问对象的接口
type ISettingDao interface {
	// GetSetting 获取用户的设置
	GetSetting(userID int) (*model.PluginTelegramPushSettingModel, error)
	// UpdateSetting 更新用户的设置
	UpdateSetting(userID int, setting *model.PluginTelegramPushSettingModel) error
	// CreateSetting 创建用户的设置
	CreateSetting(setting *model.PluginTelegramPushSettingModel) error
	// GetOrCreate 获取用户的设置，如果不存在则创建
	GetOrCreate(userID int, setting *model.PluginTelegramPushSettingModel) (*model.PluginTelegramPushSettingModel, error)
}

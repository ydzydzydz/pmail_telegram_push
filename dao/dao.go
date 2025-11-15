package dao

import "github.com/ydzydzydz/pmail_telegram_push/model"

// ISettingDao 是设置数据访问对象的接口
type ISettingDao interface {
	// GetSetting 获取用户的设置
	GetSetting(userID int) (*model.TelegramPushSetting, error)
	// UpdateSetting 更新用户的设置
	UpdateSetting(userID int, setting *model.TelegramPushSetting) error
	// CreateSetting 创建用户的设置
	CreateSetting(setting *model.TelegramPushSetting) error
	// ExistSetting 检查用户的设置是否存在
	ExistSetting(userID int) bool
}

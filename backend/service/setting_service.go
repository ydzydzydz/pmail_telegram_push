package service

import (
	"github.com/ydzydzydz/pmail_telegram_push/dao"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

// SettingService 设置服务
type SettingService struct {
	dao dao.ISettingDao
}

const (
	DefaultChatID             = ""   // 默认聊天id, 为空时, 则不发送
	DefaultShowContent        = true // 是否显示邮件内容
	DefaultSpoilerContent     = true // 是否 spoiler 显示邮件内容
	DefaultSendAttachments    = true // 是否发送附件
	DefaultDisableLinkPreview = true // 是否禁用链接预览
)

// NewSettingService 创建设置服务实例
func NewSettingService(dao dao.ISettingDao) *SettingService {
	return &SettingService{dao: dao}
}

// GetSetting 获取设置
// 如果不存在, 则创建默认设置
func (s *SettingService) GetSetting(userID int) (*model.PluginTelegramPushSettingModel, error) {
	defaultSetting := &model.PluginTelegramPushSettingModel{
		UserID:             userID,
		ChatID:             DefaultChatID,
		ShowContent:        DefaultShowContent,
		SpoilerContent:     DefaultSpoilerContent,
		SendAttachments:    DefaultSendAttachments,
		DisableLinkPreview: DefaultDisableLinkPreview,
	}
	return s.dao.GetOrCreate(userID, defaultSetting)
}

// UpdateSetting 更新设置
// 如果不存在, 则创建默认设置后更新
func (s *SettingService) UpdateSetting(userID int, setting *model.PluginTelegramPushSettingModel) error {
	defaultSetting := &model.PluginTelegramPushSettingModel{
		UserID:             userID,
		ChatID:             DefaultChatID,
		ShowContent:        DefaultShowContent,
		SpoilerContent:     DefaultSpoilerContent,
		SendAttachments:    DefaultSendAttachments,
		DisableLinkPreview: DefaultDisableLinkPreview,
	}
	_, err := s.dao.GetOrCreate(userID, defaultSetting)
	if err != nil {
		return err
	}
	return s.dao.UpdateSetting(userID, setting)
}

// CreateDefaultSetting 创建默认设置
func (s *SettingService) CreateDefaultSetting(userID int) error {
	setting := &model.PluginTelegramPushSettingModel{
		UserID:             userID,
		ChatID:             DefaultChatID,
		ShowContent:        DefaultShowContent,
		SpoilerContent:     DefaultSpoilerContent,
		SendAttachments:    DefaultSendAttachments,
		DisableLinkPreview: DefaultDisableLinkPreview,
	}
	return s.dao.CreateSetting(setting)
}

package service

import (
	"errors"
	"strings"

	"github.com/ydzydzydz/pmail_telegram_push/model"
	"xorm.io/xorm"
)

var ErrSettingNotFound = errors.New("setting not found")

// SettingService 设置服务
type SettingService struct {
	db *xorm.Engine
}

const (
	DefaultChatID             = ""   // 默认聊天id, 为空时, 则不发送
	DefaultShowContent        = true // 是否显示邮件内容
	DefaultSpoilerContent     = true // 是否 spoiler 显示邮件内容
	DefaultSendAttachments    = true // 是否发送附件
	DefaultDisableLinkPreview = true // 是否禁用链接预览
)

// NewSettingService 创建设置服务实例
func NewSettingService(db *xorm.Engine) *SettingService {
	return &SettingService{db: db}
}

// findSetting 按 user_id 查询单条设置记录
func (s *SettingService) findSetting(userID int) (*model.PluginTelegramPushSettingModel, bool, error) {
	setting := new(model.PluginTelegramPushSettingModel)
	has, err := s.db.Where("user_id = ?", userID).Get(setting)
	if err != nil {
		return nil, false, err
	}
	return setting, has, nil
}

// getOrCreateSetting 原子化的 GetOrCreate，避免并发插入时的重复键冲突
func (s *SettingService) getOrCreateSetting(userID int, defaultSetting *model.PluginTelegramPushSettingModel) (*model.PluginTelegramPushSettingModel, error) {
	existing, has, err := s.findSetting(userID)
	if err != nil {
		return nil, err
	}
	if has {
		return existing, nil
	}

	if _, createErr := s.db.Insert(defaultSetting); createErr == nil {
		return defaultSetting, nil
	} else if !isDuplicateErr(createErr) {
		return nil, createErr
	}

	setting, has, err := s.findSetting(userID)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrSettingNotFound
	}
	return setting, nil
}

// isDuplicateErr 判断错误是否为重复键冲突（MySQL / SQLite / PostgreSQL 常见信息）
func isDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate") ||
		strings.Contains(msg, "UNIQUE constraint") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique_violation")
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
	return s.getOrCreateSetting(userID, defaultSetting)
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
	if _, err := s.getOrCreateSetting(userID, defaultSetting); err != nil {
		return err
	}
	_, err := s.db.Where("user_id = ?", userID).AllCols().Update(setting)
	return err
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
	_, err := s.db.Insert(setting)
	return err
}

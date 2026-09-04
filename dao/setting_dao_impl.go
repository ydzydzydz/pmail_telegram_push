package dao

import (
	"github.com/ydzydzydz/pmail_telegram_push/db/repository"
	"github.com/ydzydzydz/pmail_telegram_push/model"
	"xorm.io/xorm"
)

// SettingDaoImpl 实现了 ISettingDao 接口
type SettingDaoImpl struct {
	db   *xorm.Engine
	repo *repository.Repository[model.PluginTelegramPushSettingModel]
}

var _ ISettingDao = (*SettingDaoImpl)(nil)

// NewSettingDaoImpl 创建一个新的 SettingDaoImpl 实例
func NewSettingDaoImpl(db *xorm.Engine) *SettingDaoImpl {
	return &SettingDaoImpl{
		db:   db,
		repo: repository.NewRepository[model.PluginTelegramPushSettingModel](db),
	}
}

// GetSetting 获取用户的设置
func (d *SettingDaoImpl) GetSetting(userID int) (*model.PluginTelegramPushSettingModel, error) {
	return d.repo.FindOne(userID)
}

// UpdateSetting 更新用户的设置
func (d *SettingDaoImpl) UpdateSetting(userID int, setting *model.PluginTelegramPushSettingModel) error {
	return d.repo.Update(userID, setting)
}

// CreateSetting 创建用户的设置
func (d *SettingDaoImpl) CreateSetting(setting *model.PluginTelegramPushSettingModel) error {
	return d.repo.Create(setting)
}

// ExistSetting 检查用户的设置是否存在
func (d *SettingDaoImpl) ExistSetting(userID int) bool {
	return d.repo.Exist(userID)
}

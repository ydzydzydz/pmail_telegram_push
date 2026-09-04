package service

import (
	"testing"

	"github.com/ydzydzydz/pmail_telegram_push/hook/dao"
	"github.com/ydzydzydz/pmail_telegram_push/hook/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

func newTestSettingService(t *testing.T) (*SettingService, func()) {
	db, err := xorm.NewEngine("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	db.ShowSQL(false)
	err = db.Sync2(new(model.PluginTelegramPushSettingModel))
	if err != nil {
		t.Fatalf("failed to sync schema: %v", err)
	}
	settingDao := dao.NewSettingDaoImpl(db)
	settingService := NewSettingService(settingDao)
	return settingService, func() {
		db.Close()
	}
}

func TestSettingService_GetSetting_NotExists(t *testing.T) {
	service, cleanup := newTestSettingService(t)
	defer cleanup()

	setting, err := service.GetSetting(1)
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}

	if setting.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", setting.UserID)
	}
	if setting.ChatID != DefaultChatID {
		t.Errorf("expected ChatID '%s', got '%s'", DefaultChatID, setting.ChatID)
	}
	if setting.ShowContent != DefaultShowContent {
		t.Errorf("expected ShowContent %v, got %v", DefaultShowContent, setting.ShowContent)
	}
}

func TestSettingService_GetSetting_AlreadyExists(t *testing.T) {
	service, cleanup := newTestSettingService(t)
	defer cleanup()

	existingSetting := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "custom_chat_id",
		ShowContent:        false,
		SpoilerContent:     false,
		SendAttachments:    false,
		DisableLinkPreview: false,
	}

	err := service.UpdateSetting(1, existingSetting)
	if err != nil {
		t.Fatalf("UpdateSetting failed: %v", err)
	}

	setting, err := service.GetSetting(1)
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}

	if setting.ChatID != "custom_chat_id" {
		t.Errorf("expected ChatID 'custom_chat_id', got '%s'", setting.ChatID)
	}
	if setting.ShowContent != false {
		t.Errorf("expected ShowContent false, got %v", setting.ShowContent)
	}
}

func TestSettingService_UpdateSetting(t *testing.T) {
	service, cleanup := newTestSettingService(t)
	defer cleanup()

	newSetting := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "new_chat_id",
		ShowContent:        true,
		SpoilerContent:     true,
		SendAttachments:    true,
		DisableLinkPreview: true,
	}

	err := service.UpdateSetting(1, newSetting)
	if err != nil {
		t.Fatalf("UpdateSetting failed: %v", err)
	}

	setting, err := service.GetSetting(1)
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}

	if setting.ChatID != "new_chat_id" {
		t.Errorf("expected ChatID 'new_chat_id', got '%s'", setting.ChatID)
	}
}

func TestSettingService_CreateDefaultSetting(t *testing.T) {
	service, cleanup := newTestSettingService(t)
	defer cleanup()

	err := service.CreateDefaultSetting(999)
	if err != nil {
		t.Fatalf("CreateDefaultSetting failed: %v", err)
	}

	setting, err := service.GetSetting(999)
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}

	if setting.UserID != 999 {
		t.Errorf("expected UserID 999, got %d", setting.UserID)
	}
	if setting.ChatID != DefaultChatID {
		t.Errorf("expected default ChatID, got '%s'", setting.ChatID)
	}
}
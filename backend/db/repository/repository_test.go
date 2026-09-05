package repository

import (
	"testing"

	"github.com/ydzydzydz/pmail_telegram_push/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

func newTestRepository(t *testing.T) (*Repository[model.PluginTelegramPushSettingModel], func()) {
	db, err := xorm.NewEngine("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	db.ShowSQL(false)
	err = db.Sync2(new(model.PluginTelegramPushSettingModel))
	if err != nil {
		t.Fatalf("failed to sync schema: %v", err)
	}
	repo := NewRepository[model.PluginTelegramPushSettingModel](db)
	return repo, func() {
		db.Close()
	}
}

func TestRepository_GetOrCreate_NotExists(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	defaultItem := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "default_chat_id",
		ShowContent:        true,
		SpoilerContent:     true,
		SendAttachments:    true,
		DisableLinkPreview: true,
	}

	result, err := repo.GetOrCreate(1, defaultItem)
	if err != nil {
		t.Fatalf("GetOrCreate failed: %v", err)
	}

	if result.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", result.UserID)
	}
	if result.ChatID != "default_chat_id" {
		t.Errorf("expected ChatID 'default_chat_id', got '%s'", result.ChatID)
	}
}

func TestRepository_GetOrCreate_AlreadyExists(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	existingItem := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "existing_chat_id",
		ShowContent:        false,
		SpoilerContent:     false,
		SendAttachments:    false,
		DisableLinkPreview: false,
	}

	err := repo.Create(existingItem)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	newItem := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "new_chat_id",
		ShowContent:        true,
		SpoilerContent:     true,
		SendAttachments:    true,
		DisableLinkPreview: true,
	}

	result, err := repo.GetOrCreate(1, newItem)
	if err != nil {
		t.Fatalf("GetOrCreate failed: %v", err)
	}

	if result.ChatID != "existing_chat_id" {
		t.Errorf("expected ChatID 'existing_chat_id', got '%s'", result.ChatID)
	}
}

func TestRepository_Create_and_FindOne(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	item := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "test_chat_id",
		ShowContent:        true,
		SpoilerContent:     false,
		SendAttachments:    true,
		DisableLinkPreview: false,
	}

	err := repo.Create(item)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	result, err := repo.FindOne(1)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}

	if result.ChatID != "test_chat_id" {
		t.Errorf("expected ChatID 'test_chat_id', got '%s'", result.ChatID)
	}
	if result.ShowContent != true {
		t.Errorf("expected ShowContent true, got %v", result.ShowContent)
	}
}

func TestRepository_Update(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	item := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "original_chat_id",
		ShowContent:        true,
		SpoilerContent:     true,
		SendAttachments:    true,
		DisableLinkPreview: true,
	}

	err := repo.Create(item)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	updatedItem := &model.PluginTelegramPushSettingModel{
		UserID:             1,
		ChatID:             "updated_chat_id",
		ShowContent:        false,
		SpoilerContent:     false,
		SendAttachments:    false,
		DisableLinkPreview: false,
	}

	err = repo.Update(1, updatedItem)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	result, err := repo.FindOne(1)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}

	if result.ChatID != "updated_chat_id" {
		t.Errorf("expected ChatID 'updated_chat_id', got '%s'", result.ChatID)
	}
}

func TestRepository_Exist(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	if repo.Exist(1) {
		t.Error("expected not exist before create")
	}

	item := &model.PluginTelegramPushSettingModel{
		UserID: 1,
		ChatID: "test",
	}

	err := repo.Create(item)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if !repo.Exist(1) {
		t.Error("expected exist after create")
	}
}

func TestRepository_FindOne_NotFound(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	_, err := repo.FindOne(999)
	if err == nil {
		t.Error("expected error for non-existent item")
	}
}

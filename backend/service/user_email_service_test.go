package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/Jinnrry/pmail/models"
	"github.com/ydzydzydz/pmail_telegram_push/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

func newTestUserEmailService(t *testing.T) (*UserEmailService, *xorm.Engine, func()) {
	db, err := xorm.NewEngine("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	db.ShowSQL(false)
	err = db.Sync2(new(models.UserEmail))
	if err != nil {
		t.Fatalf("failed to sync schema: %v", err)
	}
	userEmailService := NewUserEmailService(db)
	return userEmailService, db, func() {
		db.Close()
	}
}

func TestUserEmailService_GetUserEmailStatus_Exists(t *testing.T) {
	service, db, cleanup := newTestUserEmailService(t)
	defer cleanup()

	testCases := []struct {
		name       string
		status     int8
		wantStatus model.UserEmailStatus
	}{
		{"unsent_or_received", 0, model.StatusUnsentOrReceived},
		{"sent", 1, model.StatusSent},
		{"failed", 2, model.StatusFailed},
		{"deleted", 3, model.StatusDeleted},
		{"draft", 4, model.StatusDraft},
		{"junk", 5, model.StatusJunk},
	}

	for i, tc := range testCases {
		userID := i + 1
		record := &models.UserEmail{
			UserID:  userID,
			EmailID: 100,
			Status:  tc.status,
		}
		if _, err := db.Insert(record); err != nil {
			t.Fatalf("insert user_email failed: %v", err)
		}

		status, err := service.GetUserEmailStatus(userID, 100)
		if err != nil {
			t.Fatalf("GetUserEmailStatus(%d, %d) failed: %v", userID, 100, err)
		}
		if status != tc.wantStatus {
			t.Errorf("case %q: expected status %v, got %v", tc.name, tc.wantStatus, status)
		}
	}
}

func TestUserEmailService_GetUserEmailStatus_NotExists(t *testing.T) {
	service, _, cleanup := newTestUserEmailService(t)
	defer cleanup()

	status, err := service.GetUserEmailStatus(1, 999)
	if !errors.Is(err, ErrUserEmailNotFound) {
		t.Errorf("expected ErrUserEmailNotFound, got %v", err)
	}
	if status != model.StatusUnsentOrReceived {
		t.Errorf("expected fallback status %v, got %v", model.StatusUnsentOrReceived, status)
	}

	wantMsg := "user_id=1 email_id=999"
	if !strings.Contains(err.Error(), wantMsg) {
		t.Errorf("expected error message to contain %q, got %q", wantMsg, err.Error())
	}
}

func TestUserEmailService_GetUserEmailStatus_MatchesOnlyGivenUserAndEmail(t *testing.T) {
	service, db, cleanup := newTestUserEmailService(t)
	defer cleanup()

	record := &models.UserEmail{
		UserID:  1,
		EmailID: 100,
		Status:  int8(model.StatusSent),
	}
	if _, err := db.Insert(record); err != nil {
		t.Fatalf("insert user_email failed: %v", err)
	}

	// 不同 email_id 应查不到记录
	if _, err := service.GetUserEmailStatus(1, 200); !errors.Is(err, ErrUserEmailNotFound) {
		t.Errorf("expected ErrUserEmailNotFound for unmatched email_id, got %v", err)
	}

	// 不同 user_id 应查不到记录
	if _, err := service.GetUserEmailStatus(2, 100); !errors.Is(err, ErrUserEmailNotFound) {
		t.Errorf("expected ErrUserEmailNotFound for unmatched user_id, got %v", err)
	}

	// 匹配的 user_id + email_id 应返回正确状态
	status, err := service.GetUserEmailStatus(1, 100)
	if err != nil {
		t.Fatalf("GetUserEmailStatus failed: %v", err)
	}
	if status != model.StatusSent {
		t.Errorf("expected status %v, got %v", model.StatusSent, status)
	}
}

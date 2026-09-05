package service

import (
	"errors"
	"fmt"

	"github.com/Jinnrry/pmail/models"
	"github.com/ydzydzydz/pmail_telegram_push/model"
	"xorm.io/xorm"
)

var ErrUserEmailNotFound = errors.New("user_email not found")

type UserEmailService struct {
	db *xorm.Engine
}

func NewUserEmailService(db *xorm.Engine) *UserEmailService {
	return &UserEmailService{db: db}
}

// GetUserEmailStatus 查询 user_id + email_id 对应记录的 Status 字段值。
// 若记录不存在返回 ErrUserEmailNotFound，调用方必须显式判断并处理。
func (s *UserEmailService) GetUserEmailStatus(userID int, emailID int) (model.UserEmailStatus, error) {
	var userEmail models.UserEmail
	has, err := s.db.Where("user_id = ? AND email_id = ?", userID, emailID).Get(&userEmail)
	if err != nil {
		return model.StatusUnsentOrReceived, fmt.Errorf("query user_email status: %w", err)
	}
	if !has {
		return model.StatusUnsentOrReceived, fmt.Errorf("%w: user_id=%d email_id=%d",
			ErrUserEmailNotFound, userID, emailID)
	}
	return model.UserEmailStatus(userEmail.Status), nil
}

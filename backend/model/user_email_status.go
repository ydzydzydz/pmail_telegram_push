package model

import "fmt"

// UserEmailStatus 是 user_email.status 字段的强类型枚举。
// 注意：常量值与 github.com/Jinnrry/pmail/models.UserEmail.Status 保持一一对应，
// 上游版本升级后请同步核对。
type UserEmailStatus int8

const (
	StatusUnsentOrReceived UserEmailStatus = 0 // 未发送或收件
	StatusSent             UserEmailStatus = 1 // 已发送
	StatusFailed           UserEmailStatus = 2 // 发送失败
	StatusDeleted          UserEmailStatus = 3 // 删除
	StatusDraft            UserEmailStatus = 4 // 草稿箱
	StatusJunk             UserEmailStatus = 5 // 垃圾邮件
)

func (s UserEmailStatus) String() string {
	switch s {
	case StatusUnsentOrReceived:
		return "未发送或收件"
	case StatusSent:
		return "已发送"
	case StatusFailed:
		return "发送失败"
	case StatusDeleted:
		return "删除"
	case StatusDraft:
		return "草稿箱"
	case StatusJunk:
		return "垃圾邮件"
	default:
		return fmt.Sprintf("未知状态(%d)", s)
	}
}

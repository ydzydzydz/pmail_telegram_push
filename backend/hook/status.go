package hook

// UserEmailStatus 状态类型，用于表示邮件状态枚举
type UserEmailStatus int8

const (
	StatusUnsentOrReceived UserEmailStatus = 0 // 未发送或收件
	StatusSent             UserEmailStatus = 1 // 已发送
	StatusFailed           UserEmailStatus = 2 // 发送失败
	StatusDeleted          UserEmailStatus = 3 // 删除
	StatusDraft            UserEmailStatus = 4 // 草稿箱
	StatusJunk             UserEmailStatus = 5 // 整扰邮件(Junk)
)

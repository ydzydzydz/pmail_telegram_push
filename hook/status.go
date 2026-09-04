package hook

// Status 状态类型，用于表示邮件状态枚举
type Status int8

const (
	StatusUnsentOrReceived Status = 0 // 未发送或收件
	StatusSent             Status = 1 // 已发送
	StatusFailed           Status = 2 // 发送失败
	StatusDeleted          Status = 3 // 删除
	StatusDraft            Status = 4 // 草稿箱
	StatusJunk             Status = 5 // 整扰邮件(Junk)
)

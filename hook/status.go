package hook

type Status int8

const (
	StatusUnsentOrReceived = 0 // 未发送或收件
	StatusSent             = 1 // 已发送
	StatusFailed           = 2 // 发送失败
	StatusDeleted          = 3 // 删除
	StatusDraft            = 4 // 草稿箱
	StatusJunk             = 5 // 整扰邮件(Junk)
)

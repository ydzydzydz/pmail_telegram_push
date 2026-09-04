package hook

import (
	"fmt"

	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

// getSubjectText 获取主题文本
func getSubjectText(email *parsemail.Email) string {
	if len(email.Subject) <= 0 {
		return ""
	}
	return fmt.Sprintf("🔖 主题：<b>%s</b>\n", email.Subject)
}

// getFromText 获取发件人文本
func getFromText(email *parsemail.Email) string {
	if len(email.From.EmailAddress) <= 0 {
		return ""
	}
	return fmt.Sprintf("📤 发件：&#60;%s&#62;\n", email.From.EmailAddress)
}

// getToText 获取收件人文本
func getToText(email *parsemail.Email) string {
	if len(email.To) <= 0 {
		return ""
	}
	text := "📥 收件："
	for _, to := range email.To {
		text += fmt.Sprintf("&#60;%s&#62; ", to.EmailAddress)
	}
	text += "\n"
	return text
}

// getCcText 获取抄送人文本
func getCcText(email *parsemail.Email) string {
	if len(email.Cc) <= 0 {
		return ""
	}
	text := "📋 抄送："
	for _, cc := range email.Cc {
		text += fmt.Sprintf("&#60;%s&#62; ", cc.EmailAddress)
	}
	text += "\n"
	return text
}

// getBccText 获取密送人文本
func getBccText(email *parsemail.Email) string {
	if len(email.Bcc) <= 0 {
		return ""
	}
	text := "🕵️ 密送："
	for _, bcc := range email.Bcc {
		text += fmt.Sprintf("&#60;%s&#62; ", bcc.EmailAddress)
	}
	text += "\n"
	return text
}

// getAttachmentsText 获取附件文本
func getAttachmentsText(email *parsemail.Email) string {
	if len(email.Attachments) <= 0 {
		return ""
	}
	return fmt.Sprintf("📎 附件：%d 个\n", len(email.Attachments))
}

// getContentText 获取邮件内容文本
func getContentText(email *parsemail.Email, setting *model.PluginTelegramPushSettingModel) string {
	if !setting.ShowContent {
		return ""
	}
	if len(email.Text) > 0 {
		return string(email.Text)
	}
	if len(email.HTML) > 0 {
		return removeHTMLTag(string(email.HTML))
	}
	return ""
}

// getSpoilerText 获取spoiler文本
func getSpoilerText(text string, setting *model.PluginTelegramPushSettingModel) string {
	if !setting.SpoilerContent {
		return text
	}
	return fmt.Sprintf("<tg-spoiler>%s</tg-spoiler>", text)
}

package hook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

// TelegramTextMaxSize Telegram 文本最大长度
const TELEGRAM_TEXT_MAX_SIZE = 4096

// getSubjectText 获取主题文本
func (h *PmailTelegramPushHook) getSubjectText(email *parsemail.Email) string {
	if len(email.Subject) <= 0 {
		return ""
	}
	return fmt.Sprintf("🔖 主题：<b>%s</b>\n", email.Subject)
}

// getFromText 获取发件人文本
func (h *PmailTelegramPushHook) getFromText(email *parsemail.Email) string {
	if len(email.From.EmailAddress) <= 0 {
		return ""
	}
	return fmt.Sprintf("📤 发件：&#60;%s&#62;\n", email.From.EmailAddress)
}

// getToText 获取收件人文本
func (h *PmailTelegramPushHook) getToText(email *parsemail.Email) string {
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
func (h *PmailTelegramPushHook) getCcText(email *parsemail.Email) string {
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
func (h *PmailTelegramPushHook) getBccText(email *parsemail.Email) string {
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
func (h *PmailTelegramPushHook) getAttachmentsText(email *parsemail.Email) string {
	if len(email.Attachments) <= 0 {
		return ""
	}
	return fmt.Sprintf("📎 附件：%d 个\n", len(email.Attachments))
}

// getContentText 获取邮件内容文本
func (h *PmailTelegramPushHook) getContentText(email *parsemail.Email, setting *model.TelegramPushSetting) string {
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
func (h *PmailTelegramPushHook) getSpoilerText(text string, setting *model.TelegramPushSetting) string {
	if !setting.SpoilerContent {
		return text
	}
	return fmt.Sprintf("<tg-spoiler>%s</tg-spoiler>", text)
}

// buildSendText 构建发送文本
func (h *PmailTelegramPushHook) buildSendText(email *parsemail.Email, setting *model.TelegramPushSetting) string {
	text := "📧 有新邮件\n"
	text += h.getSubjectText(email)
	text += h.getFromText(email)
	text += h.getToText(email)
	text += h.getCcText(email)
	text += h.getBccText(email)
	text += h.getAttachmentsText(email)
	text += h.getSpoilerText(h.getContentText(email, setting), setting)
	text = removeExtraSpace(text)

	// 预留 20 个字符
	maxSizeWithPadding := TELEGRAM_TEXT_MAX_SIZE - 20
	if len(text) > maxSizeWithPadding {
		// 如果在预留长度内没有 spoiler 起始标签，直接截取
		if !strings.Contains(text[:maxSizeWithPadding], "<tg-spoiler>") {
			return text[:maxSizeWithPadding] + "..."
		}
		// 如果在预留长度内有 spoiler 起始标签，末尾添加结束标签
		return text[:maxSizeWithPadding] + "..." + "</tg-spoiler>"
	}
	return text
}

// buildPamilLinkButton 创建Pamil链接按钮
func (h *PmailTelegramPushHook) buildPamilLinkButton() *models.InlineKeyboardMarkup {
	var url string
	if h.mainConfig.HttpsEnabled > 1 {
		url = "http://" + h.mainConfig.WebDomain
	} else {
		url = "https://" + h.mainConfig.WebDomain
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "查收邮件",
					URL:  url,
				},
			},
		},
	}
}

// sendText 发送文本消息
func (h *PmailTelegramPushHook) sendText(email *parsemail.Email, setting *model.TelegramPushSetting) (msg *models.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.pluginConfig.Timeout)*time.Second)
	defer cancel()

	parmas := &bot.SendMessageParams{
		ChatID:      setting.ChatID,
		Text:        h.buildSendText(email, setting),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: h.buildPamilLinkButton(),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &setting.DisableLinkPreview,
		},
	}

	return h.bot.SendMessage(ctx, parmas)
}

// sendAttachments 发送附件消息
func (h *PmailTelegramPushHook) sendAttachments(id int, email *parsemail.Email, setting *model.TelegramPushSetting) (errs error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.pluginConfig.Timeout)*time.Second)
	defer cancel()

	// 引用消息中包含附件关键字
	params := &bot.SendDocumentParams{
		ChatID: setting.ChatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: id,
			Quote:     fmt.Sprintf("📎 附件：%d 个", len(email.Attachments)),
		},
	}

	// 逐个发送附件
	for i, attachment := range email.Attachments {
		params.Caption = fmt.Sprintf("📎 附件 %d", i+1)
		params.Document = &models.InputFileUpload{
			Filename: filepath.Base(attachment.Filename),
			Data:     bytes.NewReader(attachment.Content),
		}
		// 发送附件失败，记录错误，继续发送下一个附件
		if _, err := h.bot.SendDocument(ctx, params); err != nil {
			errs = errors.Join(err, fmt.Errorf("send document failed, err: %w", err))
			continue
		}
	}
	return
}

// sendNotification 发送通知消息
// 先发送文本消息，再发送附件消息
func (h *PmailTelegramPushHook) sendNotification(email *parsemail.Email, setting *model.TelegramPushSetting) (err error) {
	msg, err := h.sendText(email, setting)
	if err != nil {
		return err
	}
	return h.sendAttachments(msg.ID, email, setting)
}

// TODO: 合并多个附件为一个消息发送
// func (h *PmailTelegramPushHook) sendAttachmentsCombine(id int, email *parsemail.Email) (msg []*models.Message, err error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.pluginConfig.Timeout)*time.Second)
// 	defer cancel()
// 	params := &bot.SendMediaGroupParams{
// 		ChatID: h.pluginConfig.TelegramChatID,
// 		ReplyParameters: &models.ReplyParameters{
// 			MessageID: id,
// 			Quote:     fmt.Sprintf("📎 附件：%d 个", len(email.Attachments)),
// 		},
// 	}
// 	for i, attachment := range email.Attachments {
// 		params.Media = append(params.Media, &models.InputMediaDocument{
// 			Media:           filepath.Base(attachment.Filename),
// 			Caption:         fmt.Sprintf("📎 附件 %d", i+1),
// 			MediaAttachment: bytes.NewReader(attachment.Content),
// 		})
// 	}
// 	return h.bot.SendMediaGroup(ctx, params)
// }

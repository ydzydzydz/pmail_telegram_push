package sender

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

const (
	TELEGRAM_TEXT_MAX_SIZE        = 4096             // TELEGRAM_TEXT_MAX_SIZE Telegram 文本最大长度
	TELEGRAM_ATTACHMENT_MAX_COUNT = 10               // TELEGRAM_ATTACHMENT_MAX_COUNT Telegram 附件最大数量
	TELEGRAM_ATTACHMENT_MAX_SIZE  = 20 * 1024 * 1024 // TELEGRAM_ATTACHMENT_MAX_SIZE Telegram 附件最大大小
)

// TelegramBotSender Telegram机器人发送器
type TelegramBotSender struct {
	bot              *bot.Bot
	pmailWebSiteLink string
}

// NewTelegramBotSender 创建Telegram机器人发送器
func NewTelegramBotSender(cfg *config.Config) (*TelegramBotSender, error) {
	bot, err := NewBot(cfg)
	if err != nil {
		return nil, err
	}
	var link string
	if cfg.MainConfig.HttpsEnabled > 1 {
		link = "http://" + cfg.MainConfig.WebDomain
	} else {
		link = "https://" + cfg.MainConfig.WebDomain
	}
	return &TelegramBotSender{
		bot:              bot,
		pmailWebSiteLink: link,
	}, nil
}

// SendNotification 发送通知消息
func (s *TelegramBotSender) SendNotification(ctx context.Context, setting *model.PluginTelegramPushSettingModel, email *parsemail.Email) error {
	msg, err := s.sendMessage(ctx, setting, email)
	if err != nil {
		return err
	}
	return s.sendAttachmentsBatch(ctx, msg.ID, email, setting)
}

// SendTestMessage 发送测试消息
func (s *TelegramBotSender) SendTestMessage(ctx context.Context, setting *model.PluginTelegramPushSettingModel) error {
	params := &bot.SendMessageParams{
		ChatID:    setting.ChatID,
		Text:      "这是一条测试消息",
		ParseMode: models.ParseModeHTML,
	}
	_, err := s.bot.SendMessage(ctx, params)
	return err
}

// GetBot 获取Telegram机器人
func (s *TelegramBotSender) GetBot() *bot.Bot {
	return s.bot
}

// buildSendText 构建发送文本
func buildSendText(email *parsemail.Email, setting *model.PluginTelegramPushSettingModel) string {
	text := "📧 有新邮件\n"
	text += getSubjectText(email)
	text += getFromText(email)
	text += getToText(email)
	text += getCcText(email)
	text += getBccText(email)
	text += getAttachmentsText(email)
	text += getSpoilerText(getContentText(email, setting), setting)
	text = removeExtraSpace(text)

	// 预留 20 个字符
	maxSizeWithPadding := TELEGRAM_TEXT_MAX_SIZE - 20
	if len(text) > maxSizeWithPadding {
		truncated := text[:maxSizeWithPadding]
		// 检查截断位置是否在 <tg-spoiler> 标签内部
		lastOpen := strings.LastIndex(truncated, "<tg-spoiler>")
		if lastOpen != -1 {
			afterOpen := truncated[lastOpen:]
			if !strings.Contains(afterOpen, "</tg-spoiler>") {
				// 截断位置在标签内部，找到前一个完整的 spoiler 块
				beforeOpen := truncated[:lastOpen]
				if strings.Contains(beforeOpen, "</tg-spoiler>") {
					// 找到前一个完整 spoiler 的结束位置
					prevClose := strings.LastIndex(beforeOpen, "</tg-spoiler>")
					truncated = beforeOpen[:prevClose+len("</tg-spoiler>")]
				} else {
					// 没有完整的 spoiler，直接截断
					truncated = beforeOpen
				}
				return truncated + "..."
			}
		}
		return truncated + "..."
	}
	return text
}

// buildPmailLinkButton 创建Pmail链接按钮
func buildPmailLinkButton(pmailWebSiteLink string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "查收邮件",
					URL:  pmailWebSiteLink,
				},
			},
		},
	}
}

// sendMessage 发送文本消息
func (s *TelegramBotSender) sendMessage(ctx context.Context, setting *model.PluginTelegramPushSettingModel, email *parsemail.Email) (msg *models.Message, err error) {
	params := &bot.SendMessageParams{
		ChatID:      setting.ChatID,
		Text:        buildSendText(email, setting),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: buildPmailLinkButton(s.pmailWebSiteLink),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &setting.DisableLinkPreview,
		},
	}
	return s.bot.SendMessage(ctx, params)
}

// sendAttachmentsBatch 批量发送附件消息
func (s *TelegramBotSender) sendAttachmentsBatch(ctx context.Context, id int, email *parsemail.Email, setting *model.PluginTelegramPushSettingModel) (err error) {

	// 引用消息中包含附件关键字
	params := &bot.SendMediaGroupParams{
		ChatID: setting.ChatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: id,
			Quote:     fmt.Sprintf("📎 附件：%d 个", len(email.Attachments)),
		},
	}

	// 批量发送附件, 每个批次最多 TELEGRAM_ATTACHMENT_MAX_COUNT 个附件
	for i := 0; i < len(email.Attachments); i += TELEGRAM_ATTACHMENT_MAX_COUNT {
		end := min(len(email.Attachments), i+TELEGRAM_ATTACHMENT_MAX_COUNT)
		params.Media = nil
		batch := email.Attachments[i:end]
		for j, attachment := range batch {
			// 判断文件大小, 超过最大大小或为空则跳过
			if len(attachment.Content) == 0 || len(attachment.Content) > TELEGRAM_ATTACHMENT_MAX_SIZE {
				continue
			}
			// 构建 InputMediaDocument
			params.Media = append(params.Media, &models.InputMediaDocument{
				Media:           fmt.Sprintf("attach://%s", filepath.Base(attachment.Filename)),
				Caption:         fmt.Sprintf("📎 附件 %d", i+j+1),
				MediaAttachment: bytes.NewReader(attachment.Content),
			})
		}
		if len(params.Media) == 0 {
			continue
		}
		if _, err = s.bot.SendMediaGroup(ctx, params); err != nil {
			return err
		}
		// 每个批次发送后休息 1 秒, 避免触发速率限制
		time.Sleep(time.Second)
	}

	return nil
}

package hook

import (
	"strings"

	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/logger"
	"github.com/ydzydzydz/pmail_telegram_push/service"

	pconfig "github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/Jinnrry/pmail/models"
	"github.com/Jinnrry/pmail/utils/context"
	"github.com/go-telegram/bot"

	"github.com/ydzydzydz/pmail_telegram_push/db"
)

const (
	PLUGIN_NAME = "pmail_telegram_push" // 插件名称
)

// PmailTelegramPushHook 插件钩子
type PmailTelegramPushHook struct {
	bot            *bot.Bot
	mainConfig     *pconfig.Config
	pluginConfig   *config.PluginConfig
	settingService *service.SettingService
}

// NewPmailTelegramPushHook 创建插件钩子
func NewPmailTelegramPushHook(cfg *config.Config) *PmailTelegramPushHook {
	bot, err := NewBot(cfg)
	if err != nil {
		logger.PluginLogger.Fatal().Err(err).Msg("创建bot失败")
	}
	logger.PluginLogger.Info().Msg("bot初始化成功")

	dataSource, err := db.NewDataSource(cfg)
	if err != nil {
		logger.PluginLogger.Fatal().Err(err).Msg("创建数据库连接失败")
	}
	logger.PluginLogger.Info().Msg("数据库初始化成功")

	settingService := service.NewSettingService(dataSource.SettingDao())
	return &PmailTelegramPushHook{
		bot:            bot,
		mainConfig:     cfg.MainConfig,
		pluginConfig:   cfg.PluginConfig,
		settingService: settingService,
	}
}

var _ framework.EmailHook = (*PmailTelegramPushHook)(nil)

// GetName 获取插件名称
func (h *PmailTelegramPushHook) GetName(ctx *context.Context) string {
	return PLUGIN_NAME
}

// ReceiveSaveAfter 接收保存后的钩子
func (h *PmailTelegramPushHook) ReceiveSaveAfter(ctx *context.Context, email *parsemail.Email, ue []*models.UserEmail) {
	for _, u := range ue {
		// 已读邮件不处理
		if u.IsRead != 0 {
			continue
		}
		// 未读邮件不处理
		if u.Status != 0 {
			continue
		}
		// 邮件ID不存在不处理
		if email.MessageId <= 0 {
			continue
		}

		setting, err := h.settingService.GetSetting(u.UserID)
		if err != nil {
			// 获取设置失败不处理
			continue
		}
		if setting.ChatID == "" {
			// 聊天ID不存在不处理
			continue
		}

		if err = h.sendNotification(email, setting); err != nil {
			logger.PluginLogger.Error().Err(err).Int64("email_message_id", email.MessageId).Msg("发送通知失败")
			continue
		}
		logger.PluginLogger.Info().Int64("email_message_id", email.MessageId).Msg("发送通知成功")
	}
}

// ReceiveParseBefore 接收解析前的钩子
func (h *PmailTelegramPushHook) ReceiveParseBefore(ctx *context.Context, email *[]byte) {

}

// ReceiveParseAfter 接收解析后的钩子
func (h *PmailTelegramPushHook) ReceiveParseAfter(ctx *context.Context, email *parsemail.Email) {

}

// SendAfter 发送后的钩子
func (h *PmailTelegramPushHook) SendAfter(ctx *context.Context, email *parsemail.Email, err map[string]error) {

}

// SendBefore 发送前的钩子
func (h *PmailTelegramPushHook) SendBefore(ctx *context.Context, email *parsemail.Email) {

}

// SettingsHtml 获取设置 HTML
func (h *PmailTelegramPushHook) SettingsHtml(ctx *context.Context, url string, requestData string) string {
	switch {
	// 获取用户设置
	case strings.Contains(url, "getSetting"):
		return h.getSetting(ctx.UserID)
	// 获取机器人信息
	case strings.Contains(url, "getBotInfo"):
		return h.getBotInfo()
	// 更新设置
	case strings.Contains(url, "updateSetting"):
		return h.updateSetting(ctx.UserID, requestData)
	default:
		return SettingHtml
	}
}

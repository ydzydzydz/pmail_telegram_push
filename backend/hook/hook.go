package hook

import (
	"strings"
	"time"

	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/controller"
	"github.com/ydzydzydz/pmail_telegram_push/db"
	"github.com/ydzydzydz/pmail_telegram_push/logger"
	"github.com/ydzydzydz/pmail_telegram_push/model"
	"github.com/ydzydzydz/pmail_telegram_push/sender"
	"github.com/ydzydzydz/pmail_telegram_push/service"

	_ "embed"

	ccontext "context"

	pconfig "github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/Jinnrry/pmail/models"
	"github.com/Jinnrry/pmail/utils/context"
)

const (
	PLUGIN_NAME = "pmail_telegram_push" // 插件名称
)

// PmailTelegramPushHook 插件钩子
type PmailTelegramPushHook struct {
	mainConfig        *pconfig.Config
	pluginConfig      *config.PluginConfig
	settingService    *service.SettingService
	userEmailService  *service.UserEmailService
	settingController *controller.SettingController
	botController     *controller.BotController
	sender            *sender.TelegramBotSender
}

var _ framework.EmailHook = (*PmailTelegramPushHook)(nil)

// NewPmailTelegramPushHook 创建插件钩子实例
func NewPmailTelegramPushHook(cfg *config.Config) *PmailTelegramPushHook {
	dataSource, err := db.NewDataSource(cfg)
	if err != nil {
		logger.PluginLogger.Fatal().Err(err).Msg("创建数据库连接失败")
	}
	logger.PluginLogger.Info().Msg("数据库初始化成功")
	settingService := service.NewSettingService(dataSource.DB())
	userEmailService := service.NewUserEmailService(dataSource.DB())

	sender, err := sender.NewTelegramBotSender(cfg)
	if err != nil {
		logger.PluginLogger.Fatal().Err(err).Msg("创建Telegram Bot发送器失败")
	}
	logger.PluginLogger.Info().Msg("Telegram Bot发送器初始化成功")

	return &PmailTelegramPushHook{
		mainConfig:        cfg.MainConfig,
		pluginConfig:      cfg.PluginConfig,
		settingService:    settingService,
		userEmailService:  userEmailService,
		settingController: controller.NewSettingController(settingService),
		botController:     controller.NewBotController(sender),
		sender:            sender,
	}
}

// GetName 获取插件名称
func (h *PmailTelegramPushHook) GetName(ctx *context.Context) string {
	return PLUGIN_NAME
}

// ReceiveSaveAfter 接收保存后的钩子
func (h *PmailTelegramPushHook) ReceiveSaveAfter(ctx *context.Context, email *parsemail.Email, userEmails []*models.UserEmail) {
	for _, userEmail := range userEmails {
		// 已读邮件不处理
		if userEmail.IsRead != 0 {
			continue
		}
		// 邮件ID不存在不处理
		if email.MessageId <= 0 {
			continue
		}

		setting, err := h.settingService.GetSetting(userEmail.UserID)
		if err != nil {
			logger.PluginLogger.Error().Err(err).Int("user_id", userEmail.UserID).Msg("获取用户设置失败")
			continue
		}
		// 聊天ID不存在不处理
		if setting.ChatID == "" {
			continue
		}

		// 未发送或收件邮件不处理
		// if model.UserEmailStatus(userEmail.Status) != model.StatusUnsentOrReceived {
		// 	continue
		// }

		// fix: https://github.com/Jinnrry/PMail/discussions/357
		// ? 不确定是否有用，无法确定查询时状态一定更新到数据库
		// ! 根本解决需要框架层面支持
		// 从数据库查询用户邮件状态
		status, err := h.userEmailService.GetUserEmailStatus(userEmail.UserID, userEmail.EmailID)
		if err != nil {
			logger.PluginLogger.Error().Err(err).Int64("email_message_id", email.MessageId).Msg("获取用户邮件状态失败")
			continue
		}
		logger.PluginLogger.Info().Int("user_id", userEmail.UserID).Int("email_id", userEmail.EmailID).Str("status", status.String()).Msg("用户邮件状态")

		if status != model.StatusUnsentOrReceived {
			continue
		}

		// 发送通知
		cctx, cancel := ccontext.WithTimeout(ccontext.Background(), time.Duration(h.pluginConfig.Timeout)*time.Second)
		defer cancel()
		if err = h.sender.SendNotification(cctx, setting, email); err != nil {
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

var (
	//go:embed dist/index.html
	SettingHtml string // 设置页面
)

// SettingsHtml 获取设置 HTML
func (h *PmailTelegramPushHook) SettingsHtml(ctx *context.Context, url string, requestData string) string {
	switch {
	case strings.HasSuffix(url, "getSetting"):
		return h.settingController.GetSetting(ctx.UserID)
	case strings.HasSuffix(url, "getBotInfo"):
		return h.botController.GetBotInfo()
	case strings.HasSuffix(url, "updateSetting"):
		return h.settingController.UpdateSetting(ctx.UserID, requestData)
	case strings.HasSuffix(url, "testMessage"):
		return h.botController.SendTestMessage(ctx.UserID, requestData)
	default:
		return SettingHtml
	}
}

package hook

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/logger"
	"golang.org/x/net/proxy"
)

// NewBot 创建一个新的 Telegram 机器人实例
func NewBot(config *config.Config) (*bot.Bot, error) {
	// 初始化 Telegram 机器人选项
	// 检查初始化超时时间
	opts := []bot.Option{
		bot.WithCheckInitTimeout(time.Duration(config.PluginConfig.Timeout) * time.Second),
	}

	// 开启调试模式
	// 开启调试模式后，会打印所有的 API 请求和响应
	if config.PluginConfig.Debug {
		opts = append(opts,
			bot.WithDebug(),
			// 自定义调试处理函数，将调试信息打印到日志中
			bot.WithDebugHandler(func(format string, args ...any) {
				logger.BotLogger.Debug().Msgf(format, args...)
			}),
		)
	}

	if config.PluginConfig.Proxy == "" {
		return newBotWithOutProxy(config, opts...)
	}
	parsedURL, err := url.Parse(config.PluginConfig.Proxy)
	if err != nil {
		logger.PluginLogger.Panic().Err(err).Msg("代理URL解析失败")
	}

	// 根据代理协议类型创建不同的 Telegram 机器人实例
	switch strings.ToLower(parsedURL.Scheme) {
	case "socks5":
		return newBotWithSocks5Proxy(config, parsedURL, opts...)
	case "http", "https":
		return newBotWithHTTPProxy(config, parsedURL, opts...)
	default:
		return newBotWithOutProxy(config, opts...)
	}
}

// newBotWithOutProxy 创建一个没有代理的 Telegram 机器人实例
func newBotWithOutProxy(config *config.Config, options ...bot.Option) (*bot.Bot, error) {
	return bot.New(config.PluginConfig.TelegramBotToken, options...)
}

// newBotWithSocks5Proxy 创建一个 SOCKS5 代理的 Telegram 机器人实例
func newBotWithSocks5Proxy(config *config.Config, proxyURL *url.URL, options ...bot.Option) (*bot.Bot, error) {
	// 解析 SOCKS5 代理 URL
	var auth *proxy.Auth
	if proxyURL.User != nil {
		password, _ := proxyURL.User.Password()
		auth = &proxy.Auth{
			User:     proxyURL.User.Username(),
			Password: password,
		}
	}
	dialer, err := proxy.SOCKS5(
		"tcp",
		proxyURL.Host,
		auth,
		proxy.Direct,
	)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
	}
	opts := append(options, bot.WithHTTPClient(time.Duration(config.PluginConfig.Timeout)*time.Second, httpClient))
	return bot.New(config.PluginConfig.TelegramBotToken, opts...)
}

// newBotWithHTTPProxy 创建一个 HTTP 代理的 Telegram 机器人实例
func newBotWithHTTPProxy(config *config.Config, proxyURL *url.URL, options ...bot.Option) (*bot.Bot, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	options = append(options, bot.WithHTTPClient(time.Duration(config.PluginConfig.Timeout)*time.Second, httpClient))
	return bot.New(config.PluginConfig.TelegramBotToken, options...)
}

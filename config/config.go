package config

import (
	"encoding/json"
	"os"

	pconfig "github.com/Jinnrry/pmail/config"
	"github.com/ydzydzydz/pmail_telegram_push/logger"
)

const (
	MAIN_CONFIG_FILE   = "./config/config.json"
	PLUGIN_CONFIG_FILE = "./config/pmail_telegram_push_config.json"
)

// PluginConfig 插件配置
type PluginConfig struct {
	TelegramBotToken string `json:"telegram_bot_token"`    // Telegram机器人Token
	Proxy            string `json:"proxy" default:""`      // 代理地址
	Timeout          int    `json:"timeout" default:"30"`  // 超时时间
	Debug            bool   `json:"debug" default:"false"` // 调试模式
}

// Config 配置
type Config struct {
	PluginConfig *PluginConfig
	MainConfig   *pconfig.Config
}

// readMainConfig 读取主配置文件
func readMainConfig() *pconfig.Config {
	content, err := os.ReadFile(MAIN_CONFIG_FILE)
	if err != nil {
		logger.PluginLogger.Panic().Err(err).Msg("主配置文件读取失败")
	}
	var mainConfig pconfig.Config
	if err := json.Unmarshal(content, &mainConfig); err != nil {
		logger.PluginLogger.Panic().Err(err).Msg("主配置文件解析失败")
	}
	logger.PluginLogger.Info().Msg("主配置文件读取成功")
	return &mainConfig
}

// readPluginConfig 读取插件配置文件
func readPluginConfig() *PluginConfig {
	content, err := os.ReadFile(PLUGIN_CONFIG_FILE)
	if err != nil {
		logger.PluginLogger.Panic().Err(err).Msg("插件配置文件读取失败")
	}
	var pluginConfig PluginConfig
	if err := json.Unmarshal(content, &pluginConfig); err != nil {
		logger.PluginLogger.Panic().Err(err).Msg("插件配置文件解析失败")
	}
	if pluginConfig.TelegramBotToken == "" {
		logger.PluginLogger.Panic().Msg("telegram bot token is empty")
	}
	logger.PluginLogger.Info().Msg("插件配置文件读取成功")
	return &pluginConfig
}

// ReadConfig 读取配置
func ReadConfig() *Config {
	return &Config{
		PluginConfig: readPluginConfig(),
		MainConfig:   readMainConfig(),
	}
}

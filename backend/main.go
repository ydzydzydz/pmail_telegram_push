package main

import (
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/hook"
)

func main() {
	// 读取配置
	cfg := config.ReadConfig()
	// 启动插件
	framework.CreatePlugin(
		hook.PLUGIN_NAME,
		hook.NewPmailTelegramPushHook(cfg),
	).Run()
}

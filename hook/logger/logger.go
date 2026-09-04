package logger

import "github.com/phuslu/log"

// BotLogger 机器人日志
var BotLogger = log.Logger{
	Level:  log.DebugLevel,
	Caller: 1,
	Writer: &log.ConsoleWriter{
		ColorOutput:    true,
		EndWithMessage: true,
	},
}

// PluginLogger 插件日志
var PluginLogger = log.Logger{
	Level:  log.InfoLevel, // 信息日志级别
	Caller: 1,
	Writer: &log.ConsoleWriter{
		ColorOutput:    true,
		EndWithMessage: true,
	},
}

package xhlog

import (
	"fmt"
	"strings"
	"time"
)

// LoggerConf 日志的配置信息
type LoggerConf struct {
	Dir        string
	Prefix     string
	Level      string
	Console    bool
	RotateSize int64
}

// 日志打印模块，实现功能如下：
// 日志输出支持日志和终端输出，可选
// 1. 终端输出，不做切割
// 2. 文件输出，支持按时间切割处理，支持错误日志单独输出操作

type Logger struct {
	fileLogger *fileLogger
	level      Level
	Console    bool
}

var LoggerExp *Logger

// Init 初始化日志打印的句柄
func Init(conf *LoggerConf) error {
	f, err := newFileLogger(conf)
	if err != nil {
		return err
	}
	LoggerExp = &Logger{
		fileLogger: f,
		level:      fromLevelName(conf.Level),
		Console:    conf.Console,
	}
	return nil
}

// Logger.Close 日志关闭
func (l *Logger) Close() {
	l.fileLogger.close()
	return
}

func (l *Logger) Print(v ...interface{}) {
	// 暂不提供实现，空函数
}

func (l *Logger) Println(v ...interface{}) {
	// 暂不提供实现，空函数
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level >= DebugLevel {
		l.print(DebugLevel, v...)
	}
	if l.Console {
		console(DebugLevel, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level >= InfoLevel {
		l.print(InfoLevel, v...)
	}
	if l.Console {
		console(InfoLevel, v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level >= WarnLevel {
		l.print(WarnLevel, v...)
	}
	if l.Console {
		console(WarnLevel, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.level >= ErrorLevel {
		l.print(ErrorLevel, v...)
	}
	if l.Console {
		console(ErrorLevel, v...)
	}
}

func (l *Logger) print(printLevel Level, v ...interface{}) {
	strTime := time.Now().Format(TimeFormat)
	logStr := fmt.Sprintf("%s %s %s\n", levels[printLevel].RawText, strTime, fmt.Sprint(v...))
	l.fileLogger.logChan <- logStr
	return
}

type Level uint32

// The available built'n log levels, users can add or modify a level via `Levels` field.
const (
	// DisableLevel will disable the printer.
	DisableLevel Level = iota
	// ErrorLevel will print only errors.
	ErrorLevel
	// WarnLevel will print errors and warnings.
	WarnLevel
	// InfoLevel will print errors, warnings and infos.
	InfoLevel
	// DebugLevel will print on any level, fatals, errors, warnings, infos and debug logs.
	DebugLevel
)

// Levels contains the levels and their
// mapped (pointer of, in order to be able to be modified) metadata, callers
// are allowed to modify this package-level global variable
// without any loses.
var levels = map[Level]*LevelMetadata{
	DisableLevel: {
		Name:    "disable",
		RawText: "",
	},
	ErrorLevel: {
		Name:    "error",
		RawText: "[ERROR]",
	},
	WarnLevel: {
		Name:    "warn",
		RawText: "[WARN]",
	},
	InfoLevel: {
		Name:    "info",
		RawText: "[INFO]",
	},
	DebugLevel: {
		Name:    "debug",
		RawText: "[DEBUG]",
	},
}

type LevelMetadata struct {
	Name    string
	RawText string
}

// fromLevelName 将string映射为level
func fromLevelName(name string) Level {
	for level, meta := range levels {
		if strings.ToLower(meta.Name) == strings.ToLower(name) {
			return level
		}
	}
	return DisableLevel
}

const (
	// 日志时间格式
	TimeFormat = "[2006-01-02T15:04:05.000Z0700]"
)

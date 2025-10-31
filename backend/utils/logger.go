package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var (
	// 当前日志级别（可通过环境变量配置）
	currentLogLevel LogLevel = LogLevelInfo
	logFile          *os.File
)

// 初始化日志系统
func init() {
	// 从环境变量读取日志级别
	logLevelStr := os.Getenv("LOG_LEVEL")
	switch logLevelStr {
	case "DEBUG":
		currentLogLevel = LogLevelDebug
	case "WARN":
		currentLogLevel = LogLevelWarn
	case "ERROR":
		currentLogLevel = LogLevelError
	default:
		currentLogLevel = LogLevelInfo
	}

	// 可选：将日志输出到文件（生产环境建议）
	logFilePath := os.Getenv("LOG_FILE")
	if logFilePath != "" {
		var err error
		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("警告: 无法打开日志文件 %s: %v", logFilePath, err)
		}
	}
}

// LogDebug 调试日志
func LogDebug(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelDebug {
		logMessage("DEBUG", format, v...)
	}
}

// LogInfo 信息日志
func LogInfo(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelInfo {
		logMessage("INFO", format, v...)
	}
}

// LogWarn 警告日志
func LogWarn(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelWarn {
		logMessage("WARN", format, v...)
	}
}

// LogError 错误日志
func LogError(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelError {
		logMessage("ERROR", format, v...)
	}
}

// logMessage 统一的日志输出函数
func logMessage(level, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := format
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	}

	logMsg := "[" + timestamp + "] [" + level + "] " + message

	// 输出到标准输出
	log.Println(logMsg)

	// 如果配置了日志文件，同时写入文件
	if logFile != nil {
		logFile.WriteString(logMsg + "\n")
	}
}

// LogSlowQuery 慢查询日志（单独函数以便后续统计）
func LogSlowQuery(query string, duration time.Duration, threshold time.Duration) {
	if duration > threshold {
		LogWarn("慢查询检测: 耗时 %v (阈值: %v), SQL: %s", duration, threshold, query)
	}
}

// CloseLogFile 关闭日志文件（在程序退出时调用）
func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}


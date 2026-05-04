package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var appLogger *log.Logger

func init() {
	appLogger = log.New(os.Stdout, "", 0)
}

// caller 获取调用方的文件名和行号
// skip=2: formatLog -> LogXxx -> 实际调用处
func caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// formatLog 格式化日志输出
func formatLog(level string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	location := caller(3) // 跳过: formatLog -> LogXxx -> 实际调用处
	msg := fmt.Sprintf(format, args...)
	appLogger.Printf("[%s] [%s] %s %s", timestamp, level, location, msg)
}

// LogInfo 输出 Info 级别日志
func LogInfo(format string, args ...interface{}) {
	formatLog("INFO", format, args...)
}

// LogWarn 输出 Warn 级别日志
func LogWarn(format string, args ...interface{}) {
	formatLog("WARN", format, args...)
}

// LogError 输出 Error 级别日志
func LogError(format string, args ...interface{}) {
	formatLog("ERROR", format, args...)
}

// LogDebug 输出 Debug 级别日志
func LogDebug(format string, args ...interface{}) {
	formatLog("DEBUG", format, args...)
}

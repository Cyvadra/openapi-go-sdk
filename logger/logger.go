// Package logger 提供日志接口和默认实现。
//
// 定义 Logger 接口，支持注入自定义 logger。
// 默认实现使用标准库 log 包输出到 stderr。
package logger

import (
	"fmt"
	"log"
	"os"
)

// Level 日志级别
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	SetLevel(level Level)
}

// DefaultLogger 默认日志实现，使用标准库 log 包
type DefaultLogger struct {
	level  Level
	logger *log.Logger
}

// NewDefaultLogger 创建默认日志实例
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		level:  LevelInfo,
		logger: log.New(os.Stderr, "[tigeropen] ", log.LstdFlags),
	}
}

// SetLevel 设置日志级别
func (l *DefaultLogger) SetLevel(level Level) {
	l.level = level
}

// Debug 输出 DEBUG 级别日志
func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.Printf("[DEBUG] "+msg, args...)
	}
}

// Info 输出 INFO 级别日志
func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.Printf("[INFO] "+msg, args...)
	}
}

// Warn 输出 WARN 级别日志
func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.logger.Printf("[WARN] "+msg, args...)
	}
}

// Error 输出 ERROR 级别日志
func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	if l.level <= LevelError {
		l.logger.Printf("[ERROR] "+msg, args...)
	}
}

// NopLogger 空日志实现，不输出任何内容
type NopLogger struct{}

func (l *NopLogger) Debug(msg string, args ...interface{}) {}
func (l *NopLogger) Info(msg string, args ...interface{})  {}
func (l *NopLogger) Warn(msg string, args ...interface{})  {}
func (l *NopLogger) Error(msg string, args ...interface{}) {}
func (l *NopLogger) SetLevel(level Level)                  {}

// 确保接口实现
var _ Logger = (*DefaultLogger)(nil)
var _ Logger = (*NopLogger)(nil)

// 全局默认 logger
var defaultLogger Logger = NewDefaultLogger()

// SetDefault 设置全局默认 logger
func SetDefault(l Logger) {
	defaultLogger = l
}

// Default 获取全局默认 logger
func Default() Logger {
	return defaultLogger
}

// 全局便捷方法

// Debugf 输出 DEBUG 级别日志
func Debugf(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

// Infof 输出 INFO 级别日志
func Infof(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

// Warnf 输出 WARN 级别日志
func Warnf(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

// Errorf 输出 ERROR 级别日志
func Errorf(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

// 确保 fmt 包被使用
var _ = fmt.Sprintf

package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDefaultLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := &DefaultLogger{
		level:  LevelWarn,
		logger: log.New(&buf, "", 0),
	}

	l.Debug("debug msg")
	l.Info("info msg")
	l.Warn("warn msg")
	l.Error("error msg")

	output := buf.String()
	if strings.Contains(output, "DEBUG") {
		t.Error("DEBUG 消息不应输出")
	}
	if strings.Contains(output, "INFO") {
		t.Error("INFO 消息不应输出")
	}
	if !strings.Contains(output, "WARN") {
		t.Error("WARN 消息应输出")
	}
	if !strings.Contains(output, "ERROR") {
		t.Error("ERROR 消息应输出")
	}
}

func TestDefaultLogger_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	l := &DefaultLogger{
		level:  LevelDebug,
		logger: log.New(&buf, "", 0),
	}
	l.Debug("test %s", "debug")
	if !strings.Contains(buf.String(), "test debug") {
		t.Error("DEBUG 级别应输出 debug 消息")
	}
}

func TestDefaultLogger_SetLevel(t *testing.T) {
	l := NewDefaultLogger()
	l.SetLevel(LevelError)
	if l.level != LevelError {
		t.Error("SetLevel 未生效")
	}
}

func TestNopLogger(t *testing.T) {
	l := &NopLogger{}
	// 不应 panic
	l.Debug("test")
	l.Info("test")
	l.Warn("test")
	l.Error("test")
	l.SetLevel(LevelDebug)
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{Level(99), "UNKNOWN"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("Level(%d).String() = %s, want %s", tt.level, got, tt.want)
		}
	}
}

func TestGlobalLogger(t *testing.T) {
	original := Default()
	defer SetDefault(original)

	nop := &NopLogger{}
	SetDefault(nop)
	if Default() != nop {
		t.Error("SetDefault 未生效")
	}
	// 全局便捷方法不应 panic
	Debugf("test")
	Infof("test")
	Warnf("test")
	Errorf("test")
}

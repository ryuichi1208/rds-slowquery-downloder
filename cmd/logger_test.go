package cmd

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	testCases := []struct {
		name     string
		logLevel string
	}{
		{
			name:     "デバッグレベル",
			logLevel: "debug",
		},
		{
			name:     "情報レベル",
			logLevel: "info",
		},
		{
			name:     "警告レベル",
			logLevel: "warn",
		},
		{
			name:     "エラーレベル",
			logLevel: "error",
		},
		{
			name:     "不明なレベル（デフォルトは情報レベル）",
			logLevel: "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := NewLogger(tc.logLevel)

			if logger == nil {
				t.Error("NewLogger() returned nil")
			}

			// ロガーの存在確認のみ行う
			// 具体的な実装の詳細はテストせず、nilでないことだけ確認
			if logger.Handler() == nil {
				t.Error("Logger handler is nil")
			}
		})
	}
}

package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateTestLogs(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := filepath.Join(os.TempDir(), "slowquery-test-logs")
	defer os.RemoveAll(tempDir) // テスト完了後にクリーンアップ

	// ログファイルを生成
	err := GenerateTestLogs(tempDir)
	if err != nil {
		t.Fatalf("GenerateTestLogs failed: %v", err)
	}

	// 生成されたファイルが存在するか確認
	files := []string{
		"aws-slowquery.log",
		"gcp-slowquery-1.log",
		"gcp-slowquery-2.log",
	}

	for _, file := range files {
		path := filepath.Join(tempDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist, but it doesn't", path)
		}

		// ファイルサイズを確認（空でないことを確認）
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("Failed to stat file %s: %v", path, err)
			continue
		}

		if info.Size() == 0 {
			t.Errorf("File %s is empty", path)
		}
	}
}

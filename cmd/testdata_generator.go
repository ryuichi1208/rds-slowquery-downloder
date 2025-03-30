package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// 様々なSQLクエリのサンプル
var sampleQueries = []struct {
	query     string
	queryTime float64
	lockTime  float64
	rowsSent  int
	rowsExam  int
	userHost  string
	clientIP  string
	timestamp time.Time
	db        string
}{
	{
		query:     "SELECT * FROM users WHERE id > 1000 AND last_login > '2023-01-01' ORDER BY created_at DESC LIMIT 100",
		queryTime: 2.5,
		lockTime:  0.01,
		rowsSent:  100,
		rowsExam:  1000000,
		userHost:  "app[app]",
		clientIP:  "10.0.1.10",
		timestamp: time.Date(2023, 5, 10, 12, 30, 15, 0, time.UTC),
		db:        "production",
	},
	{
		query:     "SELECT articles.*, users.name FROM articles LEFT JOIN users ON users.id = articles.user_id WHERE articles.published_at > NOW() - INTERVAL 7 DAY",
		queryTime: 5.8,
		lockTime:  0.02,
		rowsSent:  250,
		rowsExam:  2500000,
		userHost:  "web[web]",
		clientIP:  "10.0.1.11",
		timestamp: time.Date(2023, 5, 10, 13, 15, 22, 0, time.UTC),
		db:        "production",
	},
	{
		query:     "UPDATE users SET last_login = NOW() WHERE id IN (SELECT user_id FROM sessions WHERE last_activity > NOW() - INTERVAL 1 HOUR)",
		queryTime: 3.2,
		lockTime:  1.1,
		rowsSent:  0,
		rowsExam:  500000,
		userHost:  "batch[batch]",
		clientIP:  "10.0.1.12",
		timestamp: time.Date(2023, 5, 10, 14, 0, 5, 0, time.UTC),
		db:        "production",
	},
	{
		query:     "DELETE FROM logs WHERE created_at < NOW() - INTERVAL 30 DAY",
		queryTime: 8.9,
		lockTime:  4.5,
		rowsSent:  0,
		rowsExam:  5000000,
		userHost:  "maintenance[maintenance]",
		clientIP:  "10.0.1.13",
		timestamp: time.Date(2023, 5, 10, 23, 30, 0, 0, time.UTC),
		db:        "production",
	},
	{
		query:     "SELECT COUNT(*) FROM orders WHERE status = 'pending' GROUP BY user_id",
		queryTime: 1.2,
		lockTime:  0.001,
		rowsSent:  5000,
		rowsExam:  800000,
		userHost:  "analytics[analytics]",
		clientIP:  "10.0.1.14",
		timestamp: time.Date(2023, 5, 11, 9, 15, 30, 0, time.UTC),
		db:        "production",
	},
	{
		query:     "SELECT products.*, categories.name FROM products JOIN categories ON categories.id = products.category_id WHERE products.stock < 10 ORDER BY products.stock ASC",
		queryTime: 0.9,
		lockTime:  0.005,
		rowsSent:  120,
		rowsExam:  250000,
		userHost:  "inventory[inventory]",
		clientIP:  "10.0.1.15",
		timestamp: time.Date(2023, 5, 11, 10, 45, 12, 0, time.UTC),
		db:        "inventory",
	},
	{
		query:     "SELECT AVG(amount), DATE(created_at) FROM transactions WHERE created_at > NOW() - INTERVAL 90 DAY GROUP BY DATE(created_at) ORDER BY DATE(created_at)",
		queryTime: 4.3,
		lockTime:  0.01,
		rowsSent:  90,
		rowsExam:  3000000,
		userHost:  "finance[finance]",
		clientIP:  "10.0.1.16",
		timestamp: time.Date(2023, 5, 11, 15, 20, 35, 0, time.UTC),
		db:        "finance",
	},
	{
		query:     "INSERT INTO audit_logs (user_id, action, entity_id, created_at) SELECT user_id, 'login', id, created_at FROM sessions WHERE created_at > NOW() - INTERVAL 1 DAY",
		queryTime: 2.1,
		lockTime:  0.8,
		rowsSent:  0,
		rowsExam:  150000,
		userHost:  "audit[audit]",
		clientIP:  "10.0.1.17",
		timestamp: time.Date(2023, 5, 12, 8, 5, 40, 0, time.UTC),
		db:        "audit",
	},
}

// インスタンス名のリスト
var instanceNames = []string{
	"mysql-instance-1",
	"mysql-instance-2",
	"gcp-mysql-prod",
	"gcp-mysql-dev",
}

// GenerateTestLogs はテスト用のスロークエリログファイルを生成します
func GenerateTestLogs(outputDir string) error {
	// テスト用ディレクトリの作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
	}

	// AWSスタイルのログファイル
	awsFile, err := os.Create(filepath.Join(outputDir, "aws-slowquery.log"))
	if err != nil {
		return fmt.Errorf("AWSスタイルログファイルの作成に失敗しました: %w", err)
	}
	defer awsFile.Close()

	// GCPスタイルのログファイル1
	gcpFile1, err := os.Create(filepath.Join(outputDir, "gcp-slowquery-1.log"))
	if err != nil {
		return fmt.Errorf("GCPスタイルログファイル1の作成に失敗しました: %w", err)
	}
	defer gcpFile1.Close()

	// GCPスタイルのログファイル2
	gcpFile2, err := os.Create(filepath.Join(outputDir, "gcp-slowquery-2.log"))
	if err != nil {
		return fmt.Errorf("GCPスタイルログファイル2の作成に失敗しました: %w", err)
	}
	defer gcpFile2.Close()

	// インスタンス名を含むログファイル
	for _, instance := range instanceNames {
		instFile, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("slowquery.%s.log", instance)))
		if err != nil {
			return fmt.Errorf("インスタンスログファイル %s の作成に失敗しました: %w", instance, err)
		}
		defer instFile.Close()

		// インスタンスごとに異なるクエリのセットを書き込む
		for i, q := range sampleQueries {
			if i%len(instanceNames) == indexOf(instanceNames, instance) {
				logEntry := formatSlowQueryLog(q)
				if _, err := instFile.WriteString(logEntry); err != nil {
					return fmt.Errorf("インスタンスファイル %s への書き込みに失敗しました: %w", instance, err)
				}
			}
		}
	}

	// ファイルに書き込む
	// AWSファイルには全クエリを書き込む
	for _, q := range sampleQueries {
		logEntry := formatSlowQueryLog(q)
		if _, err := awsFile.WriteString(logEntry); err != nil {
			return fmt.Errorf("AWSファイルへの書き込みに失敗しました: %w", err)
		}
	}

	// GCPファイル1には奇数インデックスのクエリを書き込む
	for i, q := range sampleQueries {
		if i%2 == 0 {
			logEntry := formatSlowQueryLog(q)
			if _, err := gcpFile1.WriteString(logEntry); err != nil {
				return fmt.Errorf("GCPファイル1への書き込みに失敗しました: %w", err)
			}
		}
	}

	// GCPファイル2には偶数インデックスのクエリを書き込む
	for i, q := range sampleQueries {
		if i%2 == 1 {
			logEntry := formatSlowQueryLog(q)
			if _, err := gcpFile2.WriteString(logEntry); err != nil {
				return fmt.Errorf("GCPファイル2への書き込みに失敗しました: %w", err)
			}
		}
	}

	return nil
}

// formatSlowQueryLog はクエリ情報からMySQLスロークエリログのエントリを生成します
func formatSlowQueryLog(q struct {
	query     string
	queryTime float64
	lockTime  float64
	rowsSent  int
	rowsExam  int
	userHost  string
	clientIP  string
	timestamp time.Time
	db        string
}) string {
	unixTime := q.timestamp.Unix()

	logEntry := fmt.Sprintf(`# Time: %s
# User@Host: %s @ [%s]
# Query_time: %.6f  Lock_time: %.6f Rows_sent: %d  Rows_examined: %d
use %s;
SET timestamp=%d;
%s;

`,
		q.timestamp.Format("2006-01-02T15:04:05.000000Z"),
		q.userHost,
		q.clientIP,
		q.queryTime,
		q.lockTime,
		q.rowsSent,
		q.rowsExam,
		q.db,
		unixTime,
		q.query)

	return logEntry
}

// indexOf はスライス内の要素のインデックスを返します
func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return 0
}

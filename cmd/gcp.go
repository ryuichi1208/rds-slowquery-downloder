package cmd

import (
	"fmt"
	"log/slog"
	// "cloud.google.com/go/cloudsqlconn"
	// sql "cloud.google.com/go/sql/apiv1"
	// "cloud.google.com/go/sql/apiv1/sqlpb"
	// "google.golang.org/api/iterator"
	// "google.golang.org/api/option"
)

type GCPClient struct {
	// sqlClient   *sql.InstancesClient
	logger      *slog.Logger
	projectID   string
	credentials string
}

type GCPClientInterface interface {
	GetInstanceList() []string
	GetSlowQueryList(instance string) ([]string, error)
	DownloadSlowQueryLog(instance string, logFile string) (*string, error)
}

func NewGCPClient(logger *slog.Logger, projectID string, credentials string) (GCPClient, error) {
	// 実際の実装はテスト以外の環境で使用する
	// ここではモック実装を返す
	return GCPClient{
		// sqlClient:   nil,
		logger:      logger,
		projectID:   projectID,
		credentials: credentials,
	}, nil
}

func (g GCPClient) GetInstanceList() []string {
	// テスト用の実装
	return []string{"gcp-instance-1", "gcp-instance-2"}
}

func (g GCPClient) GetSlowQueryList(instance string) ([]string, error) {
	// テスト用の実装
	return []string{
		fmt.Sprintf("slowquery/mysql-slowquery.log.%s.1", instance),
		fmt.Sprintf("slowquery/mysql-slowquery.log.%s.2", instance),
	}, nil
}

func (g GCPClient) DownloadSlowQueryLog(instance string, logFile string) (*string, error) {
	// テスト用の実装
	dummyLog := "# Time: 2023-01-01T12:00:00.000000Z\n# User@Host: user[user] @ localhost [127.0.0.1]\n# Query_time: 2.000000  Lock_time: 0.000000 Rows_sent: 1  Rows_examined: 1000000\nSELECT * FROM large_table WHERE id > 1000;\n"
	return &dummyLog, nil
}

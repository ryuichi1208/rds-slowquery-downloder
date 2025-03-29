package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	sql "cloud.google.com/go/sql/apiv1"
	"cloud.google.com/go/sql/apiv1/sqlpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GCPClient struct {
	sqlClient   *sql.InstancesClient
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
	ctx := context.Background()

	opts := []option.ClientOption{}
	if credentials != "" {
		opts = append(opts, option.WithCredentialsFile(credentials))
	}

	sqlClient, err := sql.NewInstancesClient(ctx, opts...)
	if err != nil {
		logger.Error(err.Error())
		return GCPClient{}, err
	}

	return GCPClient{
		sqlClient:   sqlClient,
		logger:      logger,
		projectID:   projectID,
		credentials: credentials,
	}, nil
}

func (g GCPClient) GetInstanceList() []string {
	var instanceList []string
	ctx := context.Background()

	req := &sqlpb.ListInstancesRequest{
		Project: g.projectID,
	}

	it := g.sqlClient.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			g.logger.Error(fmt.Sprintf("Failed to list instances: %v", err))
			break
		}
		g.logger.Info(fmt.Sprintf("DB instance %v", resp.Name))
		// GCPのインスタンス名はプロジェクト名を含むフルパスなので、最後の部分だけを抽出
		parts := strings.Split(resp.Name, "/")
		instanceName := parts[len(parts)-1]
		instanceList = append(instanceList, instanceName)
	}

	return instanceList
}

func (g GCPClient) GetSlowQueryList(instance string) ([]string, error) {
	var slowQueryList []string
	ctx := context.Background()

	req := &sqlpb.ListInstancesRequest{
		Project: g.projectID,
	}

	it := g.sqlClient.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			g.logger.Error(fmt.Sprintf("Failed to list instances: %v", err))
			return slowQueryList, err
		}

		// インスタンス名の最後の部分を抽出して比較
		parts := strings.Split(resp.Name, "/")
		instanceName := parts[len(parts)-1]
		if instanceName != instance {
			continue
		}

		// スロークエリログファイルのリストを取得
		// 注意：実際のGCP Cloud SQLではログファイルのリスト取得方法は異なります
		// ここではデモンストレーション目的のために簡略化しています
		logFiles, err := g.listLogFiles(ctx, instance)
		if err != nil {
			return slowQueryList, err
		}

		for _, logFile := range logFiles {
			if strings.Contains(logFile, "slowquery") {
				slowQueryList = append(slowQueryList, logFile)
			}
		}
	}

	return slowQueryList, nil
}

func (g GCPClient) listLogFiles(ctx context.Context, instance string) ([]string, error) {
	// 実際のGCP Cloud SQLでは、ログファイルリストの取得は異なるAPI呼び出しになります
	// ここではデモンストレーション目的で仮の実装を提供しています

	// 例として、固定のログファイル名リストを返します
	// 実際の実装では、Cloud SQLのAPIを使用してログファイルを取得する必要があります
	return []string{
		fmt.Sprintf("slowquery/mysql-slowquery.log.%s.1", instance),
		fmt.Sprintf("slowquery/mysql-slowquery.log.%s.2", instance),
		fmt.Sprintf("slowquery/mysql-slowquery.log.%s.3", instance),
	}, nil
}

func (g GCPClient) DownloadSlowQueryLog(instance string, logFile string) (*string, error) {
	// 実際のGCP Cloud SQLでは、ログファイルのダウンロードは異なるAPI呼び出しになります
	// ここではデモンストレーション目的で仮の実装を提供しています

	// 例として、固定の文字列を返します
	dummyLog := fmt.Sprintf("# Time: 2023-01-01T12:00:00.000000Z\n# User@Host: user[user] @ localhost [127.0.0.1]\n# Query_time: 2.000000  Lock_time: 0.000000 Rows_sent: 1  Rows_examined: 1000000\nSELECT * FROM large_table WHERE id > 1000;\n")
	return &dummyLog, nil
}

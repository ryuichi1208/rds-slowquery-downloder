package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type AWSClient struct {
	cfg       aws.Config
	rdsClient *rds.Client
	logger    *slog.Logger
}

type AWSClientInterface interface {
	GetInstanceList() []string
	GetSlowQueryList(instance string) ([]string, error)
	DownloadSlowQueryLog(instance string, logFile string) (*string, error)
}

func NewAWSClient(logger *slog.Logger) (AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Error(err.Error())
		return AWSClient{}, err
	}

	rdsClient := rds.NewFromConfig(cfg)

	return AWSClient{
		cfg:       cfg,
		logger:    logger,
		rdsClient: rdsClient,
	}, nil
}

func FilterInstance(a AWSClientInterface, target string) string {
	instanceList := a.GetInstanceList()

	for _, instance := range instanceList {
		if len(instance) < len(target) {
			continue
		}

		// 前方一致でフィルタリングし最初に見つかったインスタンスを返す
		if instance[:len(target)] == target {
			return instance
		}
	}
	return ""
}

func GetSlowQueryList(a AWSClientInterface, instance string) ([]string, error) {
	return a.GetSlowQueryList(instance)
}

func (a AWSClient) GetInstanceList() []string {
	var instanceList []string
	const maxInstances = 20
	a.logger.Debug(fmt.Sprintf("Let's list up to %v DB instances.", maxInstances))
	output, err := a.rdsClient.DescribeDBInstances(context.TODO(),
		&rds.DescribeDBInstancesInput{MaxRecords: aws.Int32(maxInstances)})
	if err != nil {
		a.logger.Debug(fmt.Sprintf("Couldn't list DB instances: %v", err))
		return instanceList
	}
	if len(output.DBInstances) == 0 {
		a.logger.Debug("No DB instances found.")
	} else {
		for _, instance := range output.DBInstances {
			a.logger.Error(fmt.Sprintf("DB instance %v", *instance.DBInstanceIdentifier))
			instanceList = append(instanceList, *instance.DBInstanceIdentifier)
		}
	}

	return instanceList
}

func (a AWSClient) GetSlowQueryList(instance string) ([]string, error) {
	var slowQueryList []string

	fmt.Println(instance)
	req, err := a.rdsClient.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		a.logger.Error(fmt.Sprintf("Couldn't list DB instances: %v", err))
	}

	if len(req.DBInstances) == 0 {
		return slowQueryList, fmt.Errorf("No DB instances found.")
	}

	for _, dbInstance := range req.DBInstances {
		if *dbInstance.DBInstanceIdentifier != instance {
			continue
		}

		slowQueryReq, _ := a.rdsClient.DescribeDBLogFiles(context.Background(), &rds.DescribeDBLogFilesInput{
			DBInstanceIdentifier: dbInstance.DBInstanceIdentifier,
			FilenameContains:     aws.String("slowquery/mysql-slowquery"),
		})

		for _, logFile := range slowQueryReq.DescribeDBLogFiles {
			slowQueryList = append(slowQueryList, *logFile.LogFileName)
		}
	}
	return slowQueryList, nil
}

func DownloadSlowQueryLog(a AWSClientInterface, instance, filter string, logFile []string) (*string, error) {
	var str *string
	for _, log := range logFile {
		// フィルタの文字列が含まれていない場合はスキップ
		if filter != "" && !strings.Contains(log, filter) {
			continue
		}

		str, err := a.DownloadSlowQueryLog(instance, log)
		if err != nil {
			return str, err
		}

		// logDataをファイルに書き出す
		err = WriteLogData("a.log", str)
		if err != nil {
			return str, err
		}
	}

	return str, nil
}

func WriteLogData(path string, logData *string) error {
	if logData == nil {
		return nil
	}

	// 追記モードでファイルを開く
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(*logData)
	if err != nil {
		return err
	}

	return nil
}

func (a AWSClient) DownloadSlowQueryLog(instance string, logFile string) (*string, error) {
	input := &rds.DownloadDBLogFilePortionInput{
		DBInstanceIdentifier: aws.String(instance),
		LogFileName:          aws.String(logFile),
	}

	req, err := a.rdsClient.DownloadDBLogFilePortion(context.Background(), input)

	if err != nil {
		return nil, err
	}
	return req.LogFileData, nil
}

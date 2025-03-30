# mysql-slowquery-downloder

## Description

This tool is a tool to retrieve logs of slow queries on a specified date from RDS for MySQL and Cloud SQL for MySQL in a batch. By combining it with pt-query-digest etc., you will be able to immediately analyze queries.

## Usage

```
MySQL slow query log downloader

Usage:
  mysql-slowquery-downloder [flags]

Flags:
      --credentials string   path to GCP credentials file
      --debug                debug mode
      --filter string        log filter string (default "f")
  -h, --help                 help for mysql-slowquery-downloder
      --instance string      instance name (default "i")
  -o, --output string        output file path (default "stdout")
      --project string       GCP project ID
      --provider string      cloud provider (aws or gcp) (default "aws")
```

## Cloud Providers

### AWS (Default)

For AWS RDS, the tool will use the default AWS credentials.

### GCP Cloud SQL

For GCP Cloud SQL, you need to specify:

1. `--provider gcp` to use GCP
2. `--project` with your GCP project ID
3. Optional: `--credentials` path to your service account JSON key file

## テストログ生成

このツールはテスト用のMySQLスロークエリログを生成する機能も提供しています。

```
MySQL slow query log downloader - Test Log Generator

Usage:
  mysql-slowquery-downloder testlog [flags]

Flags:
  -o, --output-dir string   生成したログファイルの出力先ディレクトリ (default "testdata")
```

以下のコマンドを実行することで、テスト用のスロークエリログを生成できます：

```
go run main.go testlog -o testdata
```

生成されるログファイル：

- `aws-slowquery.log`: AWSスタイルのスロークエリログ（すべてのクエリを含む）
- `gcp-slowquery-1.log`, `gcp-slowquery-2.log`: GCPスタイルのスロークエリログ
- インスタンス別ログファイル: `slowquery.{インスタンス名}.log`

これらのログファイルは、開発やテスト目的で使用できます。

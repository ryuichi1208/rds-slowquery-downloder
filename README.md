# mysql-slowquery-downloder

## Description

This tool is designed to retrieve logs of slow queries on a specified date from RDS for MySQL and Cloud SQL for MySQL in a batch. By combining it with pt-query-digest etc., you will be able to immediately analyze queries.

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

## Test Log Generation

This tool also provides functionality to generate MySQL slow query logs for testing purposes.

```
MySQL slow query log downloader - Test Log Generator

Usage:
  mysql-slowquery-downloder testlog [flags]

Flags:
  -o, --output-dir string   directory to output generated log files (default "testdata")
```

You can generate test slow query logs by running the following command:

```
go run main.go testlog -o testdata
```

Generated log files:

- `aws-slowquery.log`: AWS-style slow query log (containing all queries)
- `gcp-slowquery-1.log`, `gcp-slowquery-2.log`: GCP-style slow query logs
- Instance-specific log files: `slowquery.{instance-name}.log`

These log files can be used for development and testing purposes.

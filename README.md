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

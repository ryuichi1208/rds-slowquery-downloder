# mysql-slowquery-downloder

## Description

This tool is a tool to retrieve logs of slow queries on a specified date from RDS for MySQL in a batch. By combining it with pt-query-digest etc., you will be able to immediately analyze queries.

## Usage

```
RDS MySQL slow query log downloader

Usage:
  mysql-slowquery-downloder [flags]

Flags:
      --filter string     log filter string (default "f")
  -h, --help              help for mysql-slowquery-downloder
      --instance string   instance name (default "i")
  -o, --output string     output file path (default "stdout")
```

package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mysql-slowquery-downloder",
	Short: "RDS MySQL slow query log downloader",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Do(cmd, args)
	},
}

func Do(cmd *cobra.Command, args []string) error {
	var logger *slog.Logger
	if cmd.Flag("debug").Value.String() == "true" {
		logger = NewLogger("debug")
	} else {
		logger = NewLogger("info")
	}

	logger.Info("Start mysql-slowquery-downloder")

	// クラウドプロバイダーの選択
	provider := cmd.Flag("provider").Value.String()
	instance := ""
	var logList []string
	var err error

	switch provider {
	case "aws":
		var aws AWSClient
		aws, err = NewAWSClient(logger)
		if err != nil {
			return err
		}

		instance = FilterInstance(aws, cmd.Flag("instance").Value.String())
		logList, err = GetSlowQueryList(aws, instance)
		if err != nil {
			return err
		}

		for _, logFile := range logList {
			logger.Debug(fmt.Sprintf("logFile: %s", logFile))
		}

		_, err = DownloadSlowQueryLog(aws, instance, cmd.Flag("filter").Value.String(), logList)
		if err != nil {
			return err
		}
	case "gcp":
		projectID := cmd.Flag("project").Value.String()
		if projectID == "" {
			return fmt.Errorf("GCP project ID is required")
		}

		credentials := cmd.Flag("credentials").Value.String()
		var gcp GCPClient
		gcp, err = NewGCPClient(logger, projectID, credentials)
		if err != nil {
			return err
		}

		instance = FilterInstance(gcp, cmd.Flag("instance").Value.String())
		if instance == "" {
			return fmt.Errorf("No matching instance found")
		}

		logList, err = GetSlowQueryList(gcp, instance)
		if err != nil {
			return err
		}

		for _, logFile := range logList {
			logger.Debug(fmt.Sprintf("logFile: %s", logFile))
		}

		_, err = DownloadSlowQueryLog(gcp, instance, cmd.Flag("filter").Value.String(), logList)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported provider: %s. Use 'aws' or 'gcp'", provider)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("debug", "d", false, "debug mode")
	rootCmd.Flags().String("instance", "i", "instance name")
	rootCmd.Flags().String("filter", "f", "log filter string")
	rootCmd.Flags().StringP("output", "o", "stdout", "output file path")
	rootCmd.Flags().String("provider", "aws", "cloud provider (aws or gcp)")
	rootCmd.Flags().String("project", "", "GCP project ID")
	rootCmd.Flags().String("credentials", "", "path to GCP credentials file")
}

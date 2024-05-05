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

	aws, err := NewAWSClient(logger)
	if err != nil {
		return err
	}

	instance := FilterInstance(aws, cmd.Flag("instance").Value.String())

	logList, err := GetSlowQueryList(aws, instance)
	if err != nil {
		return err
	}

	for _, logFile := range logList {
		logger.Debug(fmt.Sprintf("logFile: %s", logFile))
	}

	DownloadSlowQueryLog(aws, instance, cmd.Flag("filter").Value.String(), logList)
	if err != nil {
		return err
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
}

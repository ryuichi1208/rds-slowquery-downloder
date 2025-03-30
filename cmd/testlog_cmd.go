package cmd

import (
	"github.com/spf13/cobra"
)

// testlogCmd はテスト用のスロークエリログを生成するコマンドです
var testlogCmd = &cobra.Command{
	Use:   "testlog",
	Short: "テスト用のMySQLスロークエリログを生成します",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir, _ := cmd.Flags().GetString("output-dir")
		return GenerateTestLogs(outputDir)
	},
}

func init() {
	rootCmd.AddCommand(testlogCmd)
	testlogCmd.Flags().StringP("output-dir", "o", "testdata", "生成したログファイルの出力先ディレクトリ")
}

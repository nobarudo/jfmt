package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var colorize bool

var rootCmd = &cobra.Command{
	Use:   "jsonfmt",
	Short: "標準入力からJSONを受け取り、見やすく整形します",
	Long:  `パイプ経由で渡されたJSONデータを読み込み、インデントを整えて出力するCLIツールです。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := readStdin()
		if err != nil {
			return err
		}

		var out bytes.Buffer
		if err := json.Indent(&out, input, "", "  "); err != nil {
			return fmt.Errorf("JSONの解析に失敗しました: %w", err)
		}

		result := out.String()
		if colorize {
			result = applyColor(result)
		}

		fmt.Println(result)
		return nil
	},
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
	rootCmd.Flags().BoolVarP(&colorize, "color", "c", false, "シンタックスハイライトを有効にする")
}

// readStdin はパイプからの標準入力を読み取ります
func readStdin() ([]byte, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("エラー: 標準入力が空です。 'curl ... | jsonfmt' のようにパイプを使用してください")
	}
	return io.ReadAll(os.Stdin)
}

func applyColor(j string) string {
	reKey := regexp.MustCompile(`"(.*?)":`)
	reString := regexp.MustCompile(`:\s*"(.*?)"`)
	reNum := regexp.MustCompile(`:\s*([-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?)`)
	reBool := regexp.MustCompile(`:\s*(true|false|null)`)

	cyan := "\033[36m"
	green := "\033[32m"
	magenta := "\033[35m"
	yellow := "\033[33m"
	reset := "\033[0m"

	j = reKey.ReplaceAllString(j, cyan+`"$1"`+reset+":")
	j = reString.ReplaceAllString(j, ": "+green+`"$1"`+reset)
	j = reNum.ReplaceAllString(j, ": "+magenta+`$1`+reset)
	j = reBool.ReplaceAllString(j, ": "+yellow+`$1`+reset)

	return j
}

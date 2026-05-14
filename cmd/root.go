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
	Use:   "jsonfmt [file]",
	Short: "JSONを受け取り、見やすく整形します",
	Long:  `ファイルまたはパイプ経由の標準入力からJSONデータを読み込み、インデントを整えて出力します。`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputReader io.Reader

		if len(args) > 0 {
			// ファイルから読み込む
			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("ファイルを読み込めませんでした: %w", err)
			}
			defer file.Close()
			inputReader = file
		} else {
			// 標準入力から読み込む
			inputReader = cmd.InOrStdin()
		}

		input, err := readInput(inputReader)
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

		fmt.Fprintln(cmd.OutOrStdout(), result)
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

// readInput は指定されたリーダーからデータを読み取ります
func readInput(r io.Reader) ([]byte, error) {
	// os.Stdin の場合のみ、空チェック（CharDeviceチェック）を行う
	if f, ok := r.(*os.File); ok && f == os.Stdin {
		stat, err := f.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) != 0 {
			return nil, fmt.Errorf("エラー: 入力がありません。 'jsonfmt file.json' または 'curl ... | jsonfmt' のように使用してください")
		}
	}
	return io.ReadAll(r)
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

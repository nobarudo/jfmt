package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var minifyCmd = &cobra.Command{
	Use:   "minify [file]",
	Short: "JSONを1行のテキストに圧縮します",
	Long:  `ファイルまたはパイプ経由の標準入力からJSONデータを読み込み、1行のコンパクトな文字列として出力します。`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputReader io.Reader

		if len(args) > 0 {
			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("ファイルを読み込めませんでした: %w", err)
			}
			defer file.Close()
			inputReader = file
		} else {
			inputReader = cmd.InOrStdin()
		}

		input, err := readInput(inputReader)
		if err != nil {
			return err
		}

		var out bytes.Buffer
		if err := json.Compact(&out, input); err != nil {
			return fmt.Errorf("JSONの解析に失敗しました: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), out.String())
		return nil
	},
}

func init() {
	// ルートコマンドにこのサブコマンドを登録する
	rootCmd.AddCommand(minifyCmd)
}

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var minifyCmd = &cobra.Command{
	Use:   "minify",
	Short: "JSONを1行のテキストに圧縮します",
	Long:  `JSON内の余計な空白や改行をすべて削除し、1行のコンパクトな文字列として出力します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := readInput(cmd.InOrStdin())
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

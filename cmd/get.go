package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <path> [file]",
	Short: "JSONから特定の値を抽出します",
	Long: `ドット記法と配列インデックス（例: user.name, users[0].id）を使用して、
JSONデータから特定の値を抽出して出力します。`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		var inputReader io.Reader

		if len(args) > 1 {
			file, err := os.Open(args[1])
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

		var data interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return fmt.Errorf("JSONの解析に失敗しました: %w", err)
		}

		value, err := extractValue(data, path)
		if err != nil {
			return err
		}

		outputValue(cmd.OutOrStdout(), value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func extractValue(data interface{}, path string) (interface{}, error) {
	if path == "" || path == "." {
		return data, nil
	}

	// a[0].b -> a.[0].b に変換して分割しやすくする
	normalizedPath := strings.ReplaceAll(path, "[", ".[")
	normalizedPath = strings.ReplaceAll(normalizedPath, "]", "")
	parts := strings.Split(normalizedPath, ".")

	current := data
	for _, part := range parts {
		if part == "" {
			continue
		}

		if strings.HasPrefix(part, "[") {
			// 配列アクセス
			indexStr := part[1:]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, fmt.Errorf("不正な配列インデックスです: %s", part)
			}

			slice, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("対象は配列ではありません: %s", part)
			}
			if index < 0 || index >= len(slice) {
				return nil, fmt.Errorf("インデックスが範囲外です: %d", index)
			}
			current = slice[index]
		} else {
			// オブジェクトアクセス
			obj, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("対象はオブジェクトではありません: %s", part)
			}
			val, exists := obj[part]
			if !exists {
				return nil, fmt.Errorf("キーが見つかりません: %s", part)
			}
			current = val
		}
	}

	return current, nil
}

func outputValue(w io.Writer, value interface{}) {
	switch v := value.(type) {
	case string:
		fmt.Fprintln(w, v)
	case float64:
		fmt.Fprintln(w, v)
	case bool:
		fmt.Fprintln(w, v)
	case nil:
		fmt.Fprintln(w, "null")
	default:
		// オブジェクトや配列は整形して出力
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(v)
	}
}

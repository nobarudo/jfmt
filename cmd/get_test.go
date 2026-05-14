package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetCmd(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_get*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	jsonContent := `{
		"name": "test-user",
		"age": 30,
		"active": true,
		"address": {
			"city": "Tokyo",
			"tags": ["work", "home"]
		},
		"items": [
			{"id": 1, "name": "item1"},
			{"id": 2, "name": "item2"}
		]
	}`
	if _, err := tmpFile.WriteString(jsonContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	tests := []struct {
		name        string
		input       string
		args        []string
		expectedOut string
		expectErr   bool
	}{
		{
			name:        "Simple key",
			input:       jsonContent,
			args:        []string{"get", "name"},
			expectedOut: "test-user\n",
		},
		{
			name:        "Nested key",
			input:       jsonContent,
			args:        []string{"get", "address.city"},
			expectedOut: "Tokyo\n",
		},
		{
			name:        "Array index",
			input:       jsonContent,
			args:        []string{"get", "address.tags[1]"},
			expectedOut: "home\n",
		},
		{
			name:        "Object in array",
			input:       jsonContent,
			args:        []string{"get", "items[0].name"},
			expectedOut: "item1\n",
		},
		{
			name:        "Whole object output",
			input:       jsonContent,
			args:        []string{"get", "address"},
			expectedOut: "{\n  \"city\": \"Tokyo\",\n  \"tags\": [\n    \"work\",\n    \"home\"\n  ]\n}\n",
		},
		{
			name:        "Read from file",
			input:       "",
			args:        []string{"get", "age", tmpFile.Name()},
			expectedOut: "30\n",
		},
		{
			name:      "Key not found",
			input:     jsonContent,
			args:      []string{"get", "unknown"},
			expectErr: true,
		},
		{
			name:      "Index out of range",
			input:     jsonContent,
			args:      []string{"get", "items[5]"},
			expectErr: true,
		},
		{
			name:      "Not an array access",
			input:     jsonContent,
			args:      []string{"get", "name[0]"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outBuf := new(bytes.Buffer)
			inBuf := bytes.NewBufferString(tt.input)

			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(outBuf)
			rootCmd.SetIn(inBuf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if !tt.expectErr {
				// Encode(v) は末尾に改行を追加するので、期待値と比較する際に注意
				actual := outBuf.String()
				if actual != tt.expectedOut {
					// 複数行のJSON出力の場合は、トリミングして比較するか、厳密に比較する
					// ここでは厳密に比較
					t.Errorf("expected output:\n%q\ngot:\n%q", tt.expectedOut, actual)
				}
			}
		})
	}
}

func TestExtractValue(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": []interface{}{
				map[string]interface{}{"c": "target"},
			},
		},
	}

	tests := []struct {
		path    string
		want    interface{}
		wantErr bool
	}{
		{"a.b[0].c", "target", false},
		{"a.b", data["a"].(map[string]interface{})["b"], false},
		{"a.x", nil, true},
		{"a.b[1]", nil, true},
		{"", data, false},
		{".", data, false},
	}

	for _, tt := range tests {
		got, err := extractValue(data, tt.path)
		if (err != nil) != tt.wantErr {
			t.Errorf("extractValue(%s) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !strings.Contains(fmt.Sprint(got), fmt.Sprint(tt.want)) {
			// interface{} の比較を簡略化
			t.Errorf("extractValue(%s) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

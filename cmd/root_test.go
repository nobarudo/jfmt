package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestRootCmd(t *testing.T) {
	// Create a temporary file for testing file input
	tmpFile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{"file": "test"}`
	if _, err := tmpFile.WriteString(content); err != nil {
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
			name:        "Valid JSON from Stdin",
			input:       `{"name":"test","value":1}`,
			args:        []string{},
			expectedOut: "{\n  \"name\": \"test\",\n  \"value\": 1\n}\n",
			expectErr:   false,
		},
		{
			name:        "Valid JSON from File",
			input:       "",
			args:        []string{tmpFile.Name()},
			expectedOut: "{\n  \"file\": \"test\"\n}\n",
			expectErr:   false,
		},
		{
			name:      "File Not Found",
			input:     "",
			args:      []string{"nonexistent.json"},
			expectErr: true,
		},
		{
			name:      "Invalid JSON",
			input:     `{name:"test"}`,
			args:      []string{},
			expectErr: true,
		},
		{
			name:  "Colorized JSON output",
			input: `{"key": "string", "num": 1, "bool": true}`,
			args:  []string{"--color"},
			expectedOut: "{\n  \x1b[36m\"key\"\x1b[0m: \x1b[32m\"string\"\x1b[0m,\n  \x1b[36m\"num\"\x1b[0m: \x1b[35m1\x1b[0m,\n  \x1b[36m\"bool\"\x1b[0m: \x1b[33mtrue\x1b[0m\n}\n",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			colorize = false

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
				if outBuf.String() != tt.expectedOut {
					t.Errorf("expected output: %q, got: %q", tt.expectedOut, outBuf.String())
				}
			}
		})
	}
}

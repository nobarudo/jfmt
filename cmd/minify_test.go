package cmd

import (
	"bytes"
	"testing"
)

func TestMinifyCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		args        []string
		expectedOut string
		expectErr   bool
	}{
		{
			name: "Valid JSON format",
			input: `{
  "name": "test",
  "value": 1
}`,
			args:        []string{"minify"},
			expectedOut: `{"name":"test","value":1}` + "\n",
			expectErr:   false,
		},
		{
			name:      "Invalid JSON",
			input:     `{name:"test"}`,
			args:      []string{"minify"},
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
				if outBuf.String() != tt.expectedOut {
					t.Errorf("expected output: %q, got: %q", tt.expectedOut, outBuf.String())
				}
			}
		})
	}
}

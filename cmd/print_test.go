package cmd

import (
	"bytes"
	"testing"
)

func TestPrint_InvalidArg(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"print"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

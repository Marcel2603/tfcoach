package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRoot_HelpShown(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error :%v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "Usage:") {
		t.Fatalf("expected usage help, got %q", got)
	}
}

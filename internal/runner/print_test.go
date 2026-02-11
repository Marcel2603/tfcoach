package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunPrint(t *testing.T) {
	testDir := t.TempDir()
	reportPath := filepath.Join(testDir, "report.json")
	err := os.WriteFile(reportPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatal("Setup error", err)
	}

	type args struct {
		srcReportPath string
		outputFormat  string
		allowEmojis   bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "From file OK",
			args: args{
				srcReportPath: reportPath,
				outputFormat:  "pretty",
				allowEmojis:   true,
			},
			want: 0,
		},
		{
			name: "From unknown file",
			args: args{
				srcReportPath: "/hello",
				outputFormat:  "compact",
				allowEmojis:   false,
			},
			want: 1,
		},
		{
			name: "Conversion error",
			args: args{
				srcReportPath: reportPath,
				outputFormat:  "abcd",
				allowEmojis:   true,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			returnCode := Print(tt.args.srcReportPath, w, tt.args.outputFormat, tt.args.allowEmojis)
			if returnCode != tt.want {
				t.Errorf("Print() = %v, want %v", returnCode, tt.want)
			}
		})
	}
}

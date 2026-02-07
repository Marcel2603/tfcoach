package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/Marcel2603/tfcoach/internal/format"
)

func Print(srcReportPath string, w io.Writer, outputFormat string, allowEmojis bool) int {
	var reportContent []byte
	var err error
	if srcReportPath == "-" {
		reportContent, err = io.ReadAll(os.Stdin)
	} else {
		reportContent, err = os.ReadFile(srcReportPath)
	}
	if err != nil {
		_, _ = fmt.Fprintf(w, "failed to read %q: %v", srcReportPath, err)
		return 1
	}

	err = format.ReformatResults(reportContent, w, outputFormat, allowEmojis)
	if err != nil {
		_, _ = fmt.Fprintf(w, "failed to convert report: %v", err)
		return 2
	}

	return 0
}

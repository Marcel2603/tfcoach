//revive:disable:var-naming For now it's okay to have a generic name
package format

import (
	"encoding/json"
	"io"
)

func writeJSON(issues []issueOutput, w io.Writer) error {
	output := jsonOutput{
		IssueCount: len(issues),
		Issues:     issues,
	}
	outputAsStr, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(outputAsStr)
	if err != nil {
		return err
	}
	return nil
}

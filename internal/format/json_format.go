package format

import (
	"encoding/json"
	"io"
)

type jsonOutput struct {
	IssueCount int           `json:"issue_count"`
	Issues     []issueOutput `json:"issues"`
}

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

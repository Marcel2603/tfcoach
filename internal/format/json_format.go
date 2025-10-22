package format

import (
	"encoding/json"
	"io"

	"github.com/Marcel2603/tfcoach/internal/types"
)

type jsonOutput struct {
	IssueCount int           `json:"issue_count"`
	Issues     []issueOutput `json:"issues"`
}

func writeJSON(issues []types.Issue, w io.Writer) error {
	output := jsonOutput{
		IssueCount: len(issues),
		Issues:     toIssueOutputs(issues),
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

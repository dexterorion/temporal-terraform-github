package workflows

import (
	"encoding/json"
	"errors"

	"github.com/dexterorion/meetupgo/github"
	"go.temporal.io/sdk/worker"
)

var (
	ErrWorkflowNotFound = errors.New("workflow not found")
)

func Register(w worker.Worker) {
	github.Register(w)
}

func MapInputToWorkflow(wfname string, input string) (interface{}, error) {
	switch wfname {
	case "CreateRepositoryWorkflow":
		in := github.CreateRepositoryInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DestroyRepositoryWorkflow":
		in := github.DestroyRepositoryInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "CreateBranchWorkflow":
		in := github.CreateBranchInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	default:
		return nil, ErrWorkflowNotFound
	}
}

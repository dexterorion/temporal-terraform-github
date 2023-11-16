package github

import "go.temporal.io/sdk/worker"

func Register(w worker.Worker) {
	w.RegisterWorkflow(CreateRepositoryWorkflow)
	w.RegisterActivity(CreateRepositoryActivity)

	w.RegisterWorkflow(DestroyRepositoryWorkflow)
	w.RegisterActivity(DestroyRepositoryActivity)

	w.RegisterWorkflow(CreateBranchWorkflow)
	w.RegisterActivity(CreateBranchActivity)
}

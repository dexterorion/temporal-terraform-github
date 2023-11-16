# Infra Automation

## Start Worker

First, you should set proper environment variable.


### ENV VARS
Must create a `.env` file with format:

- `TEMPORAL_HOSTPORT`="127.0.0.1:7233" (default)
- `TEMPORAL_NAMESPACE`="default" (default)
- `TEMPORAL_TASKQUEUE`="meetupgo" (default)
- `TERRAFORM_STATE_BUCKET`="meetupgo" (default)
- `GITHUB_TOKEN`="" (testing purposes)

### Running 

```
go run cmd/worker/main.go
```

## Run a workflow

#### Create Repository

```
go run cmd/orchestrator/main.go -workflow=CreateRepositoryWorkflow -taskqueue=meetupgo -input='{"Name": "test", "Description": "The Test Repo"}'
```

#### Delete Repository

```
go run cmd/orchestrator/main.go -workflow=DestroyRepositoryWorkflow -taskqueue=meetupgo -input='{"Name": "test"}'
```

#### Create Branch

```
go run cmd/orchestrator/main.go -workflow=CreateBranchWorkflow -taskqueue=meetupgo -input='{"Repository": "test", "Branch": "develop"}'
```
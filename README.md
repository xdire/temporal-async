# temporal-async
Asynchronous execution of ordinary workflow

# Why
This repository shows how to create Temporal worker with a simple processing task

# How

1) Temporal binary need to be installed in your system
```shell
temporal server start-dev --db-filename temporal.db
```

2) Run worker from this folder
```shell
go run service-worker/main.go
```

3) Run the API service to manipulate the workflow
```shell
go run service-api/main.go
```

4) Start the workflow and observe workflow in the temporal UI
```shell
curl -v -X POST 'http://localhost:8093/process/create?name=xxx&desc=testprocess'
```

5) Workflow will be awaiting next signals, that can update the Process record
```shell
curl -X PATCH 'http://localhost:8093/process/stage?id=proc_${ID_FROM_UI}&stage=execute'
```

6) Confirm that workflow can move into the final part of the execution
```shell
curl -X PATCH 'http://localhost:8093/process/stage?id=proc_${ID_FROM_UI}&stage=finalize'
```

7) Complete the workflow with the confirmation signal
```shell
curl -X PATCH 'http://localhost:8093/process/confirm?id=proc_${ID_FROM_UI}'   
```

8) Observe workflow completed and went through all the execution steps in the Event History
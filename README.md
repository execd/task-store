# Task Store
Simple redis based task store.

# Building and running
```bash
$ gofmt -w  . && megacheck ./... && go test -v -race ./... && golint -set_exit_status $(go list ./...) && go build
$ PORT=8080 ./task-store
```

To test:

```bash
$ curl -XPOST -d '{"name":"test", "image":"alpine", "init":"init.sh"}' localhost:8080/task/
task:1
```
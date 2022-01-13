include .env

.PHONY: fmt
fmt: ## Run go fmt for the whole project
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do gofmt -e -l -w $$d/*.go; done)

run:
	go run ./cmd/api/main.go
go-build:
	go build -o ./build/api ./cmd/api/main.go
run-build:
	./build/api
go-get:
	go get -v all
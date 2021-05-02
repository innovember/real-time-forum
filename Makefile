include .env
run:
	go run ./cmd/api/main.go
go-build:
	go build -o ./build/api ./cmd/api/main.go
run-build:
	./build/api
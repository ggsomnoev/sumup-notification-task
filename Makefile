lint:
	golangci-lint run
tidy:
	go mod tidy
run:
	go run .\cmd\main.go
test:
	ginkgo ./...
generate:
	go generate ./...

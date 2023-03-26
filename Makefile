local: gen-mocks tests
	go mod tidy
	go run main.go

tests:
	go test -v ./... -cover

gen-mocks:
		go generate ./...	
local: gen-mocks tests
	go mod tidy
	go run main.go

build: gen-mocks tests
	go mod tidy
	go build .	

tests:
	go test -v ./... -cover

gen-mocks:
		go generate ./...	
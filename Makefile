.PHONY: build
build:
	go build cmd/loy-certificate-api/main.go

.PHONY: test
test:
	go test -v ./...
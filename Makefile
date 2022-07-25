test:
	go test -v ./... -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out

gen: clean
	go run gen/*.go commands
	go generate ./...

commands:
	go run gen/*.go commands
	go generate commands/...

commands/%:
	go run gen/*.go -class $* commands
	go generate ./...

hostapi:
	go run gen/*.go hostapi
	go generate ./...

hostapi/%:
	go run gen/*.go -class $* hostapi
	go generate ./...

clean:
	rm -rf commands hostapi

.PHONY: gen clean test coverage

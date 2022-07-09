

gen:
	go run gen/main.go
	go generate ./...

.PHONY: gen

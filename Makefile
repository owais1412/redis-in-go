build:
	@go build -o bin/app cmd/main.go

run: build
	@./bin/app
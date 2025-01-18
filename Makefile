build:
	@go build -o bin/equity-api cmd/main.go

run: build
	@./bin/equity-api
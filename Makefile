.PHONY: build
build:
	@echo "Building spotifycli..."
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o bin/spotifycli
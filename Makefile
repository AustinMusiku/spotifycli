.PHONY: build snapshot changelog release

current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
current_commit = $(shell git rev-parse --short HEAD)
version = $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)

PKG = github.com/AustinMusiku/spotifycli/internal/version
LDFLAGS = -s -w \
  -X $(PKG).Version=$(version) \
  -X $(PKG).Commit=$(current_commit) \
  -X $(PKG).BuildDate=$(current_time)

build:
	@echo "Building spotifycli..."
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/spotifycli

snapshot:
	goreleaser release --snapshot --clean

release:
	git tag $(VERSION)
	git push origin $(VERSION)
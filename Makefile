# Usage:
#   make test            # run tests with coverage printed inline
#   make cover           # write coverage profile to coverage.out
#   make cover-html      # openable HTML report at coverage.html
#   make clean           # remove coverage artifacts

PKGS        := $(shell go list ./...)
COVERMODE   := atomic
COVERFILE   := coverage.out
HTMLFILE    := coverage.html
GOFLAGS     := -race -shuffle=on

.PHONY: test cover cover-html build clean docs-rules

test:
	go test $(GOFLAGS) ./... -cover

cover:
	go test $(GOFLAGS) ./... -covermode=$(COVERMODE) -coverprofile=$(COVERFILE)
	@echo
	@go tool cover -func=$(COVERFILE) | tail -n1

cover-html: cover
	go tool cover -html=$(COVERFILE) -o $(HTMLFILE)
	@echo "Wrote $(HTMLFILE)"

build:
	 go build .

clean:
	rm -f $(COVERFILE) $(HTMLFILE)

doc-rules:
	@go run ./tools/cmd/gen-rules-doc/main.go > docs/pages/rules/index.md

format:
	@gofmt -l -s -w .

lint: format
	@which revive > /dev/null || go install github.com/mgechev/revive@latest
	@revive -config config.toml -formatter friendly ./...
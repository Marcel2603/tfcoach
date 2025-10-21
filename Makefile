# Usage:
#   make test            # run tests with coverage printed inline
#   make cover           # write coverage profile to coverage.out
#   make cover-html      # openable HTML report at coverage.html
#   make clean           # remove coverage artifacts

PKGS        := $(shell go list ./...)
COVERMODE   := atomic
COVERFILE   := coverage.out
HTMLFILE    := coverage.html
GOFLAGS     := -race -shuffle=on -tags=test -covermode=$(COVERMODE) -coverprofile=$(COVERFILE)

.PHONY: test cover cover-html build clean generate-documentation lint

test:
	go test $(GOFLAGS) $(PKGS) -cover

cover:
	go test $(GOFLAGS) $(PKGS)
	@echo
	@go tool cover -func=$(COVERFILE) | tail -n1

cover-html: cover
	go tool cover -html=$(COVERFILE) -o $(HTMLFILE)
	@echo "Wrote $(HTMLFILE)"

build:
	 go build .

clean:
	rm -fv $(COVERFILE) $(HTMLFILE)

generate-documentation:
	@go run -tags tfcoach_tools ./tools/cmd/gen-docs

format:
	@gofmt -l -s -w .

lint: format
	@go run github.com/mgechev/revive@latest -config config.toml -formatter friendly ./...

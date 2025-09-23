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

.PHONY: test cover cover-html clean

test:
	go test $(GOFLAGS) ./... -cover

cover:
	go test $(GOFLAGS) ./... -covermode=$(COVERMODE) -coverprofile=$(COVERFILE)
	@echo
	@go tool cover -func=$(COVERFILE) | tail -n1

cover-html: cover
	go tool cover -html=$(COVERFILE) -o $(HTMLFILE)
	@echo "Wrote $(HTMLFILE)"

clean:
	rm -f $(COVERFILE) $(HTMLFILE)

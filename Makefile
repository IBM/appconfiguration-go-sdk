MODULE      := github.com/IBM/appconfiguration-go-sdk
GOTEST      := go test
GOVET       := go vet
GOFMT       := gofmt
GOIMPORTS   := goimports
GOLINT      := golangci-lint
TESTFLAGS   := -race -count=1 -timeout=120s
COVERFLAGS  := -coverprofile=coverage.out -covermode=atomic
PKGS        := ./...

all: fmt tidy vet lint test

fmt:
	$(GOFMT) -l -w .
	@which $(GOIMPORTS) > /dev/null 2>&1 && $(GOIMPORTS) -local $(MODULE) -w . || \
		echo "goimports not installed; run: go install golang.org/x/tools/cmd/goimports@latest"

tidy:
	go mod tidy
	go mod verify

vet:
	$(GOVET) $(PKGS)

lint:
	@which $(GOLINT) > /dev/null 2>&1 || \
		(echo "golangci-lint not installed; see https://golangci-lint.run/usage/install/" && exit 1)
	$(GOLINT) run --config=.golangci.yml $(PKGS)

test:
	$(GOTEST) $(TESTFLAGS) $(COVERFLAGS) $(PKGS)
	go tool cover -func=coverage.out

test-race:
	$(GOTEST) -race -count=1 -timeout=120s $(PKGS)

coverage: test
	go tool cover -html=coverage.out

build:
	go build $(PKGS)

clean:
	rm -f coverage.out

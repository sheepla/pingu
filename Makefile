NAME = pingu
BIN := bin/$(NAME)

# version e.g. v0.0.1
VERSION := $(shell git describe --tags --abbrev=0 | tr -d "v")
# commit hash of HEAD e.g. 3a913f
REVISION := $(shell git rev-parse --short HEAD)

LDFLAGS := -w \
		   -s \
		   -X "main.appVersion=$(VERSION)" \
		   -X "main.appRevision=$(REVISION)"

COVERAGE_OUT := .test/cover.out
COVERAGE_HTML := .test/cover.html

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN)

.PHONY: fmt
fmt:
	go fmt

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: test
test:
	mkdir -p .test
	go test -coverprofile=$(COVERAGE_OUT) ./...

.PHONY: coverage
coverage:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

.PHONY: clean
clean:
	rm $(BIN)
	rm $(COVERAGE_OUT)
	rm $(COVERAGE_HTML)


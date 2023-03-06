tag := $(shell git describe --exact-match --tags 2>git_describe_error.tmp; rm -f git_describe_error.tmp)
branch := $(shell git rev-parse --abbrev-ref HEAD)
commit := $(shell git rev-parse --short=8 HEAD)
glibc_version := 2.17
version := $(tag:v%=%)

MAKEFLAGS += --no-print-directory
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
HOSTGO := env -u GOOS -u GOARCH -u GOARM -- go

LDFLAGS := $(LDFLAGS) -X main.commit=$(commit) -X main.branch=$(branch) -X main.goos=$(GOOS) -X main.goarch=$(GOARCH)
ifneq ($(tag),)
	LDFLAGS += -X main.version=$(version)
endif

# Go built-in race detector works only for 64 bits architectures.
ifneq ($(GOARCH), 386)
	race_detector := -race
endif


GOFILES ?= $(shell git ls-files '*.go')
GOFMT ?= $(shell gofmt -l -s $(GOFILES))

prefix ?= /usr/local
bindir ?= $(prefix)/bin
sysconfdir ?= $(prefix)/etc
localstatedir ?= $(prefix)/var
pkgdir ?= build/dist

.PHONY: all
all:
	@$(MAKE) deps
	@$(MAKE) wsb

.PHONY: help
help:
	@echo 'Targets:'
	@echo '  all        - download dependencies and compile wsb binary'
	@echo '  deps       - download dependencies'
	@echo '  wsb        - compile wsb binary'
	@echo '  test       - run short unit tests'
	@echo '  fmt        - format source files'
	@echo '  tidy       - tidy go modules'
	@echo '  lint       - run linter'
	@echo '  check-deps - check docs/LICENSE_OF_DEPENDENCIES.md'
	@echo '  clean      - delete build artifacts'
	@echo ''

.PHONY: deps
deps:
	go mod download -x

.PHONY: wsb
wsb:
	go build -ldflags "$(LDFLAGS)" -o ./wsb/cmd/wsb .

.PHONY: test
test:
	go test -short $(race_detector) ./...

.PHONY: fmt
fmt:
	@gofmt -s -w $(GOFILES)

.PHONY: fmtcheck
fmtcheck:
	@if [ ! -z "$(GOFMT)" ]; then \
		echo "[ERROR] gofmt has found errors in the following files:"  ; \
		echo "$(GOFMT)" ; \
		echo "" ;\
		echo "Run make fmt to fix them." ; \
		exit 1 ;\
	fi

.PHONY: test-windows
test-windows:
	go test -short ./...

.PHONY: vet
vet:
	@echo 'go vet $$(go list ./...)'
	@go vet $$(go list ./...) ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "go vet has found suspicious constructs. Please remediate any reported errors"; \
		echo "to fix them before submitting code for review."; \
		exit 1; \
	fi

.PHONY: lint
lint:
ifeq (, $(shell which golangci-lint))
	$(info golangci-lint can't be found, please install it: https://golangci-lint.run/usage/install/)
	exit 1
endif

	golangci-lint -v run

.PHONY: tidy
tidy:
	go mod verify
	go mod tidy
	@if ! git diff --quiet go.mod go.sum; then \
		echo "please run go mod tidy and check in changes"; \
		exit 1; \
	fi

.PHONY: check
check: fmtcheck vet

.PHONY: test-all
test-all: fmtcheck vet
	go test $(race_detector) ./...

.PHONY: clean
clean:
	rm -f wsb
	rm -f wsb.exe
	rm -rf build


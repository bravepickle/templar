# Helper methods for working with GO application

# Color settings
CL_RED := \033[31m
CL_GREEN := \033[32m
CL_YELLOW := \033[33m
CL_RESET := \033[0m

# Disable coloring output
ifeq ($(NOCOLOR),1)
	CL_RED =
    CL_GREEN =
    CL_YELLOW =
    CL_RESET =
endif

# App settings
OUT_DIR ?= $(shell pwd)
APP_CONFIGS_DIR ?= $(OUT_DIR)

# Output for test coverage
OUT_TESTS_COVER ?= ./tests_coverprofile.out

# build flags
GIT_COMMIT ?= $(shell git rev-list -1 HEAD | cut -c 1-7)
LDFLAGS ?= -X github.com/bravepickle/templar.GitCommitHash=$(GIT_COMMIT) \
	-X github.com/bravepickle/templar.AppVersion=$(APP_VERSION) \
	-X github.com/bravepickle/templar.AppConfigsDir=$(APP_CONFIGS_DIR)

# skip staticheck linting
SKIP_STATICCHECK ?=

# skip goimports tool
SKIP_GOIMPORTS ?=

# interactive mode for commands
APP_INTERACTIVE ?= 1

# skip pulling changes form git on make update command
APP_SKIP_GIT_UPDATE ?= 0

# skip updating GO dependencies
APP_SKIP_PKG_UPDATE ?= 1

# skip tests on some actions
APP_SKIP_TESTS ?= 0

# skip check for vulnerabilities in some cases
APP_SKIP_VULNCHECK ?= 0

# path to listen to for displaying documentation as a web server
APP_DOC_SERVER ?= localhost:8080

GOCMD ?= go

APP_RELEASE_SUFFIX=-$(shell $(GOCMD) env GOARCH)-$(shell $(GOCMD) env GOOS)

.DEFAULT_GOAL := build

# prepare and run application
.PHONY: all
all: test build

.PHONY: setup
setup:
	$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GOCMD) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GOCMD) install golang.org/x/pkgsite/cmd/pkgsite@latest
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest

	$(GOCMD) mod download

# format all GO files
fmt:
	@echo "$(CL_YELLOW)=> Code formatting...$(CL_RESET)"
	$(GOCMD) fmt ./...
ifeq ($(APP_SKIP_GOIMPORTS),1)
	@echo Skipping goimports...
else
	goimports -w .
endif

# static analysis (aka lint)
.PHONY: lint
lint: fmt
	@echo "$(CL_YELLOW)=> Code analysis...$(CL_RESET)"
	$(GOCMD) vet ./...

ifeq ($(SKIP_STATICCHECK),1)
	@echo Skipping staticcheck linting...
else
	staticcheck ./...
endif

# build binary from source code - for development mostly. May lack extra meta data which can be found in "release" task
.PHONY: build
build: lint
	@echo "$(CL_YELLOW)=> Building app...$(CL_RESET)"
	$(GOCMD) build .

# build binary from source code and add extra data for release
.PHONY: release
release: lint
	@echo "$(CL_YELLOW)=> Building app...$(CL_RESET)"
ifneq ($(APP_SKIP_VULNCHECK),1)
	govulncheck ./...
endif

# build & install binary from source code
.PHONY: install
install: lint
	@echo "$(CL_YELLOW)=> Installing app...$(CL_RESET)"

ifneq ($(APP_SKIP_VULNCHECK),1)
	govulncheck ./...
endif

# project cleanup - caches, temporary files etc.
.PHONY: clean
clean:
	@echo "$(CL_YELLOW)App cleanup...$(CL_RESET)"
	$(GOCMD) clean
	rm -f $(OUT_TESTS_COVER)

# run tests
.PHONY: test
test: lint
	@echo "$(CL_YELLOW)Running tests...$(CL_RESET)"

	$(GOCMD) test ./... $(TEST_ARGS)

# run benchmarks
.PHONY: bench
bench: lint
	@echo "$(CL_YELLOW)Running benchmarks...$(CL_RESET)"

	$(GOCMD) test -bench=. -benchmem -cpu ./...

# run test coverage
.PHONY: cover
cover: lint
	@echo "$(CL_YELLOW)Running tests with code coverage...$(CL_RESET)"
	$(GOCMD) test -coverprofile $(OUT_TESTS_COVER) ./...
	$(GOCMD) tool cover -func $(OUT_TESTS_COVER)

# run test coverage in HTML view
.PHONY: cover_web
cover_web: lint
	@echo "$(CL_YELLOW)Running tests with code coverage with WEB view...$(CL_RESET)"
	$(GOCMD) test -coverprofile $(OUT_TESTS_COVER) ./...
	$(GOCMD) tool cover -html $(OUT_TESTS_COVER)

# show documentation from source code
.PHONY: doc
doc: lint
	@echo "$(CL_YELLOW)Show docs...$(CL_RESET)"
	$(GOCMD) doc -all -C .

# display documentation in web UI (aka browser)
.PHONY: doc_web
doc_web: lint
	@echo "$(CL_YELLOW)Show docs in WEB...$(CL_RESET)"

ifeq ($(shell which pkgsite),)
	@echo "$(CL_GREEN)pkgsite$(CL_RESET) application not found. Try checking paths or installing it. E.g. \"$(CL_YELLOW)go install golang.org/x/pkgsite/cmd/pkgsite@latest$(CL_RESET)\""
else
	@echo "Starting DOCs server at $(CL_GREEN)$(APP_DOC_SERVER)$(CL_RESET)..."
	pkgsite -http $(APP_DOC_SERVER)
endif

# list available Makefile targets
.PHONY: help
help:
	@echo "$(CL_YELLOW)Available targets:$(CL_RESET)"
	@awk '/^[a-zA-Z_-]+:/ {print substr($$1, 1, length($$1)-1)}' $(MAKEFILE_LIST)

# update dependencies & application
.PHONY: update
update:
ifneq ($(APP_SKIP_GIT_UPDATE),1)
	git status
	git pull
endif

	$(GOCMD) mod tidy

ifneq ($(APP_SKIP_PKG_UPDATE),1)
# update external apps
	$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GOCMD) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GOCMD) install golang.org/x/pkgsite/cmd/pkgsite@latest
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest

# update package dependencies
#	$(GOCMD) get -u ./...
endif

	$(GOCMD) mod vendor

# skip tests after update check
ifneq ($(APP_SKIP_TESTS),1)
	$(MAKE) test
endif



TARGET      = xkcdpwd

PACKAGE     = github.com/wfscheper/xkcdpwd
DATE       ?= $(shell date +%Y-%m-%d)
COMMIT      = $(shell git rev-parse --short HEAD 2>/dev/null)
VERSION    ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(CURDIR)/.version 2> /dev/null || echo v0)
BIN         = $(GOPATH)/bin
PKGS        = $(or $(PKG),$(shell $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS    = $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))

PLATFORMS = linux darwin windows
ARCHES    = amd64 386

# Allow tags to be set on command-line, but don't set them
# by default
override TAGS := $(and $(TAGS),-tags $(TAGS))

GOLDFLAGS  := '-X main.version=$(VERSION) -X main.buildDate=$(DATE) -X main.commitHash=$(COMMIT)'
GO      = go
GOBUILD = $(GO) build -v
GOTEST  = $(GO) test -tags test
GODOC   = godoc
GOFMT   = gofmt
DOCKER  = docker
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: generate fmt lint vendor ; $(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build $(TAGS) \
		-ldflags $(GOLDFLAGS) \
		-o bin/$(TARGET) cmd/xkcdpwd/main.go

.PHONY: dist
dist: generate fmt lint vendor ; $(info $(M) building distributions...)
	$Q [[ -d dist ]] && rm -rf dist/ || : ; \
	mkdir dist ; \
	for GOOS in $(PLATFORMS) ; do \
		for GOARCH in $(ARCHES) ; do \
			binout=$(TARGET)-$$GOOS-$$GOARCH ; \
			GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) -a $(TAGS) \
				-ldflags $(GOLDFLAGS) \
				-o dist/$$binout ; \
			sha256sum dist/$$binout > dist/$$binout.sha256 ; \
		done ; \
	done

# Tools

GODEP = $(BIN)/dep
$(BIN)/dep: ; $(info $(M) building go dep…)
	$Q go get github.com/golang/dep/cmd/dep

GOLINT = $(BIN)/golint
$(BIN)/golint: ; $(info $(M) building golint…)
	$Q go get github.com/golang/lint/golint

GOCOV = $(BIN)/gocov
$(BIN)/gocov: ; $(info $(M) building gocov…)
	$Q go get github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: ; $(info $(M) building gocov-xml…)
	$Q go get github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: ; $(info $(M) building go2xunit…)
	$Q go get github.com/tebeka/go2xunit

GOBINDATA = $(BIN)/go-bindata
$(BIN)/go-bindata: ; $(info $(M) building go-bindata...)
	$Q go get github.com/shuLhan/go-bindata/...

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: generate fmt lint vendor ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GOTEST) -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint vendor | $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q 2>&1 $(GOTEST) -timeout 20s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/cover.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test
test-coverage: fmt lint vendor test-coverage-tools ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)
	$Q $(GOTEST) \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_PROFILE)" \
			$(TESTPKGS)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: test-update
test-update: fmt lint vendor ; $(info $(M) updating golden files...)
	$Q $(GOTEST) ./cmd/xkcdpwd/... -update

.PHONY: lint
lint: vendor | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

.PHONY: generate generate-tools
generate-tools: | $(GOBINDATA)
generate: generate-tools internal/langs/languages.go ; $(info $(M) running go generate...)
internal/langs/languages.go: internal/langs/languages/en
	$Q go generate ./...

# Dependency management

Gopkg.toml: | $(GODEP); $(info $(M) generating Gopkg.toml…)
	$Q $(GODEP) init

Gopkg.lock: | Gopkg.toml $(GODEP); $(info $(M) updating Gopkg.lock…)
	$Q $(GODEP) ensure -no-vendor

vendor: Gopkg.lock | $(GODEP) ; $(info $(M) retrieving dependencies…)
	$Q $(GODEP) ensure
	@touch $@
.PHONY: vendor-update
vendor-update: | $(GODEP)
ifeq "$(origin PKG)" "command line"
	$(info $(M) updating $(PKG) dependency…)
	$Q $(GODEP) ensure -update $(PKG)
else
	$(info $(M) updating all dependencies…)
	$Q $(GODEP) ensure -update
endif
	@touch vendor

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -f internal/langs/languages.go
	@rm -rf bin/ dist/
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)

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
GO          = go
GOBUILD     = $(GO) build -v
GOTEST      = $(GO) test
GODOC       = godoc
GOFMT       = $(GO) fmt
DOCKER      = docker
TIMEOUT     = 15

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

# Environment
GOOS   = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

.PHONY: all
all: vendor generate fmt lint ; $(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build $(TAGS) \
		-ldflags $(GOLDFLAGS) \
		-o bin/$(TARGET) cmd/xkcdpwd/main.go

.PHONY: dist
dist: vendor generate fmt lint ; $(info $(M) building distributions...)
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

TOOLS := tools
$(TOOLS):
	@mkdir -p $@

GODEP := $(TOOLS)/dep
$(GODEP): | $(TOOLS) ; $(info $(M) building go dep…)
	$Q curl -fsSL https://github.com/golang/dep/releases/download/v0.4.1/dep-$(GOOS)-$(GOARCH) -o $@
	$Q echo "$$(curl -fsSL https://github.com/golang/dep/releases/download/v0.4.1/dep-$(GOOS)-$(GOARCH).sha256 | awk '{print $$1}')  $@" | sha256sum -c - >/dev/null
	$Q chmod +x $@

GOLINT := $(TOOLS)/golint
$(GOLINT): | $(TOOLS) ; $(info $(M) building golint…)
	$Q $(GOBUILD) -o $@ ./vendor/github.com/golang/lint/golint

GOCOV := $(TOOLS)/gocov
$(GOCOV): | $(TOOLS) ; $(info $(M) building gocov…)
	$Q $(GOBUILD) -o $@ ./vendor/github.com/axw/gocov/gocov

GOCOVXML := $(TOOLS)/gocov-xml
$(GOCOVXML): | $(TOOLS) ; $(info $(M) building gocov-xml…)
	$Q $(GOBUILD) -o $@ ./vendor/github.com/AlekSi/gocov-xml

GO2XUNIT := $(TOOLS)/go2xunit
$(GO2XUNIT): ; $(info $(M) building go2xunit…)
	$Q $(GOBUILD) -o $@ ./vendor/github.com/tebeka/go2xunit

GOBINDATA = $(TOOLS)/go-bindata
$(GOBINDATA): ; $(info $(M) building go-bindata...)
	$Q $(GOBUILD) -o $@ ./vendor/github.com/shuLhan/go-bindata/cmd/go-bindata

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

test-xml: vendor fmt lint | $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
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
test-update: vendor fmt lint ; $(info $(M) updating golden files...)
	$Q $(GOTEST) ./cmd/xkcdpwd/... -update

.PHONY: lint
lint: vendor | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) running go fmt…) @ ## Run go fmt on all packages
	$Q $(GOFMT) $(PKGS)

.PHONY: generate
generate: | $(GOBINDATA) ; $(info $(M) running go generate...)
	$Q PATH=$(TOOLS):$$PATH && $(GO) generate ./dictionary.go

# Dependency management

Gopkg.lock: Gopkg.toml | $(GODEP); $(info $(M) updating Gopkg.lock…)
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
	@rm -rf bin/ dist/ tools/
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)

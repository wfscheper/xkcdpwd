TARGET   = bin/xkcdpwd
SOURCES := $(shell find . -name '*.go' -a -not -name '*_test.go')

# commands
GOTAGGER      = bin/gotagger
GOLANGCILINT  = bin/golangci-lint
GORELEASER    = bin/goGORELEASER
STENTOR       = bin/stentor
GOTESTSUM     = bin/gotestsum

# variables
BUILDDATE       := $(shell date +%Y-%m-%d)
COMMIT          := $(shell git rev-parse HEAD)
GORELEASERFLAGS  = $(if $(filter false,$(DRYRUN)),,--snapshot) --rm-dist
GOTAGGERFLAGS    = $(if $(filter false,$(DRYRUN)),-push,)
REPORTDIR        = reports
STENTORFLAGS     = $(if $(filter false,$(DRYRUN)),-release)

# output controls
override Q = $(if $(filter 1,$(V)),,@)
override M = ▶

.PHONY: all
all: lint build test

.PHONY: help
help:
	@printf "Available make targets:\
	\n  all:         Run the lint, build, and test targets. Default target.\
	\n  build:       Compile $(TARGET).\
	\n  changelog:   Generate changelog updates. Set DRYRUN=false to update CHANGELOG.md.\
	\n  clean:       Delete generated files.\
	\n  coverage:    Run tests and generate coverage reports.\
	\n  dist:        Build distributable archives. Set DRYRUN=false to publish to github.\
	\n  fmt format:  Format source code.\
	\n  lint:        Lint source code.\
	\n  release:     Tag repository with current version if DRYRUN is set to false.\
	\n  show:        Print make variables, eg. 'make show VALUE=TARGET'.\
	\n  test tests:  Run tests.\
	\n  test-report: Run tests and generate a xUnit-style test report.\
	\n  test-watch:  Run tests and watch for changes (hit CTRL-C to stop).\
	\n  version:     Print the version of xkcdpwd.\
	\n"

.PHONY: build
build: $(TARGET)

.PHONY: changelog
changelog: | $(STENTOR) $(GOTAGGER) ; $(info $(M) generating changelog)
	$Q $(STENTOR) $(STENTORFLAGS) $(shell bin/gotagger) $(shell git tag --list --sort=-version:refname | head -n1)

.PHONY: clean
clean:
	$Q $(RM) -r bin/ dist/ reports/

.PHONY: dist
dist: | $(GORELEASER) ; $(info $(M) building dist…)
	$(GORELEASER) release $(GORELEASERFLAGS)

.PHONY: fmt format ; $(info $(M) formatting…)
fmt format: GOLANGCILINTFLAGS += --fix
fmt format: lint | $(GOLANGCILINT)

.PHONY: lint
lint: | $(GOLANGCILINT) ; $(info $(M) linting…)
	$Q $(GOLANGCILINT) run $(GOLANGCILINTFLAGS)

.PHONY: test tests coverage
coverage: TESTFLAGS += -cover -covermode=atomic
coverage: tests
test tests: | $(GOTESTSUM) ; $(info $(M) running tests…)
	$Q $(GOTESTSUM) $(GOTESTSUMFLAGS) -- $(TESTFLAGS) ./...

.PHONY: test-report
test-report: GOTESTSUMFLAGS += --junitfile reports/junit.xml
test-report: TESTFLAGS += -cover -covermode=atomic -coverprofile=reports/cover.out
test-report: $(REPORTDIR)
test-report: test

.PHONY: test-watch
test-watch: GOTESTSUMFLAGS += --watch
test-watch: test

.PHONY: release
release: | $(GOTAGGER); $(info $(M) tagging release…)
	$Q $(GOTAGGER) $(GOTAGGERFLAGS)

.PHONY: show
show:
	@echo $(VALUE)=$($(VALUE))

.PHONY: version
version: | $(GOTAGGER)
	@$(GOTAGGER)

$(GOTAGGER): tools/go.mod tools/go.sum ; $(info $(M) building $(GOTAGGER)…)
	cd tools && GOBIN=$(CURDIR)/bin go install github.com/sassoftware/gotagger/cmd/gotagger

$(GOLANGCILINT): tools/go.mod tools/go.sum ; $(info $(M) building $(GOLANGCILINT)…)
	cd tools && GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GORELEASER): tools/go.mod tools/go.sum ; $(info $(M) building $(GORELEASER)…)
	cd tools && GOBIN=$(CURDIR)/bin go install github.com/goGORELEASER/goGORELEASER

$(REPORTDIR):
	@mkdir -p $@

$(STENTOR): tools/go.mod tools/go.sum ; $(info $(M) building $(STENTOR)…)
	cd tools && GOBIN=$(CURDIR)/bin go install github.com/wfscheper/stentor/cmd/stentor

$(TARGET): $(SOURCES) go.mod go.sum | $(GOTAGGER); $(info $(M) building $(TARGET)…)
	go build -o $@ -mod=readonly -ldflags "-X main.buildDate=$(BUILDDATE) -X main.commit=$(COMMIT) -X main.version=$$(bin/gotagger)" ./cmd/xkcdpwd/

$(GOTESTSUM): tools/go.mod tools/go.sum ; $(info $(M) building $(GOTESTSUM)…)
	cd tools && GOBIN=$(CURDIR)/bin go install gotest.tools/gotestsum

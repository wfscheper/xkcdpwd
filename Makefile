BINARY   = bin/xkcdpwd

# commands
GOTAGGER       = tools/gotagger
GOLANGCILINT   = tools/golangci-lint
GORELEASER     = tools/goreleaser
GOTESTSUM      = tools/gotestsum
STENTOR        = tools/stentor
TOOLS          = $(GOTAGGER) $(GOLANGCILINT) $(GORELEASER) $(GOTESTSUM) $(STENTOR)

# variables
BUILDDATE   := $(shell date +%Y-%m-%d)
COMMIT      := $(shell git rev-parse HEAD)
RELEASEFLAGS = $(if $(filter false,$(DRYRUN)),,--snapshot)
STENTORFLAGS = $(if $(filter false,$(DRYRUN)),-release)
TESTFLAGS    = -cover -covermode=atomic

# output controls
override Q = $(if $(filter 1,$(V)),,@)
override M = ▶

.PHONY: all
all: lint build test

.PHONY: build
build: $(BINARY)

.PHONY: changelog
changelog: | $(STENTOR) $(GOTAGGER) ; $(info $(M) generating changelog)
	$Q $(STENTOR) $(STENTORFLAGS) $(shell $(GOTAGGER)) $(shell git tag --list --sort=-version:refname | head -n1)

.PHONY: clean
clean:
	$Q $(RM) -r bin/ dist/ cover.out junit.xml $(TOOLS)

.PHONY: release
release: | $(GORELEASER) ; $(info $(M) building release dist…)
	$(GORELEASER) release

.PHONY: snapshot
snapshot: | $(GORELEASER) ; $(info $(M) building snapshot dist…)
	$(GORELEASER) release --snapshot --clean

.PHONY: fmt format ; $(info $(M) formatting…)
fmt format: GOLANGCILINTFLAGS += --fix
fmt format: lint | $(GOLANGCILINT)

.PHONY: lint
lint: | $(GOLANGCILINT) ; $(info $(M) linting…)
	$Q $(GOLANGCILINT) run $(GOLANGCILINTFLAGS)

.PHONY: test tests
test tests: | $(GOTESTSUM) ; $(info $(M) running tests…)
	$Q $(GOTESTSUM) $(GOTESTSUMFLAGS) -- $(TESTFLAGS) ./...

.PHONY: test-report
test-report: GOTESTSUMFLAGS += --junitfile junit.xml
test-report: TESTFLAGS += -coverprofile=cover.out
test-report: test

.PHONY: test-watch
test-watch: GOTESTSUMFLAGS += --watch
test-watch: test

.PHONY: show
show:
	@echo $(value)=$($(value))

.PHONY: version
version: | $(GOTAGGER)
	@$(GOTAGGER)

# real targets
.PHONY: FORCE
$(BINARY): FORCE | $(GOTAGGER); $(info $(M) building $(BINARY)…)
	go build -o $@ -mod=readonly -ldflags "-X main.buildDate=$(BUILDDATE) -X main.commit=$(COMMIT) -X main.version=$$($(GOTAGGER))" ./cmd/xkcdpwd/

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod

define build_tool
$(1): tools/go.mod tools/go.sum ; $$(info $$(M) building $(1)…)
	cd tools && go build -mod=readonly $(2)
endef

FOO = $(call build_tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint)
$(eval $(call build_tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint))
$(eval $(call build_tool,$(GORELEASER),github.com/goreleaser/goreleaser))
$(eval $(call build_tool,$(GOTAGGER),github.com/sassoftware/gotagger/cmd/gotagger))
$(eval $(call build_tool,$(GOTESTSUM),gotest.tools/gotestsum))
$(eval $(call build_tool,$(STENTOR),github.com/wfscheper/stentor/cmd/stentor))

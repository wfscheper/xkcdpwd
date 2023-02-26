//go:build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/sassoftware/gotagger/cmd/gotagger"
	_ "github.com/wfscheper/stentor/cmd/stentor"
	_ "gotest.tools/gotestsum"
)

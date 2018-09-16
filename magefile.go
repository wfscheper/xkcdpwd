// Copyright © 2017 Walter Scheper <walter.scheper@gmal.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build mage

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"

	// tools
	_ "github.com/shuLhan/go-bindata"
	_ "golang.org/x/lint"
)

const (
	targetName  = "xkcdpwd"
	packageName = "github.com/wfscheper/" + targetName
)

var (
	// Default is the default mage target
	Default = Build

	ldflags = "-X main.version=$VERSION -X main.buildDate=$DATE -X main.commitHash=$COMMIT"
	goexe   = "go"

	// commands
	gofmt     = sh.RunCmd(goexe, "fmt")
	gotest    = sh.RunCmd(goexe, "test", "-timeout", "15s")
	govet     = sh.RunCmd(goexe, "vet")
	goveralls = filepath.Join("tools", "goveralls")
	golint    = filepath.Join("tools", "golint")

	// distribution targets
	platforms = []string{"darwin", "linux", "windows"}
	arches    = []string{"386", "amd64"}
)

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// Force use of go modules
	os.Setenv("GO111MODULES", "on")
}

func buildEnvs() map[string]string {
	commit, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	version, _ := sh.Output("git", "describe", "--tags", "--always", "--dirty", "--match=v*")
	return map[string]string{
		"COMMIT":  commit,
		"DATE":    time.Now().UTC().Format(time.RFC3339),
		"VERSION": version,
	}
}

func buildTags() string {
	if tags := os.Getenv("BUILD_TAGS"); tags != "" {
		return tags
	}
	return "none"
}

func toolsEnv() map[string]string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tools := filepath.Join(cwd, "tools")
	path := strings.Join([]string{
		tools,
		os.Getenv("PATH"),
	}, string(os.PathListSeparator))
	return map[string]string{
		"GOBIN": tools,
		"PATH":  path,
	}
}

// Build builds the xkcdpwd binary
func Build(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate, Fmt, Lint, Vet)
	fmt.Println("building " + targetName + "…")
	return sh.RunWith(buildEnvs(), goexe, "build", "-v", "-tags", buildTags(), "-ldflags", ldflags, "-o",
		filepath.Join("bin", targetName), filepath.Join("cmd", "xkcdpwd", "main.go"))
}

// Dist prepare a release
func Dist(ctx context.Context) (err error) {
	mg.CtxDeps(ctx, TestRace)
	fmt.Println("building distribution…")
	for _, goos := range platforms {
		for _, goarch := range arches {
			binname := strings.Join([]string{targetName, goos, goarch}, "-")
			if "windows" == goos {
				binname += ".exe"
			}
			env := buildEnvs()
			env["GOOS"] = goos
			env["GOARCH"] = goarch
			fmt.Println("building " + binname + "…")
			err = sh.RunWith(env, goexe, "build", "-v", "-a",
				"-tags", buildTags(),
				"-ldflags", ldflags,
				"-o", filepath.Join("dist", binname),
				filepath.Join("cmd", "xkcdpwd", "main.go"))
			if err != nil {
				break
			}
			err = checksum(binname)
			if err != nil {
				break
			}
		}
	}
	return err
}

func checksum(filename string) (err error) {
	os.Chdir("dist")
	defer os.Chdir("..")
	var checksum string
	switch runtime.GOOS {
	case "darwin":
		checksum, err = sh.Output("shasum", "-a", "256", filename)
	default:
		checksum, err = sh.Output("sha256sum", filename)
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename+".sha256", []byte(checksum), 0644)
}

// Fmt runs go fmt
func Fmt(ctx context.Context) error {
	fmt.Println("running go fmt…")
	return gofmt("./...")
}

// Generate generates dictionary files
func Generate(ctx context.Context) error {
	rebuild, err := target.Dir(filepath.Join(".", "internal", "langs", "languages.go"), filepath.Join(".", "internal", "langs", "languages"))
	if err != nil {
		return err
	}
	if rebuild {
		mg.CtxDeps(ctx, getGobindata)
		fmt.Println("running go generate…")
		return sh.RunWith(toolsEnv(), goexe, "generate", filepath.Join(".", "dictionary.go"))
	}
	return nil
}

func getGobindata(ctx context.Context) error {
	rebuild, err := target.Path(filepath.Join("tools", "go-bindata"))
	if err != nil {
		return err
	}
	if rebuild {
		fmt.Println("getting go-bindata…")
		return sh.RunWith(toolsEnv(), goexe, "install", "github.com/shuLhan/go-bindata/cmd/go-bindata")
	}
	return nil
}

// Lint runs golint
func Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, getGolint)
	fmt.Println("running golint…")
	return sh.Run(golint, "./...")
}

func getGolint(ctx context.Context) error {
	if rebuild, err := target.Path(golint); err != nil {
		return err
	} else if rebuild {
		fmt.Println("getting golint…")
		return sh.RunWith(toolsEnv(), goexe, "install", "golang.org/x/lint/golint")
	}
	return nil
}

// Vet runs go vet
func Vet(ctx context.Context) error {
	fmt.Println("running go vet…")
	return govet("./...")
}

// Test runs the test suite
func Test(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate, Fmt, Vet, Lint)
	return runTest()
}

// TestRace runs the test suite with race detection
func TestRace(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate, Fmt, Vet, Lint)
	return runTest("-race")
}

// TestShort runs only tests marked as short
func TestShort(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate, Fmt, Vet, Lint)
	return runTest("-short")
}

// Benchmark runs the benchmark suite
func Benchmark(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate, Fmt, Vet, Lint)
	return runTest("-run=__absolutelynothing__", "-bench")
}

func runTest(testType ...string) error {
	var space string
	if len(testType) > 1 {
		space = " "
	}
	fmt.Printf("running go test%s%s…\n", space, strings.Join(testType, " "))
	testType = append(testType, "./...")
	return gotest(testType...)
}

// Coverage generates coverage reports
func Coverage(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate, Fmt, Vet, Lint)
	sh.Run("mkdir", "-p", "coverage")
	mode := os.Getenv("COVERAGE_MODE")
	if mode == "" {
		mode = "atomic"
	}
	if err := runTest("-cover", "-covermode", mode, "-coverprofile=coverage/cover.out"); err != nil {
		return err
	}
	if err := sh.Run(goexe, "tool", "cover", "-html=coverage/cover.out", "-o", "coverage/index.html"); err != nil {
		return err
	}
	return nil
}

// Coveralls uploads coverage report
func Coveralls(ctx context.Context) error {
	// only do something if within travis
	if os.Getenv("TRAVIS_HOME") == "" {
		return nil
	}
	mg.CtxDeps(ctx, getGoveralls)
	fmt.Println("running goveralls…")
	return sh.Run(goveralls, "-coverprofile=coverage/cover.out", "-service=travis-ci")
}

func getGoveralls(ctx context.Context) error {
	if rebuild, err := target.Path(goveralls); err != nil {
		return err
	} else if rebuild {
		fmt.Println("getting goveralls…")
		return sh.RunWith(toolsEnv(), goexe, "install", "github.com/mattn/goveralls")
	}
	return nil
}

func Clean() error {
	return sh.Run("rm", "-rf", "bin/", "dist/", "tools/", "coverage/")
}

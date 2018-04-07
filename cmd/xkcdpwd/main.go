// Copyright © 2017 Walter Scheper <walter.scheper@gmal.com>
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

package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"

	dict "github.com/wfscheper/xkcdpwd"
)

const (
	successExitCode = 0
	errorExitCode   = 1
)

var (
	appName    = "xkcdpwd"
	version    = "deveL"
	buildDate  string
	commitHash string
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get working directory", err)
		os.Exit(errorExitCode)
	}

	exc := &Xkcdpwd{
		Args:       os.Args,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		WorkingDir: wd,
		Env:        os.Environ(),
	}
	exitCode := exc.Run()
	os.Exit(exitCode)
}

// Xkcdpwd specifies an execution of xkcdpwd
type Xkcdpwd struct {
	WorkingDir     string    // Where to execute
	Args           []string  // command-line arguments
	Env            []string  // os environment
	Stdout, Stderr io.Writer // output writers
}

// Run executes xkcdpwd
func (x *Xkcdpwd) Run() int {
	var (
		// flags
		showVersion bool
	)

	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(x.Stderr)
	_ = flags.Bool("v", false, "be more verbose")

	// register global flags
	flags.BoolVar(&showVersion, "version", false, "show version information")

	// wrap stdout and stderr in loggers
	outLogger := log.New(x.Stdout, "", 0)
	errLogger := log.New(x.Stderr, "", 0)

	setUsage(errLogger, flags)
	if err := flags.Parse(x.Args[1:]); err != nil {
		return errorExitCode
	}

	if showVersion {
		outLogger.Printf(`%s
 version     : %s
 build date  : %s
 git hash    : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, appName, version, buildDate, commitHash, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
		return successExitCode
	}
	d := dict.GetDict("en")
	l := big.NewInt(int64(d.Length() - 1))
	words := make([]string, 4, 4)
	for i := 0; i < 10; i++ {
		for j := 0; j < 4; j++ {
			idx, err := rand.Int(rand.Reader, l)
			if err != nil {
				errLogger.Println("error: cannot generate random words:", err)
				return errorExitCode
			}
			words[j] = d.Word(int(idx.Int64()))
		}
		outLogger.Println(strings.Join(words, " "))
	}
	return successExitCode
}

func setUsage(logger *log.Logger, fs *flag.FlagSet) {
	var flagsUsage bytes.Buffer
	tw := tabwriter.NewWriter(&flagsUsage, 0, 4, 2, ' ', 0)
	fs.VisitAll(func(f *flag.Flag) {
		switch f.DefValue {
		case "":
			fmt.Fprintf(tw, "\t-%s\t%s\n", f.Name, f.Usage)
		default:
			fmt.Fprintf(tw, "\t-%s\t%s (default: %s)\n", f.Name, f.Usage, f.DefValue)
		}
	})
	tw.Flush()
	fs.Usage = func() {
		logger.Printf("Usage: %s [OPTIONS]\n", appName)
		logger.Println()
		logger.Printf("%s is a passphrase generator based on XKCD comic #936\n", appName)
		logger.Println()
		logger.Println("Flags:")
		logger.Println()
		logger.Println(flagsUsage.String())
	}
}

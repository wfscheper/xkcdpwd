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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/pelletier/go-toml"
	dict "github.com/wfscheper/xkcdpwd"
	"github.com/wfscheper/xkcdpwd/internal/userinfo"
)

const (
	successExitCode = 0
	errorExitCode   = 1
)

var (
	appName    = "xkcdpwd"
	version    = "devel"
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
	// wrap stdout and stderr in loggers
	outLogger := log.New(x.Stdout, "", 0)
	errLogger := log.New(x.Stderr, "", 0)

	cfgfileDefault, err := userinfo.DefaultConfigFile(appName)
	if err != nil {
		errLogger.Printf("Cannot determine default config file location: %s", err)
		return errorExitCode
	}

	var cfgfile string
	var cfg *toml.Tree
	cfgFlags := flag.NewFlagSet(appName, flag.ContinueOnError)
	cfgFlags.SetOutput(io.Discard)
	cfgFlags.StringVar(&cfgfile, "cfgfile", "", "path to config file")
	if err := cfgFlags.Parse(x.Args[1:]); err != nil {
		switch err {
		case flag.ErrHelp:
		default:
			cfgfile = cfgfileDefault
		}
	}
	cfg, err = toml.LoadFile(cfgfile)
	if err != nil {
		if !os.IsNotExist(err) {
			errLogger.Printf("cannot read config file '%s': %s", cfgfile, err)
			return errorExitCode
		}
	}
	if cfg == nil {
		cfg, _ = toml.TreeFromMap(map[string]interface{}{})
	}

	// register global flags
	var (
		// flags
		capitalize      string
		lang            string
		maxWordLength   int
		minWordLength   int
		passphraseCount int
		separator       string
		showVersion     bool
		wordCount       int
	)
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(x.Stderr)

	_ = flags.Bool("v", false, "be more verbose")
	_ = flags.String("cfgfile", cfgfile, "path to config file")
	flags.BoolVar(&showVersion, "version", false, "show version information")

	var capitalizeDefault = cfg.GetDefault(appName+".capitalize", "none").(string)
	flags.StringVar(&capitalize, "capitalize", capitalizeDefault, "capitalize letters in passphrase")

	var langDefault = cfg.GetDefault(appName+".lang", "").(string)
	var langHelpDefault string
	if langDefault == "" {
		langHelpDefault = " (default: en)"
	}
	flags.StringVar(&lang, "lang", langDefault, fmt.Sprintf("language to use, a valid IETF language tag%s", langHelpDefault))

	var maxWordLengthDefault = cfg.GetDefault(appName+".max-length", int64(0)).(int64)
	flags.IntVar(&maxWordLength, "max-length", int(maxWordLengthDefault), "maximum word length")

	var minWordLengthDefault = cfg.GetDefault(appName+".min-length", int64(0)).(int64)
	flags.IntVar(&minWordLength, "min-length", int(minWordLengthDefault), "minimum word length")

	var passphraseCountDefault = cfg.GetDefault(appName+".phrases", int64(10)).(int64)
	flags.IntVar(&passphraseCount, "phrases", int(passphraseCountDefault), "the number of passphrases")

	var separatorDefault = cfg.GetDefault(appName+".separator", " ").(string)
	flags.StringVar(&separator, "separator", separatorDefault, "passphrase separator")

	var wordCountDefault = cfg.GetDefault(appName+".words", int64(4)).(int64)
	flags.IntVar(&wordCount, "words", int(wordCountDefault), "the number of words in each passphrase")

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

	// check that words is valid
	if wordCount <= 0 {
		errLogger.Printf("error: words must be greater than 0")
		return errorExitCode
	}

	// check that phrases is valid
	if passphraseCount <= 0 {
		errLogger.Printf("error: phrases must be greater than 0")
		return errorExitCode
	}

	// check that separator is valid
	if !checkSeparator(separator) {
		errLogger.Printf("error: invalid separator '%s'\n", separator)
		return errorExitCode
	}

	// check that capitalize is valid
	switch capitalize {
	case "all", "first", "none", "random":
	default:
		errLogger.Printf("error: invalid capitalization strategy '%s'", capitalize)
		return errorExitCode
	}

	// Source lang from the environment, but prefer the command line if set
	if envLang, ok := os.LookupEnv("LANG"); ok && lang == "" {
		lang = envLang
	}
	d := dict.GetDict(lang)
	d.SetCapitalize(capitalize)
	d.SetMaxWordLength(maxWordLength)
	d.SetMinWordLength(minWordLength)
	for i := 0; i < passphraseCount; i++ {
		words, err := d.Passphrase(wordCount)
		if err != nil {
			errLogger.Printf("error: %v\n", err)
			return errorExitCode
		}
		outLogger.Println(strings.Join(words, separator))
	}
	return successExitCode
}

func checkSeparator(sep string) bool {
	switch sep {
	case "", " ", ".", "-", "_", "=":
		return true
	default:
		return false
	}
}

func setUsage(logger *log.Logger, fs *flag.FlagSet) {
	var flagsUsage bytes.Buffer
	tw := tabwriter.NewWriter(&flagsUsage, 0, 4, 2, ' ', 0)
	fs.VisitAll(func(f *flag.Flag) {
		switch f.DefValue {
		case "":
			fmt.Fprintf(tw, "\t-%s\t%s\n", f.Name, f.Usage)
		case " ":
			fmt.Fprintf(tw, "\t-%s\t%s (default: '%s')\n", f.Name, f.Usage, f.DefValue)
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

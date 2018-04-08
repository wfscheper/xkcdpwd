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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wfscheper/xkcdpwd/internal/test"
)

// Entry point for running integration tests.
func TestIntegration(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var relPath = "testdata"
	filepath.Walk(relPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal("error walking testdata")
		}

		if filepath.Base(path) != "testcase.json" {
			return nil
		}

		segments := strings.Split(path, string(filepath.Separator))
		// testName is the everything after "testdata/", excluding "testcase.json"
		testName := strings.Join(segments[1:len(segments)-1], "/")
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			t.Run("external", runTest(testName, relPath, wd, execCmd))
			t.Run("internal", runTest(testName, relPath, wd, runMain))
		})
		return nil
	})
}

func runTest(name, relPath, wd string, run test.RunFunc) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		testCase := test.NewCase(t, filepath.Join(wd, relPath), name)

		// Skip tests
		if testCase.Skip {
			t.Skipf("skipping %s", name)
		}

		testEnv := test.NewEnvironment(t, wd, run)

		var err error
		for i, args := range testCase.Commands {
			err = testEnv.Run(args)
			if err != nil && i < len(testCase.Commands)-1 {
				t.Fatalf("cmd '%s' raised an unexpected error: %s", strings.Join(args, " "), err.Error())
			}
		}

		if *test.UpdateGolden {
			testCase.UpdateStderr(testEnv.GetStderr())
			testCase.UpdateStdout(testEnv.GetStdout())
		} else {
			testCase.CompareError(err, testEnv.GetStderr())
			testCase.CompareOutput(testEnv.GetStdout())
		}
	}
}

func execCmd(prog string, args []string, stdout, stderr io.Writer, dir string, env []string) error {
	cmd := exec.Command(prog, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Dir = dir
	cmd.Env = env
	return cmd.Run()
}

func runMain(prog string, args []string, stdout, stderr io.Writer, dir string, env []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = r
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	exc := &Xkcdpwd{
		Args:       append([]string{prog}, args...),
		Stdout:     stdout,
		Stderr:     stderr,
		WorkingDir: dir,
		Env:        env,
	}
	if exitCode := exc.Run(); exitCode != 0 {
		err = fmt.Errorf("exit status %d", exitCode)
	}
	return
}

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

package userinfo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigDir(t *testing.T) {
	actual := DefaultConfigDir("foo")
	if "foo/.config" != actual {
		t.Errorf("wrong config dir")
	}
}

func TestDefaultConfigFile(t *testing.T) {
	actual, err := DefaultConfigFile("foo")
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join("/home", os.Getenv("USER"), ".config", "foo", "foo.conf")
	if actual != expected {
		t.Errorf("wrong config file")
	}
}

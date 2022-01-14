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

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigDir(t *testing.T) {
	got := DefaultConfigDir("foo")
	assert.Equal(t, "foo/.config", got)
}

func TestDefaultConfigFile(t *testing.T) {
	got, err := DefaultConfigFile("foo")
	if assert.NoError(t, err) {
		want := filepath.Join("/home", os.Getenv("USER"), ".config", "foo", "foo.conf")
		assert.Equal(t, want, got)
	}
}

func TestDefaultConfigDirEnv(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", "/foo")
	t.Cleanup(func() {
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	})

	got := DefaultConfigDir("bar")
	assert.Equal(t, "/foo", got)
}

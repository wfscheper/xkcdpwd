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
)

// DefaultConfigDir returns the path to the OS specific user config directory.
func DefaultConfigDir(homeDir string) string {
	if cfgDir, ok := os.Lookup("APPDATA"); ok {
		return cfgDir
	}
	return filepath.Join(homeDir, "AppData", "Roaming")
}

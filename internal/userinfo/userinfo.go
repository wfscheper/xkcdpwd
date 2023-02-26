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

//go:build darwin || linux || windows
// +build darwin linux windows

package userinfo

import (
	"errors"
	"os/user"
	"path/filepath"
)

// DefaultConfigFile returns the path to default location for a config file,
// based on the underlying OS.
func DefaultConfigFile(appName string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	if u.HomeDir == "" {
		return "", errors.New("cannot determine user specific home directory")
	}
	cfgDir := DefaultConfigDir(u.HomeDir)
	return filepath.Join(cfgDir, appName, appName+".conf"), nil
}

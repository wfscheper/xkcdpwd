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

import "testing"

func Test_checkSeparatro(t *testing.T) {
	// valid
	valid := []string{"", " ", "-", ".", "_"}
	for _, sep := range valid {
		if !checkSeparator(sep) {
			t.Errorf("valid separator '%s' failed check", sep)
		}
	}
	invalid := []string{"ajfjke;ja;", "     ", "----", "....", "____", "a", "z", "A", "Z", "ü"}
	for _, sep := range invalid {
		if checkSeparator(sep) {
			t.Errorf("invalid separator '%s' passed check", sep)
		}
	}
}

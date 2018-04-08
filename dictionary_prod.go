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

// +build !test

//go:generate go-bindata -prefix internal/langs/ -o internal/langs/languages.go -pkg langs internal/langs/languages/...

package xkcdpwd

import (
	"bytes"

	"github.com/jlucktay/xkcdpwd/internal/langs"
)

// GetDict returns the dictionary associated with the language code lang.
func GetDict(lang string) *Dictionary {
	switch lang {
	case "en":
		data := langs.MustAsset("languages/en")
		return NewDictionary(bytes.NewBuffer(data))
	default:
		return nil
	}
}

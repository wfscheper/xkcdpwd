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

//go:generate go run -tags=dev languages_generate.go

package langs

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/text/language"
)

var matcher = language.NewMatcher([]language.Tag{
	language.English,
})

func GetLanguage(lang string) ([]byte, error) {
	tag, _ := language.MatchStrings(matcher, lang)
	switch tag {
	case language.English:
		file, err := Languages.Open("/en")
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return ioutil.ReadAll(file)
	default:
		return nil, fmt.Errorf("No language file found for %s", lang)
	}
}

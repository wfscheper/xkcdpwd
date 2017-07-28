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

package dictionary

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDictionary(t *testing.T) {
	is := assert.New(t)
	r := bytes.NewBufferString("word\n")
	d := NewDictionary(r)
	if is.Len(d.words, 1) {
		is.Equal([]string{"word"}, d.words)
	}
}

func TestLength(t *testing.T) {
	is := assert.New(t)
	d := Dictionary{words: []string{"word"}}
	is.Nil(d.length)
	is.Equal(1, d.Length())
	is.Equal(1, *d.length)

	d = Dictionary{words: []string{"another", "word"}}
	is.Nil(d.length)
	is.Equal(2, d.Length())
	is.Equal(2, *d.length)
}

func TestNewDictionaryComment(t *testing.T) {
	r := bytes.NewBufferString("word\n# comment\n  # comment2\nanother # inline\n")
	d := NewDictionary(r)
	if assert.Len(t, d.words, 2) {
		assert.Equal(t, []string{"word", "another"}, d.words)
	}
}

func TestNewDictionaryNoTrailingNewLine(t *testing.T) {
	r := bytes.NewBufferString("word")
	d := NewDictionary(r)
	if assert.Len(t, d.words, 1) {
		assert.Equal(t, []string{"word"}, d.words)
	}
	r = bytes.NewBufferString("word\nanother")
	d = NewDictionary(r)
	if assert.Len(t, d.words, 2) {
		assert.Equal(t, []string{"word", "another"}, d.words)
	}
}

func TestWord(t *testing.T) {
	tests := []struct {
		Index    int
		Expected string
	}{
		{
			Index:    0,
			Expected: "word",
		},
		{
			Index:    1,
			Expected: "another",
		},
		{
			Index:    -1,
			Expected: "",
		},
		{
			Index:    -15,
			Expected: "",
		},
		{
			Index:    2,
			Expected: "",
		},
		{
			Index:    303,
			Expected: "",
		},
	}
	is := assert.New(t)
	d := Dictionary{words: []string{"word", "another"}}
	for idx, test := range tests {
		actual := d.Word(test.Index)
		is.Equalf(test.Expected, actual, "%d: Expected '%s', got '%s'", idx, test.Expected, actual)
	}
}

func TestGetDict(t *testing.T) {
	is := assert.New(t)
	d := GetDict("en")
	is.Equal(d.Word(0), "able")
	is.Nil(GetDict("foo"))
}

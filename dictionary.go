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

package xkcdpwd

import (
	"bufio"
	"io"
	"strings"
)

// Dictionary wraps a word list and its length.
type Dictionary struct {
	words  []string
	length *int
}

// Length returns the number of words in the Dictionary.
func (d *Dictionary) Length() int {
	if d.length == nil {
		l := len(d.words)
		d.length = &l
	}
	return *d.length
}

// Word returns the word at index idx. If idx is less than 0, or greater than
// or equal to the number of words in the dictionary, then Word returns an
// empty string.
func (d *Dictionary) Word(idx int) string {
	if idx < 0 || idx >= d.Length() {
		return ""
	}
	return d.words[idx]
}

// NewDictionary scans r line-by-line and returns a Dictionary. Each line in r
// should be a word in the dictionary. Lines beginning with a #-character are
// considred comments and are ignored.
func NewDictionary(r io.Reader) *Dictionary {
	d := &Dictionary{words: make([]string, 0, 10)}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		w := scanner.Text()
		i := strings.IndexRune(w, '#')
		switch {
		case i < 0:
			w = strings.TrimSpace(w)
		case i > 0:
			w = strings.TrimSpace(w[:i])
		default:
			w = ""
		}
		if w != "" {
			d.words = append(d.words, w)
		}
	}
	return d
}

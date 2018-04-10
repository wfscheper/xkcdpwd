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

//go:generate go-bindata -prefix internal/langs/ -o internal/langs/languages.go -pkg langs internal/langs/languages/...

package xkcdpwd

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"

	"golang.org/x/text/language"

	"github.com/wfscheper/xkcdpwd/internal/langs"
)

// Dictionary wraps a word list and its length.
type Dictionary struct {
	words  []string
	length *int
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
			// Insert words according to length, so that we can more easily
			// filter them later
			d.words = insertByLength(d.words, w)
		}
	}
	return d
}

func insertByLength(s []string, w string) []string {
	if len(s) == 0 {
		// first word, append it
		return append(s, w)
	}
	i, j := 0, len(s)-1
	for {
		switch {
		case len(w) < len(s[i]):
			if i == 0 {
				return append([]string{w}, s...)
			}
			return insert(s, w, i)
		case len(w) > len(s[j]):
			if j == len(s)-1 {
				return append(s, w)
			}
			return insert(s, w, j)
		case i >= j:
			return append(s, w)
		default:
			i++
			j--
		}
	}
}

func insert(s []string, w string, i int) []string {
	s = append(s, "")
	copy(s[i+1:], s[i:])
	s[i] = w
	return s
}

// Entropy returns the number of bits of entropy the current dictionary configuration can support.
func (d *Dictionary) Entropy(n int) float64 {
	return float64(n) * math.Log2(float64(d.Length()))
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

// Words returns a slice of n randomly chosen words. If the dictionary cannot
// support the minimum entropy, then an error is returned.
func (d *Dictionary) Words(n int) ([]string, error) {
	if d.Entropy(n) < 30.0 {
		return nil, fmt.Errorf("dictionary cannot support more than 30 bits of entropy")
	}
	dictLengh := big.NewInt(int64(d.Length() - 1))
	words := make([]string, n, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, dictLengh)
		if err != nil {
			return nil, fmt.Errorf("cannot generate random words: %s", err)
		}
		words[i] = d.Word(int(idx.Int64()))
	}
	return words, nil
}

var matcher = language.NewMatcher([]language.Tag{
	language.English,
})

// GetDict returns the dictionary associated with the language code lang.
func GetDict(lang string) *Dictionary {
	tag, _ := language.MatchStrings(matcher, lang)
	switch tag {
	case language.English:
		data := langs.MustAsset("languages/en")
		return NewDictionary(bytes.NewBuffer(data))
	default:
		return nil
	}
}

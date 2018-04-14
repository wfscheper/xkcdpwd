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
	"sort"
	"strings"

	"golang.org/x/text/language"

	"github.com/wfscheper/xkcdpwd/internal/langs"
)

const minEntropy = 30.0

// Dictionary wraps a word list and its length.
type Dictionary struct {
	capitalize    string
	minWordLength int
	maxWordLength int
	randReader    io.Reader
	words         []string
	start         int
	stop          int
}

// NewDictionary scans r line-by-line and returns a Dictionary. Each line in r
// should be a word in the dictionary. Lines beginning with a #-character are
// considred comments and are ignored.
func NewDictionary(r io.Reader) *Dictionary {
	d := &Dictionary{words: []string{}, randReader: rand.Reader}
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
			wLength := len(w)
			if d.MaxWordLength() < wLength {
				d.SetMaxWordLength(wLength)
			}
			if d.MinWordLength() > wLength || d.MinWordLength() == 0 {
				d.SetMinWordLength(wLength)
			}
		}
	}
	// sort words according to length, so that we can more easily
	// filter them later
	sort.Slice(d.words, func(i, j int) bool {
		if len(d.words[i]) < len(d.words[j]) {
			return true
		}
		return false
	})
	return d
}

// Capitalize returns the current capitalizaton strategy.
func (d *Dictionary) Capitalize() string {
	return d.capitalize
}

// SetCapitalize sets the capitalization strategy. If the string passed in is not a recognized strategy, then 'none' is used.
func (d *Dictionary) SetCapitalize(s string) {
	switch s {
	case "all", "random", "first":
		d.capitalize = s
	default:
		d.capitalize = "none"
	}
}

// MaxWordLength returns the current max word length
func (d *Dictionary) MaxWordLength() int {
	return d.maxWordLength
}

// SetMaxWordLength updates the maximum word length of the dictionary. Values of 0 or less are taken to mean no limit.
func (d *Dictionary) SetMaxWordLength(n int) {
	d.maxWordLength = n
	d.updateStop()
}

func (d *Dictionary) updateStop() {
	if d.maxWordLength <= 0 {
		d.stop = len(d.words)
		return
	}
	end := len(d.words) - 1
	for i := range d.words {
		word := d.words[end-i]
		if len(word) <= d.maxWordLength {
			d.stop = end - i + 1
			return
		}
	}
	d.stop = 0
}

// MinWordLength returns the current minimum word length
func (d *Dictionary) MinWordLength() int {
	return d.minWordLength
}

// SetMinWordLength updates the minimum word length of the dictionary. Values of 0 or less are taken to mean no limit.
func (d *Dictionary) SetMinWordLength(n int) {
	d.minWordLength = n
	d.updateStart()
}

func (d *Dictionary) updateStart() {
	if d.minWordLength <= 0 {
		d.start = 0
		return
	}
	for idx, word := range d.words {
		if len(word) >= d.minWordLength {
			d.start = idx
			return
		}
	}
	d.start = len(d.words)
}

// Entropy returns the number of bits of entropy the current dictionary configuration can support.
func (d *Dictionary) Entropy(n int) float64 {
	return float64(n) * math.Log2(float64(d.Length()))
}

// Length returns the number of words in the Dictionary.
func (d *Dictionary) Length() int {
	if d.start >= len(d.words) || d.stop <= 0 {
		return 0
	}
	return len(d.words[d.start:d.stop])
}

// Word returns the word at index idx. If idx is less than 0, or greater than
// or equal to the number of words in the dictionary, then Word returns an
// empty string.
func (d *Dictionary) Word(idx int) string {
	if idx < 0 || d.start+idx >= d.stop {
		return ""
	}
	return d.words[d.start+idx]
}

// Passphrase returns a slice of n randomly chosen words. If the dictionary
// cannot support the minimum entropy, then an error is returned.
func (d *Dictionary) Passphrase(n int) ([]string, error) {
	if d.Entropy(n) < minEntropy {
		return nil, fmt.Errorf("dictionary cannot support more than %0.0f bits of entropy", minEntropy)
	}
	dictLengh := big.NewInt(int64(d.Length() - 1))
	words := make([]string, n, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(d.randReader, dictLengh)
		if err != nil {
			return nil, fmt.Errorf("cannot generate random words: %s", err)
		}
		word := d.words[d.start+int(idx.Int64())]
		switch d.capitalize {
		case "all":
			word = strings.ToUpper(word)
		case "first":
			word = strings.ToUpper(word[0:1]) + word[1:]
		case "random":
			for j := range word {
				choice, err := rand.Int(d.randReader, big.NewInt(int64(9)))
				if err != nil {
					return []string{}, err
				}
				if int(choice.Int64()) < 5 {
					switch j {
					case 0:
						word = strings.ToUpper(word[0:1]) + word[1:]
					case len(word) - 1:
						word = word[0:j] + strings.ToUpper(word[j:])
					default:
						word = word[0:j] + strings.ToUpper(word[j:j+1]) + word[j+1:]
					}
				}
			}
		}
		words[i] = word
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

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
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

// newDictionary returns a Dictionary of the slice words. words is assumed to
// already be sorted by word length, and the Dictionary will be configured with
// a non-random randReader.
func newDictionary(words []string) (d *Dictionary) {
	d = &Dictionary{
		randReader: constantReader(0),
		words:      words,
	}
	d.SetMaxWordLength(len(words[len(words)-1]))
	d.SetMinWordLength(len(words[0]))
	return
}

type constantReader int

func (c constantReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(c)
	}
	return len(p), nil
}

func TestNewDictionary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		data                  string
		expectedWords         []string
		expectedMinWordLength int
		expectedMaxWordLength int
	}{
		{
			data:                  "word",
			expectedWords:         []string{"word"},
			expectedMaxWordLength: 4,
			expectedMinWordLength: 4,
		},
		{
			data:                  "word\n",
			expectedWords:         []string{"word"},
			expectedMaxWordLength: 4,
			expectedMinWordLength: 4,
		},
		{
			data:                  "word\nanother",
			expectedWords:         []string{"word", "another"},
			expectedMaxWordLength: 7,
			expectedMinWordLength: 4,
		},
		{
			data:                  "word\nanother\n",
			expectedWords:         []string{"word", "another"},
			expectedMaxWordLength: 7,
			expectedMinWordLength: 4,
		},
		{
			data:                  "word\n# comment\n  # comment2\nanother # inline\n",
			expectedWords:         []string{"word", "another"},
			expectedMaxWordLength: 7,
			expectedMinWordLength: 4,
		},
	}
	for idx, test := range tests {
		r := bytes.NewBufferString(test.data)
		d := NewDictionary(r)
		t.Run(fmt.Sprint(idx+1), func(t *testing.T) {
			if reflect.DeepEqual(test.expectedWords, d.words) {
				if test.expectedMaxWordLength != d.MaxWordLength() {
					t.Errorf("expected max word length %d, got %d", test.expectedMaxWordLength, d.MaxWordLength())
				}
				if test.expectedMinWordLength != d.MinWordLength() {
					t.Errorf("expected min word length %d, got %d", test.expectedMinWordLength, d.MinWordLength())
				}
			} else {
				t.Errorf("expected %v words, got %v", test.expectedWords, d.words)
			}
		})
	}
}

func TestNewDictionaryOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"f"}, []string{"f"}},
		{[]string{"f", "b"}, []string{"f", "b"}},
		{[]string{"b", "fo"}, []string{"b", "fo"}},
		{[]string{"fo", "b"}, []string{"b", "fo"}},
		{[]string{"b", "foo"}, []string{"b", "foo"}},
		{[]string{"foo", "b"}, []string{"b", "foo"}},
		{[]string{"b", "fo", "bar"}, []string{"b", "fo", "bar"}},
		{[]string{"b", "bar", "fo"}, []string{"b", "fo", "bar"}},
		{[]string{"fo", "b", "bar"}, []string{"b", "fo", "bar"}},
		{[]string{"fo", "bar", "b"}, []string{"b", "fo", "bar"}},
		{[]string{"bar", "b", "fo"}, []string{"b", "fo", "bar"}},
		{[]string{"bar", "fo", "b"}, []string{"b", "fo", "bar"}},
		{[]string{"f", "g", "h", "i"}, []string{"f", "g", "h", "i"}},
		{[]string{"1", "12", "1234", "1234", "12345"}, []string{"1", "12", "1234", "1234", "12345"}},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprint(idx+1), func(t *testing.T) {
			input := strings.Join(test.input, "\n")
			actual := NewDictionary(bytes.NewBufferString(input))
			if !reflect.DeepEqual(test.expected, actual.words) {
				t.Errorf("expected %v, got %v", test.expected, actual.words)
			}
		})
	}
}

func TestLength(t *testing.T) {
	d := newDictionary([]string{"word"})
	if d.Length() != 1 {
		t.Error("expected length 1")
	}

	d = newDictionary([]string{"word", "another"})
	if d.Length() != 2 {
		t.Error("expected length 2")
	}
}

func TestSetCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"all", "all"},
		{"first", "first"},
		{"none", "none"},
		{"random", "random"},
		{"", "none"},
		{"fjekajfe", "none"},
		{"alll", "none"},
	}
	d := newDictionary([]string{"word", "another"})
	for idx, test := range tests {
		d.SetCapitalize(test.input)
		if d.Capitalize() != test.expected {
			t.Errorf("%d: expected '%s', got '%s'", idx+1, test.expected, d.capitalize)
		}
	}
}

func TestSetMaxWordLength(t *testing.T) {
	tests := []struct {
		maxLength      int
		expectedStop   int
		expectedLength int
	}{
		{maxLength: -10, expectedStop: 5, expectedLength: 5},
		{maxLength: -1, expectedStop: 5, expectedLength: 5},
		{maxLength: 0, expectedStop: 5, expectedLength: 5},
		{maxLength: 1, expectedStop: 1, expectedLength: 1},
		{maxLength: 2, expectedStop: 2, expectedLength: 2},
		{maxLength: 3, expectedStop: 2, expectedLength: 2},
		{maxLength: 4, expectedStop: 4, expectedLength: 4},
		{maxLength: 5, expectedStop: 5, expectedLength: 5},
		{maxLength: 6, expectedStop: 5, expectedLength: 5},
		{maxLength: 16, expectedStop: 5, expectedLength: 5},
	}
	d := newDictionary([]string{"1", "12", "1234", "1234", "12345"})
	for idx, test := range tests {
		d.SetMaxWordLength(test.maxLength)
		if test.expectedStop == d.stop {
			if test.expectedLength != d.Length() {
				t.Errorf("%d: expected length %d, got %d", idx+1, test.expectedLength, d.Length())
			}
		} else {
			t.Errorf("%d: expected stop %d, got %d", idx+1, test.expectedStop, d.stop)
		}
	}
}

func TestSetMinWordLength(t *testing.T) {
	t.Parallel()
	tests := []struct {
		minLength      int
		expectedStart  int
		expectedLength int
	}{
		{minLength: -10, expectedStart: 0, expectedLength: 5},
		{minLength: -1, expectedStart: 0, expectedLength: 5},
		{minLength: 0, expectedStart: 0, expectedLength: 5},
		{minLength: 1, expectedStart: 0, expectedLength: 5},
		{minLength: 2, expectedStart: 1, expectedLength: 4},
		{minLength: 3, expectedStart: 2, expectedLength: 3},
		{minLength: 4, expectedStart: 2, expectedLength: 3},
		{minLength: 5, expectedStart: 4, expectedLength: 1},
		{minLength: 6, expectedStart: 5, expectedLength: 0},
		{minLength: 16, expectedStart: 5, expectedLength: 0},
	}
	d := newDictionary([]string{"1", "12", "1234", "1234", "12345"})
	for idx, test := range tests {
		t.Run(fmt.Sprint(idx+1), func(t *testing.T) {
			d.SetMinWordLength(test.minLength)
			if test.expectedStart == d.start {
				if test.expectedLength != d.Length() {
					t.Errorf("expected length %d, got %d", test.expectedLength, d.Length())
				}
			} else {
				t.Errorf("expected start %d, got %d", test.expectedStart, d.start)
			}
		})
	}
}

func TestWord(t *testing.T) {
	t.Parallel()
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
	d := newDictionary([]string{"word", "another"})
	for idx, test := range tests {
		t.Run(fmt.Sprint(idx+1), func(t *testing.T) {
			actual := d.Word(test.Index)
			if test.Expected != actual {
				t.Errorf("%d: Expected '%s', got '%s'", idx, test.Expected, actual)
			}
		})
	}
}

func TestWordWithMax(t *testing.T) {
	t.Parallel()
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
			Expected: "",
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
	d := newDictionary([]string{"word", "another"})
	d.SetMaxWordLength(5)
	for idx, test := range tests {
		t.Run(fmt.Sprint(idx+1), func(t *testing.T) {
			actual := d.Word(test.Index)
			if test.Expected != actual {
				t.Errorf("%d: Expected '%s', got '%s'", idx, test.Expected, actual)
			}
		})
	}
}

func TestWordWithMin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Index    int
		Expected string
	}{
		{
			Index:    0,
			Expected: "another",
		},
		{
			Index:    1,
			Expected: "",
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
	d := newDictionary([]string{"word", "another"})
	d.SetMinWordLength(5)
	for idx, test := range tests {
		t.Run(fmt.Sprint(idx+1), func(t *testing.T) {
			actual := d.Word(test.Index)
			if test.Expected != actual {
				t.Errorf("%d: Expected '%s', got '%s'", idx, test.Expected, actual)
			}
		})
	}
}

func TestGetDict(t *testing.T) {
	d := GetDict("en")
	actual := d.Word(0)
	if d != nil {
		if "able" != actual {
			t.Errorf("expected 'able', got '%s'", actual)
		}
	}
	// matcher will try *really* hard to give a match
	d = GetDict("foo")
	if d != nil {
		if "able" != actual {
			t.Errorf("expected 'able', got '%s'", actual)
		}
	}
}

func TestCapitalizeAll(t *testing.T) {
	d := newDictionary(strings.Split(strings.Repeat("able,", 1000), ","))
	d.SetCapitalize("all")
	p, err := d.Passphrase(4)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"ABLE", "ABLE", "ABLE", "ABLE"}, p) {
		t.Errorf("expected [ABLE ABLE ABLE ABLE], got %v", p)
	}
}

func TestCapitalizeRandom(t *testing.T) {
	d := newDictionary(strings.Split(strings.Repeat("able,", 1000), ","))
	// will get all caps since our randReader only returns 0
	d.SetCapitalize("random")
	p, err := d.Passphrase(4)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"ABLE", "ABLE", "ABLE", "ABLE"}, p) {
		t.Errorf("expected [ABLE ABLE ABLE ABLE], got %v", p)
	}
}

func TestCapitalizeFirst(t *testing.T) {
	d := newDictionary(strings.Split(strings.Repeat("able,", 1000), ","))
	d.SetCapitalize("first")
	p, err := d.Passphrase(4)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"Able", "Able", "Able", "Able"}, p) {
		t.Errorf("expected [Able Able Able Able], got %v", p)
	}
}

func BenchmarkNewDictionaySorted(b *testing.B) {
	data := make([]string, 5000)
	length := 3
	wordsPerLength := len(data) / 10
	for i := 0; i < len(data); i++ {
		data[i] = strings.Repeat("1", length+(i/wordsPerLength))
	}
	benchmarkNewDictionary(b, data)
}

func BenchmarkNewDictionayReversed(b *testing.B) {
	data := make([]string, 5000)
	length := 12
	wordsPerLength := len(data) / 10
	for i := 0; i < len(data); i++ {
		data[i] = strings.Repeat("1", length-(i/wordsPerLength))
	}
	benchmarkNewDictionary(b, data)
}

func BenchmarkNewDictionayRandom(b *testing.B) {
	data := make([]string, 5000)
	for i := 0; i < len(data); i++ {
		length := rand.Intn(11) + 1
		data[i] = strings.Repeat("1", length)
	}
	benchmarkNewDictionary(b, data)
}

func benchmarkNewDictionary(b *testing.B, words []string) {
	r := bytes.NewBufferString(strings.Join(words, "\n"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewDictionary(r)
	}
}

// Copyright 2018 Andrew Merenbach
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pasc

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"github.com/merenbach/goldbug/internal/masc"
	"github.com/merenbach/goldbug/internal/stringutil"
)

// TabulaRecta holds a tabula recta.
type TabulaRecta struct {
	Strict   bool
	Caseless bool

	PtAlphabet  string
	CtAlphabet  string
	KeyAlphabet string

	DictFunc func(s string, i int) (*masc.Tableau, error)
}

// NewTabulaRecta creates a new tabula recta from multiple invocations of a MASC tableau generation function.
// NewTabulaRecta is the canonical method to generate a typical tabula recta.
func NewTabulaRecta(ptAlphabet string, keyAlphabet string, f func(s string, i int) (*masc.Tableau, error)) (*TabulaRecta, error) {
	if ptAlphabet == "" {
		ptAlphabet = Alphabet
	}

	if keyAlphabet == "" {
		keyAlphabet = ptAlphabet
	}

	t := TabulaRecta{
		PtAlphabet:  ptAlphabet,
		KeyAlphabet: keyAlphabet,
		DictFunc:    f,
	}
	return &t, nil
}

func (tr *TabulaRecta) makedictsfromfunc() (ReciprocalTable, error) {
	ptAlphabet := tr.PtAlphabet

	keyAlphabet := tr.KeyAlphabet
	if keyAlphabet == "" {
		keyAlphabet = ptAlphabet
	}

	f := tr.DictFunc
	if f == nil {
		return tr.makereciprocaltable()
	}

	m := make(map[rune]*masc.Tableau)

	keyRunes := []rune(keyAlphabet)

	if len(keyRunes) != len(keyAlphabet) {
		return nil, errors.New("Row headers must have same rune length as rows slice")
	}

	for i, r := range keyRunes {
		t, err := f(ptAlphabet, i)
		if err != nil {
			return nil, err
		}
		t.Caseless = tr.Caseless
		t.Strict = tr.Strict
		m[r] = t
	}

	return m, nil
}

// MakeTabulaRecta creates a standard Caesar shift tabula recta.
func (tr *TabulaRecta) makereciprocaltable() (ReciprocalTable, error) {
	ptAlphabet := tr.PtAlphabet

	ctAlphabet := tr.CtAlphabet
	if ctAlphabet == "" {
		ctAlphabet = ptAlphabet
	}

	keyAlphabet := tr.KeyAlphabet
	if keyAlphabet == "" {
		keyAlphabet = ptAlphabet
	}

	ctAlphabets := make([]string, utf8.RuneCountInString(tr.KeyAlphabet))
	ctAlphabetLen := utf8.RuneCountInString(ctAlphabet)

	// Cast to []rune to increase index without gaps
	for y := range ctAlphabets {
		ii := make([]int, ctAlphabetLen)
		for x := range ii {
			ii[x] = (x + y) % ctAlphabetLen
		}

		out, err := stringutil.Backpermute(ctAlphabet, ii)
		if err != nil {
			return nil, err
		}
		ctAlphabets[y] = out
	}

	m := make(map[rune]*masc.Tableau)

	keyRunes := []rune(keyAlphabet)
	if len(keyRunes) != len(keyAlphabet) {
		return nil, errors.New("Row headers must have same rune length as rows slice")
	}

	for i, r := range keyRunes {
		t, err := masc.NewTableau(ptAlphabet, func(string) (string, error) {
			return ctAlphabets[i], nil
		})
		if err != nil {
			return nil, err
		}
		t.Strict = tr.Strict
		// t.Caseless = tr.Caseless // TODO
		m[r] = t
	}

	return m, nil
}

func (tr *TabulaRecta) String() string {
	return fmt.Sprintf("%+v", map[string]interface{}{
		"k":  tr.KeyAlphabet,
		"pt": tr.PtAlphabet,
		"ct": tr.CtAlphabet,
	})
}

// Printable representation of this tabula recta.
func (tr *TabulaRecta) Printable() (string, error) {
	ptAlphabet := tr.PtAlphabet

	keyAlphabet := tr.KeyAlphabet
	if keyAlphabet == "" {
		keyAlphabet = ptAlphabet
	}

	rt, err := tr.makedictsfromfunc()
	if err != nil {
		return "", err
	}

	var b strings.Builder

	w := tabwriter.NewWriter(&b, 4, 1, 3, ' ', 0)

	formatForPrinting := func(s string) string {
		spl := strings.Split(s, "")
		return strings.Join(spl, " ")
	}

	fmt.Fprintf(w, "\t%s\n", formatForPrinting(ptAlphabet))
	for _, r := range []rune(keyAlphabet) {
		if tableau, ok := rt[r]; ok {
			out, err := tableau.Encipher(ptAlphabet)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(w, "\n%c\t%s", r, formatForPrinting(out))
		}
	}

	w.Flush()
	return b.String(), nil
}

// // Encipher a plaintext rune with a given key alphabet rune.
// // Encipher will return (-1, false) if the key rune is invalid.
// // Encipher will return (-1, true) if the key rune is valid but the message rune is not.
// // Encipher will otherwise return the transcoded rune as the first argument and true as the second.
// func (tr *TabulaRecta) encipher(r rune, k rune) (rune, bool) {
// 	c, ok := tr.pt2ct[k]
// 	if !ok {
// 		return (-1), false
// 	}

// 	if o, ok := c[r]; ok {
// 		return o, true
// 	}
// 	return (-1), true
// }

// // Decipher a ciphertext rune with a given key alphabet rune.
// // Decipher will return (-1, false) if the key rune is invalid.
// // Decipher will return (-1, true) if the key rune is valid but the message rune is not.
// // Decipher will otherwise return the transcoded rune as the first argument and true as the second.
// func (tr *TabulaRecta) decipher(r rune, k rune) (rune, bool) {
// 	c, ok := tr.ct2pt[k]
// 	if !ok {
// 		return (-1), false
// 	}

// 	if o, ok := c[r]; ok {
// 		return o, true
// 	}
// 	return (-1), true
// }

// Encipher a string.
func (tr *TabulaRecta) Encipher(s string, k string, onSuccess func(rune, rune, *[]rune)) (string, error) {
	rt, err := tr.makedictsfromfunc()
	if err != nil {
		return "", err
	}
	return rt.Encipher(s, k, onSuccess)
}

// Decipher a string.
func (tr *TabulaRecta) Decipher(s string, k string, onSuccess func(rune, rune, *[]rune)) (string, error) {
	rt, err := tr.makedictsfromfunc()
	if err != nil {
		return "", err
	}
	return rt.Decipher(s, k, onSuccess)
}

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

package vigenere

import (
	"github.com/merenbach/goldbug/internal/pasc"
)

// An AutokeyOption determines the autokey setting for the cipher.
type autokeyOption uint8

const (
	// NoAutokey signifies not to autokey the cipher.
	NoAutokey autokeyOption = iota

	// TextAutokey denotes a text-autokey mechanism.
	TextAutokey

	// KeyAutokey denotes a key-autokey mechanism.
	KeyAutokey
)

// Cipher implements a Vigenere cipher.
type Cipher struct {
	Alphabet string
	Autokey  autokeyOption
	Key      string
	Strict   bool
}

func (c *Cipher) maketableau() (*pasc.TabulaRecta, error) {
	alphabet := c.Alphabet
	if alphabet == "" {
		alphabet = pasc.Alphabet
	}

	return &pasc.TabulaRecta{
		PtAlphabet:  alphabet,
		CtAlphabet:  alphabet,
		KeyAlphabet: alphabet,
		Strict:      c.Strict,
	}, nil
}

// Encipher a message.
func (c *Cipher) Encipher(s string) (string, error) {
	t, err := c.maketableau()
	if err != nil {
		return "", err
	}
	return t.Encipher(s, c.Key, func(original rune, translated rune) rune {
		switch c.Autokey {
		case TextAutokey:
			return original
		case KeyAutokey:
			return translated
		}
		return (-1)
	})
}

// Decipher a message.
func (c *Cipher) Decipher(s string) (string, error) {
	t, err := c.maketableau()
	if err != nil {
		return "", err
	}
	return t.Decipher(s, c.Key, func(original rune, translated rune) rune {
		switch c.Autokey {
		case TextAutokey:
			return translated
		case KeyAutokey:
			return original
		}
		return (-1)
	})
}

// Tableau for encipherment and decipherment.
func (c *Cipher) tableau() (string, error) {
	t, err := c.maketableau()
	if err != nil {
		return "", err
	}
	return t.Printable()
}

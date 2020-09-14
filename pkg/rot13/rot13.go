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

package rot13

import (
	"github.com/merenbach/goldbug/internal/masc"
	"github.com/merenbach/goldbug/pkg/caesar"
)

// Shift alphabet by this constant number of characters.
const shift = 13

// Cipher implements a ROT13 cipher.
type Cipher struct {
	Caseless bool
	Strict   bool
}

func (c *Cipher) maketableau() *caesar.Cipher {
	return &caesar.Cipher{
		Alphabet: "",
		Caseless: c.Caseless,
		Shift:    shift,
		Strict:   c.Strict,
	}
}

// Encipher a message.
func (c *Cipher) Encipher(s string) (string, error) {
	return c.maketableau().Encipher(s)
}

// Decipher a message.
func (c *Cipher) Decipher(s string) (string, error) {
	return c.maketableau().Decipher(s)
}

// Tableau for encipherment and decipherment.
func (c *Cipher) Tableau() (*masc.Tableau, error) {
	return c.maketableau().Tableau()
}

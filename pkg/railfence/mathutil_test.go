// Copyright 2019 Andrew Merenbach
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

package railfence

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestAbs(t *testing.T) {
	testdata, err := ioutil.ReadFile(filepath.Join("testdata", "abs.json"))
	if err != nil {
		t.Fatal("Could not read testdata fixture:", err)
	}

	var tables []struct {
		Input  int
		Output int
	}
	if err := json.Unmarshal(testdata, &tables); err != nil {
		t.Fatal("Could not unmarshal testdata:", err)
	}

	for _, table := range tables {
		if out := abs(table.Input); out != table.Output {
			t.Errorf("Expected abs(%d)=%d, but instead got abs(%d)=%d", table.Input, table.Output, table.Input, out)
		}
	}
}

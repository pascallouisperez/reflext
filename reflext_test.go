// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reflext

import (
	. "gopkg.in/check.v1"
)

func (_ *ReflextSuite) TestString(c *C) {
	examples := map[string]string{
		"int":                                 "",
		"[2]bool":                             "",
		"[]bool":                              "",
		"[][]rune":                            "[][]int32",
		"*string":                             "",
		"map[byte]error":                      "map[uint8]error",
		"chan int":                            "",
		"chan <- int":                         "",
		"<- chan int":                         "",
		"kind[uint8]":                         "",
		"struct":                              "kind[struct]",
		"alias[chan uint8]":                   "",
		"_":                                   "",
		"{int}":                               "",
		"map[{_}]*{_}":                        "",
		"{{int | {_}}}":                       "",
		"func(int) byte":                      "func(int) uint8",
		"func(int) (int, int)":                "",
		"func(int | uint, bool) (int, int)":   "",
		"int | uint":                          "",
		"int | kind[int] | uint | kind[uint]": "",
	}
	for s, expected := range examples {
		c.Log(s)
		actual, err := parse(s)
		c.Assert(err, IsNil)
		if expected == "" {
			c.Assert(actual.String(), Equals, s)
		} else {
			c.Assert(actual.String(), Equals, expected)
		}
	}
}

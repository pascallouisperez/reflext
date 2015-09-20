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
	"reflect"

	. "gopkg.in/check.v1"
)

func (_ *ReflextSuite) TestParser_good(c *C) {
	examples := map[string]expression{
		"int":               &exact{baseTypes["int"]},
		"[2]bool":           &arrayOf{2, &exact{baseTypes["bool"]}},
		"[]bool":            &sliceOf{&exact{baseTypes["bool"]}},
		"[][]rune":          &sliceOf{&sliceOf{&exact{baseTypes["rune"]}}},
		"*string":           &ptrOf{&exact{baseTypes["string"]}},
		"map[byte]error":    &mapOf{&exact{baseTypes["byte"]}, &exact{baseTypes["error"]}},
		"chan int":          &chanOf{&exact{baseTypes["int"]}, reflect.BothDir},
		"chan <- int":       &chanOf{&exact{baseTypes["int"]}, reflect.SendDir},
		"<- chan int":       &chanOf{&exact{baseTypes["int"]}, reflect.RecvDir},
		"kind[uint8]":       &kindOf{kinds["uint8"]},
		"struct":            &kindOf{kinds["struct"]},
		"alias[chan uint8]": &aliasOf{&chanOf{&exact{baseTypes["uint8"]}, reflect.BothDir}},
		"_":                 &any{},
		"int | uint":        &firstOf{[]expression{&exact{baseTypes["int"]}, &exact{baseTypes["uint"]}}},
		"int | kind[int] | uint | kind[uint]": &firstOf{[]expression{
			&exact{baseTypes["int"]}, &kindOf{kinds["int"]},
			&exact{baseTypes["uint"]}, &kindOf{kinds["uint"]},
		}},
		"{int}":         &captureOf{&exact{baseTypes["int"]}, 0},
		"map[{_}]*{_}":  &mapOf{&captureOf{&any{}, 0}, &ptrOf{&captureOf{&any{}, 1}}},
		"{{int | {_}}}": &captureOf{&captureOf{&firstOf{[]expression{&exact{baseTypes["int"]}, &captureOf{&any{}, 2}}}, 1}, 0},
	}
	for s, expected := range examples {
		c.Log(s)
		actual, err := parse(s)
		c.Assert(err, IsNil)
		c.Assert(actual, DeepEquals, expected)
	}
}

func (_ *ReflextSuite) TestParser_func(c *C) {
	examples := map[string]expression{
		"func(int) byte": &funcOf{
			[]expression{&exact{baseTypes["int"]}},
			[]expression{&exact{baseTypes["byte"]}},
		},
		"func(int) (int, int)": &funcOf{
			[]expression{&exact{baseTypes["int"]}},
			[]expression{&exact{baseTypes["int"]}, &exact{baseTypes["int"]}},
		},
		"func(int | uint, bool) (int, int)": &funcOf{
			[]expression{&firstOf{[]expression{&exact{baseTypes["int"]}, &exact{baseTypes["uint"]}}}, &exact{baseTypes["bool"]}},
			[]expression{&exact{baseTypes["int"]}, &exact{baseTypes["int"]}},
		},
	}
	for s, expected := range examples {
		c.Log(s)
		actual, err := parse(s)
		c.Assert(err, IsNil)
		c.Assert(actual, DeepEquals, expected)
	}
}

func (_ *ReflextSuite) TestParser_concrete(c *C) {
	examples := map[string]expression{
		"%T": nil,
		// func(%T, %T) %T
		// etc.
	}
	for s, expected := range examples {
		c.Log(s)
		actual, err := parse(s)
		c.Assert(err, IsNil)
		c.Assert(actual, DeepEquals, expected)
	}
}

func (_ *ReflextSuite) TestParser_bad(c *C) {
	examples := []string{
		"int int",
		"[ int",
		"] int",
		"[w]int",
	}
	for _, s := range examples {
		_, err := parse(s)
		c.Assert(err, NotNil)
	}
}

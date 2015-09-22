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

type someInterface interface{}

func (_ *ReflextSuite) TestParser_good(c *C) {
	examples := map[string]expression{
		"int":            &exact{types["int"]},
		"[2]bool":        &arrayOf{2, &exact{types["bool"]}},
		"[]bool":         &sliceOf{&exact{types["bool"]}},
		"[][]rune":       &sliceOf{&sliceOf{&exact{types["rune"]}}},
		"*string":        &ptrOf{&exact{types["string"]}},
		"map[byte]error": &mapOf{&exact{types["byte"]}, &implements{types["error"]}},
		"chan int":       &chanOf{&exact{types["int"]}, reflect.BothDir},
		"chan <- int":    &chanOf{&exact{types["int"]}, reflect.SendDir},
		"<- chan int":    &chanOf{&exact{types["int"]}, reflect.RecvDir},
		"kind[uint8]":    &kindOf{kinds["uint8"]},
		"struct":         &kindOf{kinds["struct"]},
		"_":              &any{},
		"int | uint":     &firstOf{[]expression{&exact{types["int"]}, &exact{types["uint"]}}},
		"int | kind[int] | uint | kind[uint]": &firstOf{[]expression{
			&exact{types["int"]}, &kindOf{kinds["int"]},
			&exact{types["uint"]}, &kindOf{kinds["uint"]},
		}},
		"map[_]_ | _":   &firstOf{[]expression{&mapOf{&any{}, &any{}}, &any{}}},
		"*_ | _":        &firstOf{[]expression{&ptrOf{&any{}}, &any{}}},
		"chan _ | _":    &firstOf{[]expression{&chanOf{&any{}, reflect.BothDir}, &any{}}},
		"func() _ | _":  &firstOf{[]expression{&funcOf{nil, []expression{&any{}}}, &any{}}},
		"{int}":         &captureOf{&exact{types["int"]}, 0},
		"map[{_}]*{_}":  &mapOf{&captureOf{&any{}, 0}, &ptrOf{&captureOf{&any{}, 1}}},
		"{{int | {_}}}": &captureOf{&captureOf{&firstOf{[]expression{&exact{types["int"]}, &captureOf{&any{}, 2}}}, 1}, 0},

		// alias
		"alias[string]":     &aliasOf{&convertibleTo{types["string"]}},
		"alias[chan uint8]": &aliasOf{&chanOf{&exact{types["uint8"]}, reflect.BothDir}},

		// func
		"func(int)": &funcOf{
			[]expression{&exact{types["int"]}},
			nil,
		},
		"func() int": &funcOf{
			nil,
			[]expression{&exact{types["int"]}},
		},
		"func(int) byte": &funcOf{
			[]expression{&exact{types["int"]}},
			[]expression{&exact{types["byte"]}},
		},
		"func(int) (int, int)": &funcOf{
			[]expression{&exact{types["int"]}},
			[]expression{&exact{types["int"]}, &exact{types["int"]}},
		},
		"func(int | uint, bool) (int, int)": &funcOf{
			[]expression{&firstOf{[]expression{&exact{types["int"]}, &exact{types["uint"]}}}, &exact{types["bool"]}},
			[]expression{&exact{types["int"]}, &exact{types["int"]}},
		},

		// Concrete
		"%T":           &exact{types["int"]},
		"%T | %T | %T": &firstOf{[]expression{&exact{types["int"]}, &exact{types["bool"]}, &exact{types["string"]}}},
		"map[%T]*%T":   &mapOf{&exact{types["int"]}, &ptrOf{&exact{types["bool"]}}},
	}
	for s, expected := range examples {
		c.Log(s)

		// Parse
		actual, _, err := parse(s, 0, true, "")
		c.Assert(err, IsNil)
		c.Assert(actual, DeepEquals, expected)

		// Check round trip (expression > string > expression) is identity function
		actual2, _, err := parse(expected.String(), 0, true, "")
		c.Assert(err, IsNil)
		c.Assert(actual2, DeepEquals, expected)
	}
}

func (_ *ReflextSuite) TestParser_concreteWithInterface(c *C) {
	actual, _, err := parse("%T", reflect.TypeOf((*someInterface)(nil)).Elem())
	c.Assert(err, IsNil)
	c.Assert(actual, DeepEquals, &implements{reflect.TypeOf((*someInterface)(nil)).Elem()})
	c.Assert(actual.String(), Equals, "reflext.someInterface")
}

func (_ *ReflextSuite) TestParser_bad(c *C) {
	examples := []string{
		"int int",
		"[ int",
		"] int",
		"[w]int",
		"%t",
	}
	for _, s := range examples {
		_, _, err := parse(s)
		c.Assert(err, NotNil)
	}
}

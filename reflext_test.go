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
	"errors"
	"reflect"

	. "gopkg.in/check.v1"
)

func (_ *ReflextSuite) TestMatch(c *C) {
	examples := map[string]struct {
		matches, doesnt []interface{}
	}{
		"int": {
			matches: []interface{}{int(0)},
			doesnt:  []interface{}{int8(0)},
		},
		"error": {
			matches: []interface{}{errors.New(""), &myError{}},
			doesnt:  []interface{}{int8(0)},
		},
		"[2]int": {
			matches: []interface{}{[2]int{}},
			doesnt:  []interface{}{[3]int{}, [2]bool{}},
		},
		"[]int": {
			matches: []interface{}{[]int{}},
			doesnt:  []interface{}{[]bool{}, [2]int{}},
		},
		"*int": {
			matches: []interface{}{new(int)},
			doesnt:  []interface{}{new(bool)},
		},
		"map[string]bool": {
			matches: []interface{}{map[string]bool{}},
			doesnt:  []interface{}{map[string]string{}, map[bool]bool{}, new(bool)},
		},
		"chan int": {
			matches: []interface{}{make(chan int)},
			doesnt:  []interface{}{make(chan bool), make(chan<- int), make(<-chan int), new(int)},
		},
		"chan<- int": {
			matches: []interface{}{make(chan<- int)},
			doesnt:  []interface{}{make(chan bool), make(chan int), make(<-chan int), new(int)},
		},
		"<-chan int": {
			matches: []interface{}{make(<-chan int)},
			doesnt:  []interface{}{make(chan bool), make(chan<- int), make(chan int), new(int)},
		},
		"kind[struct]": {
			matches: []interface{}{struct{}{}, struct{ int }{0}},
			doesnt:  []interface{}{map[int]bool{}},
		},
		"alias[string]": {
			matches: []interface{}{stringAlias("")},
			doesnt:  []interface{}{""},
		},
		"alias[chan int]": {
			matches: []interface{}{make(chanIntAlias)},
			doesnt:  []interface{}{make(chan int)},
		},
		"_": {
			matches: []interface{}{0, true, "", map[int]int{}},
			doesnt:  []interface{}{},
		},
		"{int}": {
			matches: []interface{}{0},
			doesnt:  []interface{}{false},
		},
		"int | bool": {
			matches: []interface{}{0, true},
			doesnt:  []interface{}{6.7, ""},
		},
		"func(int) int": {
			matches: []interface{}{func(int) int { return 0 }},
			doesnt: []interface{}{
				func(int) (int, int) { return 0, 0 },
				func(int) bool { return true },
				func() int { return 0 },
				func(bool) int { return 0 },
			},
		},
	}
	for s, cases := range examples {
		r := MustCompile(s)
		for _, value := range cases.matches {
			c.Logf("%s matches %T", s, value)
			c.Assert(r.Match(value), Equals, true)
		}
		for _, value := range cases.doesnt {
			c.Logf("%s does not match %T", s, value)
			c.Assert(r.Match(value), Equals, false)
		}
	}
}

func (_ *ReflextSuite) TestMatchType(c *C) {
	examples := map[string]struct {
		matches, doesnt []interface{}
	}{
		"int": {
			matches: []interface{}{int(0)},
			doesnt:  []interface{}{int8(0)},
		},
		"error": {
			matches: []interface{}{errors.New(""), &myError{}},
			doesnt:  []interface{}{int8(0)},
		},
		"[2]int": {
			matches: []interface{}{[2]int{}},
			doesnt:  []interface{}{[3]int{}, [2]bool{}},
		},
		"[]int": {
			matches: []interface{}{[]int{}},
			doesnt:  []interface{}{[]bool{}, [2]int{}},
		},
		"*int": {
			matches: []interface{}{new(int)},
			doesnt:  []interface{}{new(bool)},
		},
		"map[string]bool": {
			matches: []interface{}{map[string]bool{}},
			doesnt:  []interface{}{map[string]string{}, map[bool]bool{}, new(bool)},
		},
		"chan int": {
			matches: []interface{}{make(chan int)},
			doesnt:  []interface{}{make(chan bool), make(chan<- int), make(<-chan int), new(int)},
		},
		"chan<- int": {
			matches: []interface{}{make(chan<- int)},
			doesnt:  []interface{}{make(chan bool), make(chan int), make(<-chan int), new(int)},
		},
		"<-chan int": {
			matches: []interface{}{make(<-chan int)},
			doesnt:  []interface{}{make(chan bool), make(chan<- int), make(chan int), new(int)},
		},
		"kind[struct]": {
			matches: []interface{}{struct{}{}, struct{ int }{0}},
			doesnt:  []interface{}{map[int]bool{}},
		},
		"alias[string]": {
			matches: []interface{}{stringAlias("")},
			doesnt:  []interface{}{""},
		},
		"alias[chan int]": {
			matches: []interface{}{make(chanIntAlias)},
			doesnt:  []interface{}{make(chan int)},
		},
		"_": {
			matches: []interface{}{0, true, "", map[int]int{}},
			doesnt:  []interface{}{},
		},
		"{int}": {
			matches: []interface{}{0},
			doesnt:  []interface{}{false},
		},
		"int | bool": {
			matches: []interface{}{0, true},
			doesnt:  []interface{}{6.7, ""},
		},
		"func(int) int": {
			matches: []interface{}{func(int) int { return 0 }},
			doesnt: []interface{}{
				func(int) (int, int) { return 0, 0 },
				func(int) bool { return true },
				func() int { return 0 },
				func(bool) int { return 0 },
			},
		},
	}
	for s, cases := range examples {
		r := MustCompile(s)
		for _, value := range cases.matches {
			t := reflect.TypeOf(value)
			c.Logf("%s matches %s", s, t)
			c.Assert(r.MatchType(t), Equals, true)
		}
		for _, value := range cases.doesnt {
			t := reflect.TypeOf(value)
			c.Logf("%s does not match %s", s, t)
			c.Assert(r.MatchType(t), Equals, false)
		}
	}
}

func (_ *ReflextSuite) TestFindAll(c *C) {
	examples := map[string][]struct {
		value    interface{}
		expected []reflect.Type
	}{
		"{_}": {
			{0, []reflect.Type{types["int"]}},
			{true, []reflect.Type{types["bool"]}},
		},
		"{int} | {bool}": {
			{0, []reflect.Type{types["int"], nil}},
			{true, []reflect.Type{nil, types["bool"]}},
		},
		"map[{_}]chan {_}": {
			{make(map[int]chan string), []reflect.Type{types["int"], types["string"]}},
		},
		"map[string]int | uint": {
			{make(map[string]int), []reflect.Type{}},
			{uint(0), []reflect.Type{}},
		},
		"map[string]{int | uint}": {
			{make(map[string]int), []reflect.Type{types["int"]}},
			{make(map[string]uint), []reflect.Type{types["uint"]}},
		},
		"func({_})": {
			{func(rune) {}, []reflect.Type{types["rune"]}},
		},
		"func() {_}": {
			{func() int { return 0 }, []reflect.Type{types["int"]}},
		},
		"{%T}": {
			{&myError{}, []reflect.Type{reflect.TypeOf(&myError{})}},
		},
	}
	for s, cases := range examples {
		r := MustCompile(s, reflect.TypeOf((*error)(nil)).Elem())
		for _, eg := range cases {
			c.Logf("%s with %T", s, eg.value)
			captures, ok := r.FindAll(eg.value)
			c.Assert(ok, Equals, true)
			c.Assert(captures, DeepEquals, eg.expected)
		}
	}
}

func (_ *ReflextSuite) TestString(c *C) {
	r := MustCompile("map[int]bool")
	c.Assert(r.String(), Equals, "map[int]bool")
}

func (_ *ReflextSuite) TestString_expression(c *C) {
	examples := map[string]string{
		"int":                                 "",
		"[2]bool":                             "",
		"[]bool":                              "",
		"[][]rune":                            "[][]int32",
		"*string":                             "",
		"map[byte]error":                      "map[uint8]error",
		"chan int":                            "",
		"chan<- int":                          "",
		"<-chan int":                          "",
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
		actual, _, err := parse(s)
		c.Assert(err, IsNil)
		if expected == "" {
			c.Assert(actual.String(), Equals, s)
		} else {
			c.Assert(actual.String(), Equals, expected)
		}
	}
}

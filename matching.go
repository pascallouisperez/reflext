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
	"strconv"
	"strings"
)

type expression interface {
	Match(reflect.Type) bool
	String() string
}

type exact struct {
	typ reflect.Type
}

func (m *exact) Match(typ reflect.Type) bool {
	return m.typ == typ
}

func (m *exact) String() string {
	return m.typ.String()
}

type implements struct {
	typ reflect.Type
}

func (m *implements) Match(typ reflect.Type) bool {
	return typ.Implements(m.typ)
}

func (m *implements) String() string {
	return m.typ.String()
}

type sliceOf struct {
	exp expression
}

func (m *sliceOf) Match(typ reflect.Type) bool {
	if typ.Kind() != reflect.Slice {
		return false
	}
	return m.exp.Match(typ.Elem())
}

func (m *sliceOf) String() string {
	return "[]" + m.exp.String()
}

type arrayOf struct {
	size int
	exp  expression
}

func (m *arrayOf) Match(typ reflect.Type) bool {
	if typ.Kind() != reflect.Array {
		return false
	}
	if m.size != typ.Len() {
		return false
	}
	return m.exp.Match(typ.Elem())
}

func (m *arrayOf) String() string {
	return "[" + strconv.Itoa(m.size) + "]" + m.exp.String()
}

type ptrOf struct {
	exp expression
}

func (m *ptrOf) Match(typ reflect.Type) bool {
	if typ.Kind() != reflect.Ptr {
		return false
	}
	return m.exp.Match(typ.Elem())
}

func (m *ptrOf) String() string {
	return "*" + m.exp.String()
}

type mapOf struct {
	key   expression
	value expression
}

func (m *mapOf) Match(typ reflect.Type) bool {
	if typ.Kind() != reflect.Map {
		return false
	}
	if !m.key.Match(typ.Key()) {
		return false
	}
	if !m.value.Match(typ.Elem()) {
		return false
	}
	return true
}

func (m *mapOf) String() string {
	return "map[" + m.key.String() + "]" + m.value.String()
}

type chanOf struct {
	exp expression
	dir reflect.ChanDir
}

func (m *chanOf) Match(typ reflect.Type) bool {
	if typ.Kind() != reflect.Chan {
		return false
	}
	if typ.ChanDir() != m.dir {
		return false
	}
	return m.exp.Match(typ.Elem())
}

func (m *chanOf) String() string {
	switch m.dir {
	case reflect.BothDir:
		return "chan " + m.exp.String()
	case reflect.SendDir:
		return "chan<- " + m.exp.String()
	case reflect.RecvDir:
		return "<-chan " + m.exp.String()
	}
	panic("unreachable")
}

type funcOf struct {
	arguments []expression
	returns   []expression
}

func (m *funcOf) Match(typ reflect.Type) bool {
	if typ.Kind() != reflect.Func {
		return false
	}
	if len(m.arguments) != typ.NumIn() {
		return false
	}
	if len(m.returns) != typ.NumOut() {
		return false
	}
	for i, a := range m.arguments {
		if !a.Match(typ.In(i)) {
			return false
		}
	}
	for i, r := range m.returns {
		if !r.Match(typ.Out(i)) {
			return false
		}
	}
	return true
}

func (m *funcOf) String() string {
	var (
		args []string
		rets []string
		ret  string
	)
	for _, a := range m.arguments {
		args = append(args, a.String())
	}
	if len(m.returns) == 1 {
		ret = m.returns[0].String()
	} else {
		for _, r := range m.returns {
			rets = append(rets, r.String())
		}
		ret = "(" + strings.Join(rets, ", ") + ")"
	}
	return "func(" + strings.Join(args, ", ") + ") " + ret
}

type kindOf struct {
	kind reflect.Kind
}

func (m *kindOf) Match(typ reflect.Type) bool {
	return m.kind == typ.Kind()
}

func (m *kindOf) String() string {
	return "kind[" + m.kind.String() + "]"
}

type aliasOf struct {
	exp expression
}

func (m *aliasOf) Match(typ reflect.Type) bool {
	if typ.Name() == "" {
		return false
	}
	return m.exp.Match(typ)
}

func (m *aliasOf) String() string {
	return "alias[" + m.exp.String() + "]"
}

type convertibleTo struct {
	typ reflect.Type
}

func (m *convertibleTo) Match(typ reflect.Type) bool {
	return m.typ != typ && typ.ConvertibleTo(m.typ)
}

func (m *convertibleTo) String() string {
	return m.typ.String()
}

type any struct{}

func (m *any) Match(reflect.Type) bool {
	return true
}

func (m *any) String() string {
	return "_"
}

type firstOf struct {
	exps []expression
}

func (m *firstOf) Match(typ reflect.Type) bool {
	for _, exp := range m.exps {
		if exp.Match(typ) {
			return true
		}
	}
	return false
}

func (m *firstOf) String() string {
	var e []string
	for _, exp := range m.exps {
		e = append(e, exp.String())
	}
	return strings.Join(e, " | ")
}

type captureOf struct {
	exp   expression
	index int
}

func (m *captureOf) Match(typ reflect.Type) bool {
	return m.exp.Match(typ)
}

func (m *captureOf) String() string {
	return "{" + m.exp.String() + "}"
}

// Assert that all matches implement the expression interface.
var _ = []expression{
	&exact{},
	&implements{},
	&sliceOf{},
	&arrayOf{},
	&ptrOf{},
	&mapOf{},
	&funcOf{},
	&kindOf{},
	&aliasOf{},
	&convertibleTo{},
	&any{},
	&firstOf{},
	&captureOf{},
}

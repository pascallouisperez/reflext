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
)

type expression interface {
	Match(reflect.Type) bool
}

type exact struct {
	exp reflect.Type
}

func (m *exact) Match(reflect.Type) bool { return false }

type sliceOf struct {
	exp expression
}

func (m *sliceOf) Match(reflect.Type) bool { return false }

type arrayOf struct {
	size int
	exp  expression
}

func (m *arrayOf) Match(reflect.Type) bool { return false }

type ptrOf struct {
	exp expression
}

func (m *ptrOf) Match(reflect.Type) bool { return false }

type mapOf struct {
	key   expression
	value expression
}

func (m *mapOf) Match(reflect.Type) bool { return false }

type chanOf struct {
	exp expression
	dir reflect.ChanDir
}

func (m *chanOf) Match(reflect.Type) bool { return false }

type funcOf struct {
	arguments []expression
	returns   []expression
}

func (m *funcOf) Match(reflect.Type) bool { return false }

type kindOf struct {
	kind reflect.Kind
}

func (m *kindOf) Match(reflect.Type) bool { return false }

type aliasOf struct {
	exp expression
}

func (m *aliasOf) Match(reflect.Type) bool { return false }

type any struct{}

func (m *any) Match(reflect.Type) bool { return false }

type firstOf struct {
	exps []expression
}

func (m *firstOf) Match(reflect.Type) bool { return false }

type captureOf struct {
	exp   expression
	index int
}

func (m *captureOf) Match(reflect.Type) bool { return false }

// Assert that all matches implement the expression interface.
var _ = []expression{
	&exact{},
	&sliceOf{},
	&arrayOf{},
	&ptrOf{},
	&mapOf{},
	&funcOf{},
	&kindOf{},
	&aliasOf{},
	&any{},
	&firstOf{},
	&captureOf{},
}

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
	"testing"
)

type someStruct struct{}

func exampleFunc(someInterface, *someStruct, map[int]bool) error { return nil }

func BenchmarkWithReflext(b *testing.B) {
	r := MustCompile(
		"func(%T, *struct, map[int]bool) error",
		reflect.TypeOf((*someInterface)(nil)).Elem())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ok := r.Match(exampleFunc); !ok {
			b.Error("must match on every iteration")
		}
	}
}

func BenchmarkByHand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if ok := matchByHand(exampleFunc); !ok {
			b.Error("must match on every iteration")
		}
	}
}

var (
	errorType         = reflect.TypeOf((*error)(nil)).Elem()
	someInterfaceType = reflect.TypeOf((*someInterface)(nil)).Elem()
	intType           = types["int"]
	boolType          = types["bool"]
)

func matchByHand(v interface{}) bool {
	value := reflect.ValueOf(v)
	typ := value.Type()
	if typ.Kind() != reflect.Func {
		return false
	}

	// Arguments
	if typ.NumIn() != 3 {
		return false
	}
	// arg0: someInterface
	if typ.In(0).Kind() != reflect.Interface && typ.In(0).Implements(someInterfaceType) {
		return false
	}
	// arg1: *struct
	if typ.In(1).Kind() != reflect.Ptr && typ.In(1).Elem().Kind() != reflect.Struct {
		return false
	}
	// arg2: map[int]bool
	if typ.In(2).Kind() != reflect.Map && typ.In(2).Elem().Key() != intType && typ.In(2).Elem() != boolType {
		return false
	}

	// Return
	if typ.NumOut() != 1 {
		return false
	}
	// return0: error
	if typ.Out(0).Kind() != reflect.Interface && typ.Out(0).Implements(errorType) {
		return false
	}

	return true
}

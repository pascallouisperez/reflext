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

type Reflext struct {
	expression
	numGroup int
}

func Compile(s string, args ...interface{}) (*Reflext, error) {
	exp, numGroup, err := parse(s, args)
	if err != nil {
		return nil, err
	}
	return &Reflext{exp, numGroup}, nil
}

func MustCompile(s string, args ...interface{}) *Reflext {
	r, err := Compile(s, args...)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Reflext) Match(value interface{}) bool {
	return r.expression.Match(reflect.TypeOf(value), nil)
}

func (r *Reflext) FindAll(value interface{}) ([]reflect.Type, bool) {
	captured := make([]reflect.Type, r.numGroup, r.numGroup)
	if ok := r.expression.Match(reflect.TypeOf(value), &captured); !ok {
		return nil, false
	}
	return captured, true
}

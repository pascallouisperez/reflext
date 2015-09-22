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
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

var types = map[string]reflect.Type{
	"bool":       reflect.TypeOf(true),
	"byte":       reflect.TypeOf(byte(0)),
	"complex64":  reflect.TypeOf(complex64(0)),
	"complex128": reflect.TypeOf(complex128(0)),
	"error":      reflect.TypeOf(new(error)).Elem(),
	"float32":    reflect.TypeOf(float32(0)),
	"float64":    reflect.TypeOf(float64(0)),
	"int":        reflect.TypeOf(int(0)),
	"int8":       reflect.TypeOf(int8(0)),
	"int16":      reflect.TypeOf(int16(0)),
	"int32":      reflect.TypeOf(int32(0)),
	"int64":      reflect.TypeOf(int64(0)),
	"rune":       reflect.TypeOf(rune(0)),
	"string":     reflect.TypeOf(string("")),
	"uint":       reflect.TypeOf(uint(0)),
	"uint8":      reflect.TypeOf(uint8(0)),
	"uint16":     reflect.TypeOf(uint16(0)),
	"uint32":     reflect.TypeOf(uint32(0)),
	"uint64":     reflect.TypeOf(uint64(0)),
	"uintptr":    reflect.TypeOf(uintptr(0)),
}

var kinds = map[string]reflect.Kind{
	"bool":       reflect.Bool,
	"int":        reflect.Int,
	"int8":       reflect.Int8,
	"int16":      reflect.Int16,
	"int32":      reflect.Int32,
	"int64":      reflect.Int64,
	"uint":       reflect.Uint,
	"uint8":      reflect.Uint8,
	"uint16":     reflect.Uint16,
	"uint32":     reflect.Uint32,
	"uint64":     reflect.Uint64,
	"uintptr":    reflect.Uintptr,
	"float32":    reflect.Float32,
	"float64":    reflect.Float64,
	"complex64":  reflect.Complex64,
	"complex128": reflect.Complex128,
	"array":      reflect.Array,
	"chan":       reflect.Chan,
	"func":       reflect.Func,
	"interface":  reflect.Interface,
	"map":        reflect.Map,
	"ptr":        reflect.Ptr,
	"slice":      reflect.Slice,
	"string":     reflect.String,
	"struct":     reflect.Struct,
}

var stop = map[string]bool{
	"":  true,
	")": true,
	",": true,
	"|": true,
}

type token struct {
	text   string
	offset int
}

type parser struct {
	tokens []token
	index  int

	args      []interface{}
	argsIndex int

	group int
}

func parse(expr string, args ...interface{}) (expression, int, error) {
	p := &parser{
		tokens: tokenize(expr),
		index:  0,
		args:   args,
	}
	exp, ok := p.parseExp()
	if !ok || p.index != len(p.tokens) {
		return nil, 0, fmt.Errorf("unable to parse %s", expr)
	}
	return exp, p.group, nil
}

func (p *parser) parseExpList() ([]expression, bool) {
	var exps []expression
	if text, ok := p.peek(); !ok {
		return nil, false
	} else if stop[text] {
		return exps, true
	}
	exp, ok := p.parseExp()
	if !ok {
		return nil, false
	}
	exps = append(exps, exp)
	var done bool
	for !done {
		if text, ok := p.peek(); !ok {
			done = true
		} else if text == "," {
			p.next()
			exp, ok := p.parseExp()
			if !ok {
				return nil, false
			}
			exps = append(exps, exp)
		} else {
			done = true
		}
	}
	return exps, true
}

func (p *parser) parseExp() (expression, bool) {
	var exps []expression
	exp, ok := p.parseSubExp()
	if !ok {
		return nil, false
	}
	exps = append(exps, exp)
	var done bool
	for !done {
		if text, ok := p.peek(); !ok {
			done = true
		} else if text == "|" {
			p.next()
			exp, ok := p.parseSubExp()
			if !ok {
				return nil, false
			}
			exps = append(exps, exp)
		} else {
			done = true
		}
	}
	if len(exps) == 1 {
		return exps[0], true
	} else {
		return &firstOf{exps}, true
	}
}

func (p *parser) parseSubExp() (expression, bool) {
	text, ok := p.next()
	if !ok {
		return nil, false
	}
	switch text {

	case "[":
		text, ok := p.next()
		if !ok {
			return nil, false
		}
		if text == "]" {
			exp, ok := p.parseExp()
			if !ok {
				return nil, false
			}
			return &sliceOf{exp}, true
		} else if size, err := strconv.ParseInt(text, 10, 32); err == nil {
			if !p.consume("]") {
				return nil, false
			}
			exp, ok := p.parseExp()
			if !ok {
				return nil, false
			}
			return &arrayOf{int(size), exp}, true
		}
		return nil, false

	case "*":
		exp, ok := p.parseSubExp()
		if !ok {
			return nil, false
		}
		return &ptrOf{exp}, true

	case "_":
		return &any{}, true

	case "{":
		group := p.group
		p.group++
		exp, ok := p.parseExp()
		if !ok {
			return nil, false
		}
		if ok := p.consume("}"); !ok {
			return nil, false
		}
		return &captureOf{exp, group}, true

	case "map":
		if ok := p.consume("["); !ok {
			return nil, false
		}
		exp1, ok := p.parseExp()
		if !ok {
			return nil, false
		}
		if ok := p.consume("]"); !ok {
			return nil, false
		}
		exp2, ok := p.parseSubExp()
		if !ok {
			return nil, false
		}
		return &mapOf{exp1, exp2}, true

	case "kind":
		if ok := p.consume("["); !ok {
			return nil, false
		}
		k, ok := p.next()
		if !ok {
			return nil, false
		}
		kind, ok := kinds[k]
		if !ok {
			return nil, false
		}
		if ok := p.consume("]"); !ok {
			return nil, false
		}
		return &kindOf{kind}, true

	case "struct":
		return &kindOf{reflect.Struct}, true

	case "alias":
		if !p.consume("[") {
			return nil, false
		}
		exp, ok := p.parseExp()
		if !ok {
			return nil, false
		}
		if !p.consume("]") {
			return nil, false
		}
		if e, ok := exp.(*exact); ok {
			return &aliasOf{&convertibleTo{e.typ}}, true
		} else {
			return &aliasOf{exp}, true
		}

	case "chan":
		dir := reflect.BothDir
		if text, ok := p.peek(); !ok {
			return nil, false
		} else if text == "<" {
			p.next()
			if ok := p.consume("-"); !ok {
				return nil, false
			}
			dir = reflect.SendDir
		}
		exp, ok := p.parseSubExp()
		if !ok {
			return nil, false
		}
		return &chanOf{exp, dir}, true

	case "<":
		if ok := p.consume("-"); !ok {
			return nil, false
		}
		if ok := p.consume("chan"); !ok {
			return nil, false
		}
		exp, ok := p.parseSubExp()
		if !ok {
			return nil, false
		}
		return &chanOf{exp, reflect.RecvDir}, true

	case "func":
		if ok := p.consume("("); !ok {
			return nil, false
		}
		argsExp, ok := p.parseExpList()
		if !ok {
			return nil, false
		}
		if ok := p.consume(")"); !ok {
			return nil, false
		}
		var returnsExp []expression
		if text, _ := p.peek(); text == "(" {
			p.next()
			returnsExp, ok = p.parseExpList()
			if !ok {
				return nil, false
			}
			if ok := p.consume(")"); !ok {
				return nil, false
			}
		} else if !stop[text] {
			returnExp, ok := p.parseSubExp()
			if !ok {
				return nil, false
			}
			returnsExp = append(returnsExp, returnExp)
		}
		return &funcOf{argsExp, returnsExp}, true

	case "%":
		argsIndex := p.argsIndex
		p.argsIndex++
		if ok := p.consume("T"); !ok {
			return nil, false
		}
		if len(p.args) <= argsIndex {
			return nil, false
		}
		return exactOrImplements(reflect.TypeOf(p.args[argsIndex])), true

	default:
		if typ, ok := types[text]; ok {
			return exactOrImplements(typ), true
		}
		return nil, false

	}
	panic("unreachable")
}

func (p *parser) next() (string, bool) {
	if len(p.tokens) <= p.index {
		return "", false
	}
	tok := p.tokens[p.index]
	p.index++
	return tok.text, true
}

func (p *parser) peek() (string, bool) {
	if len(p.tokens) <= p.index {
		return "", false
	}
	tok := p.tokens[p.index]
	return tok.text, true
}

func (p *parser) consume(match string) bool {
	if text, ok := p.next(); !ok {
		return false
	} else {
		return text == match
	}
}

func tokenize(expr string) []token {
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	var tok rune
	var tokens []token
	for tok != scanner.EOF {
		tok = s.Scan()
		text := s.TokenText()
		if text == "" {
			return tokens
		} else {
			tokens = append(tokens, token{
				text:   text,
				offset: s.Pos().Offset,
			})
		}
	}
	panic("unreachable")
}

func exactOrImplements(typ reflect.Type) expression {
	if typ.Kind() == reflect.Interface {
		return &implements{typ}
	} else {
		return &exact{typ}
	}
}

package main

import (
	"fmt"
	"strings"
)

type Symbol interface {
	String() string
	Value() any
	Type() string
}

// ----------------------------
// StructSymbol
// ----------------------------

type StructSymbol struct {
	Environment *Env
	TypeName    string
}

func (s *StructSymbol) String() string {
	return "Struct: " + s.Environment.String()
}

func (s *StructSymbol) Value() any  { return s.Environment }
func (s *StructSymbol) Type() string { return s.TypeName }

// ----------------------------
// FuncSymbol
// ----------------------------

type FuncSymbol struct {
	Body   *BlockNode
	Params []string
	TypeName string
	Env    *Env
	NaviteFn func(e *Evaluator, self *Env, args []any, pos Position) any
}

func (f *FuncSymbol) String() string { return "!!Function!!" }
func (f *FuncSymbol) Value() any     { return f }
func (f *FuncSymbol) Type() string   { return f.TypeName }

// ----------------------------
// VarSymbol
// ----------------------------

type VarSymbol struct {
	value any
	typeName string
}

func (v *VarSymbol) String() string { return fmt.Sprintf("Var: %v", v.value) }
func (v *VarSymbol) Value() any     { return v.value }
func (v *VarSymbol) Type() string   { return v.typeName }

// ----------------------------
// Env (Environment)
// ----------------------------

type Env struct {
	Type    string
	Symbols map[string]Symbol
	Parent  *Env
}

func NewEnv(parent *Env, envType string) *Env {
	return &Env{
		Type:    envType,
		Symbols: make(map[string]Symbol),
		Parent:  parent,
	}
}

func (e *Env) String() string {
	var lines []string
	for name, val := range e.Symbols {
		lines = append(lines, fmt.Sprintf("Name: %s; Type(%T): %v", name, val, val.String()))
	}
	return strings.Join(lines, "\n")
}

// ----------------------------
// Symbol Management
// ----------------------------

func (e *Env) SymbolExistsInCurrent(name string) bool {
	_, exists := e.Symbols[name]
	return exists
}

func (e *Env) SymbolExists(name string) bool {
	for env := e; env != nil; env = env.Parent {
		if _, exists := env.Symbols[name]; exists {
			return true
		}
		if env.Type == "block" && env.Parent != nil && env.Parent.Type == "global" {
			break
		}
	}
	return false
}

func (e *Env) IsSymbolFunc(name string) bool {
	for env := e; env != nil; env = env.Parent {
		if sym, ok := env.Symbols[name]; ok {
			_, isFunc := sym.(*FuncSymbol)
			return isFunc
		}
	}
	return false
}

func (e *Env) IsSymbolArray(name string) bool {
	for env := e; env != nil; env = env.Parent {
		if sym, ok := env.Symbols[name]; ok {
			_, isArray := sym.Value().([]any)
			return isArray
		}
	}
	return false
}

func (e *Env) AddVarSymbol(name, typ string, val any) {
	e.Symbols[name] = &VarSymbol{value: val, typeName: typ}
}

func (e *Env) AddFuncSymbol(name string, params []string, body *BlockNode, newEnv *Env) {
	e.Symbols[name] = &FuncSymbol{
		Body:     body,
		Params:   params,
		TypeName: "nan",
		Env:      newEnv,
	}
}

func (e *Env) AddStructSymbol(name string, structEnv *Env) {
	e.Symbols[name] = &StructSymbol{
		TypeName:    name,
		Environment: structEnv,
	}
}

func (e *Env) AddStructMethod(
	structName,
	methodName string,
	params []string,
	body []Node,
	nativeFn any) int {
	for env := e; env != nil; env = env.Parent {
		sym, exists := env.Symbols[structName]
		if !exists {
			continue
		}

		structSym, ok := sym.(*StructSymbol)
		if !ok {
			return -1
		}

		if structSym.Environment.Symbols == nil {
			structSym.Environment.Symbols = make(map[string]Symbol)
		}

		if _, exists := structSym.Environment.Symbols[methodName]; exists {
			return -1
		}

		funcSym := &FuncSymbol{
			Body:     &BlockNode{Statements: body},
			Params:   params,
			TypeName: "nan",
			Env:      NewEnv(structSym.Environment, "method"),
		}

		structSym.Environment.Symbols[methodName] = funcSym
		env.Symbols[structName] = structSym
		return 0
	}
	return 1
}

func (e *Env) GetSymbolType(name string) string {
	for env := e; env != nil; env = env.Parent {
		if sym, ok := env.Symbols[name]; ok {
			return sym.Type()
		}
	}
	return ""
}

func (e *Env) FindSymbol(name string) any {
	for env := e; env != nil; env = env.Parent {
		if sym, ok := env.Symbols[name]; ok {
			return sym
		}
		if env.Type == "block" && env.Parent != nil && env.Parent.Type == "global" {
			break
		}
	}
	return nil
}

func (e *Env) FindStructSymbol(name string) *Env {
	for env := e; env != nil; env = env.Parent {
		sym, ok := env.Symbols[name]
		if !ok {
			continue
		}
		if structSym, ok := sym.(*StructSymbol); ok {
			if structEnv, ok := structSym.Value().(*Env); ok {
				return structEnv
			}
		}
	}
	return nil
}

func (e *Env) FindStructMember(structName, memberName string) Symbol {
	for env := e; env != nil; env = env.Parent {
		sym, ok := env.Symbols[structName]
		if !ok {
			continue
		}
		structSym, ok := sym.(*StructSymbol)
		if !ok {
			return nil
		}
		if member, ok := structSym.Environment.Symbols[memberName]; ok {
			return member
		}
		return nil
	}
	return nil
}

func (e *Env) RemoveSymbol(name string) {
	delete(e.Symbols, name)
}

func (e *Env) UpdateSymbol(name string, newValue any, newType string) {
	for env := e; env != nil; env = env.Parent {
		if sym, ok := env.Symbols[name]; ok {
			if env.IsSymbolFunc(name) {
				break
			}
			varType := sym.Type()
			if len(newType) != 0 {
				varType = newType
			}
			env.Symbols[name] = &VarSymbol{
				value: newValue,
				typeName: varType,
			}
		}
	}
}

func (e *Env) ChangeArrayValue(name string, index int, value any) {
	for env := e; env != nil; env = env.Parent {
		sym, ok := env.Symbols[name]
		if !ok {
			continue
		}
		varSym, ok := sym.(*VarSymbol)
		if !ok {
			return
		}
		arr, ok := varSym.Value().([]any)
		if !ok || index < 0 || index >= len(arr) {
			return
		}
		arr[index] = value
		env.Symbols[name] = &VarSymbol{value: arr, typeName: varSym.typeName}
		return
	}
}

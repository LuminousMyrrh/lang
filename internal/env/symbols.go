package env

import (
	"fmt"
	"lang/internal/eval"
	"lang/internal/parser"
)

// ----------------------------
// StructSymbol
// ----------------------------

type StructSymbol struct {
	Environment *Env
	TypeName    string
}

func (s *StructSymbol) Name() string {
	return "Struct: " + s.Environment.String()
}

func (s *StructSymbol) Value() any   { return s.Environment }
func (s *StructSymbol) Type() string { return s.TypeName }

// ----------------------------
// FuncSymbol
// ----------------------------

type FuncSymbol struct {
	Body       *parser.BlockNode
	Params     []string
	TypeName   string
	Env        *Env
	NativeFunc func(e *eval.Evaluator, self *Env, args []any, pos parser.Position) any
}

func (f *FuncSymbol) Name() string { return "!!Function!!" }
func (f *FuncSymbol) Value() any     { return f }
func (f *FuncSymbol) Type() string   { return f.TypeName }

// ----------------------------
// VarSymbol
// ----------------------------

type VarSymbol struct {
	value    any
	typeName string
}

func (v *VarSymbol) Name() string { return fmt.Sprintf("Var value: %v\n", v.value) }
func (v *VarSymbol) Value() any     { return v.value }
func (v *VarSymbol) Type() string   { return v.typeName }

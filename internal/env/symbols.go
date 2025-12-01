package env

import (
	"lang/internal/core"
	"lang/internal/parser"
)

// ----------------------------
// StructSymbol
// ----------------------------

type StructSymbol struct {
	Environment *Env
	TypeName    string
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
	NativeFunc func(e core.Evaluator, self *Env, args []any, pos parser.Position) any
}

func (f *FuncSymbol) Value() any     { return f }
func (f *FuncSymbol) Type() string   { return f.TypeName }

// ----------------------------
// VarSymbol
// ----------------------------

type VarSymbol struct {
	value    any
	typeName string
}

func (v *VarSymbol) Value() any     { return v.value }
func (v *VarSymbol) Type() string   { return v.typeName }

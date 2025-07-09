package main

type Symbol interface {
	String() string
	Value() any
}

type FuncSymbol struct {
	Body *BlockNode
	Params []string
	Type string
	Enf *Env
}

func (f FuncSymbol) String() string {
	return ""
}

func (f FuncSymbol) Value() any {
	return f
}

type VarSymbol struct {
	Val any
	Type string
}

func (v VarSymbol) Value() any {
	return v
}

func (v VarSymbol) String() string {
	return ""
}

type Env struct {
	Symbols map[string]Symbol
	Parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env {
		Symbols: make(map[string]Symbol),
		Parent: parent,
	}
}

func (e *Env) AddVarSymbol(name string, typ string, val any) {
	e.Symbols[name] = VarSymbol{Val: val, Type: typ}
}

func (e *Env) AddFuncSymbol(
	name string,
	params []string,
	body *BlockNode,
	NewEnv *Env) {

	e.Symbols[name] = FuncSymbol{
		Body: body,
		Params: params,
		Type: "",
		Enf: NewEnv,
	}
}

func (e *Env) FindSymbol(name string) any {
	env := e
	for env != nil {
		if sym, ok := env.Symbols[name]; ok {
			return sym.Value()
		}
		env = env.Parent
	}
	return nil
}

func (e *Env) RemoveSymbol(name string) {
	delete(e.Symbols, name)
}

func (e *Env) UpdateSymbol(name string, newValue any) {
	env := e
	for env != nil {
		if _, ok := env.Symbols[name]; ok {
			//env.Symbols[name] = Symbol{Value: newNode, Type: sym.Type}
			break
		}
		env = env.Parent
	}
}

package env

type Symbol interface {
	Name() string
	Value() any
	Type() string
}

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


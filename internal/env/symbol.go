package env

import "lang/internal/core"

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

func (e *Env) GetSymbolType(name string) string {
	for env := e; env != nil; env = env.Parent {
		if sym, ok := env.Symbols[name]; ok {
			return sym.Type()
		}
	}
	return ""
}

func (e *Env) FindSymbol(name string) core.Symbol {
	for env := e; env != nil; env = env.Parent {
		if symValue, ok := env.Symbols[name]; ok {
			return symValue
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

func (e *Env) FindStructMember(structName, memberName string) core.Symbol {
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


package env

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

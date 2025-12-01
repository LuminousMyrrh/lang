package env

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

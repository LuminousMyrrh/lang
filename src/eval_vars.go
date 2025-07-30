package main

import ("fmt")

func (e *Evaluator) evalVarDef(stmt *VarDefNode) any {
	if e.currentEnv.SymbolExists(stmt.Name) {
		e.genError(fmt.Sprintf(
			"Var '%s' already exists", stmt.Name),
			stmt.Position)
		return nil
	}
	value := e.eval(stmt.Value)
	var_type := e.resolveType(value, stmt.Position)
	e.currentEnv.AddVarSymbol(stmt.Name, var_type, value)
	return value
}

func (e *Evaluator) evalNil(stmt *NilNode) any {
	return nilValue{}
}

func (e *Evaluator) evalIdentifier(id *IdentifierNode) any {
    i := e.currentEnv.FindSymbol(id.Name)
	if i == nil {
		e.genError(fmt.Sprintf(
			"Unknown identifier '%s'", id.Name),
			id.Position)
		return nil
	}
	return i.(Symbol).Value()
}

func (e *Evaluator) evalAssignment(a *AssignmentNode) any {
    switch target := a.Name.(type) {
	case *IdentifierNode:
		if !e.currentEnv.SymbolExists(target.Name) {
			e.genError(fmt.Sprintf(
				"Variable '%s' does not exist",
				target.Name),
				a.Position)
			return nil
		}

		value := e.eval(a.Value)
		existing := e.currentEnv.FindSymbol(target.Name)
		if existing == nil {
			e.genError(fmt.Sprintf(
				"Variable '%s' not found", target.Name),
				target.Position,
				)
			return nil
		}

		sym, ok := existing.(Symbol)
		if !ok {
			e.genError(fmt.Sprintf(
					"'%s' is not a valid symbol",
					target.Name),
				target.Position,
				)
			return nil
		}

		currentVal := sym.Value()

		switch a.Op {
		case "+=":
			switch curr := currentVal.(type) {
			case int:
				if val, ok := value.(int); ok {
					e.currentEnv.UpdateSymbol(target.Name, curr+val)
					return curr + val
				}
			case string:
				if val, ok := value.(string); ok {
					e.currentEnv.UpdateSymbol(target.Name, curr+val)
					return curr + val
				}
			}
			e.genError(fmt.Sprintf(
				"Unsupported types for '+=': %T and %T", currentVal, value),
				target.Position,
				)
			return nil

		case "-=":
			curr, ok1 := currentVal.(int)
			val, ok2 := value.(int)
			if ok1 && ok2 {
				e.currentEnv.UpdateSymbol(target.Name, curr-val)
				return curr - val
			}
			e.genError(fmt.Sprintf(
				"Unsupported types for '-=': %T and %T", currentVal, value),
				target.Position,
				)
			return nil

		case "=":
			e.currentEnv.UpdateSymbol(target.Name, value)
			return value

		default:
			e.genError(fmt.Sprintf(
				"Unsupported assignment operator: '%s'", a.Op),
				target.Position,
				)
			return nil
		}

    case *StructMethodCall:
        // Struct field assignment: obj.field = ...
        if !target.IsField || len(target.Args) > 0 {
            e.genError(
				"Assignment target must be a field access, not a method call",
				target.Position,
				)
            return nil
        }
        // Caller must be an identifier (e.g., self, point)
        callerIdent, ok := target.Caller.(*IdentifierNode)
        if !ok {
            e.genError(
				"Struct field assignment target must be an identifier",
				target.Position,
				)
            return nil
        }
        sym := e.currentEnv.FindSymbol(callerIdent.Name)
        varSym, ok := sym.(*VarSymbol)
        if !ok {
            e.genError(fmt.Sprintf(
				"'%s' is not a variable", callerIdent.Name),
				target.Position,
				)
            return nil
        }
        instanceEnv, ok := varSym.Value().(*Env)
        if !ok {
            e.genError(fmt.Sprintf(
				"'%s' is not a struct instance", callerIdent.Name),
				target.Position,
				)
            return nil
        }
        if !instanceEnv.SymbolExistsInCurrent(target.MethodName) {
            e.genError(fmt.Sprintf(
					"Field '%s' does not exist in struct '%s'",
					target.MethodName,
					callerIdent.Name),
				target.Position,
				)
            return nil
        }
        value := e.eval(a.Value)
        instanceEnv.UpdateSymbol(target.MethodName, value)
		return value

	case *ArrayAccessNode: {
		// Unwrap the identifier
		arrNameNode := target.Target
		indexValue := e.eval(target.Index)
		indexInt, ok := indexValue.(int)
		if !ok {
			e.genError("Array index must be an integer",
				target.Position,
				)
			return nil
		}
		arrIdent, ok := arrNameNode.(*IdentifierNode)
		if !ok {
			e.genError(
				"Array assignment target must be an identifier",
				target.Position,
				)
			return nil
		}
		arrSym := e.currentEnv.FindSymbol(arrIdent.Name)
		varSym, ok := arrSym.(*VarSymbol)
		if !ok {
			e.genError(fmt.Sprintf(
				"Variable '%s' is not a VarSymbol",
				arrIdent.Name),
				target.Position,
				)
			return nil
		}
		arr, ok := varSym.Value().([]any)
		if !ok {
			e.genError(fmt.Sprintf(
					"Variable '%s' is not an array. Got: %T",
					arrIdent.Name,
					varSym.Value()),
				target.Position,
				)
			return nil
		}
		if indexInt < 0 {
			e.genError("Negative array index", target.Position)
			return nil
		}

		value := e.eval(a.Value)
		if indexInt < len(arr) {
			// Normal case: overwrite
			arr[indexInt] = value
		} else if indexInt == len(arr) {
			// Append/grow array by one position
			arr = append(arr, value)
		} else {
			// Trying to set beyond the next element: error
			e.genError(fmt.Sprintf(
				"Index %d is out of range. You can insert only at len(arr)=%d",
				indexInt, len(arr)),
				target.Position,
				)
			return nil
		}

		// Store the grown array back in the environment
		e.currentEnv.UpdateSymbol(arrIdent.Name, arr)
		return value
	}
    default:
        e.genError("Invalid assignment target", a.Position)
        return nil
    }
}

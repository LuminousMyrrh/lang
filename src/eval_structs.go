package main

import ("fmt")

func (e *Evaluator) evalStructDef(stmt *StructDefNode) any {
	structEnv := NewEnv(e.currentEnv, stmt.Name)

	for _, field := range stmt.Fields {
		if field.Value != nil {
			value := e.eval(field.Value)
			structEnv.AddVarSymbol(
				field.Name,
				e.resolveType(value, stmt.Position),
				value)
			continue
		}
		structEnv.AddVarSymbol(
			field.Name,
			"nil",
			nil)
	}

	e.currentEnv.AddStructSymbol(
		stmt.Name,
		structEnv,
		)

	return structEnv
}

func (e *Evaluator) evalStructMethodDef(stmt *StructMethodDef) any {
	structEnv := e.currentEnv.FindStructSymbol(stmt.StructName)
	if structEnv == nil {
		e.genError(fmt.Sprintf(
			"Struct '%s' doesn't exists", stmt.StructName),
			stmt.Position)
		return nil
	}

	e.currentEnv.AddStructMethod(
		stmt.StructName,
		stmt.MethodName,
		stmt.Parameters,
		stmt.Body.Statements,
		)

	return structEnv
}

func (e *Evaluator) evalStructInit(stmt *StructInitNode) any {
    sym := e.currentEnv.FindSymbol(stmt.Name)
    structSym, ok := sym.(*StructSymbol)
    if !ok {
        e.genError(fmt.Sprintf(
			"Struct type '%s' not found", stmt.Name),
			stmt.Position)
        return nil
    }

    instanceEnv := NewEnv(structSym.Environment, stmt.Name)

	for fieldName, fieldSym := range structSym.Environment.Symbols {
		if varSym, ok := fieldSym.(*VarSymbol); ok {
			instanceEnv.AddVarSymbol(fieldName, varSym.Type(), varSym.Value())
		}
	}

    for _, fieldAssign := range stmt.InitFields {
        assign, ok := fieldAssign.(*AssignmentNode)
        if !ok {
            e.genError(
				"Invalid field assignment in struct initialization",
				stmt.Position)
            return nil
        }
        val := e.eval(assign.Value)
		name, ok := assign.Name.(*IdentifierNode)
		if !instanceEnv.SymbolExistsInCurrent(name.Name) {
			e.genError(fmt.Sprintf(
				"Field '%s' is not defined in struct '%s'",
				name.Name,
				stmt.Name),
				assign.Position)
			return nil
		}

		if !ok {
			e.genError(
				"Struct field assignment must use an identifier",
				stmt.Position,
				)
			return nil
		}
        instanceEnv.UpdateSymbol(name.Name,
			val, e.resolveType(val, assign.Position))
    }

    return instanceEnv
}

func (e *Evaluator) evalStructMethodCall(structEnv *Env,
	methodName string,
	args []any,
	pos Position) any {
    if structEnv.Parent == nil {
        e.genError(fmt.Sprintf(
			"Struct type environment for method '%s' not found",
			methodName),
			pos,
			)
        return nil
    }
    methodSym, ok := structEnv.Parent.Symbols[methodName]
    if !ok {
        e.genError(fmt.Sprintf(
			"Method '%s' not found in struct",
			methodName),
			pos)
        return nil
    }
    method, ok := methodSym.(*FuncSymbol)
    if !ok {
        e.genError(fmt.Sprintf("'%s' is not a method", methodName), pos)
        return nil
    }
    if len(args) != len(method.Params) {
        e.genError(fmt.Sprintf(
			"Method '%s' expects %d args, got %d",
			methodName,
			len(method.Params),
			len(args)), pos)
        return nil
    }
    callEnv := NewEnv(structEnv, "function")
    callEnv.AddVarSymbol("self", structEnv.Type, structEnv)
    for i, val := range args {
        callEnv.AddVarSymbol(method.Params[i],
			e.resolveType(val, method.Body.Position), val)
    }
    prevEnv := e.currentEnv
    e.currentEnv = callEnv
    result := e.evalFuncBlock(method.Body)
    e.currentEnv = prevEnv
    if ret, ok := result.(returnValue); ok {
        return ret.value
    }
    return result
}

func (e *Evaluator) evalStructMemberAccess(stmt *StructMethodCall) any {
    // Get the variable holding the struct instance
	callerIdent, ok := stmt.Caller.(*IdentifierNode)
	if !ok {
		e.genError(fmt.Sprintf(
			"Struct method call: Caller must be an identifier '%T'",
			stmt.Caller,
			), stmt.Position)
		return nil
	}
	sym := e.currentEnv.FindSymbol(callerIdent.Name)
    varSym, ok := sym.(*VarSymbol)
    if !ok {
        e.genError(fmt.Sprintf("'%s' is not a variable",
			callerIdent.Name),
			stmt.Position)
        return nil
    }

    instanceEnv, ok := varSym.Value().(*Env)
    if !ok {
        e.genError(fmt.Sprintf(
			"Variable '%s' is not a struct instance",
			callerIdent.Name),
			stmt.Position)
        return nil
    }

    if stmt.IsField {
        if fieldSym, ok := instanceEnv.Symbols[stmt.MethodName]; ok {
            return fieldSym.Value()
        }
        e.genError(fmt.Sprintf("Field '%s' not found in struct instance",
			stmt.MethodName),
			stmt.Position)
        return nil
    }

    // Otherwise: method call
    argValues := make([]any, len(stmt.Args))
    for i, arg := range stmt.Args {
        argValues[i] = e.eval(arg)
    }

    return e.evalStructMethodCall(
		instanceEnv, stmt.MethodName, argValues, stmt.Position)
}

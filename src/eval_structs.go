package main

import ("fmt")

func (e *Evaluator) initBuiltintClasses() {
	stringEnv := NewEnv(nil, "string")
	stringEnv.AddVarSymbol(
		"value",
		"string",
		nil)
	e.currentEnv.AddStructSymbol("string", stringEnv)
	

	intEnv := NewEnv(nil, "int")
	intEnv.AddVarSymbol(
		"value",
		"int",
		nil)
	e.currentEnv.AddStructSymbol("int", intEnv)

	floatEnv := NewEnv(nil, "float")
	floatEnv.AddVarSymbol(
		"value",
		"float",
		nil)
	e.currentEnv.AddStructSymbol("float", floatEnv)
}

func (e *Evaluator) initBuiltintMethods() int {
	stringEnv := e.currentEnv.FindStructSymbol("string")
	if stringEnv != nil {
		stringEnv.Symbols["substring"] = &FuncSymbol{
			NaviteFn: stringSubstring,
			TypeName: "string",
		}
		stringEnv.Symbols["capitalize"] = &FuncSymbol{
			NaviteFn: stringCapitalize,
			TypeName: "string",
		}
		stringEnv.Symbols["contains"] = &FuncSymbol{
			NaviteFn: stringContains,
			TypeName: "string",
		}
		stringEnv.Symbols["empty"] = &FuncSymbol{
			NaviteFn: stringEmpty,
			TypeName: "string",
		}
		stringEnv.Symbols["isDigit"] = &FuncSymbol{
			NaviteFn: stringIsDigit,
			TypeName: "string",
		}
		stringEnv.Symbols["isAlph"] = &FuncSymbol{
			NaviteFn: stringIsAlph,
			TypeName: "string",
		}
	} else {
		return -1
	}
	return 0
}

func (e *Evaluator) evalStructDef(stmt *StructDefNode) any {
	if (e.currentEnv.SymbolExists(stmt.Name)) {
		e.genError(fmt.Sprintf(
			"Class '%s' already exists", stmt.Name), stmt.Position)
		return nil
	}
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

	if e.currentEnv.AddStructMethod(
		stmt.StructName,
		stmt.MethodName,
		stmt.Parameters,
		stmt.Body.Statements,
		nil,
		) == -1 {
		e.genError(fmt.Sprintf(
			"Method '%s' already exists in class '%s'",
			stmt.MethodName, stmt.StructName), stmt.Position)
		return nil
	}

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

func (e *Evaluator) evalStructMethodCall(
	self *Env,
	methodName string,
	args []any,
	pos Position) any {

    if self.Parent == nil {
        e.genError(fmt.Sprintf(
			"Struct type environment for method '%s' not found",
			methodName),
			pos,
			)
        return nil
    }
    methodSym, ok := self.Parent.Symbols[methodName]
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

	if method.NaviteFn != nil {
		return method.NaviteFn(e, self, args, pos)
	}

	if len(args) != len(method.Params) {
		e.genError(fmt.Sprintf(
			"Method '%s' expects %d args, got %d",
			methodName,
			len(method.Params),
			len(args)), pos)
        return nil
    }
    callEnv := NewEnv(self, "function")
    callEnv.AddVarSymbol("self", self.Type, self)

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
    // Evaluate caller expression, expecting a struct instance Env
    callerValue := e.eval(stmt.Caller)

    instanceEnv, ok := callerValue.(*Env)
    if !ok {
        e.genError(fmt.Sprintf(
           "Caller is not a struct instance but %T", callerValue),
           stmt.Position)
        return nil
    }

    if stmt.IsField {
        if fieldSym, ok := instanceEnv.Symbols[stmt.MethodName]; ok {
            return fieldSym.Value()
        }
        e.genError(fmt.Sprintf("Field '%s' not found in struct instance",
            stmt.MethodName), stmt.Position)
        return nil
    }

    // Method call: evaluate arguments
    argValues := make([]any, len(stmt.Args))
    for i, arg := range stmt.Args {
        argValues[i] = e.eval(arg)
    }

    return e.evalStructMethodCall(
        instanceEnv, stmt.MethodName, argValues, stmt.Position)
}

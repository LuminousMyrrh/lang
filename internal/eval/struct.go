package eval

import (
	"fmt"
	"lang/internal/env"
	"lang/internal/parser"
)

func (e *Evaluator) evalStructDef(stmt *parser.StructDefNode) any {
	if e.currentEnv.SymbolExists(stmt.Name) {
		e.genError(fmt.Sprintf(
			"Class '%s' already exists", stmt.Name), stmt.Position)
		return nil
	}
	structEnv := env.NewEnv(e.currentEnv, stmt.Name)

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
			nilValue{})
	}

	e.currentEnv.AddStructSymbol(
		stmt.Name,
		structEnv,
	)

	return structEnv
}

func (e *Evaluator) evalStructMethodDef(stmt *parser.StructMethodDef) any {
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

func (e *Evaluator) evalStructInit(stmt *parser.StructInitNode) any {
	sym := e.currentEnv.FindSymbol(stmt.Name)
	structSym, ok := sym.(*env.StructSymbol)
	if !ok {
		e.genError(fmt.Sprintf(
			"Struct type '%s' not found", stmt.Name),
			stmt.Position)
		return nil
	}

	instanceEnv := env.NewEnv(structSym.Environment, stmt.Name)

	for fieldName, fieldSym := range structSym.Environment.Symbols {
		if varSym, ok := fieldSym.(*env.VarSymbol); ok {
			instanceEnv.AddVarSymbol(fieldName, varSym.Type(), varSym.Value())
		}
	}

	for _, fieldAssign := range stmt.InitFields {
		assign, ok := fieldAssign.(*parser.AssignmentNode)
		if !ok {
			e.genError(
				"Invalid field assignment in struct initialization",
				stmt.Position)
			return nil
		}
		val := e.eval(assign.Value)
		name, ok := assign.Name.(*parser.IdentifierNode)
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
	self *env.Env,
	methodName string,
	args []any,
	pos parser.Position) any {

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
	method, ok := methodSym.(*env.FuncSymbol)
	if !ok {
		e.genError(fmt.Sprintf("'%s' is not a method", methodName), pos)
		return nil
	}

	if method.NativeFunc != nil {
		return method.NativeFunc(e, self, args, pos)
	}

	if len(args) != len(method.Params) {
		e.genError(fmt.Sprintf(
			"Method '%s' expects %d args, got %d",
			methodName,
			len(method.Params),
			len(args)), pos)
		return nil
	}
	callEnv := env.NewEnv(self, "function")
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

func (e *Evaluator) evalStructMemberAccess(stmt *parser.StructMethodCall) any {
	// Evaluate caller expression, expecting a struct instance Env
	callerValue := e.eval(stmt.Caller)

	instanceEnv, ok := callerValue.(*env.Env)
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

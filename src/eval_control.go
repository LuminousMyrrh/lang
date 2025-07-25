package main

func (e *Evaluator) evalBlock(block *BlockNode) any {
	blockEnv := NewEnv(e.currentEnv, "block")
	prevEnv := e.currentEnv
	e.currentEnv = blockEnv
	var result any
	for _, stmt := range block.Statements {
		result = e.eval(stmt)
		if _, ok := stmt.(*ReturnNode); ok {
			e.currentEnv = prevEnv
			return result
		}
	}
	e.currentEnv = prevEnv
	return result
}

func (e *Evaluator) evalIf(stmt *IfNode) any {
	cond := e.evalCondition(stmt.Condition)
	if cond == nil {
		return nil
	}
	if cond == true && stmt.ThenBranch != nil {
		return e.eval(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return e.eval(stmt.ElseBranch)
	} else {
		return 0
	}
}

func (e *Evaluator) evalWhile(stmt *WhileNode) any {
	var result any
	for {
		cond := e.eval(stmt.Condition)
		b, ok := cond.(bool)
		if !ok {
			e.genError("While loop condition should return bool",
				stmt.Position)
			return nil
		}
		if !b {
			break
		}
		bodyResult := e.evalLoopBlock(stmt.Body)
		if _, isBreak := bodyResult.(breakSignal); isBreak {
			break
		}
		if ret, isReturn := bodyResult.(returnValue); isReturn {
			return ret
		}
		result = bodyResult
	}
	return result
}

func (e *Evaluator) evalFor(stmt *ForNode) any {
	var result any
	switch v := stmt.Init.(type) {
	case *VarDefNode: {
		if v != nil {
			val := e.eval(v.Value)
			e.currentEnv.AddVarSymbol(
					v.Name,
					e.resolveType(val, v.Position),
					val)
		}
	}
	case *AssignmentNode: {
		if v != nil {
			e.eval(v)
		}
	}
	}
	for {
		cond := e.eval(stmt.Condition)
		b, ok := cond.(bool)
		if !ok {
			e.genError("For loop condition should return bool",
				stmt.Position)
			return nil
		}
		if !b {
			break
		}
		bodyResult := e.evalLoopBlock(stmt.Body)
		if _, isBreak := bodyResult.(breakSignal); isBreak {
			break
		}
		if ret, isRet := bodyResult.(returnValue); isRet {
			return ret
		}
		result = bodyResult
		e.eval(stmt.Post)
	}
	return result
}

func (e *Evaluator) evalReturn(ret *ReturnNode) any {
	val := e.eval(ret.Value)
	return returnValue{val}
}

func (e *Evaluator) evalLoopBlock(block *BlockNode) any {
	blockEnv := NewEnv(e.currentEnv, "block")
	prevEnv := e.currentEnv
	e.currentEnv = blockEnv

	var result any
	for _, stmt := range block.Statements {
		result := e.eval(stmt)
		if _, ok := stmt.(*BreakNode); ok {
			e.currentEnv = prevEnv
			return breakSignal{}
		}
		if ret, ok := result.(returnValue); ok {
			e.currentEnv = prevEnv
			return ret
		}
	}
	e.currentEnv = prevEnv
	return result
}


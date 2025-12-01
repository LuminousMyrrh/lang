package eval

import (
	"lang/internal/core"
	"lang/internal/env"
	"lang/internal/parser"
)

func (e *Evaluator) evalBlock(block *parser.BlockNode) any {
	blockEnv := env.NewEnv(e.currentEnv, "block")
	prevEnv := e.currentEnv
	e.currentEnv = blockEnv
	var result any
	for _, stmt := range block.Statements {
		result = e.EvalNode(stmt)
		if _, ok := result.(core.ReturnValue); ok || result != nil {
			e.currentEnv = prevEnv
			return result
		}
	}
	e.currentEnv = prevEnv
	return result
}

func (e *Evaluator) evalIf(stmt *parser.IfNode) any {
	cond := e.evalCondition(stmt.Condition)
	if cond == nil {
		e.GenError("Condition doesn't exist", stmt.Position)
		return nil
	}
	if cond == true && stmt.ThenBranch != nil {
		res := e.EvalNode(stmt.ThenBranch)
		if ret, ok := res.(core.ReturnValue); ok {
			return ret
		}
		return res
	} else if stmt.ElseBranch != nil {
		res := e.EvalNode(stmt.ElseBranch)
		if ret, ok := res.(core.ReturnValue); ok {
			return ret
		}
		return res
	} else {
		return nil
	}
}

func (e *Evaluator) evalWhile(stmt *parser.WhileNode) any {
	var result any
	for {
		cond := e.EvalNode(stmt.Condition)
		b, ok := cond.(bool)
		if !ok {
			e.GenError("While loop condition should return bool",
				stmt.Position)
			return nil
		}
		if !b {
			break
		}
		bodyResult := e.evalLoopBlock(stmt.Body)
		if _, isBreak := bodyResult.(core.BreakSignal); isBreak {
			break
		}
		if ret, isReturn := bodyResult.(core.ReturnValue); isReturn {
			return ret
		}
		result = bodyResult
	}
	return result
}

func (e *Evaluator) evalFor(stmt *parser.ForNode) any {
	var result any
	switch v := stmt.Init.(type) {
	case *parser.VarDefNode:
		{
			if v != nil {
				val := e.EvalNode(v.Value)
				e.currentEnv.AddVarSymbol(
					v.Name,
					e.resolveType(val, v.Position),
					val)
			}
		}
	case *parser.AssignmentNode:
		{
			if v != nil {
				e.EvalNode(v)
			}
		}
	}
	for {
		cond := e.EvalNode(stmt.Condition)
		b, ok := cond.(bool)
		if !ok {
			e.GenError("For loop condition should return bool",
				stmt.Position)
			return nil
		}
		if !b {
			break
		}
		bodyResult := e.evalLoopBlock(stmt.Body)
		if bodyResult == nil {
			return nil
		}
		if _, isBreak := bodyResult.(core.BreakSignal); isBreak {
			break
		}
		if ret, isRet := bodyResult.(core.ReturnValue); isRet {
			return ret
		}
		result = bodyResult
		e.EvalNode(stmt.Post)
	}
	return result
}

func (e *Evaluator) evalReturn(ret *parser.ReturnNode) any {
	val := e.EvalNode(ret.Value)
	return core.ReturnValue{Value: val}
}

func (e *Evaluator) evalLoopBlock(block *parser.BlockNode) any {
	blockEnv := env.NewEnv(e.currentEnv, "block")
	prevEnv := e.currentEnv
	e.currentEnv = blockEnv

	var result any = core.NilValue{}
	for _, stmt := range block.Statements {
		result := e.EvalNode(stmt)
		if _, ok := stmt.(*parser.BreakNode); ok {
			e.currentEnv = prevEnv
			return core.BreakSignal{}
		}
		if ret, ok := result.(core.ReturnValue); ok {
			e.currentEnv = prevEnv
			return ret
		}
	}
	e.currentEnv = prevEnv
	return result
}

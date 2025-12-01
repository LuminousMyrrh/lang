package eval

import (
	"fmt"
	"lang/internal/core"
	"lang/internal/parser"
)

func (e *Evaluator) evalArray(arr *parser.ArrayNode) any {
	var values []any

	for _, el := range arr.Elements {
		values = append(values, e.EvalNode(el))
	}

	return values
}

func (e *Evaluator) evalArrayAccess(stmt *parser.ArrayAccessNode) any {
	// Recursively evaluate the target to get the array value
	arr := unwrapBuiltinValue(e.EvalNode(stmt.Target))
	if arr == nil {
		e.GenError("Array does not exist", stmt.Position)
		return nil
	}

	// Evaluate the index expression
	idx := unwrapBuiltinValue(e.EvalNode(stmt.Index))
	i, ok := idx.(int)
	if !ok {
		e.GenError("Array index must be an integer", stmt.Position)
		return nil
	}

	// Check if arr is actually an array
	switch a := arr.(type) {
	case []any:
		{
			if i < 0 || i >= len(a) {
				e.GenError(fmt.Sprintf("Index %d out of bounds", i),
					stmt.Position)
				return nil
			}
			return a[i]
		}
	case string:
		{
			runes := []rune(a)
			if i < 0 || i >= len(runes) {
				e.GenError(fmt.Sprintf("Index %d out of bounds", i), stmt.Position)
				return nil
			}
			return e.createString(string(runes[i]))
		}
	default:
		e.GenError(fmt.Sprintf(
			"Target has incorrect type. It should be array: %T",
			a,
		), stmt.Position)
		return nil
	}

}

func (e *Evaluator) evalArrayAssign(stmt *parser.ArrayAssign) any {
	// stmt.Target is an ArrayAccessNode (possibly nested)
	// Descend to the parent array and index
	var parent any
	var index int

	// Unwind nested ArrayAccessNodes to get to the parent array and final index
	target := stmt.Target
	for {
		if access, ok := target.(*parser.ArrayAccessNode); ok {
			// If the target itself is an ArrayAccessNode, keep going
			if inner, ok := access.Target.(*parser.ArrayAccessNode); ok {
				target = inner
				continue
			}
			// Now, access.Target is the base IdentifierNode
			parent = e.EvalNode(access.Target)
			idx := e.EvalNode(access.Index)
			index, ok = idx.(int)
			if !ok {
				e.GenError("Array index must be an integer",
					stmt.Position)
				return nil
			}
			break
		} else {
			e.GenError("Invalid array assignment target", stmt.Position)
			return nil
		}
	}

	// parent should be an array
	arr, ok := parent.([]any)
	if !ok {
		e.GenError(fmt.Sprintf(
			"Target has incorrect type. It should be array: %T",
			parent,
		), stmt.Position)
		return nil
	}
	if index < 0 || index >= len(arr) {
		e.GenError(fmt.Sprintf("Index %d out of bounds", index),
			stmt.Position)
		return nil
	}

	// Assign the value
	val := e.EvalNode(stmt.Value)
	arr[index] = val
	e.currentEnv.UpdateSymbol(stmt.Target.String(),
		arr, e.resolveType(arr, stmt.Position))

	return core.NilValue{}
}

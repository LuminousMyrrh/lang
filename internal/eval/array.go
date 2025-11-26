package eval

import (
	"fmt"
)

func (e *Evaluator) evalArray(arr *ArrayNode) any {
	var values []any

	for _, el := range arr.Elements {
		values = append(values, e.eval(el))
	}

	return values;
}

func (e *Evaluator) evalArrayAccess(stmt *ArrayAccessNode) any {
	// Recursively evaluate the target to get the array value
	arr := unwrapBuiltinValue(e.eval(stmt.Target))
	if arr == nil {
		e.genError("Array does not exist", stmt.Position)
		return nil
	}

	// Evaluate the index expression
	idx := unwrapBuiltinValue(e.eval(stmt.Index))
	i, ok := idx.(int)
	if !ok {
		e.genError("Array index must be an integer", stmt.Position)
		return nil
	}

	// Check if arr is actually an array
	switch a := arr.(type) {
	case []any: {
		if i < 0 || i >= len(a) {
			e.genError(fmt.Sprintf("Index %d out of bounds", i),
				stmt.Position)
			return nil
		}
		return a[i]
	}
	case string: {
		runes := []rune(a)
		if i < 0 || i >= len(runes) {
			e.genError(fmt.Sprintf("Index %d out of bounds", i), stmt.Position)
			return nil
		}
		return e.createString(string(runes[i]))
	}
	default:
		e.genError(fmt.Sprintf(
			"Target has incorrect type. It should be array: %T",
			a,
			), stmt.Position)
		return nil
	}

}

func (e *Evaluator) evalArrayAssign(stmt *ArrayAssign) any {
	// stmt.Target is an ArrayAccessNode (possibly nested)
	// Descend to the parent array and index
	var parent any
	var index int

	// Unwind nested ArrayAccessNodes to get to the parent array and final index
	target := stmt.Target
	for {
		if access, ok := target.(*ArrayAccessNode); ok {
			// If the target itself is an ArrayAccessNode, keep going
			if inner, ok := access.Target.(*ArrayAccessNode); ok {
				target = inner
				continue
			}
			// Now, access.Target is the base IdentifierNode
			parent = e.eval(access.Target)
			idx := e.eval(access.Index)
			index, ok = idx.(int)
			if !ok {
				e.genError("Array index must be an integer",
					stmt.Position)
				return nil
			}
			break
		} else {
			e.genError("Invalid array assignment target", stmt.Position)
			return nil
		}
	}

	// parent should be an array
	arr, ok := parent.([]any)
	if !ok {
		e.genError(fmt.Sprintf(
			"Target has incorrect type. It should be array: %T",
			parent,
			), stmt.Position)
		return nil
	}
	if index < 0 || index >= len(arr) {
		e.genError(fmt.Sprintf("Index %d out of bounds", index),
			stmt.Position)
		return nil
	}

	// Assign the value
	val := e.eval(stmt.Value)
	arr[index] = val
	e.currentEnv.UpdateSymbol(stmt.Target.String(),
		arr, e.resolveType(arr, stmt.Position))

	return nilValue{}
}

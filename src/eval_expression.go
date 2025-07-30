package main

import (
	"fmt"
)

func (e *Evaluator) evalBinary(expr *BinaryOpNode) any {
	left := unwrapBuiltinValue(e.eval(expr.Left))
	if ret, ok := left.(returnValue); ok {
		left = ret.value
	}
	right := unwrapBuiltinValue(e.eval(expr.Right))
	if ret, ok := right.(returnValue); ok {
		right = ret.value
	}

    switch expr.Op {
    case "+":
        switch l := left.(type) {
        case int:
            switch r := right.(type) {
            case int:
                return l + r
            case float64:
                return float64(l) + float64(r)
            default:
                e.genError(
                    "Right operand of '+' must be int or float64 if left is int",
                    expr.Position)
                return nil
            }
        case float64:
            switch r := right.(type) {
            case int:
                return float64(l) + float64(r)
            case float64:
                return l + r
            default:
                e.genError(
                    "Right operand of '+' must be int or float64 if left is float64",
                    expr.Position)
                return nil
            }
        case string:
            r, ok := right.(string)
            if !ok {
                e.genError(
                    "Right operand of '+' must be string if left is string",
                    expr.Position)
                return nil
            }
            return l + r // String concatenation
        default:
            e.genError("Unsupported type for '+' operator",
                expr.Position)
            return nil
        }
    case "-", "*", "/":
        // Convert both operands to float64 if either is float64,
        // else use int arithmetic
        var lFloat, rFloat float64
        var lInt, rInt int
        lIsInt := false
        rIsInt := false

        switch l := left.(type) {
        case int:
            lInt = l
            lIsInt = true
            lFloat = float64(l)
        case float64:
            lFloat = l
        default:
            e.genError(fmt.Sprintf(
                "Left operand of '%s' must be int or float64", expr.Op),
                expr.Position)
            return nil
        }

        switch r := right.(type) {
        case int:
            rInt = r
            rIsInt = true
            rFloat = float64(r)
        case float64:
            rFloat = r
        default:
            e.genError(fmt.Sprintf(
                "Right operand of '%s' must be int or float64", expr.Op),
                expr.Position)
            return nil
        }

        if lIsInt && rIsInt {
            // Both int, do integer math
            switch expr.Op {
            case "-":
                return lInt - rInt
            case "*":
                return lInt * rInt
            case "/":
                if rInt == 0 {
                    e.genError("Division by zero!", expr.Position)
                    return nil
                }
                return lInt / rInt
            }
        } else {
            // Float math
            switch expr.Op {
            case "-":
                return lFloat - rFloat
            case "*":
                return lFloat * rFloat
            case "/":
                if rFloat == 0 {
                    e.genError("Division by zero!", expr.Position)
                    return nil
                }
                return lFloat / rFloat
            }
        }
    default:
        return e.evalLogical(expr)
    }
    e.genError("Unknown", expr.Position)
    return nil
}

func (e *Evaluator) evalUnary(node *UnaryOpNode) any {
    value := unwrapBuiltinValue(e.eval(node.Expr))
    if _, ok := value.(nilValue); ok {
        e.genError("Value with unary shouldn't be nil!", node.Position)
        return nil
    }
    switch node.Op {
    case "++", "--":
        if node.Op == "++" {
            return e.handlePostfixPP(node)
        } else {
            return e.handlePostfixMM(node)
        }
    case "-":
        // Negation for numbers int or float64
        switch v := value.(type) {
        case int:
            return -v
        case float64:
            return -v
        }
        e.genError(fmt.Sprintf(
            "Unary '-' not supported for type %T", value),
            node.Position)
        return nil
    case "!":
        // Logical NOT for booleans
        if v, ok := value.(bool); ok {
            return !v
        }
        e.genError(fmt.Sprintf(
            "Unary '!' not supported for type %T", value),
            node.Position)
        return nil
    default:
        return nil
    }
}

func (e *Evaluator) evalCondition(condition Node) any {
	switch t := condition.(type) {
	case *UnaryOpNode:
		return e.evalUnary(t)
	case *BinaryOpNode:
		return e.evalBinary(t)
	case *FunctionCallNode:
		return e.evalFunctionCall(t)
	case *IdentifierNode:
		return e.evalIdentifier(t)
	case *LiteralNode:
		return e.evalLiteral(t)
	case *NilNode:
		return nil
	case *StructMethodCall:
		return e.evalStructMemberAccess(t)
	default:
		e.genError(fmt.Sprintf(
			"Unsupported type: %T", condition),
			Position{-1, -1})
		return nil
	}
}

func (e *Evaluator) evalLiteral(lit *LiteralNode) any {
	if e.resolveType(lit.Value, lit.Position) == "string" {
		stringEnv := e.currentEnv.FindStructSymbol("string")
		if stringEnv == nil {
			return nil
		}
		instEnv := NewEnv(stringEnv, "string")
		instEnv.AddVarSymbol("value", "string", lit.Value)
		return instEnv
	}
	return lit.Value
}

func (e *Evaluator) evalTrue() any {
	return true
}

func (e *Evaluator) evalFalse() any {
	return false
}

func (e *Evaluator) evalLogical(node *BinaryOpNode) any {
	switch node.Op {
	case "&&": {
		left := unwrapBuiltinValue(e.eval(node.Left))
		lBool, ok := left.(bool)
		if !ok {
			e.genError(fmt.Sprintf(
				"Left operand of '&&' must be bool, got %T", left),
				node.Position)
			return nil
		}
		if !lBool {
			return false
		}

		right := unwrapBuiltinValue(e.eval(node.Right))
		rBool, ok := right.(bool)
		if !ok {
			e.genError(fmt.Sprintf(
				"Right operand of '&&' must be bool, got %T", right),
				node.Position)
			return nil
		}
		return rBool

	}
	case "||": {
		left := unwrapBuiltinValue(e.eval(node.Left))
		lBool, ok := left.(bool)
		if !ok {
			e.genError(fmt.Sprintf(
				"Left operand of '||' must be bool, got %T", left),
				node.Position)
			return nil
		}
		if lBool {
			// Short-circuit: true || _ == true
			return true
		}
		right := unwrapBuiltinValue(e.eval(node.Right))
		rBool, ok := right.(bool)
		if !ok {
			e.genError(fmt.Sprintf(
				"Right operand of '||' must be bool, got %T", right),
				node.Position)
			return nil
		}
		return rBool
	}
	}
	left := unwrapBuiltinValue(e.eval(node.Left))
	right := unwrapBuiltinValue(e.eval(node.Right))

	switch node.Op {
	case ">", "<", ">=", "<=":
		if lStr, lok := left.(string); lok {
			if rStr, rok := right.(string); rok {
				switch node.Op {
				case ">":
					return lStr > rStr
				case "<":
					return lStr < rStr
				case ">=":
					return lStr >= rStr
				case "<=":
					return lStr <= rStr
				}
			} else {
				e.genError(
					"Rigth operand must be a string for string comp",
					node.Position)
				return nil
			}
		}
		lInt, lok := left.(int)
		rInt, rok := right.(int)
		if !lok || !rok {
			e.genError(fmt.Sprintf(
				"Operator '%s' requires integer operands", node.Op),
				node.Position)
			return nil
		}
		switch node.Op {
		case ">":
			return lInt > rInt
		case "<":
			return lInt < rInt
		case ">=":
			return lInt >= rInt
		case "<=":
			return lInt <= rInt
		}
	case "==", "!=":
		if left == nil || right == nil {
			e.genError("Cannot compare nil operands", node.Position)
			return nil
		}

		// if reflect.TypeOf(left) != reflect.TypeOf(right) {
		// 	e.genError(fmt.Sprintf(
		// 		"Operator '%s' requires operands of the same type", node.Op),
		// 		node.Position)
		// 	return nil
		// }
		switch node.Op {
		case "==":
			// Allow nil == nil, and nil == any other is false
			return isNilValue(left) && isNilValue(right) || left == right
		case "!=":
			return !isNilValue(left) && isNilValue(right) || isNilValue(left) && !isNilValue(right) || left != right
		}
	default:
		e.genError(fmt.Sprintf(
			"Unknown operator: %s", node.Op),
			node.Position)
		return nil
	}
	return nil
}


func (e *Evaluator) handlePostfixPP(node *UnaryOpNode) any {
	ident, ok := node.Expr.(*IdentifierNode)
	if !ok {
		e.genError("++ operator requires an identifier",
			node.Position)
		return nil
	}

	val := e.currentEnv.FindSymbol(ident.Name)
	intVal, ok := val.(int)
	if !ok {
		e.genError("++ operator requires integer value",
			node.Position)
		return nil
	}

	origVal := val

	e.currentEnv.UpdateSymbol(ident.Name, intVal+1, "int")

	return origVal
}

func (e *Evaluator) handlePostfixMM(node *UnaryOpNode) any {
	ident, ok := node.Expr.(*IdentifierNode)
	if !ok {
		e.genError("-- operator requires an identifier",
			node.Position)
		return nil
	}

	val := e.currentEnv.FindSymbol(ident.Name)
	intVal, ok := val.(int)
	if !ok {
		e.genError("-- operator requires integer value",
			node.Position)
		return nil
	}

	origVal := val

	e.currentEnv.UpdateSymbol(ident.Name, intVal-1, "int")

	return origVal
}

func isNilValue(x any) bool {
    _, ok := x.(nilValue)
    return ok
}

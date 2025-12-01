package eval

import (
	"fmt"
	"lang/internal/core"
	"lang/internal/parser"
)

func (e *Evaluator) evalBinary(expr *parser.BinaryOpNode) any {
	left := unwrapBuiltinValue(e.EvalNode(expr.Left))
	if ret, ok := left.(core.ReturnValue); ok {
		left = ret.Value
	}
	right := unwrapBuiltinValue(e.EvalNode(expr.Right))
	if ret, ok := right.(core.ReturnValue); ok {
		right = ret.Value
	}
	if (left == nil || right == nil) {
		e.GenError("Expression operand cannot bet nil", expr.Position)
		return nil
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
                e.GenError(
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
                e.GenError(
                    "Right operand of '+' must be int or float64 if left is float64",
                    expr.Position)
                return nil
            }
        case string:
            r, ok := right.(string)
            if !ok {
                e.GenError(
                    "Right operand of '+' must be string if left is string",
                    expr.Position)
                return nil
            }
            return e.createString(l + r)
        default:
            e.GenError("Unsupported type for '+' operator",
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
            e.GenError(fmt.Sprintf(
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
            e.GenError(fmt.Sprintf(
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
                    e.GenError("Division by zero!", expr.Position)
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
                    e.GenError("Division by zero!", expr.Position)
                    return nil
                }
                return lFloat / rFloat
            }
        }
    default:
        return e.evalLogical(expr)
    }
    e.GenError("Unknown", expr.Position)
    return nil
}

func (e *Evaluator) evalUnary(node *parser.UnaryOpNode) any {
    value := unwrapBuiltinValue(e.EvalNode(node.Expr))
    if _, ok := value.(core.NilValue); ok {
        e.GenError("Value with unary shouldn't be nil!", node.Position)
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
        e.GenError(fmt.Sprintf(
            "Unary '-' not supported for type %T", value),
            node.Position)
        return nil
    case "!":
        // Logical NOT for booleans
        if v, ok := value.(bool); ok {
            return !v
        }
        e.GenError(fmt.Sprintf(
            "Unary '!' not supported for type %T", value),
            node.Position)
        return nil
    default:
		e.GenError("Unsupported unary operator", node.Position)
        return nil
    }
}

func (e *Evaluator) evalCondition(condition parser.Node) any {
	switch t := condition.(type) {
	case *parser.UnaryOpNode:
		return e.evalUnary(t)
	case *parser.BinaryOpNode:
		return e.evalBinary(t)
	case *parser.FunctionCallNode:
		return e.evalFunctionCall(t)
	case *parser.IdentifierNode:
		return e.evalIdentifier(t)
	case *parser.LiteralNode:
		return e.evalLiteral(t)
	case *parser.StructMethodCall:
		return e.evalStructMemberAccess(t)
	default:
		e.GenError(fmt.Sprintf(
			"Unsupported type: %T", condition),
			parser.Position{Row: -1, Column: -1})
		return nil
	}
}

func (e *Evaluator) evalLiteral(lit *parser.LiteralNode) any {
	if e.resolveType(lit.Value, lit.Position) == "string" {
		if s, ok := lit.Value.(string); ok {
			fmt.Printf("Creating string: '%s' \n", s)
			return e.createString(s)
		} else {
			e.GenError("Failed to parse string", lit.Position)
			return nil
		}
	}
	return lit.Value
}

func (e *Evaluator) evalTrue() any {
	return true
}

func (e *Evaluator) evalFalse() any {
	return false
}

func (e *Evaluator) evalLogical(node *parser.BinaryOpNode) any {
	switch node.Op {
	case "&&": {
		left := unwrapBuiltinValue(e.EvalNode(node.Left))
		lBool, ok := left.(bool)
		if !ok {
			e.GenError(fmt.Sprintf(
				"Left operand of '&&' must be bool, got %T", left),
				node.Position)
			return nil
		}
		if !lBool {
			return false
		}

		right := unwrapBuiltinValue(e.EvalNode(node.Right))
		rBool, ok := right.(bool)
		if !ok {
			e.GenError(fmt.Sprintf(
				"Right operand of '&&' must be bool, got %T", right),
				node.Position)
			return nil
		}
		return rBool

	}
	case "||": {
		left := unwrapBuiltinValue(e.EvalNode(node.Left))
		lBool, ok := left.(bool)
		if !ok {
			e.GenError(fmt.Sprintf(
				"Left operand of '||' must be bool, got %T", left),
				node.Position)
			return nil
		}
		if lBool {
			// Short-circuit: true || _ == true
			return true
		}
		right := unwrapBuiltinValue(e.EvalNode(node.Right))
		rBool, ok := right.(bool)
		if !ok {
			e.GenError(fmt.Sprintf(
				"Right operand of '||' must be bool, got %T", right),
				node.Position)
			return nil
		}
		return rBool
	}
	}
	left := unwrapBuiltinValue(e.EvalNode(node.Left))
	right := unwrapBuiltinValue(e.EvalNode(node.Right))

	if (left == nil || right == nil) {
		e.GenError("Failed to get value", node.Position)
		return nil
	}

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
				e.GenError(
					"Rigth operand must be a string for string comp",
					node.Position)
				return nil
			}
		}
		lInt, lok := left.(int)
		rInt, rok := right.(int)
		if !lok || !rok {
			e.GenError(fmt.Sprintf(
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
			e.GenError("Cannot compare nil operands", node.Position)
			return nil
		}

		// if reflect.TypeOf(left) != reflect.TypeOf(right) {
		// 	e.GenError(fmt.Sprintf(
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
		e.GenError(fmt.Sprintf(
			"Unknown operator: %s", node.Op),
			node.Position)
		return nil
	}
	e.GenError("Uncaught error happen", node.Position)
	return nil
}


func (e *Evaluator) handlePostfixPP(node *parser.UnaryOpNode) any {
	ident, ok := node.Expr.(*parser.IdentifierNode)
	if !ok {
		e.GenError("++ operator requires an identifier",
			node.Position)
		return nil
	}

	val := e.currentEnv.FindSymbol(ident.Name)
	intVal, ok := val.(int)
	if !ok {
		e.GenError("++ operator requires integer value",
			node.Position)
		return nil
	}

	origVal := val

	e.currentEnv.UpdateSymbol(ident.Name, intVal+1, "int")

	return origVal
}

func (e *Evaluator) handlePostfixMM(node *parser.UnaryOpNode) any {
	ident, ok := node.Expr.(*parser.IdentifierNode)
	if !ok {
		e.GenError("-- operator requires an identifier",
			node.Position)
		return nil
	}

	val := e.currentEnv.FindSymbol(ident.Name)
	intVal, ok := val.(int)
	if !ok {
		e.GenError("-- operator requires integer value",
			node.Position)
		return nil
	}

	origVal := val

	e.currentEnv.UpdateSymbol(ident.Name, intVal-1, "int")

	return origVal
}

func isNilValue(x any) bool {
    _, ok := x.(core.NilValue)
    return ok
}

package main

import (
	"fmt"
)

func (e *Evaluator) evalBinary(expr *BinaryOpNode) any {
	left := e.eval(expr.Left)
	right := e.eval(expr.Right)

	switch expr.Op {
	case "+":
		switch l := left.(type) {
		case int:
			r, ok := right.(int)
			if !ok {
				e.genError(
					"Right operand of '+' must be int if left is int",
					expr.Position)
				return nil
			}
			return l + r
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
		lInt, lOk := left.(int)
		rInt, rOk := right.(int)
		if !lOk || !rOk {
			e.genError(fmt.Sprintf(
				"Operator '%s' requires integer operands", expr.Op),
				expr.Position)
			return nil
		}
		switch expr.Op {
		case "-": return lInt - rInt
		case "*": return lInt * rInt
		case "/": {
			if rInt == 0 {
				e.genError("Division by zero!", expr.Position)
				return nil
			}
			return lInt / rInt
		}
		}
	default:
		return e.evalLogical(expr)
	}
	e.genError("Unknown", expr.Position)
	return nil
}

func (e *Evaluator) evalUnary(node *UnaryOpNode) any {
	value := e.eval(node.Expr)
	if _, ok := value.(nilValue); ok {
		e.genError("Value with unary shoundn't be nil!",
			node.Position)
		return nil
	}
	switch node.Op {
	case "-":
		// Negation for numbers
		if v, ok := value.(int); ok {
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
	default:
		e.genError(fmt.Sprintf(
			"Unsupported type: %T", condition),
			Position{-1, -1})
		return nil
	}
}

func (e *Evaluator) evalLiteral(lit *LiteralNode) any {
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
		left := e.eval(node.Left)
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

		right := e.eval(node.Right)
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
		left := e.eval(node.Left)
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
		right := e.eval(node.Right)
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
	left := e.eval(node.Left)
	right := e.eval(node.Right)

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


func isNilValue(x any) bool {
    _, ok := x.(nilValue)
    return ok
}

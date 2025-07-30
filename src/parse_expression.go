package main

import (
	"fmt"
)

func (p *Parser) parseExpression(precedence int) Node {
    tok := p.currentToken()
    var left Node

    // Parse prefix (numbers, unary minus, parentheses)
    switch tok.TType {
    case Int, Float:
		left = p.parseDigit()
	case Identifier:
		left = p.parseIdentifier()
	case StringTok:
		left = p.parseString()
    case LParen:
        p.advance()
        left = p.parseExpression(0)
        if p.currentToken() == nil || p.currentToken().TType != RParen {
            p.genError("expected ')'")
            return nil
        }
        p.advance()
	case Nil:
		p.advance()
		return &NilNode{
			Position{
				Row: tok.Line,
				Column: tok.Column,
			},
		}
    case Minus:
        p.advance()
        expr := p.parseExpression(100) // high precedence for unary minus
        left = &UnaryOpNode{
			Position: Position {
				Row: p.currentToken().Line,
				Column: p.currentToken().Column,
			},
			Op: "-", Expr: expr}
	case Bang:
		p.advance()
		expr := p.parseExpression(100) // high precedence for unary ops
		left = &UnaryOpNode{
			Position: Position {
				Row: p.currentToken().Line,
				Column: p.currentToken().Column,
			},
			Op: "!", Expr: expr}
	case PlusPlus, MinusMinus:
		op := ""
		if tok.TType == PlusPlus {
			op = "++"
		} else {
			op = "--"
		}
		p.advance()

		var expr Node
		switch p.currentToken().TType {
		case Identifier:
			expr = p.parseIdentifier()
		default:
			p.genError(fmt.Sprintf("expected identifier after '%s'", op))
			return nil
		}

		left = &UnaryOpNode{
			Position: Position{
				Row: p.currentToken().Line,
				Column: p.currentToken().Column,
			},
			Op:   op,
			Expr: expr,
		}
	case True: {
		tok := p.currentToken()
		p.advance()
		left = &TrueNode{
			Position {
				Row: tok.Line,
				Column: tok.Column,
			},
		}
	}
	case False: {
		tok := p.currentToken()
		p.advance()
		left = &FalseNode{
			Position {
				Row: tok.Line,
				Column: tok.Column,
			},
		}
	}
	default:
		p.genError(fmt.Sprintf("Unexpected token in expression: %v", tok.Lexeme))
		return nil
	}

    // Parse infix (binary) operators
	for {
		next := p.currentToken()
		if next == nil {
			break
		}
		opPrec := precedenceOf(next.TType)

		if opPrec <= precedence {
			break
		}

		// handle postfix ++ and -- specially, no right expression
		if next.TType == PlusPlus || next.TType == MinusMinus {
			p.genError("'++' or '--' currently doesn't supported")
			return nil
			op := ""
			if next.TType == PlusPlus {
				op = "++"
			} else {
				op = "--"
			}
			p.advance()
			// Create a UnaryOpNode with left as the Expr for postfix
			left = &UnaryOpNode{
				Position: Position{
					Row: next.Line,
					Column: next.Column,
				},
				Op: op,
				Expr: left,
			}
			continue // continue loop for potentially chained postfix ops
		}

		op := next.TType
		p.advance()
		right := p.parseExpression(opPrec)
		left = &BinaryOpNode{
			Position: Position{
				Row: p.currentToken().Line,
				Column: p.currentToken().Column,
			},
			Op:    tokenTypeToString(op),
			Left:  left,
			Right: right,
		}
	}
    return left
}

func precedenceOf(tok TokenType) int {
    switch tok {
    case Or:
        return 1
    case And:
        return 2
    case Equals, NotEquals:
        return 3
    case Less, More, LessEq, MoreEq:
        return 4
    case Plus, Minus:
        return 5
    case Star, Slash:
        return 6
	case PlusPlus, MinusMinus:
		return 7
    }
    return 0
}


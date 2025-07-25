package main

import (
	"fmt"
)

func (p *Parser) parseExpression(precedence int) Node {
    tok := p.currentToken()
    var left Node

    // Parse prefix (numbers, unary minus, parentheses)
    switch tok.TType {
    case Digit:
		left = p.parseDigit()
	case Identifier:
		left = p.parseIdentifier()
	case String:
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
    default:
        p.genError(fmt.Sprintf("unexpected token in expression: %v", tok.Lexeme))
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

        op := next.TType
        p.advance()
        right := p.parseExpression(opPrec)
        left = &BinaryOpNode{
			Position: Position {
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
    }
    return 0
}


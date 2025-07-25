package main

import (
	"errors"
	"fmt"
)

func (p *Parser) genError(message string) {
	msg := fmt.Sprintf("Parse error in %d:%d: %s",
		p.currentToken().Line,
		p.currentToken().Column,
		message)
	p.Errors = append(p.Errors, errors.New(msg))
}

func tokenTypeToString(t TokenType) string {
    switch t {
    case Plus:      return "+"
    case Minus:     return "-"
    case Star:      return "*"
    case Slash:     return "/"
    case And:       return "&&"
    case Or:        return "||"
    case Equals:    return "=="
    case NotEquals: return "!="
    case Less:      return "<"
    case More:      return ">"
    case LessEq:    return "<="
    case MoreEq:    return ">="
    case Assign:    return "="
    default:        return ""
    }
}

func (p *Parser) expectAndAdvance(tType TokenType) bool{
	if p.currentToken().TType == tType {
		p.advance()
		return true
	} else {
		p.genError(fmt.Sprintf("Expected %T; But got: %V",
			tType, p.currentToken().Lexeme))
		return false
	}
}

func (p *Parser) expect(tType TokenType) bool {
	if p.currentToken().TType == tType {
		return true
	}
	return false
}

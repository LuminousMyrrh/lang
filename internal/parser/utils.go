package parser

import (
	"errors"
	"fmt"
	"lang/internal/token"
)

func (p *Parser) genError(message string) {
	msg := fmt.Sprintf("Parse error in %d:%d: %s",
		p.currentToken().Line,
		p.currentToken().Column,
		message)
	p.Errors = append(p.Errors, errors.New(msg))
}

func tokenTypeToString(t token.TokenType) string {
    switch t {
    case token.Plus:      return "+"
    case token.Minus:     return "-"
    case token.Star:      return "*"
    case token.Slash:     return "/"
    case token.And:       return "&&"
    case token.Or:        return "||"
    case token.Equals:    return "=="
    case token.NotEquals: return "!="
    case token.Less:      return "<"
    case token.More:      return ">"
    case token.LessEq:    return "<="
    case token.MoreEq:    return ">="
    case token.Assign:    return "="
    default:        return ""
    }
}

func (p *Parser) expectAndAdvance(tType token.TokenType) bool {
	if p.currentToken().TType == tType {
		p.advance()
		return true
	} else {
		p.genError(fmt.Sprintf("Expected %v; But got: %v",
			tType, p.currentToken().Lexeme))
		return false
	}
}

func (p *Parser) expect(tType token.TokenType) bool {
	if p.currentToken().TType == tType {
		return true
	}
	return false
}

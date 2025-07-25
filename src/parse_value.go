package main

import "strconv"

func (p *Parser) parseDigit() *LiteralNode {
	val, err := strconv.Atoi(p.currentToken().Lexeme)
	if err != nil {
		p.genError(err.Error())
		return nil
	}

	p.advance()
	return &LiteralNode{
		Position: Position {
			Row: p.currentToken().Line,
			Column: p.currentToken().Column,
		},
		Value: val,
	}
}

func (p *Parser) parseString() *LiteralNode {
	str := LiteralNode{Value: p.currentToken().Lexeme}
	p.advance()
	return &str
}


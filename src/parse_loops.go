package main

func (p *Parser) parseWhile() *WhileNode {
	initTok := p.currentToken()
	p.advance()

	if p.currentToken().TType == LCurly {
		// while {}
		body := p.parseBlock()

		return &WhileNode{
			Condition: nil,
			Body: body,
		}
	} else {
		// while () ...
		if p.currentToken().TType != LParen {
			p.genError("Expected '(' after 'while'")
			return nil
		}
		p.advance()

		cond := p.parseExpression(0)
		if p.currentToken().TType != RParen {
			p.genError("Expected ')' after condition")
			return nil
		}
		p.advance()
		p.advance()

		body := p.parseBlock()
		return &WhileNode{
			Position: Position {
				Row: initTok.Line,
				Column: initTok.Column,
			},
			Condition: cond,
			Body: body,
		}
	}
}

func (p *Parser) parseFor() *ForNode {
	p.advance()
	if !p.expectAndAdvance(LParen) {
		return nil
	}
	var Init *VarDefNode
	if p.currentToken().TType == Semicolon {
		Init = nil
		p.advance()
	} else {
		Init = p.parseVarDef()
		if !p.expectAndAdvance(Semicolon) {
			return nil
		}
	}
	var Condition Node
	if p.currentToken().TType == Semicolon {
		Condition = nil
		p.advance()
	} else {
		Condition = p.parseExpression(0)
	}
	var Post Node
	if p.currentToken().TType == RParen {
		Post = nil
	} else {
		Post = p.parseExpression(0)
		if !p.expectAndAdvance(RParen) {
			return nil
		}
	}
	return nil
}

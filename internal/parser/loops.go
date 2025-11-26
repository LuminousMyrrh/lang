package parser

import "lang/internal/token"

func (p *Parser) parseWhile() *WhileNode {
	initTok := p.currentToken()
	p.advance()

	if p.currentToken().TType == token.LCurly {
		// while {}
		body := p.parseBlock()

		return &WhileNode{
			Condition: nil,
			Body: body,
		}
	} else {
		// while () ...
		if p.currentToken().TType != token.LParen {
			p.genError("Expected '(' after 'while'")
			return nil
		}
		p.advance()

		cond := p.parseExpression(0)
		if p.currentToken().TType != token.RParen {
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

func (p *Parser) parseForLoop() *ForNode {
	p.advance()
	initToken := p.currentToken()
	if !p.expectAndAdvance(token.LParen) {
		return nil
	}
	var Init *VarDefNode = nil
	if p.currentToken().TType == token.Semicolon {
		p.advance()
	} else {
		Init = p.parseVarDef()
		if Init == nil {
			return nil
		}
	}
	var Condition Node = nil
	if p.currentToken().TType == token.Semicolon {
		p.advance()
	} else {
		Condition = p.parseExpression(0)
		if !p.expectAndAdvance(token.Semicolon) {
			return nil
		}
	}
	var Post Node = nil
	if p.currentToken().TType == token.RParen {
		p.advance()
	} else {
		Post = p.parseExpression(0)
		if !p.expectAndAdvance(token.RParen) {
			return nil
		}
	}

	if !p.expectAndAdvance(token.LCurly) {
		return nil
	}

	body := p.parseBlock()

	node := &ForNode{
		Position {
			Row: initToken.Line,
			Column: initToken.Column,
		},
		Init,
		Condition,
		Post,
		body,
	}

	return node
}

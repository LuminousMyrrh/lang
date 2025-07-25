package main

import "fmt"

func (p *Parser) parseFuncDef() *FunctionDefNode {
	p.advance()
	nameTok := p.currentToken();
	if nameTok == nil || nameTok.TType != Identifier {
		p.genError(fmt.Sprintf(
			"expected function name after 'func', got %v", nameTok))
		return nil
	}
	name := nameTok.Lexeme

	p.advance()
	if p.currentToken().TType != LParen {
		p.genError(fmt.Sprintf(
			"expected function parameters after function name, got %v", nameTok))
		return nil
	}
	p.advance()
	var params []string
	for p.currentToken().TType != RParen {
		if p.currentToken().TType == Comma {
			p.advance()
		}
		params = append(params, p.currentToken().Lexeme)
		p.advance()
	}
	p.advance()
	p.advance()
	body := p.parseBlock()

	return &FunctionDefNode{
		Position: Position {
			Row: nameTok.Line,
			Column: nameTok.Column,
		},
		Name: name,
		Parameters: params,
		Body: body,
	}
}

func (p *Parser) parseFuncCall(node Node) *FunctionCallNode {
	initTok := p.currentToken()
	p.advance() // skip '('
	var args []Node
	for p.currentToken() != nil && p.currentToken().TType != RParen {
		if p.currentToken().TType == Comma {
			p.advance()
			continue
		}
		args = append(args, p.parseExpression(0))
	}
	if p.currentToken() != nil && p.currentToken().TType == RParen {
		p.advance() // skip ')'
	} else {
		p.genError("Expected ')' after function call arguments")
		return nil
	}
	return &FunctionCallNode{
		Position: Position {
			Row: initTok.Line,
			Column: initTok.Column,
		},
		Name: node,
		Args: args,
	}
}

func (p *Parser) parseBlock() *BlockNode {
	body := &BlockNode{
		Position: Position {
			Row: p.currentToken().Line,
			Column: p.currentToken().Column,
		},
	}

	for p.currentToken().TType != RCurly {
		stmt := p.parseStatement()
		if stmt == nil {
			return nil
		}
		if _, ok := stmt.(*SemicolonNode); !ok {
			body.Statements = append(body.Statements, stmt)
		}
	}
	p.advance()

	return body
}

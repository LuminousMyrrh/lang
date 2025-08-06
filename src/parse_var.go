package main

import "fmt"

func (p *Parser) parseVarDef() *VarDefNode {
	p.advance()
	nameTok := p.currentToken()
	if nameTok == nil || nameTok.TType != Identifier {
		p.genError(fmt.Sprintf(
			"expected identifier after 'var', got %v", nameTok))
		return nil
	}
	name := nameTok.Lexeme
	p.advance()

	assignTok := p.currentToken()
	if assignTok == nil || assignTok.TType != Assign {
		p.advance()
		return &VarDefNode{
			Name: name,
			Value: nil,
		}
	}
	p.advance()

	value := p.parseValue()

	if !p.expectAndAdvance(Semicolon) {
		return nil
	}
	return &VarDefNode{
		Position: Position {
			Row: nameTok.Line,
			Column: nameTok.Column,
		},
		Name: name, Value: value}
}

func (p *Parser) parseIdentifier() Node {
    if p.currentToken() == nil {
        p.genError("Nil token")
        return nil
    }
    // Parse the identifier name and advance
    id := p.currentToken()
    p.advance()

    // Start with a plain identifier node
    var node Node = &IdentifierNode{
		Position: Position {
			Row: id.Line,
			Column: id.Column,
		},
		Name: id.Lexeme,
	}

    // Handle any number of array accesses: x[i][j][k]
	for {
		if p.currentToken() != nil && p.currentToken().TType == LBrace {
			node = p.parseArrayAccess(node)
			if node == nil {
				return nil
			}
		} else if p.currentToken() != nil && p.currentToken().TType == Dot {
			node = p.parseStructMethodCall(node)
			if node == nil {
				return nil
			}
		} else if p.currentToken() != nil && p.currentToken().TType == LParen {
			node = p.parseFuncCall(node)
			if node == nil {
				return nil
			}
		} else {
			break
		}
	}

    // Assignment: x = ... or x[i][j] = ...
    if p.currentToken() != nil && (p.currentToken().TType == Assign ||
		p.currentToken().TType == PlusEq || p.currentToken().TType == MinusEq ) {
		op := p.currentToken().Lexeme
		p.advance()
        value := p.parseValue()
        return &AssignmentNode{
			Position: Position {
				Row: id.Line,
				Column: id.Column,
			},
            Name:  node,
            Value: value,
			Op: op,
        }
    }

    return node
}

func (p *Parser) parseValue() Node {
	currTok := p.currentToken()
	if currTok.TType == Nil {
		p.advance()
		return &NilNode{
			Position{
				Row: currTok.Line,
				Column: currTok.Column,
			},
		}
	}
	if currTok.TType == LBrace {
		return p.parseArray()
	}
	if currTok.TType == Identifier &&
		p.nextToken().TType == LCurly {
		return p.parseStructInit()
	}
	if currTok.TType == True {
		tok := p.currentToken()
		p.advance()
		return &TrueNode{
			Position: Position {
				Row: tok.Line,
				Column: tok.Column,
			},
		}
	} else if currTok.TType == False {
		tok := p.currentToken()
		p.advance()
		return &FalseNode{
			Position: Position {
				Row: tok.Line,
				Column: tok.Column,
			},
		}
	}
	

	return p.parseExpression(0)
}

package parser

import (
	"fmt"
	"lang/internal/token"
)

func (p *Parser) parseVarDef() *VarDefNode {
	p.advance()
	nameTok := p.currentToken()
	if nameTok == nil || nameTok.TType != token.Identifier {
		p.genError(fmt.Sprintf(
			"expected identifier after 'var', got %v", nameTok))
		return nil
	}
	name := nameTok.Lexeme
	p.advance()

	assignTok := p.currentToken()
	if assignTok == nil || assignTok.TType != token.Assign {
		p.advance()
		return &VarDefNode{
			Name: name,
			Value: nil,
		}
	}
	p.advance()

	value := p.parseValue()

	if !p.expectAndAdvance(token.Semicolon) {
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
		if p.currentToken() != nil && p.currentToken().TType == token.LBrace {
			node = p.parseArrayAccess(node)
			if node == nil {
				return nil
			}
		} else if p.currentToken() != nil && p.currentToken().TType == token.Dot {
			node = p.parseStructMethodCall(node)
			if node == nil {
				return nil
			}
		} else if p.currentToken() != nil && p.currentToken().TType == token.LParen {
			node = p.parseFuncCall(node)
			if node == nil {
				return nil
			}
		} else {
			break
		}
	}

    // Assignment: x = ... or x[i][j] = ...
    if p.currentToken() != nil && (p.currentToken().TType == token.Assign ||
		p.currentToken().TType == token.PlusEq || p.currentToken().TType == token.MinusEq ) {
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
	if currTok.TType == token.Nil {
		p.advance()
		return &NilNode{
			Position{
				Row: currTok.Line,
				Column: currTok.Column,
			},
		}
	}
	if currTok.TType == token.LBrace {
		return p.parseArray()
	}
	if currTok.TType == token.Identifier &&
		p.nextToken().TType == token.LCurly {
		return p.parseStructInit()
	}
	if currTok.TType == token.True {
		tok := p.currentToken()
		p.advance()
		return &TrueNode{
			Position: Position {
				Row: tok.Line,
				Column: tok.Column,
			},
		}
	} else if currTok.TType == token.False {
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

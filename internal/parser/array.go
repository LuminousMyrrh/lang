package parser

import "lang/internal/token"

func (p *Parser) parseArray() *ArrayNode {
    p.advance()
	initTok := p.currentToken()
    var elems []Node

    for p.currentToken() != nil && p.currentToken().TType != token.RBrace {
        if p.currentToken().TType == token.Comma {
            p.advance()
            continue
        }
        elems = append(elems, p.parseValue())
        if p.currentToken() != nil &&
			p.currentToken().TType != token.RBrace &&
			p.currentToken().TType != token.Comma {

            p.advance()
        }
    }

    if p.currentToken() != nil && p.currentToken().TType == token.RBrace {
        p.advance()
    } else {
        p.genError("Expected ']' at end of array literal")
    }

    return &ArrayNode{
		Position: Position {
			Row: initTok.Line,
			Column: initTok.Column,
		},
        Elements: elems,
    }
}

func (p *Parser) parseArrayAccess(node Node) *ArrayAccessNode {
	var retNode *ArrayAccessNode
    for p.currentToken() != nil && p.currentToken().TType == token.LBrace {
        p.advance() // skip '['
        index := p.parseValue()
        if p.currentToken() == nil || p.currentToken().TType != token.RBrace {
            p.genError("Expected ']' after array index")
            return nil
        }
        p.advance() // skip ']'
        retNode = &ArrayAccessNode{
			Position: Position {
				Row: p.currentToken().Line,
				Column: p.currentToken().Column,
			},
            Target: node, // The previous node (IdentifierNode or ArrayAccessNode)
            Index: index,
        }
    }

	return retNode
}

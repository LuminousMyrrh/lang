package main

func (p *Parser) parseArray() *ArrayNode {
    // Assume current token is LBracket
    p.advance() // move past '['
	initTok := p.currentToken()
    var elems []Node

    for p.currentToken() != nil && p.currentToken().TType != RBrace {
        if p.currentToken().TType == Comma {
            p.advance()
            continue
        }
        elems = append(elems, p.parseValue())
        // Only advance if the next token is not a comma or closing bracket
        if p.currentToken() != nil && p.currentToken().TType != RBrace && p.currentToken().TType != Comma {
            p.advance()
        }
    }

    if p.currentToken() != nil && p.currentToken().TType == RBrace {
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
    for p.currentToken() != nil && p.currentToken().TType == LBrace {
        p.advance() // skip '['
        index := p.parseExpression(0)
        if p.currentToken() == nil || p.currentToken().TType != RBrace {
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

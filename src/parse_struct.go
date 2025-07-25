package main

func (p *Parser) parseStructDef() *StructDefNode {
	p.advance() // skip 'struct'
	nameTok := p.currentToken()
	if nameTok == nil || nameTok.TType != Identifier {
		p.genError("Expected struct name")
		return nil
	}
	name := nameTok.Lexeme
	p.advance() // skip struct name

	if p.currentToken().TType != LCurly {
		p.genError("Expected '{' after struct name")
		return nil
	}
	p.advance() // skip '{'

	var fields []*StructField

	for p.currentToken() != nil && p.currentToken().TType != RCurly {
		if p.currentToken().TType == Comma {
			p.advance()
			continue
		}
		field := p.parseStructField()
		fields = append(fields, field)
	}

	if p.currentToken() == nil || p.currentToken().TType != RCurly {
		p.genError("Expected '}' at end of struct")
		return nil
	}
	p.advance() // skip '}'

	return &StructDefNode{
		Position: Position {
			Row: nameTok.Line,
			Column: nameTok.Column,
		},
		Name: name,
		Fields: fields,
	}
}


func (p *Parser) parseStructMethodCall(node Node) *StructMethodCall {
	var retNode *StructMethodCall
	for p.currentToken() != nil && p.currentToken().TType == Dot {
		p.advance() // skip '.'
		if p.currentToken() == nil || p.currentToken().TType != Identifier {
			p.genError("Expected identifier after '.'")
			return nil
		}
		member := p.currentToken()
		memberName := member.Lexeme
		p.advance()

		// Check for method call: obj.method(...)
		if p.currentToken() != nil && p.currentToken().TType == LParen {
			p.advance() // skip '('
			var args []Node
			for p.currentToken() != nil && p.currentToken().TType != RParen {
				if p.currentToken().TType == Comma {
					p.advance()
					continue
				}
				args = append(args, p.parseValue())
			}
			if p.currentToken() != nil && p.currentToken().TType == RParen {
				p.advance() // skip ')'
			} else {
				p.genError("Expected ')' after method call arguments")
				return nil
			}
			retNode = &StructMethodCall{
				Position: Position {
					Row: member.Line,
					Column: member.Column,
				},
				Caller: node,
				MethodName: memberName,
				IsField: false,
				Args: args,
			}
		} else {
			// Field access: obj.field
			retNode = &StructMethodCall{
				Position: Position {
					Row: member.Line,
					Column: member.Column,
				},
				Caller: node,
				MethodName: memberName,
				IsField: true,
				Args: nil,
			}
		}
	}

	return retNode
}

func (p *Parser) parseStructField() *StructField {
	var isPub bool = true
	if p.currentToken().TType == Private {
		isPub = false
	}
	p.advance()
	nameTok := p.currentToken()
	name := nameTok.Lexeme
	p.advance()
	if p.currentToken().TType == Assign {
		p.advance()
		value := p.parseValue()
		return &StructField{
			Name: name,
			Value: value,
			IsPublic: isPub,
		}
	} 

	return &StructField{
		Position: Position {
			Row: nameTok.Line,
			Column: nameTok.Column,
		},
		Name: name,
		Value: nil,
		IsPublic: isPub,
	}
}

func (p *Parser) parseStructMethodDef() *StructMethodDef {
	isPub := true
	if p.currentToken() != nil && p.currentToken().TType == Private {
		isPub = false
	}
	p.advance()
	if p.currentToken() == nil {
		p.genError("Expected struct name")
		return nil
	}
	structName := p.currentToken().Lexeme
	p.advance()
	if p.currentToken() == nil || p.currentToken().TType != RightArrow {
		p.genError("Expected '->' after struct name")
		return nil
	}
	p.advance()
	if p.currentToken() == nil {
		p.genError("Expected method name after '->'")
		return nil
	}
	nameTok := p.currentToken()
	methodName := nameTok.Lexeme
	p.advance()
	if p.currentToken() == nil || p.currentToken().TType != LParen {
		p.genError("Expected '(' after method name")
		return nil
	}
	p.advance()
	var params []string
	for p.currentToken() != nil && p.currentToken().TType != RParen {
		if p.currentToken().TType == Comma {
			p.advance()
			continue
		}
		params = append(params, p.currentToken().Lexeme)
		p.advance()
	}
	if p.currentToken() == nil || p.currentToken().TType != RParen {
		p.genError("Expected ')' after parameter list")
		return nil
	}
	p.advance() // skip ')'
	if p.currentToken() == nil || p.currentToken().TType != LCurly {
		p.genError("Expected '{' to start method body")
		return nil
	}
	p.advance()
	body := p.parseBlock()
	return &StructMethodDef{
		Position: Position {
			Row: nameTok.Line,
			Column: nameTok.Column,
		},
		IsPub: isPub,
		StructName: structName,
		MethodName: methodName,
		Parameters: params,
		Body: body,
	}
}

func (p *Parser) parseStructInit() *StructInitNode {
	structName := p.currentToken().Lexeme
	p.advance() // skip name
	p.advance() // skip '{'

	var fieldsInit []Node
	for p.currentToken() != nil &&
		p.currentToken().TType != RCurly {
		if p.currentToken().TType == Comma {
			p.advance()
			continue
		}
		fieldName := p.currentToken().Lexeme
		p.advance()
		if p.currentToken().TType != Colon {
			p.genError("Expected ':' and field name")
			return nil
		}
		p.advance() // skip '='
		value := p.parseValue()
		fieldAssign := &AssignmentNode{
			Name: &IdentifierNode{Name: fieldName},
			Value: value,
		}
		fieldsInit = append(fieldsInit, fieldAssign)
	}
	if p.currentToken().TType != RCurly {
		p.genError("Expected '}' after struct init")
		return nil
	}
	p.advance()
	return &StructInitNode{
		Name: structName,
		InitFields: fieldsInit,
	}
}

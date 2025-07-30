package main

import (
	"fmt"
)

type Parser struct {
	Tokens []*Token
	TokensLength int
	MainNode ProgramNode
	Errors []error
	pos int
}

func NewParser(toks []*Token) *Parser {
	return &Parser {
		Tokens: toks,
		TokensLength: len(toks),
		MainNode: ProgramNode{},
		Errors: make([]error, 0),
		pos: 0,
	}
}

func (p *Parser) Parse() (*ProgramNode, []error) {
	p.genAst()

	return &p.MainNode, p.Errors
}

func (p *Parser) genAst() {
	for p.pos < p.TokensLength {
		node := p.parseStatement()
		if node != nil {
			if _, ok := node.(*SemicolonNode); !ok {
				p.MainNode.Nodes = append(p.MainNode.Nodes, node)
			}
		} else {
			return
		}
	}
}

func (p *Parser) parseStatement() Node {
	if p.currentToken() == nil {
		p.genError("Returning nil statement")
		return nil
	}

	switch p.currentToken().TType {
	case Int, Float, LParen, Minus:
		expr := p.parseExpression(0)

		return expr
    case Var: {
		node := p.parseVarDef()
		return node
	}
	case Func: {
		node := p.parseFuncDef()
		return node
	}
	case If:
		node := p.parseIf()
		return node
	case Struct: {
		node := p.parseStructDef()
		return node
	}
	case Public, Private: {
		node := p.parseStructMethodDef()
		return node
	}
	case Identifier: {
		node := p.parseIdentifier()
		return node
	}
	case StringTok, EmptyStringTok: {
		node := p.parseString()
		return node
	}
	case While: {
		node := p.parseWhile()
		return node
	}
	case For: {
		node := p.parseForLoop()
		return node
	}
	case Semicolon:
		p.advance()
		return &SemicolonNode{}
	case Import:
		node := p.parseImport()
		return node

	case Return:
		// skip 'return'
		initTok := p.currentToken()
		p.advance()
		if p.currentToken().TType == Semicolon {
			p.advance()
			return &ReturnNode {
				Position {
					Row: initTok.Line,
					Column: initTok.Column,
				},
				nil,
			}
		}
		value := p.parseExpression(0)
		// expecting to advace ';'
		if p.currentToken().TType != Semicolon {
			p.genError("Expected ';' after value in return block")
			return nil
		}
		p.advance()
		return &ReturnNode{
			Value: value,
		}

	default:
        // Try to parse as an expression statement
        expr := p.parseExpression(0)
        if expr == nil {
            p.genError(fmt.Sprintf("Unknown token: %v", p.currentToken().Lexeme))
			return nil
        }
        // Expect semicolon after expression statement
        if p.currentToken() == nil || p.currentToken().TType != Semicolon {
            p.genError(fmt.Sprintf("expected ';' after expression statement", p.currentToken().Lexeme))
			return nil
        }
        p.advance()
        return &ExpressionStatementNode{Expr: expr}
	}
}

func (p *Parser) parseImport() *ImportNode {
	p.advance()
	file := p.currentToken().Lexeme
	p.advance()
	if p.currentToken().TType == Semicolon {
		p.advance()
		return &ImportNode{
			File: file,
		}
	} else if p.currentToken().TType == More {
		p.advance()
		var symbols []string

		for p.currentToken().TType != Semicolon {
			if p.currentToken().TType == Comma {
				p.advance()
				continue
			}

			symbols = append(symbols, p.currentToken().Lexeme)

			p.advance()
		}

		return &ImportNode{
			File: file,
			Symbols: symbols,
		}
	} 
	p.genError("Unknown")
	return nil
}

func (p *Parser) parseIf() *IfNode {
	initTok := p.currentToken()
	p.advance()

    if p.currentToken().TType != LParen {
        p.genError("Expected '(' after 'if'")
        return nil
    }
    p.advance() // consume '('

	condition := p.parseExpression(0)

    if p.currentToken().TType != RParen {
        p.genError("Expected ')' after if condition")
        return nil
    }
    p.advance() // consume ')'

    p.advance() // consume '{'
	thenBranch := p.parseBlock()

	var elseBranch *BlockNode = nil

    if p.currentToken() != nil && p.currentToken().TType == Else {
        p.advance() // consume 'else'
        if p.currentToken() != nil && p.currentToken().TType == If {
            // else if: recursively parse as a nested IfNode in the else branch
            elseIfNode := p.parseIf()
            // Wrap the else-if IfNode in a block (so AST is consistent)
            elseBranch = &BlockNode{
				Position: Position {
					Row: p.currentToken().Line,
					Column: p.currentToken().Column,
				},
				Statements: []Node{elseIfNode}}
        } else {
			p.advance()
            elseBranch = p.parseBlock()
        }
    }

    return &IfNode{
		Position: Position {
			Row: initTok.Line,
			Column: initTok.Column,
		},
        Condition:  condition,
        ThenBranch: thenBranch,
        ElseBranch: elseBranch,
    }
}

func (p *Parser) currentToken() *Token {
	if p.pos < p.TokensLength {
		return p.Tokens[p.pos]
	}
	return nil
}

func (p *Parser) nextToken() *Token {
	if p.pos + 1 < p.TokensLength {
		return p.Tokens[p.pos + 1]
	}
	return nil
}

func (p *Parser) advance() {
	p.pos++
}

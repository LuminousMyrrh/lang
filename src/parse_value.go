package main

import (
	"strconv"
)

func (p *Parser) parseDigit() *LiteralNode {
    var val any
    var err error

    tok := p.currentToken()

    if tok.TType == Int {
        val, err = strconv.Atoi(tok.Lexeme)
        if err != nil {
            p.genError(err.Error())
            return nil
        }
    } else {
        val, err = strconv.ParseFloat(tok.Lexeme, 64)
        if err != nil {
            p.genError(err.Error())
            return nil
        }
    }

    pos := Position{
        Row:    tok.Line,
        Column: tok.Column,
    }

    p.advance()

    return &LiteralNode{
        Position: pos,
        Value:    val,
    }
}

func (p *Parser) parseString() *LiteralNode {
	str := LiteralNode{Value: p.currentToken().Lexeme}
	p.advance()
	return &str
}


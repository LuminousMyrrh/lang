package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Node interface {
    String() string
}

// Root node for the entire program
type ProgramNode struct {
    Nodes []Node // List of statements or declarations
}

func (p ProgramNode) Print() {
	for _, node := range p.Nodes {
		fmt.Println(node.String())
	}
}

// Literals
type LiteralNode struct {
    Value any // Can be int, string, bool, etc.
}

func (n LiteralNode) String() string {
	switch val := n.Value.(type) {
	case string:
		return fmt.Sprintf("\"%v\"", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Identifiers (variable or function names)
type IdentifierNode struct {
    Name string
}

func (i IdentifierNode) String() string {
	return fmt.Sprintf("@%v", i.Name)
}

// Binary operations (e.g., a + b)
type BinaryOpNode struct {
    Op    string
    Left  Node
    Right Node
}

func (n BinaryOpNode) String() string {
	return fmt.Sprintf("#%v %v %v#",n.Left.String(), n.Op, n.Right.String())
}

// Unary operations (e.g., -x, !flag)
type UnaryOpNode struct {
    Op   string
    Expr Node
}

func (n UnaryOpNode) String() string {
	return fmt.Sprintf("%v %v", n.Op, n.Expr.String())
}

type StructField struct {
	Name string
	Access bool // 0 - private 1 - pub
}

type StructDefNode struct {
	Name string
	Fields []StructField
	Methods []*FunctionDefNode
}

func (s *StructDefNode) String() string {
	var str string
	str += s.Name + " "
	for _, field := range s.Fields {
		str += field.Name + " "
	}
	str += "\n"

	for _, method := range s.Methods {
		str += method.String() + "\n"
	}

	return str
}

// Variable declaration (e.g., var x = 5)
type VarDefNode struct {
    Name  string
    Value Node
}

func (n VarDefNode) String() string {
	return fmt.Sprintf("@%v: %v", n.Name, n.Value.String())
}

// Assignment (e.g., x = 10)
type AssignmentNode struct {
    Name  string
    Value Node
}

func (a AssignmentNode) String() string {
	return fmt.Sprintf("@%v: %v", a.Name, a.Value.String())
}

// Function definition (e.g., func foo(a, b) { ... })
type FunctionDefNode struct {
    Name       string
    Parameters []string
    Body       *BlockNode
}

func (f FunctionDefNode) String() string {
	var str string
	str += f.Name + " "
	for _, param := range f.Parameters {
		str += param + " "
	}
	str += "\n"

	for _, node := range f.Body.Statements {
		str += node.String() + "\n"
	}

	return str
}

// Function call (e.g., foo(a, b))
type FunctionCallNode struct {
    Name string
    Args []Node
}

func (f FunctionCallNode) String() string {
	var str string

	str += f.Name
	for _, arg := range f.Args {
		str += arg.String()
	}

	return str
}

// Block of statements (e.g., { ... })
type BlockNode struct {
    Statements []Node
}

func (b BlockNode) String() string {
	var str string
	for _, node := range b.Statements {
		str += node.String()
	}

	return str
}

// If/Else conditional
type IfNode struct {
    Condition   Node
    ThenBranch  *BlockNode
    ElseBranch  *BlockNode // Can be nil if no else
}

// While loop
type WhileNode struct {
    Condition Node
    Body      *BlockNode
}

// For loop (basic: for var i = 0; i < 10; i = i + 1 { ... })
type ForNode struct {
    Init      Node      // e.g., VarDefNode or AssignmentNode
    Condition Node
    Post      Node      // e.g., AssignmentNode
    Body      *BlockNode
}

// Return statement
type ReturnNode struct {
    Value Node // Can be nil for "return"
}

func (r ReturnNode) String() string {
	//return r.Value.String()
	return ""
}

// Expression statement (e.g., a function call as a statement)
type ExpressionStatementNode struct {
    Expr Node
}

func (e ExpressionStatementNode) String() string {
	return e.Expr.String()
}

type Parser struct {
	Tokens []*Token
	TokensLength int
	MainNode ProgramNode
	Errors []error
	GlobalEnv *Env
	currentEnv *Env
	pos int
}

func NewParser(env *Env, toks []*Token) *Parser {
	return &Parser {
		Tokens: toks,
		TokensLength: len(toks),
		MainNode: ProgramNode{},
		Errors: make([]error, 0),
		GlobalEnv: env,
		currentEnv: env,
		pos: 0,
	}
}

func (p *Parser) Parse(tokens []*Token) (*ProgramNode, []error) {
	p.genAst()

	return &p.MainNode, p.Errors
}

func (p *Parser) genAst() {
	for p.pos < p.TokensLength {
		node, err := p.parseStatement()
		if err != nil {
			p.genError(err.Error())
			return
		}
		if node != nil {
			p.MainNode.Nodes = append(p.MainNode.Nodes, node)
		}
	}
}

func (p *Parser) parseStatement() (Node, error){
	if p.currentToken() == nil {
		return nil, nil
	}

	switch p.currentToken().TType {
	case Digit, LParen, Minus:
		expr := p.parseExpression(0)

		return expr, nil
    case Var: {
        p.advance()
        nameTok := p.currentToken()
        if nameTok == nil || nameTok.TType != Identifier {
            return nil, fmt.Errorf("expected identifier after 'var', got %v", nameTok)
        }
        name := nameTok.Lexeme
        p.advance()

        assignTok := p.currentToken()
        if assignTok == nil || assignTok.TType != Assign {
            return nil, fmt.Errorf("expected '=' after identifier in var declaration")
        }
        p.advance()

        value, err := p.parseStatement()
        if err != nil {
            return nil, err
        }

        semiTok := p.currentToken()
        if semiTok == nil || semiTok.TType != Semicolon {
			return nil, fmt.Errorf("expected ';' after var declaration: %v", semiTok.Lexeme)
        }
        p.advance()
        return VarDefNode{Name: name, Value: value}, nil
	}
	case Func: {
		p.advance()
		nameTok := p.currentToken();
        if nameTok == nil || nameTok.TType != Identifier {
            return nil, fmt.Errorf(
				"expected function name after 'func', got %v", nameTok)
        }
        name := nameTok.Lexeme

		p.advance()
		if p.currentToken().TType != LParen {
            return nil, fmt.Errorf(
				"expected function parameters after function name, got %v", nameTok)
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

		return FunctionDefNode{
			Name: name,
			Parameters: params,
			Body: body,
		}, nil
	}
	case Struct: {
		node := p.parseStruct()
		return node, nil
	}
	case Identifier: {
		node := p.parseIdentifier()
		return node, nil
	}
	case String: {
		node := LiteralNode{Value: p.currentToken().Lexeme}
		p.advance()
		return node, nil
	}
	case Semicolon:
		p.advance()
		return nil, nil
	case Return:
		// skip 'return'
		p.advance()
		value := p.parseExpression(0)
		// expecting to advace ';'
		p.advance()
		return ReturnNode{
			Value: value,
		}, nil

	default:
        // Try to parse as an expression statement
        expr := p.parseExpression(0)
        if expr == nil {
            return nil, fmt.Errorf("Unknown token: %v", p.currentToken().Lexeme)
        }
        // Expect semicolon after expression statement
        if p.currentToken() == nil || p.currentToken().TType != Semicolon {
            return nil, fmt.Errorf("expected ';' after expression statement")
        }
        p.advance()
        return ExpressionStatementNode{Expr: expr}, nil
	}
}

func (p *Parser) parseBlock() *BlockNode {
	body := &BlockNode{}
	newEnv := NewEnv(p.currentEnv)
	prevEnv := p.currentEnv
	p.currentEnv = newEnv

	for p.currentToken().TType != RCurly {
		stmt, err := p.parseStatement()
		if err != nil {
			p.genError(err.Error())
		}

		body.Statements = append(body.Statements, stmt)
	}
	p.advance()
	p.currentEnv = prevEnv

	return body
}

func (p *Parser) parseIdentifier() Node {
	if p.currentToken() == nil {
		p.genError("Nil token")
		return nil
	}
	id := p.currentToken().Lexeme
	next := p.nextToken()
	if next == nil {
		p.genError("Nil token")
		return nil
	}

	switch next.TType {
	case LParen: {
		// func call
		p.advance()
		p.advance()
		var args []Node
		for p.currentToken().TType != RParen {
			if p.currentToken().TType == Comma {
				p.advance()
				continue
			}

			stmt := p.parseExpression(0)

			args = append(args, stmt)
		}
		if p.currentToken().TType == RParen {
			p.advance()
		} else {
			return nil
		}

		return FunctionCallNode{
			Name: id,
			Args: args,
		}
	}
	case Assign: {
		p.advance()
		p.advance()

		value, err := p.parseStatement()
		if err != nil {
			p.genError(err.Error())
			return nil
		}

		return AssignmentNode{
			Name: id,
			Value: value,
		}
	}
	default: {
		node := IdentifierNode {
			id,
		}
		p.advance()
		return node
	}
	}

}

func (p *Parser) parseStruct() Node {
	p.advance()
	nameTok := p.currentToken()
	if nameTok == nil || nameTok.TType != Identifier {
		p.genError(
			fmt.Sprintf("Expected struct name, got '%v'", nameTok.Lexeme))
		return nil
	}
	name := nameTok.Lexeme
	p.advance()
	p.advance()
	var fields []StructField
	var methods []*FunctionDefNode

	for p.currentToken().TType != RCurly {
		switch p.currentToken().TType {
		case Private: {

		}
		case Identifier: {
			// id := p.currentToken().Lexeme
			p.advance()
			switch p.currentToken().TType {
			case Identifier:
			}
		}
		}
	}

	return &StructDefNode{
		Name: name,
		Fields: fields,
		Methods: methods,
	}
}

func (p *Parser) parseExpression(precedence int) Node {
    tok := p.currentToken()
	fmt.Println("Parsing token: ", tok.Lexeme)
    var left Node

    // Parse prefix (numbers, unary minus, parentheses)
    switch tok.TType {
    case Digit:
		val, err := strconv.Atoi(tok.Lexeme)
		if err != nil {
			p.genError(err.Error())
			return nil
		}

        left = LiteralNode{Value: val}
        p.advance()
	case Identifier:
		left = p.parseIdentifier()
    case LParen:
        p.advance()
        left = p.parseExpression(0)
        if p.currentToken() == nil || p.currentToken().TType != RParen {
            p.genError("expected ')'")
            return nil
        }
        p.advance()
    case Minus:
        p.advance()
        expr := p.parseExpression(100) // high precedence for unary minus
        left = UnaryOpNode{Op: "-", Expr: expr}
    default:
        p.genError(fmt.Sprintf("unexpected token in expression: %v", tok.Lexeme))
        return nil
    }

    // Parse infix (binary) operators
    for {
        next := p.currentToken()
        if next == nil {
            break
        }
        opPrec := precedenceOf(next.TType)
        if opPrec <= precedence {
            break
        }

        op := next.TType
        p.advance()
        right := p.parseExpression(opPrec)
        left = BinaryOpNode{
            Op:    tokenTypeToString(op),
            Left:  left,
            Right: right,
        }
    }
    return left
}

func precedenceOf(tok TokenType) int {
    switch tok {
    case Plus, Minus:
        return 10
    case Star, Slash:
        return 20
    }
    return 0
}

func tokenTypeToString(t TokenType) string {
    switch t {
    case Plus:
        return "+"
    case Minus:
        return "-"
    case Star:
        return "*"
    case Slash:
        return "/"
    default:
        return ""
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

func (p *Parser) genError(message string) {
	msg := fmt.Sprintf("Error in %d:%d: %s",
		p.currentToken().Line,
		p.currentToken().Column,
		message)
	p.Errors = append(p.Errors, errors.New(msg))
}

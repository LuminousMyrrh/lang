package main

import (
	"errors"
	"fmt"
)

type Position struct {
	Row int
	Column int
}

type Node interface {
    String() string
}

// Root node for the entire program
type ProgramNode struct {
    Nodes []Node // List of statements or declarations
}

func (p ProgramNode) Find(name string) (Node, error) {
	for _, stmt := range p.Nodes {
		switch s := stmt.(type) {
		case *VarDefNode:
			if s.Name == name {
				return s, nil
			}
		case *FunctionDefNode:
			if s.Name == name {
				return s, nil
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("Symbol '%v' not found", name))
}

func (p *ProgramNode) Print() {
	for _, node := range p.Nodes {
		fmt.Println(node.String())
	}
}

type SemicolonNode struct {}

func (s *SemicolonNode) String() string {
	return ";";
}

// Literals
type LiteralNode struct {
	Position
    Value any // Can be int, string, bool, etc.
}

func (n *LiteralNode) String() string {
	switch val := n.Value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"\n", val)
	default:
		return fmt.Sprintf("(%v)\n", val)
	}
}

type TrueNode struct {
	Position
}

func (t *TrueNode) String() string {
	return "true";
}

type FalseNode struct {
	Position
}

func (f *FalseNode) String() string {
	return "false";
}

// Identifiers (variable or function names)
type IdentifierNode struct {
	Position
    Name string
}

func (i *IdentifierNode) String() string {
	return fmt.Sprintf("%v", i.Name)
}

type ArrayNode struct {
	Position
	Elements []Node
}

func (a *ArrayNode) String() string {
	str := ""
	for _, el := range a.Elements  {
		str += el.String()
	}

	return fmt.Sprintf("[%s]", str)
}

type ArrayAccessNode struct {
	Position
	Target Node
	Index Node
}

func (a *ArrayAccessNode) String() string {
	return fmt.Sprintf("%v[%v]", a.Target, a.Index)
}

type ArrayAssign struct {
	Position
	Target Node
	Value Node
}

func (a *ArrayAssign) String() string {
	return fmt.Sprintf("%v = %v",
		a.Target, a.Value)
}

type LogicalExprNode struct {
	Position
    Op    string
    Left  Node
    Right Node
}

func (l *LogicalExprNode) String() string {
	return fmt.Sprintf("& :%v '%s' :%v &\n", l.Left, l.Op, l.Right)
}

// Binary operations (e.g., a + b)
type BinaryOpNode struct {
	Position
    Op    string
    Left  Node
    Right Node
}

func (n *BinaryOpNode) String() string {
	return fmt.Sprintf("# :%v '%v' :%v #\n",n.Left.String(), n.Op, n.Right.String())
}

// Unary operations (e.g., -x, !flag)
type UnaryOpNode struct {
	Position
    Op   string
    Expr Node
}

func (n *UnaryOpNode) String() string {
	return fmt.Sprintf("%v %v\n", n.Op, n.Expr.String())
}

type StructDefNode struct {
	Position
	Name string
	Fields []*StructField
}

func (s *StructDefNode) String() string {
	var str string
	str += s.Name + " "
	for _, field := range s.Fields {
		str += field.Name + " "
	}
	str += "\n"

	return str
}

type StructInitNode struct {
	Position
	Name string
	InitFields []Node
}

func (s *StructInitNode) String() string {
	return fmt.Sprintf("%s - %v", s.Name, s.InitFields)
}

type StructField struct {
	Position
	Name string
	Value Node
	IsPublic bool // 0 - private 1 - pub
}

func (s *StructField) String() string {
	if s.Value != nil {
		return fmt.Sprintf(
			"%s(%b) - %v",
			s.Name,
			s.IsPublic,
			s.Value.String(),
			)
	}
	return fmt.Sprintf(
		"%s(%b) - (nil)",
		s.Name,
		s.IsPublic,
		)
}


type StructMethodDef struct {
	Position
	IsPub bool
	StructName string
	MethodName string
	Parameters []string
	Body *BlockNode
}

func (s *StructMethodDef) String() string {
    str := ""
    if s.IsPub {
        str += "pub "
    } else {
        str += "priv "
    }

    str += s.StructName + "->"
    str += s.MethodName + " "

    for _, p := range s.Parameters {
        str += p + " "
    }

    str += "\n"
    
    // Always check for nil before dereferencing
    if s.Body != nil {
        str += s.Body.String()
    } else {
        str += "<nil body>"
    }

    return str
}

// could be field call as well
// self.x
// self.y()
type StructMethodCall struct {
	Position
	Caller Node
	MethodName string
	IsField bool
	Args []Node
}

func (s *StructMethodCall) String() string {
	if s.IsField {
		return fmt.Sprintf(
			"(%s.%s)", s.Caller, s.MethodName)
	}
	return fmt.Sprintf(
		"(%s.%s: (%v))", s.Caller, s.MethodName, s.Args)
}

// Variable declaration (e.g., var x = 5)
type VarDefNode struct {
	Position
    Name  string
    Value Node
}

func (n *VarDefNode) String() string {
	return fmt.Sprintf("@%v: %v", n.Name, n.Value.String())
}

// Assignment (e.g., x = 10)
type AssignmentNode struct {
	Position
    Name  Node
    Value Node
	Op string
}

func (a *AssignmentNode) String() string {
	return fmt.Sprintf("@%v (%s): %v", a.Name, a.Op, a.Value.String())
}

// Function definition (e.g., func foo(a, b) { ... })
type FunctionDefNode struct {
	Position
    Name       string
    Parameters []string
    Body       *BlockNode
}

func (f *FunctionDefNode) String() string {
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
	Position
    Name Node
    Args []Node
}

func (f *FunctionCallNode) String() string {
	var str string

	str += f.Name.String() + " ("
	for _, arg := range f.Args {
		str += arg.String()
	}
	str += ")"

	return str
}

// Block of statements (e.g., { ... })
type BlockNode struct {
	Position
    Statements []Node
}

func (b *BlockNode) String() string {
	var str string
	for _, node := range b.Statements {
		str += node.String()
	}

	return str
}

// If/Else conditional
type IfNode struct {
	Position
    Condition   Node
    ThenBranch  *BlockNode
    ElseBranch  *BlockNode // Can be nil if no else
}

func (i *IfNode) String() string {
	elseBranch := i.ElseBranch
	if elseBranch != nil {
		return fmt.Sprintf("%v %v %v",
			i.Condition.String(),
			i.ThenBranch.String(),
			i.ElseBranch.String(),
			)
	} else {
		return fmt.Sprintf("%v %v",
			i.Condition.String(),
			i.ThenBranch.String(),
			)
	}
}

// While loop
type WhileNode struct {
	Position
    Condition Node
    Body      *BlockNode
}

func (w *WhileNode) String() string {
	return fmt.Sprintf("%v %v", w.Condition.String(), w.Body.String())
}

// For loop (basic: for var i = 0; i < 10; i = i + 1 { ... })
type ForNode struct {
	Position
    Init      Node      // e.g., VarDefNode or AssignmentNode
    Condition Node
    Post      Node      // e.g., AssignmentNode
    Body      *BlockNode
}

func (f *ForNode) String() string {
	return "";
}

type BreakNode struct {
	Position
}

func (b *BreakNode) String() string {
	return "break";
}

// Return statement
type ReturnNode struct {
	Position
    Value Node // Can be nil for "return"
}

func (r *ReturnNode) String() string {
	//return r.Value.String()
	return r.Value.String()
}

// Expression statement (e.g., a function call as a statement)
type ExpressionStatementNode struct {
	Position
    Expr Node
}

func (e *ExpressionStatementNode) String() string {
	return e.Expr.String()
}

// import filename > symbol1, symbol2, ...;
type ImportNode struct {
	Position
	File string
	Symbols []string
}

func (i *ImportNode) String() string {
	return i.File
}

type NilNode struct {
	Position
}

func (n *NilNode) String() string {
	return "nil";
}

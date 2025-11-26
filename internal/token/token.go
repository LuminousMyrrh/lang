package token

type TokenType int

const (
    Func TokenType = iota
    If
    Else
    While
    For
    Foreach
    Struct
    Class
    Interface
    Var
    Const
    Return
    Identifier
    StringTok  // renamed from String to StringTok to avoid conflict with built-in type
	EmptyStringTok

    Import
    Private
    Public

    True
    False
    Nil

    Plus
    PlusPlus // ++
    PlusEq   // +=
    Minus
    MinusMinus // --
    MinusEq    // -=
    Slash
    OpSlash
    Star

    LeftArrow  // <-
    RightArrow // ->

    LParen
    RParen

    LBrace
    RBrace

    LCurly
    RCurly

    And // &&
    Or  // ||
    Assign
    Bang
    NotEquals
    Equals
    Less
    More
    LessEq
    MoreEq
    Question
    Semicolon
    Colon

    Comma
    Dot

    Comment

    Int
	Float
)

func (t TokenType) String() string {
    switch t {
    case Func:
        return "Func"
    case If:
        return "If"
    case Else:
        return "Else"
    case While:
        return "While"
    case For:
        return "For"
    case Foreach:
        return "Foreach"
    case Struct:
        return "Struct"
    case Class:
        return "Class"
    case Interface:
        return "Interface"
    case Var:
        return "Var"
    case Const:
        return "Const"
    case Return:
        return "Return"
    case Identifier:
        return "Identifier"
    case StringTok:
        return "String"
	case EmptyStringTok:
        return "EmptyString"
    case Import:
        return "Import"
    case Private:
        return "Private"
    case Public:
        return "Public"
    case True:
        return "True"
    case False:
        return "False"
    case Nil:
        return "Nil"
    case Plus:
        return "Plus"
    case PlusPlus:
        return "PlusPlus"
    case PlusEq:
        return "PlusEq"
    case Minus:
        return "Minus"
    case MinusMinus:
        return "MinusMinus"
    case MinusEq:
        return "MinusEq"
    case Slash:
        return "Slash"
    case OpSlash:
        return "OpSlash"
    case Star:
        return "Star"
    case LeftArrow:
        return "LeftArrow"
    case RightArrow:
        return "RightArrow"
    case LParen:
        return "LParen"
    case RParen:
        return "RParen"
    case LBrace:
        return "LBrace"
    case RBrace:
        return "RBrace"
    case LCurly:
        return "LCurly"
    case RCurly:
        return "RCurly"
    case And:
        return "And"
    case Or:
        return "Or"
    case Assign:
        return "Assign"
    case Bang:
        return "Bang"
    case NotEquals:
        return "NotEquals"
    case Equals:
        return "Equals"
    case Less:
        return "Less"
    case More:
        return "More"
    case LessEq:
        return "LessEq"
    case MoreEq:
        return "MoreEq"
    case Question:
        return "Question"
    case Semicolon:
        return "Semicolon"
    case Colon:
        return "Colon"
    case Comma:
        return "Comma"
    case Dot:
        return "Dot"
    case Comment:
        return "Comment"
    case Int:
        return "Int"
    case Float:
        return "Float"
    default:
        return "Unknown"
    }
}


type Token struct {
	Lexeme string
	TType  TokenType
	Line   int
	Column int
}


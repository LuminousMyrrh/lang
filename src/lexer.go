package main

import (
	"fmt"
	"unicode"
)

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

type Lexer struct {
	source        []rune
	currentLine   int
	currentColumn int
}

func (l *Lexer) Read(s string) ([]*Token, error) {
	l.source = []rune(s)
	length := len(l.source)
	var tokens []*Token
	l.currentLine = 1
	l.currentColumn = 1

	for i := 0; i < length; i++ {
		ch := l.source[i]

		// Handle newlines first
		if ch == '\n' {
			l.currentLine++
			l.currentColumn = 1
			continue
		}

		// Skip other whitespace (except newline) but update column
		if unicode.IsSpace(ch) {
			l.currentColumn++
			continue
		}

		startColumn := l.currentColumn

		// Identifiers/Keywords
		if unicode.IsLetter(ch) {
			start := i
			for i < length && (unicode.IsLetter(l.source[i]) || unicode.IsDigit(l.source[i]) || l.source[i] == '_') {
				i++
				l.currentColumn++
			}
			word := string(l.source[start:i])
			tokenType := l.getWordType(word)
			tokens = append(tokens, l.genTokenAtPosition(word, tokenType, l.currentLine, startColumn))
			i-- // Adjust for the for loop increment
			continue
		}

		// Numbers
		if unicode.IsDigit(ch) {
			start := i
			hasDot := false
			for i < length {
				c := l.source[i]
				if unicode.IsDigit(c) {
					i++
					l.currentColumn++
				} else if c == '.' && !hasDot {
					// Check for decimal point if none encountered before
					// Also, make sure next char is digit to avoid e.g. "12."
					if i+1 < length && unicode.IsDigit(l.source[i+1]) {
						hasDot = true
						i++
						l.currentColumn++
					} else {
						break
					}
				} else {
					break
				}
			}
			word := string(l.source[start:i])
			if hasDot {
				tokens = append(tokens, l.genTokenAtPosition(word, Float, l.currentLine, startColumn))
			} else {
				tokens = append(tokens, l.genTokenAtPosition(word, Int, l.currentLine, startColumn))
			}
			i--  // adjust for the for loop increment
			continue
		}

		if ch == '"' {
			i++              
			l.currentColumn++
			var str []rune

			for i < length {
				c := l.source[i]

				if c == '"' {
					break
				}

				if c == '\n' {
					l.currentLine++
					l.currentColumn = 1
					i++
					continue
				}

				if c == '\\' && i+1 < length {
					nextChar := l.source[i+1]
					switch nextChar {
					case '"', '\\', 'n', 't', 'r':
						if nextChar == 'n' {
							str = append(str, '\n')
						} else if nextChar == 't' {
							str = append(str, '\t')
						} else if nextChar == 'r' {
							str = append(str, '\r')
						} else {
							str = append(str, nextChar)
						}
						i += 2
						l.currentColumn += 2
						continue
					}
				}

				str = append(str, c)
				i++
				l.currentColumn++
			}

			if i >= length || l.source[i] != '"' {
				return tokens, fmt.Errorf("unterminated string literal at line %d, column %d", l.currentLine, l.currentColumn)
			}
			l.currentColumn++

			if len(str) == 0 {
				tokens = append(tokens, l.genTokenAtPosition("", EmptyStringTok, l.currentLine, startColumn))
			} else {
				tokens = append(tokens, l.genTokenAtPosition(string(str), StringTok, l.currentLine, startColumn))
			}
			continue
		}

		// Handle comments (// ... to end of line)
		if ch == '/' && l.Next(i) == '/' {
			startColumn := l.currentColumn
			i += 2
			l.currentColumn += 2
			var commentRunes []rune

			for i < length && l.source[i] != '\n' {
				commentRunes = append(commentRunes, l.source[i])
				i++
				l.currentColumn++
			}
			tokens = append(tokens, l.genTokenAtPosition(string(commentRunes), Comment, l.currentLine, startColumn))
			i-- // adjust for outer loop
			continue
		}

		// Handle multi-char and single-char tokens
		switch ch {
		case '+':
			if l.Next(i) == '+' {
				tokens = append(tokens, l.genTokenAtPosition("++", PlusPlus, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("+=", PlusEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("+", Plus, l.currentLine, startColumn))
			l.currentColumn++
		case '-':
			if l.Next(i) == '-' {
				tokens = append(tokens, l.genTokenAtPosition("--", MinusMinus, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '>' {
				tokens = append(tokens, l.genTokenAtPosition("->", RightArrow, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("-=", MinusEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("-", Minus, l.currentLine, startColumn))
			l.currentColumn++
		case '/':
			if l.Next(i) == '/' {
				// Already handled comment above, but could include here if needed
				// Skipping here because handled above
				l.currentColumn++
			} else {
				tokens = append(tokens, l.genTokenAtPosition("/", Slash, l.currentLine, startColumn))
				l.currentColumn++
			}
		case '*':
			tokens = append(tokens, l.genTokenAtPosition("*", Star, l.currentLine, startColumn))
			l.currentColumn++
		case '<':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("<=", LessEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '-' {
				tokens = append(tokens, l.genTokenAtPosition("<-", LeftArrow, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("<", Less, l.currentLine, startColumn))
			l.currentColumn++
		case '>':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition(">=", MoreEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition(">", More, l.currentLine, startColumn))
			l.currentColumn++
		case '=':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("==", Equals, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("=", Assign, l.currentLine, startColumn))
			l.currentColumn++
		case '&':
			if l.Next(i) == '&' {
				tokens = append(tokens, l.genTokenAtPosition("&&", And, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("&", And, l.currentLine, startColumn))
			l.currentColumn++
		case '|':
			if l.Next(i) == '|' {
				tokens = append(tokens, l.genTokenAtPosition("||", Or, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("|", Or, l.currentLine, startColumn))
			l.currentColumn++
		case '!':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("!=", NotEquals, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("!", Bang, l.currentLine, startColumn))
			l.currentColumn++
		case '?':
			tokens = append(tokens, l.genTokenAtPosition("?", Question, l.currentLine, startColumn))
			l.currentColumn++
		case ';':
			tokens = append(tokens, l.genTokenAtPosition(";", Semicolon, l.currentLine, startColumn))
			l.currentColumn++
		case ':':
			tokens = append(tokens, l.genTokenAtPosition(":", Colon, l.currentLine, startColumn))
			l.currentColumn++
		case ',':
			tokens = append(tokens, l.genTokenAtPosition(",", Comma, l.currentLine, startColumn))
			l.currentColumn++
		case '.':
			tokens = append(tokens, l.genTokenAtPosition(".", Dot, l.currentLine, startColumn))
			l.currentColumn++
		case '(':
			tokens = append(tokens, l.genTokenAtPosition("(", LParen, l.currentLine, startColumn))
			l.currentColumn++
		case ')':
			tokens = append(tokens, l.genTokenAtPosition(")", RParen, l.currentLine, startColumn))
			l.currentColumn++
		case '{':
			tokens = append(tokens, l.genTokenAtPosition("{", LCurly, l.currentLine, startColumn))
			l.currentColumn++
		case '}':
			tokens = append(tokens, l.genTokenAtPosition("}", RCurly, l.currentLine, startColumn))
			l.currentColumn++
		case '[':
			tokens = append(tokens, l.genTokenAtPosition("[", LBrace, l.currentLine, startColumn))
			l.currentColumn++
		case ']':
			tokens = append(tokens, l.genTokenAtPosition("]", RBrace, l.currentLine, startColumn))
			l.currentColumn++
		default:
			fmt.Printf("Unknown character '%c' at line %d, column %d\n", ch, l.currentLine, l.currentColumn)
			l.currentColumn++ // skip unknown char to avoid infinite loop
		}
	}

	return tokens, nil
}

func (l *Lexer) genTokenAtPosition(lexeme string, ttype TokenType, line int, column int) *Token {
	return &Token{
		Lexeme: lexeme,
		TType:  ttype,
		Line:   line,
		Column: column,
	}
}

func (l *Lexer) Next(i int) rune {
	if i+1 < len(l.source) {
		return l.source[i+1]
	}
	return 0
}

func (l *Lexer) getWordType(word string) TokenType {
	switch word {
	case "true":
		return True
	case "false":
		return False
	case "func":
		return Func
	case "if":
		return If
	case "else":
		return Else
	case "while":
		return While
	case "for":
		return For
	case "foreach":
		return Foreach
	case "struct":
		return Struct
	case "class":
		return Struct
	case "interface":
		return Interface
	case "var":
		return Var
	case "const":
		return Const
	case "return":
		return Return
	case "import":
		return Import
	case "pri", "private":
		return Private
	case "pub", "public":
		return Public
	case "nil":
		return Nil
	default:
		return Identifier
	}
}

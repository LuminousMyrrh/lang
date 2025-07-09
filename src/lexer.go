package main

import (
	"fmt"
	"unicode"
)

type TokenType int
const (
	Func TokenType = iota
	If 
	While 
	For 
	Foreach 
	Struct 
	Interface 
	Var 
	Const 
	Return 
	Identifier 
	String

	Import 
	Private

	Plus 
	PlusPlus // ++
	Minus 
	MinusMinus // --
	Slash 
	OpSlash 
	Star 

	LParen 
	RParen 

	LBrace 
	RBrace 

	LCurly 
	RCurly 

	And 
	Or 
	Assign  // =
	Bang 
	NotEquals  // != 
	Equals  // ==
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

	Digit
)

type Token struct {
	Lexeme string
	TType TokenType 
	Line int
	Column int
}

type Lexer struct {
	source []rune
	current_line int
	current_column int
}

func (l *Lexer) Read(s string) ([]*Token, error) {
	source := []rune(s)
	var tokens []*Token
	l.source = source
	length := len(l.source) - 1

	for i := 0; i < len(l.source); i++ {
		ch := source[i]
		l.current_column = i

		switch {
		case unicode.IsSpace(ch):
			continue
		case unicode.IsLetter(ch): {
			start := i
			for i < length && (unicode.IsLetter(source[i]) || unicode.IsDigit(source[i]) || source[i] == '_') {
				i++
			}
			current_word := string(l.source[start:i])
			wordType := l.getWordType(current_word)
			tokens = append(tokens, l.genToken(current_word, wordType))
			i--
		}
		case unicode.IsDigit(ch): {
			start := i
			for i < length && unicode.IsDigit(source[i]) {
				i++
			}
			current_word := string(l.source[start:i])
			tokens = append(tokens, l.genToken(current_word, Digit))
			i--
		}
		default: 
			switch ch {
			case '\n':
				l.current_line++
				l.current_column = 0
			case '(':
				tokens = append(tokens, l.genToken(string(ch), LParen))
			case ')':
				tokens = append(tokens, l.genToken(string(ch), RParen))
			case '{':
				tokens = append(tokens, l.genToken(string(ch),LCurly))
			case '}':
				tokens = append(tokens, l.genToken(string(ch),RCurly))
			case '[':
				tokens = append(tokens, l.genToken(string(ch),LBrace))
			case ']':
				tokens = append(tokens, l.genToken(string(ch),RBrace))

			case ';':
				tokens = append(tokens, l.genToken(string(ch),Semicolon))
			case ':':
				tokens = append(tokens, l.genToken(string(ch),Colon))

			case '+':
				if l.Next(i) == '+' {
					current_word := "++"
					tokens = append(tokens, l.genToken(current_word,PlusPlus))
					i++
					continue
				} else {
					tokens = append(tokens, l.genToken(string(ch), Plus))
				}
			case '-':
				if l.Next(i) == '-' {
					current_word := "--"
					tokens = append(tokens, l.genToken(current_word, MinusMinus))
					i++
					continue
				} else {
					tokens = append(tokens, l.genToken(string(ch),Minus))
				}
			case '/':
				if l.Next(i) == '/' {
					current_word := "//"
					tokens = append(tokens, l.genToken(current_word, Comment))
					i++
					continue
				}
				tokens = append(tokens, l.genToken(string(ch), Slash))
			case '*':
				tokens = append(tokens, l.genToken(string(ch),Star))

			case '<':
				if l.Next(i) == '=' {
					current_word := "<="
					i++
					tokens = append(tokens, l.genToken(current_word, LessEq))
					continue
				}
				tokens = append(tokens, l.genToken(string(ch), Less))

			case '>':
				if l.Next(i) == '=' {
					current_word := ">="
					i++
					tokens = append(tokens, l.genToken(current_word, MoreEq))
					continue
				}
				tokens = append(tokens, l.genToken(string(ch),More))

			case '=':
				if l.Next(i) == '=' {
					current_word := "=="
					i++
					tokens = append(tokens, l.genToken(current_word, Equals))
					continue
				}
				tokens = append(tokens, l.genToken(string(ch),Assign))
			case '?':
				tokens = append(tokens, l.genToken(string(ch),Question))
			case '.':
				tokens = append(tokens, l.genToken(string(ch),Dot))
			case ',':
				tokens = append(tokens, l.genToken(string(ch),Comma))
			case '"': {
				i++
				var str []rune
				for i < length {
					if source[i] == '"' {
						break
					}
					if source[i] == '\\' && i+1 < length {
						// Handle escape sequences
						nextChar := source[i+1]
						if nextChar == '"' {
							str = append(str, '"')
							i += 2
							continue
						} else if nextChar == '\\' {
							str = append(str, '\\')
							i += 2
							continue
						}
						// Add more escape handling as needed (\n, \t, etc.)
					}
					str = append(str, source[i])
					i++
				}
				if i >= length {
					fmt.Println("Unterminated string literal")
					break
				}
				tokens = append(tokens, l.genToken(string(str), String))
			}
			default: 
				fmt.Println("unknown character: ", ch)
				continue
			}

		}
	}

	return tokens, nil
}

func (l *Lexer) genToken(word string, t_type TokenType) *Token {
	return &Token {
		word,
		t_type,
		l.current_line,
		l.current_column,
	}
}

func (l *Lexer) Next(curr_idx int) rune {
	if curr_idx + 1 < len(l.source) {
		return l.source[curr_idx + 1]
	}
	return '\a';
}

func (l *Lexer) getWordType(word string) TokenType {
	var wordType TokenType
	switch word {
	case "func":
		wordType = Func
	case "if":
		wordType = If
	case "while":
		wordType = While
	case "for":
		wordType = For
	case "foreach":
		wordType = Foreach
	case "struct":
		wordType = Struct
	case "interface":
		wordType = Interface
	case "var":
		wordType = Var
	case "const":
		wordType = Const
	case "return":
		wordType = Return
	case "import":
		wordType = Import
	case "private":
		wordType = Private
	default: 
		wordType = Identifier
	}

	return wordType
}

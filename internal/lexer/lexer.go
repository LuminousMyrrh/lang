package lexer

import (
	"fmt"
	"lang/internal/token"
	"unicode"
)

type Lexer struct {
	source        []rune
	currentLine   int
	currentColumn int
}

func (l *Lexer) Read(s string) ([]*token.Token, error) {
	l.source = []rune(s)
	length := len(l.source)
	var tokens []*token.Token
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
				tokens = append(tokens, l.genTokenAtPosition(word, token.Float, l.currentLine, startColumn))
			} else {
				tokens = append(tokens, l.genTokenAtPosition(word, token.Int, l.currentLine, startColumn))
			}
			i-- // adjust for the for loop increment
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
				tokens = append(tokens, l.genTokenAtPosition("", token.EmptyStringTok, l.currentLine, startColumn))
			} else {
				tokens = append(tokens, l.genTokenAtPosition(string(str), token.StringTok, l.currentLine, startColumn))
			}
			continue
		}

		if ch == '/' && l.Next(i) == '/' {
			i += 2
			l.currentColumn += 2

			for i < length && l.source[i] != '\n' {
				i++
				l.currentColumn++
			}

			i--
			continue
		}

		// Handle multi-char and single-char tokens
		switch ch {
		case '+':
			if l.Next(i) == '+' {
				tokens = append(tokens, l.genTokenAtPosition("++", token.PlusPlus, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("+=", token.PlusEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("+", token.Plus, l.currentLine, startColumn))
			l.currentColumn++
		case '-':
			if l.Next(i) == '-' {
				tokens = append(tokens, l.genTokenAtPosition("--", token.MinusMinus, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '>' {
				tokens = append(tokens, l.genTokenAtPosition("->", token.RightArrow, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("-=", token.MinusEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("-", token.Minus, l.currentLine, startColumn))
			l.currentColumn++
		case '/':
			if l.Next(i) == '/' {
				// Already handled comment above, but could include here if needed
				// Skipping here because handled above
				l.currentColumn++
			} else {
				tokens = append(tokens, l.genTokenAtPosition("/", token.Slash, l.currentLine, startColumn))
				l.currentColumn++
			}
		case '*':
			tokens = append(tokens, l.genTokenAtPosition("*", token.Star, l.currentLine, startColumn))
			l.currentColumn++
		case '<':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("<=", token.LessEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			} else if l.Next(i) == '-' {
				tokens = append(tokens, l.genTokenAtPosition("<-", token.LeftArrow, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("<", token.Less, l.currentLine, startColumn))
			l.currentColumn++
		case '>':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition(">=", token.MoreEq, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition(">", token.More, l.currentLine, startColumn))
			l.currentColumn++
		case '=':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("==", token.Equals, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("=", token.Assign, l.currentLine, startColumn))
			l.currentColumn++
		case '&':
			if l.Next(i) == '&' {
				tokens = append(tokens, l.genTokenAtPosition("&&", token.And, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("&", token.And, l.currentLine, startColumn))
			l.currentColumn++
		case '|':
			if l.Next(i) == '|' {
				tokens = append(tokens, l.genTokenAtPosition("||", token.Or, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("|", token.Or, l.currentLine, startColumn))
			l.currentColumn++
		case '!':
			if l.Next(i) == '=' {
				tokens = append(tokens, l.genTokenAtPosition("!=", token.NotEquals, l.currentLine, startColumn))
				i++
				l.currentColumn += 2
				continue
			}
			tokens = append(tokens, l.genTokenAtPosition("!", token.Bang, l.currentLine, startColumn))
			l.currentColumn++
		case '?':
			tokens = append(tokens, l.genTokenAtPosition("?", token.Question, l.currentLine, startColumn))
			l.currentColumn++
		case ';':
			tokens = append(tokens, l.genTokenAtPosition(";", token.Semicolon, l.currentLine, startColumn))
			l.currentColumn++
		case ':':
			tokens = append(tokens, l.genTokenAtPosition(":", token.Colon, l.currentLine, startColumn))
			l.currentColumn++
		case ',':
			tokens = append(tokens, l.genTokenAtPosition(",", token.Comma, l.currentLine, startColumn))
			l.currentColumn++
		case '.':
			tokens = append(tokens, l.genTokenAtPosition(".", token.Dot, l.currentLine, startColumn))
			l.currentColumn++
		case '(':
			tokens = append(tokens, l.genTokenAtPosition("(", token.LParen, l.currentLine, startColumn))
			l.currentColumn++
		case ')':
			tokens = append(tokens, l.genTokenAtPosition(")", token.RParen, l.currentLine, startColumn))
			l.currentColumn++
		case '{':
			tokens = append(tokens, l.genTokenAtPosition("{", token.LCurly, l.currentLine, startColumn))
			l.currentColumn++
		case '}':
			tokens = append(tokens, l.genTokenAtPosition("}", token.RCurly, l.currentLine, startColumn))
			l.currentColumn++
		case '[':
			tokens = append(tokens, l.genTokenAtPosition("[", token.LBrace, l.currentLine, startColumn))
			l.currentColumn++
		case ']':
			tokens = append(tokens, l.genTokenAtPosition("]", token.RBrace, l.currentLine, startColumn))
			l.currentColumn++
		default:
			fmt.Printf("Unknown character '%c' at line %d, column %d\n", ch, l.currentLine, l.currentColumn)
			l.currentColumn++ // skip unknown char to avoid infinite loop
		}
	}

	return tokens, nil
}

func (l *Lexer) genTokenAtPosition(lexeme string, ttype token.TokenType, line int, column int) *token.Token {
	return &token.Token{
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

func (l *Lexer) getWordType(word string) token.TokenType {
	switch word {
	case "true":
		return token.True
	case "false":
		return token.False
	case "func":
		return token.Func
	case "if":
		return token.If
	case "else":
		return token.Else
	case "while":
		return token.While
	case "for":
		return token.For
	case "foreach":
		return token.Foreach
	case "struct":
		return token.Struct
	case "class":
		return token.Struct
	case "interface":
		return token.Interface
	case "var":
		return token.Var
	case "const":
		return token.Const
	case "return":
		return token.Return
	case "import":
		return token.Import
	case "pri", "private":
		return token.Private
	case "pub", "public":
		return token.Public
	case "nil":
		return token.Nil
	default:
		return token.Identifier
	}
}

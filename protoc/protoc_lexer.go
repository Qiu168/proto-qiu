// 词法分析器
package protoc

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

type TokenType int

const (
	// TokenEOF is returned by [Lexer.NextToken] when the end of the input is reached.
	TokenEOF TokenType = iota
	// TokenIdent is returned by [Lexer.NextToken] when an identifier is encountered.
	// 标识符 命名
	TokenIdent
	// TokenNumber is returned by [Lexer.NextToken] when a number is encountered.
	TokenNumber
	// TokenString is returned by [Lexer.NextToken] when a string is encountered.
	// such as "hello world"
	TokenString
	// TokenSymbol is returned by [Lexer.NextToken] when a symbol is encountered.
	// such as '+', '-', '*', '/', etc.
	TokenSymbol // 如 '{', ';', '=', etc.
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	reader *bufio.Reader
	pos    int
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{reader: bufio.NewReader(r)}
}

func (l *Lexer) readRune() (rune, error) {
	r, _, err := l.reader.ReadRune()
	if err == nil {
		l.pos++
	}
	return r, err
}

func (l *Lexer) unreadRune() {
	_ = l.reader.UnreadRune()
	l.pos--
}

func (l *Lexer) skipWhitespace() {
	for {
		r, err := l.readRune()
		if err != nil {
			return
		}

		// 跳过单行注释 "// ..."
		if r == '/' {
			nextR, _ := l.readRune()
			if nextR == '/' {
				for {
					r, err := l.readRune()
					if err != nil || r == '\n' {
						break
					}
				}
				continue
			} else {
				l.unreadRune()
			}
		}

		// 跳过多行注释 "/* ... */"
		if r == '/' {
			nextR, _ := l.readRune()
			if nextR == '*' {
				for {
					r, err := l.readRune()
					if err != nil {
						return
					}
					if r == '*' {
						nextR, _ := l.readRune()
						if nextR == '/' {
							break
						}
						l.unreadRune()
					}
				}
				continue
			} else {
				l.unreadRune()
			}
		}

		if !unicode.IsSpace(r) {
			l.unreadRune()
			return
		}
	}
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	r, err := l.readRune()
	if err != nil {
		if err == io.EOF {
			return Token{Type: TokenEOF}, nil
		}
		return Token{}, err
	}

	switch {
	case isIdentStart(r):
		return l.readIdentifier(r)
	case unicode.IsDigit(r):
		return l.readNumber(r)
	case r == '"' || r == '\'':
		return l.readString(r)
	default:
		return Token{Type: TokenSymbol, Value: string(r)}, nil
	}
}

func isIdentStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (l *Lexer) readIdentifier(first rune) (Token, error) {
	var builder strings.Builder
	builder.WriteRune(first)

	for {
		r, err := l.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return Token{}, err
		}

		if !isIdentPart(r) {
			l.unreadRune()
			break
		}
		builder.WriteRune(r)
	}

	return Token{Type: TokenIdent, Value: builder.String()}, nil
}

func isIdentPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.'
}

func (l *Lexer) readNumber(first rune) (Token, error) {
	var builder strings.Builder
	builder.WriteRune(first)

	for {
		r, err := l.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return Token{}, err
		}

		if !unicode.IsDigit(r) {
			l.unreadRune()
			break
		}
		builder.WriteRune(r)
	}

	return Token{Type: TokenNumber, Value: builder.String()}, nil
}

func (l *Lexer) readString(quote rune) (Token, error) {
	var builder strings.Builder

	for {
		r, err := l.readRune()
		if err != nil {
			return Token{}, errors.New("unterminated string")
		}
		if r == quote {
			return Token{Type: TokenString, Value: builder.String()}, nil
		}
		builder.WriteRune(r)
	}
}

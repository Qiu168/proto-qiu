// 词法分析器
package main

import (
	"bufio"
	"errors"
	"io"
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
	case unicode.IsLetter(r) || r == '_' || r == '.':
		ident := string(r)
		for {
			r, err = l.readRune()
			if err != nil || (!unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '.') {
				if err == nil {
					l.unreadRune()
				}
				return Token{Type: TokenIdent, Value: ident}, nil
			}
			ident += string(r)
		}

	case unicode.IsDigit(r):
		number := string(r)
		for {
			r, err = l.readRune()
			if err != nil || !unicode.IsDigit(r) {
				if err == nil {
					l.unreadRune()
				}
				return Token{Type: TokenNumber, Value: number}, nil
			}
			number += string(r)
		}

	case r == '"' || r == '\'':
		quote := r
		str := ""
		for {
			r, err = l.readRune()
			if err != nil {
				return Token{}, errors.New("unterminated string")
			}
			if r == quote {
				return Token{Type: TokenString, Value: str}, nil
			}
			str += string(r)
		}

	default:
		return Token{Type: TokenSymbol, Value: string(r)}, nil
	}
}

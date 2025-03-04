package protoc

import (
	"io"
	"strings"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
		wantErr  bool
	}{
		{
			name:  "read identifier",
			input: "message",
			expected: []Token{
				{Type: TokenIdent, Value: "message"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "read number",
			input: "123",
			expected: []Token{
				{Type: TokenNumber, Value: "123"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "read string with double quotes",
			input: `"hello"`,
			expected: []Token{
				{Type: TokenString, Value: "hello"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "read string with single quotes",
			input: `'world'`,
			expected: []Token{
				{Type: TokenString, Value: "world"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "read symbols",
			input: "{}=;",
			expected: []Token{
				{Type: TokenSymbol, Value: "{"},
				{Type: TokenSymbol, Value: "}"},
				{Type: TokenSymbol, Value: "="},
				{Type: TokenSymbol, Value: ";"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "read mixed tokens",
			input: "message User { id = 123; }",
			expected: []Token{
				{Type: TokenIdent, Value: "message"},
				{Type: TokenIdent, Value: "User"},
				{Type: TokenSymbol, Value: "{"},
				{Type: TokenIdent, Value: "id"},
				{Type: TokenSymbol, Value: "="},
				{Type: TokenNumber, Value: "123"},
				{Type: TokenSymbol, Value: ";"},
				{Type: TokenSymbol, Value: "}"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "handle whitespace",
			input: "  \t\n  message  \n",
			expected: []Token{
				{Type: TokenIdent, Value: "message"},
				{Type: TokenEOF},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []Token{{Type: TokenEOF}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))

			for _, expected := range tt.expected {
				token, err := l.NextToken()

				if tt.wantErr {
					if err == nil {
						t.Errorf("expected error but got none")
					}
					return
				}

				if err != nil && err != io.EOF {
					t.Errorf("unexpected error: %v", err)
					return
				}

				if token.Type != expected.Type {
					t.Errorf("token type mismatch - got: %v, want: %v", token.Type, expected.Type)
				}

				if token.Value != expected.Value {
					t.Errorf("token value mismatch - got: %v, want: %v", token.Value, expected.Value)
				}
			}

			// 确保没有多余的token
			extraToken, _ := l.NextToken()
			if extraToken.Type != TokenEOF {
				t.Errorf("expected EOF but got token: %v", extraToken)
			}
		})
	}
}

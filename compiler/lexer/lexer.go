package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// Lexer tokenizes BASIC source code
type Lexer struct {
	input   string
	pos     int
	line    int
	col     int
	current rune
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
		col:   1,
	}
	l.readChar()
	return l
}

// readChar advances the lexer position
func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.current = 0
	} else {
		l.current = rune(l.input[l.pos])
	}
	l.pos++
	l.col++
}

// peekChar looks at the next character without advancing
func (l *Lexer) peekChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	for {
		l.skipWhitespace()
		// Line comment: //
		if l.current == '/' && l.peekChar() == '/' {
			l.skipLineComment()
			continue
		}
		// Block comment: /* */
		if l.current == '/' && l.peekChar() == '*' {
			l.skipBlockComment()
			continue
		}
		break
	}

	switch l.current {
	case 0:
		tok = Token{Type: TokenEOF, Line: l.line, Col: l.col}
	case '=':
		if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: TokenNotEqual, Value: "<>", Line: l.line, Col: l.col - 2}
		} else {
			tok = Token{Type: TokenAssign, Value: "=", Line: l.line, Col: l.col - 1}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenLessEqual, Value: "<=", Line: l.line, Col: l.col - 2}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: TokenNotEqual, Value: "<>", Line: l.line, Col: l.col - 2}
		} else {
			tok = Token{Type: TokenLess, Value: "<", Line: l.line, Col: l.col - 1}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenGreaterEqual, Value: ">=", Line: l.line, Col: l.col - 2}
		} else {
			tok = Token{Type: TokenGreater, Value: ">", Line: l.line, Col: l.col - 1}
		}
	case '+':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenPlusAssign, Value: "+=", Line: l.line, Col: l.col - 2}
		} else {
			tok = Token{Type: TokenPlus, Value: "+", Line: l.line, Col: l.col - 1}
		}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenMinusAssign, Value: "-=", Line: l.line, Col: l.col - 2}
		} else {
			tok = Token{Type: TokenMinus, Value: "-", Line: l.line, Col: l.col - 1}
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenStarAssign, Value: "*=", Line: l.line, Col: l.col - 2}
		} else {
			tok = Token{Type: TokenMultiply, Value: "*", Line: l.line, Col: l.col - 1}
		}
	case '/':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenSlashAssign, Value: "/=", Line: l.line, Col: l.col - 2}
		} else {
			tok = Token{Type: TokenDivide, Value: "/", Line: l.line, Col: l.col - 1}
		}
	case '%':
		tok = Token{Type: TokenMod, Value: "%", Line: l.line, Col: l.col - 1}
	case '^':
		tok = Token{Type: TokenPower, Value: "^", Line: l.line, Col: l.col - 1}
	case '\\':
		tok = Token{Type: TokenIntDiv, Value: "\\", Line: l.line, Col: l.col - 1}
	case '(':
		tok = Token{Type: TokenLeftParen, Value: "(", Line: l.line, Col: l.col - 1}
	case ')':
		tok = Token{Type: TokenRightParen, Value: ")", Line: l.line, Col: l.col - 1}
	case '[':
		tok = Token{Type: TokenLeftBracket, Value: "[", Line: l.line, Col: l.col - 1}
	case ']':
		tok = Token{Type: TokenRightBracket, Value: "]", Line: l.line, Col: l.col - 1}
	case '{':
		tok = Token{Type: TokenLeftBrace, Value: "{", Line: l.line, Col: l.col - 1}
	case '}':
		tok = Token{Type: TokenRightBrace, Value: "}", Line: l.line, Col: l.col - 1}
	case ',':
		tok = Token{Type: TokenComma, Value: ",", Line: l.line, Col: l.col - 1}
	case ':':
		tok = Token{Type: TokenColon, Value: ":", Line: l.line, Col: l.col - 1}
	case ';':
		tok = Token{Type: TokenSemicolon, Value: ";", Line: l.line, Col: l.col - 1}
	case '.':
		tok = Token{Type: TokenDot, Value: ".", Line: l.line, Col: l.col - 1}
	case '\n':
		tok = Token{Type: TokenNewLine, Value: "\n", Line: l.line, Col: l.col - 1}
		l.line++
		l.col = 1
	case '"':
		tok.Type = TokenString
		tok.Value = l.readString()
		tok.Line = l.line
		tok.Col = l.col - len(tok.Value) - 2
		return tok
	case '&':
		colStart := l.col - 1
		if num, ok := l.readAmpersandNumber(); ok {
			tok = Token{Type: TokenNumber, Value: num, Line: l.line, Col: colStart}
			return tok
		}
		tok = Token{Type: TokenUnknown, Value: "&", Line: l.line, Col: colStart}
	default:
		if unicode.IsDigit(l.current) {
			num := l.readNumber()
			if strings.Contains(num, ".") {
				tok = Token{Type: TokenNumber, Value: num, Line: l.line, Col: l.col - len(num)}
			} else {
				tok = Token{Type: TokenNumber, Value: num, Line: l.line, Col: l.col - len(num)}
			}
			return tok
		} else if unicode.IsLetter(l.current) {
			ident := l.readIdentifier()
			canon := strings.ToLower(ident)
			tokType, exists := KeywordMap[strings.ToUpper(ident)]
			if exists {
				tok = Token{Type: tokType, Value: canon, Line: l.line, Col: l.col - len(ident)}
			} else {
				// Case-insensitive: canonical form is lowercase (MyVar and myvar -> "myvar")
				tok = Token{Type: TokenIdentifier, Value: canon, Line: l.line, Col: l.col - len(ident)}
			}
			return tok
		} else {
			tok = Token{Type: TokenUnknown, Value: string(l.current), Line: l.line, Col: l.col - 1}
		}
	}

	l.readChar()
	return tok
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.pos - 1
	for unicode.IsLetter(l.current) || unicode.IsDigit(l.current) || l.current == '_' {
		l.readChar()
	}
	return l.input[position : l.pos-1]
}

// readNumber reads a number: decimal, hex (0x/0X or &H/&h), or binary (0b/0B or &B/&b).
// Returns the value as a decimal string for the rest of the pipeline.
func (l *Lexer) readNumber() string {
	// Hex: 0x / 0X
	if l.current == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar()
		l.readChar()
		hex := l.readHexDigits()
		return l.parseIntToDecimal(hex, 16)
	}
	if l.current == '0' && (l.peekChar() == 'b' || l.peekChar() == 'B') {
		l.readChar()
		l.readChar()
		bin := l.readBinaryDigits()
		return l.parseIntToDecimal(bin, 2)
	}
	// Decimal
	position := l.pos - 1
	for unicode.IsDigit(l.current) {
		l.readChar()
	}
	if l.current == '.' {
		l.readChar()
		for unicode.IsDigit(l.current) {
			l.readChar()
		}
	}
	return l.input[position : l.pos-1]
}

func (l *Lexer) readHexDigits() string {
	var b strings.Builder
	for l.current != 0 && (unicode.IsDigit(l.current) || (l.current >= 'a' && l.current <= 'f') || (l.current >= 'A' && l.current <= 'F')) {
		b.WriteRune(l.current)
		l.readChar()
	}
	return b.String()
}

func (l *Lexer) readBinaryDigits() string {
	var b strings.Builder
	for l.current == '0' || l.current == '1' {
		b.WriteRune(l.current)
		l.readChar()
	}
	return b.String()
}

// readAmpersandNumber reads &H... (hex) or &B... (binary). Caller has seen '&'.
// Returns (decimal string, true) or ("", false) if not a number. Does not consume & on false.
func (l *Lexer) readAmpersandNumber() (string, bool) {
	next := l.peekChar()
	if next == 'H' || next == 'h' {
		l.readChar() // &
		l.readChar() // H/h
		hex := l.readHexDigits()
		return l.parseIntToDecimal(hex, 16), true
	}
	if next == 'B' || next == 'b' {
		l.readChar() // &
		l.readChar() // B/b
		bin := l.readBinaryDigits()
		return l.parseIntToDecimal(bin, 2), true
	}
	return "", false
}

// parseIntToDecimal parses s as base (16 or 2) and returns decimal string.
func (l *Lexer) parseIntToDecimal(s string, base int) string {
	if s == "" {
		return "0"
	}
	var n int64
	for _, r := range s {
		n *= int64(base)
		if base == 2 {
			n += int64(r - '0')
		} else {
			switch {
			case r >= '0' && r <= '9':
				n += int64(r - '0')
			case r >= 'a' && r <= 'f':
				n += int64(r-'a') + 10
			case r >= 'A' && r <= 'F':
				n += int64(r-'A') + 10
			}
		}
	}
	return fmt.Sprintf("%d", n)
}

// readString reads a string literal and resolves escape sequences.
// Supports: \", \\, \n, \t, \r, and \xHH (two hex digits).
func (l *Lexer) readString() string {
	var b strings.Builder
	l.readChar() // consume opening "
	for l.current != '"' && l.current != 0 {
		if l.current == '\\' {
			l.readChar()
			switch l.current {
			case '"':
				b.WriteRune('"')
			case '\\':
				b.WriteRune('\\')
			case 'n':
				b.WriteRune('\n')
			case 't':
				b.WriteRune('\t')
			case 'r':
				b.WriteRune('\r')
			case 'x', 'X':
				l.readChar()
				hex := ""
				for l.current != 0 && (unicode.IsDigit(l.current) || (l.current >= 'a' && l.current <= 'f') || (l.current >= 'A' && l.current <= 'F')) && len(hex) < 2 {
					hex += string(l.current)
					l.readChar()
				}
				if len(hex) > 0 {
					var code int
					fmt.Sscanf(hex, "%x", &code)
					b.WriteRune(rune(code & 0xFF))
				}
				continue
			default:
				b.WriteRune(l.current)
			}
			l.readChar()
			continue
		}
		b.WriteRune(l.current)
		l.readChar()
	}
	l.readChar() // consume closing "
	return b.String()
}

// skipWhitespace skips whitespace characters (except newlines)
func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\r' {
		l.readChar()
	}
}

// skipLineComment skips from // to end of line
func (l *Lexer) skipLineComment() {
	for l.current != 0 && l.current != '\n' {
		l.readChar()
	}
}

// skipBlockComment skips from /* to */
func (l *Lexer) skipBlockComment() {
	l.readChar() // /
	l.readChar() // *
	for l.current != 0 {
		if l.current == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			break
		}
		l.readChar()
	}
}

// Tokenize returns all tokens from the input
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == TokenEOF {
			break
		}

		if tok.Type == TokenUnknown {
			return nil, fmt.Errorf("unknown token at line %d, col %d: %s", tok.Line, tok.Col, tok.Value)
		}
	}

	return tokens, nil
}

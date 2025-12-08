// internal/lexer/lexer.go

// Package lexer provides tokenization for numio expressions.
package lexer

import (
	"strings"
	"unicode"

	"github.com/0xsj/numio/internal/token"
	"github.com/0xsj/numio/internal/types"
)

// Lexer tokenizes input strings.
type Lexer struct {
	input   string
	pos     int  // Current position in input (points to current char)
	readPos int  // Current reading position (after current char)
	ch      rune // Current character under examination
	line    int  // Current line number (for error reporting)
	col     int  // Current column number
}

// New creates a new Lexer for the given input.
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
		col:   0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = rune(l.input[l.readPos])
		// Handle multi-byte UTF-8 characters
		if l.ch >= 0x80 {
			r, size := decodeRune(l.input[l.readPos:])
			l.ch = r
			l.readPos += size - 1 // -1 because we'll add 1 below
		}
	}
	l.pos = l.readPos
	l.readPos++
	l.col++

	// Track newlines
	if l.ch == '\n' {
		l.line++
		l.col = 0
	}
}

// decodeRune decodes a UTF-8 rune from a string.
// Simple implementation without unicode/utf8 package.
func decodeRune(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}

	b := s[0]

	// Single byte (ASCII)
	if b < 0x80 {
		return rune(b), 1
	}

	// Multi-byte
	var r rune
	var size int

	if b&0xE0 == 0xC0 {
		// 2-byte sequence
		size = 2
		r = rune(b & 0x1F)
	} else if b&0xF0 == 0xE0 {
		// 3-byte sequence
		size = 3
		r = rune(b & 0x0F)
	} else if b&0xF8 == 0xF0 {
		// 4-byte sequence
		size = 4
		r = rune(b & 0x07)
	} else {
		return unicode.ReplacementChar, 1
	}

	if len(s) < size {
		return unicode.ReplacementChar, 1
	}

	for i := 1; i < size; i++ {
		if s[i]&0xC0 != 0x80 {
			return unicode.ReplacementChar, 1
		}
		r = r<<6 | rune(s[i]&0x3F)
	}

	return r, size
}

// peekChar returns the next character without advancing.
func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	ch := rune(l.input[l.readPos])
	if ch >= 0x80 {
		r, _ := decodeRune(l.input[l.readPos:])
		return r
	}
	return ch
}

// peekCharN returns the character N positions ahead.
func (l *Lexer) peekCharN(n int) rune {
	pos := l.readPos + n - 1
	if pos >= len(l.input) {
		return 0
	}
	ch := rune(l.input[pos])
	if ch >= 0x80 {
		r, _ := decodeRune(l.input[pos:])
		return r
	}
	return ch
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	startPos := l.pos

	// Check for EOF
	if l.ch == 0 {
		return token.New(token.EOF, "", startPos)
	}

	// Check for comments
	if l.ch == '#' || (l.ch == '/' && l.peekChar() == '/') {
		return l.readComment(startPos)
	}

	// Check for currency symbols (must be before operators)
	if types.IsCurrencySymbolRune(l.ch) || types.IsCryptoSymbolRune(l.ch) {
		return l.readCurrencySymbol(startPos)
	}

	// Check for numbers (including negative and decimals starting with .)
	if isDigit(l.ch) || (l.ch == '.' && isDigit(l.peekChar())) {
		return l.readNumber(startPos)
	}

	// Check for operators and punctuation
	switch l.ch {
	case '+':
		l.readChar()
		return token.New(token.PLUS, "+", startPos)

	case '-':
		// Could be minus or negative number
		// If followed by digit and previous was operator or start, treat as number
		if isDigit(l.peekChar()) && l.isStartOfExpression() {
			return l.readNumber(startPos)
		}
		l.readChar()
		return token.New(token.MINUS, "-", startPos)

	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			l.readChar()
			return token.New(token.POWER, "**", startPos)
		}
		l.readChar()
		return token.New(token.STAR, "*", startPos)

	case '/':
		l.readChar()
		return token.New(token.SLASH, "/", startPos)

	case '^':
		l.readChar()
		return token.New(token.CARET, "^", startPos)

	case '(':
		l.readChar()
		return token.New(token.LPAREN, "(", startPos)

	case ')':
		l.readChar()
		return token.New(token.RPAREN, ")", startPos)

	case '=':
		l.readChar()
		return token.New(token.EQUALS, "=", startPos)

	case ',':
		l.readChar()
		return token.New(token.COMMA, ",", startPos)

	case '%':
		l.readChar()
		return token.New(token.PERCENT, "%", startPos)

	case '\n':
		l.readChar()
		return token.New(token.NEWLINE, "\n", startPos)
	}

	// Check for identifiers and keywords
	if isLetter(l.ch) || l.ch == '_' {
		return l.readIdentifier(startPos)
	}

	// Unknown character
	ch := l.ch
	l.readChar()
	return token.New(token.ILLEGAL, string(ch), startPos)
}

// Tokenize returns all tokens from the input.
func (l *Lexer) Tokenize() []token.Token {
	var tokens []token.Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}

// skipWhitespace skips spaces and tabs (but not newlines).
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// isStartOfExpression returns true if we're at a position where
// a negative number would make sense (vs subtraction operator).
func (l *Lexer) isStartOfExpression() bool {
	// At the very beginning
	if l.pos <= 1 {
		return true
	}

	// Look back at previous non-whitespace character
	for i := l.pos - 2; i >= 0; i-- {
		ch := rune(l.input[i])
		if ch == ' ' || ch == '\t' {
			continue
		}
		// After operator or open paren, it's start of expression
		return ch == '+' || ch == '-' || ch == '*' || ch == '/' ||
			ch == '^' || ch == '(' || ch == '=' || ch == ','
	}

	return true
}

// readNumber reads a number token (integer, decimal, or with thousands separators).
func (l *Lexer) readNumber(startPos int) token.Token {
	var sb strings.Builder

	// Handle leading negative sign
	if l.ch == '-' {
		sb.WriteRune(l.ch)
		l.readChar()
	}

	// Read integer part (with possible comma separators)
	hasDigits := false
	for isDigit(l.ch) || l.ch == ',' {
		if l.ch == ',' {
			// Validate comma placement (should have digits after)
			if !isDigit(l.peekChar()) {
				break
			}
			// Skip comma in output (or keep for parsing later)
			l.readChar()
			continue
		}
		sb.WriteRune(l.ch)
		hasDigits = true
		l.readChar()
	}

	// Read decimal part
	if l.ch == '.' && isDigit(l.peekChar()) {
		sb.WriteRune(l.ch)
		l.readChar()

		for isDigit(l.ch) {
			sb.WriteRune(l.ch)
			l.readChar()
		}
	} else if l.ch == '.' && !hasDigits {
		// Leading decimal: .5
		sb.WriteRune('0')
		sb.WriteRune(l.ch)
		l.readChar()

		for isDigit(l.ch) {
			sb.WriteRune(l.ch)
			l.readChar()
		}
	}

	// Read exponent (scientific notation)
	if l.ch == 'e' || l.ch == 'E' {
		sb.WriteRune(l.ch)
		l.readChar()

		// Optional sign
		if l.ch == '+' || l.ch == '-' {
			sb.WriteRune(l.ch)
			l.readChar()
		}

		// Exponent digits
		for isDigit(l.ch) {
			sb.WriteRune(l.ch)
			l.readChar()
		}
	}

	// Check for immediate percent sign (20% as single token)
	if l.ch == '%' {
		literal := sb.String()
		l.readChar()
		return token.New(token.PERCENT, literal+"%", startPos)
	}

	return token.New(token.NUMBER, sb.String(), startPos)
}

// readIdentifier reads an identifier or keyword.
func (l *Lexer) readIdentifier(startPos int) token.Token {
	var sb strings.Builder

	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		sb.WriteRune(l.ch)
		l.readChar()
	}

	literal := sb.String()
	lower := strings.ToLower(literal)

	// Check for keywords
	if tokType := token.LookupIdentifier(lower); tokType != token.IDENTIFIER {
		return token.New(tokType, literal, startPos)
	}

	// Check for multi-word keywords by peeking ahead
	// e.g., "turkish lira", "hong kong dollar"
	if multiWord := l.tryReadMultiWordIdentifier(literal); multiWord != "" {
		return token.New(token.IDENTIFIER, multiWord, startPos)
	}

	return token.New(token.IDENTIFIER, literal, startPos)
}

// tryReadMultiWordIdentifier tries to read a multi-word identifier.
// Returns the full identifier if found, empty string otherwise.
func (l *Lexer) tryReadMultiWordIdentifier(first string) string {
	lower := strings.ToLower(first)

	// Save state in case we need to backtrack
	savedPos := l.pos
	savedReadPos := l.readPos
	savedCh := l.ch
	savedCol := l.col

	// Check for known multi-word patterns
	multiWordPrefixes := map[string][]string{
		"turkish":   {"lira"},
		"hong":      {"kong", "dollar"},
		"new":       {"zealand", "dollar"},
		"south":     {"african", "rand", "korean", "won"},
		"saudi":     {"riyal"},
		"swiss":     {"franc", "francs"},
		"british":   {"pound", "pounds"},
		"us":        {"dollar", "dollars"},
		"mexican":   {"peso"},
		"brazilian": {"real"},
		"indian":    {"rupee", "rupees"},
		"square":    {"meter", "meters", "foot", "feet", "mile", "miles", "kilometer", "kilometers"},
		"cubic":     {"meter", "meters"},
		"fluid":     {"ounce", "ounces"},
		"troy":      {"ounce", "ounces"},
		"nautical":  {"mile", "miles"},
	}

	expectedWords, ok := multiWordPrefixes[lower]
	if !ok {
		return ""
	}

	// Try to read the expected words
	var words []string
	words = append(words, first)

	for _, expected := range expectedWords {
		l.skipWhitespace()

		if !isLetter(l.ch) {
			break
		}

		wordStart := l.pos
		var sb strings.Builder
		for isLetter(l.ch) || l.ch == '_' {
			sb.WriteRune(l.ch)
			l.readChar()
		}
		word := sb.String()

		if strings.ToLower(word) == expected {
			words = append(words, word)
		} else {
			// Backtrack this word
			l.pos = wordStart
			l.readPos = wordStart + 1
			if wordStart < len(l.input) {
				l.ch = rune(l.input[wordStart])
			} else {
				l.ch = 0
			}
			break
		}
	}

	// If we only got the first word, backtrack completely
	if len(words) == 1 {
		l.pos = savedPos
		l.readPos = savedReadPos
		l.ch = savedCh
		l.col = savedCol
		return ""
	}

	return strings.Join(words, " ")
}

// readComment reads a comment until end of line.
func (l *Lexer) readComment(startPos int) token.Token {
	var sb strings.Builder

	// Include the comment marker
	if l.ch == '#' {
		sb.WriteRune(l.ch)
		l.readChar()
	} else if l.ch == '/' {
		sb.WriteRune(l.ch)
		l.readChar()
		if l.ch == '/' {
			sb.WriteRune(l.ch)
			l.readChar()
		}
	}

	// Read until end of line or EOF
	for l.ch != '\n' && l.ch != 0 {
		sb.WriteRune(l.ch)
		l.readChar()
	}

	return token.New(token.COMMENT, sb.String(), startPos)
}

// readCurrencySymbol reads a currency symbol token.
func (l *Lexer) readCurrencySymbol(startPos int) token.Token {
	symbol := string(l.ch)
	l.readChar()

	// Determine specific token type
	tokType := token.LookupCurrencySymbol(rune(symbol[0]))

	return token.New(tokType, symbol, startPos)
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_' ||
		unicode.IsLetter(ch)
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// Tokenize is a convenience function to tokenize a string.
func Tokenize(input string) []token.Token {
	return New(input).Tokenize()
}

// TokenizeNoComments returns tokens excluding comments.
func TokenizeNoComments(input string) []token.Token {
	all := Tokenize(input)
	filtered := make([]token.Token, 0, len(all))

	for _, tok := range all {
		if tok.Type != token.COMMENT {
			filtered = append(filtered, tok)
		}
	}

	return filtered
}
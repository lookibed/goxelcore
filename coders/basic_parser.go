package coders

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsingError corresponds to C++ parsing_error in voxelcore/src/coders/BasicParser.hpp
type ParsingError struct {
	Filename string
	LineNum  uint
	Message  string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("file %s: %s at line %d", e.Filename, e.Message, e.LineNum)
}

// BasicParser is a generic character-based parser.
// Corresponds to C++ BasicParser<char> class in voxelcore/src/coders/BasicParser.hpp
type BasicParser struct {
	filename   string
	source     []rune // Using []rune for easier Unicode handling
	pos        uint
	line       uint
	linestart  uint
	hashComment bool
	clikeComment bool
}

// NewBasicParser creates a new BasicParser instance.
// Corresponds to C++ BasicParser constructor.
func NewBasicParser(filename string, source string) *BasicParser {
	return &BasicParser{
		filename:   filename,
		source:     []rune(source),
		pos:        0,
		line:       1,
		linestart:  0,
		hashComment: false,
		clikeComment: false,
	}
}

func (p *BasicParser) error(message string) error {
	return &ParsingError{
		Filename: p.filename,
		LineNum:  p.line,
		Message:  message,
	}
}

// hasNext checks if there are more characters to parse.
// Corresponds to C++ BasicParser::hasNext()
func (p *BasicParser) HasNext() bool {
	return p.pos < uint(len(p.source))
}

// peek returns the next character without advancing the position.
// Corresponds to C++ BasicParser::peek()
func (p *BasicParser) Peek() rune {
	if !p.HasNext() {
		return 0 // Null character
	}
	return p.source[p.pos]
}

// peekInLine returns the next character in the current line without advancing.
// Corresponds to C++ BasicParser::peekInLine()
func (p *BasicParser) PeekInLine() rune {
	if !p.HasNext() {
		return 0
	}
	if p.source[p.pos] == '\n' || p.source[p.pos] == '\r' {
		return 0
	}
	return p.source[p.pos]
}

// peekNoJump returns the next character without advancing the position.
// Corresponds to C++ BasicParser::peekNoJump()
func (p *BasicParser) PeekNoJump() rune {
	return p.Peek()
}


// nextChar returns the next character and advances the position.
// Corresponds to C++ BasicParser::nextChar()
func (p *BasicParser) NextChar() rune {
	if !p.HasNext() {
		return 0
	}
	char := p.source[p.pos]
	p.pos++
	if char == '\n' {
		p.line++
		p.linestart = p.pos
	}
	return char
}

// skip advances the position by n characters.
// Corresponds to C++ BasicParser::skip()
func (p *BasicParser) Skip(n uint) {
	for i := uint(0); i < n && p.HasNext(); i++ {
		p.NextChar()
	}
}

// skipLine advances the position to the end of the current line.
// Corresponds to C++ BasicParser::skipLine()
func (p *BasicParser) SkipLine() {
	for p.HasNext() {
		char := p.NextChar()
		if char == '\n' {
			break
		}
	}
}

// skipWhitespaceBasic skips basic whitespace characters.
func (p *BasicParser) skipWhitespaceBasic(newline bool) {
	for p.HasNext() {
		char := p.Peek()
		if char == ' ' || char == '\t' || (newline && (char == '\n' || char == '\r')) {
			p.NextChar()
		} else {
			break
		}
	}
}

// skipWhitespaceHashComment skips whitespace and single-line comments starting with #.
func (p *BasicParser) skipWhitespaceHashComment(newline bool) {
	for {
		p.skipWhitespaceBasic(newline)
		if p.HasNext() && p.Peek() == '#' {
			p.SkipLine()
		} else {
			break
		}
	}
}

// skipWhitespaceCLikeComment skips whitespace and C-style comments /* */.
func (p *BasicParser) skipWhitespaceCLikeComment(newline bool) {
	for {
		p.skipWhitespaceBasic(newline)
		if p.HasNext() && p.Peek() == '/' && p.pos+1 < uint(len(p.source)) && p.source[p.pos+1] == '*' {
			p.Skip(2) // Skip "/*"
			for p.HasNext() {
				if p.Peek() == '*' && p.pos+1 < uint(len(p.source)) && p.source[p.pos+1] == '/' {
					p.Skip(2) // Skip "*/"
					break
				}
			p.NextChar()
			}
		} else {
			break
		}
	}
}

// SkipWhitespace skips whitespace and comments based on parser flags.
// Corresponds to C++ BasicParser::skipWhitespace()
func (p *BasicParser) SkipWhitespace(newline bool) {
	if p.hashComment {
		p.skipWhitespaceHashComment(newline)
	} else if p.clikeComment {
		p.skipWhitespaceCLikeComment(newline)
	} else {
		p.skipWhitespaceBasic(newline)
	}
}

// skipEmptyLines skips empty lines.
// Corresponds to C++ BasicParser::skipEmptyLines()
func (p *BasicParser) SkipEmptyLines() {
	for p.HasNext() && (p.Peek() == '\n' || p.Peek() == '\r') {
		p.NextChar()
	}
}

// Expect expects a specific character.
// Corresponds to C++ BasicParser::expect(CharT expected)
func (p *BasicParser) Expect(expected rune) error {
	if p.Peek() != expected {
		return p.error(fmt.Sprintf("'%%c' expected", expected))
	}
	p.NextChar()
	return nil
}

// ExpectString expects a specific substring.
// Corresponds to C++ BasicParser::expect(const StringT& substring)
func (p *BasicParser) ExpectString(expected string) error {
	if !p.IsNext(expected) {
		return p.error(fmt.Sprintf("'%%s' expected", expected))
	}
	p.Skip(uint(len(expected)))
	return nil
}

// IsNext checks if the next characters match a substring without advancing.
// Corresponds to C++ BasicParser::isNext()
func (p *BasicParser) IsNext(s string) bool {
	if p.pos+uint(len(s)) > uint(len(p.source)) {
		return false
	}
	return string(p.source[p.pos:p.pos+uint(len(s))]) == s
}

// ExpectNewLine expects a new line.
// Corresponds to C++ BasicParser::expectNewLine()
func (p *BasicParser) ExpectNewLine() error {
	if p.Peek() == '\r' {
		p.Skip(1)
	}
	if p.Peek() != '\n' {
		return p.error("'\n' expected")
	}
	p.Skip(1)
	return nil
}

// GoBack moves the position back by count characters.
// Corresponds to C++ BasicParser::goBack()
func (p *BasicParser) GoBack(count uint) {
	if p.pos < count {
		p.pos = 0
	} else {
		p.pos -= count
	}
	// Recalculate line and linestart if necessary (complex, simplify for now)
	p.line = 1
	p.linestart = 0
	for i := uint(0); i < p.pos; i++ {
		if p.source[i] == '\n' {
			p.line++
			p.linestart = i + 1
		}
	}
}

// ReadUntil reads until a specific character is found.
// Corresponds to C++ BasicParser::readUntil(CharT c)
func (p *BasicParser) ReadUntil(c rune) string {
	start := p.pos
	for p.HasNext() && p.Peek() != c {
		p.NextChar()
	}
	return string(p.source[start:p.pos])
}

// ReadUntilString reads until a specific substring is found.
// Corresponds to C++ BasicParser::readUntil(StringViewT s, bool nothrow)
func (p *BasicParser) ReadUntilString(s string, nothrow bool) (string, error) {
	start := p.pos
	for p.HasNext() && !p.IsNext(s) {
		p.NextChar()
	}
	if !p.IsNext(s) && !nothrow {
		return "", p.error(fmt.Sprintf("'%%s' expected", s))
	}
	return string(p.source[start:p.pos]), nil
}

// ReadUntilWhitespace reads until a whitespace character is found.
// Corresponds to C++ BasicParser::readUntilWhitespace()
func (p *BasicParser) ReadUntilWhitespace() string {
	start := p.pos
	for p.HasNext() {
		char := p.Peek()
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			break
		}
		p.NextChar()
	}
	return string(p.source[start:p.pos])
}

// ReadUntilEOL reads until the end of the line.
// Corresponds to C++ BasicParser::readUntilEOL()
func (p *BasicParser) ReadUntilEOL() string {
	start := p.pos
	for p.HasNext() {
		char := p.Peek()
		if char == '\n' || char == '\r' {
			break
		}
		p.NextChar()
	}
	return string(p.source[start:p.pos])
}

// ParseName parses an identifier (name).
// Corresponds to C++ BasicParser::parseName()
func (p *BasicParser) ParseName() string {
	start := p.pos
	if !p.HasNext() {
		return ""
	}
	// First character must be a letter or underscore
	char := p.Peek()
	if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_') {
		return ""
	}
	p.NextChar()

	// Subsequent characters can be letters, numbers, or underscores
	for p.HasNext() {
		char = p.Peek()
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' {
			p.NextChar()
		} else {
			break
		}
	}
	return string(p.source[start:p.pos])
}

// ParseXmlName parses an XML-style name (can include hyphens).
// Corresponds to C++ BasicParser::parseXmlName()
func (p *BasicParser) ParseXmlName() string {
	start := p.pos
	if !p.HasNext() {
		return ""
	}
	char := p.Peek()
	if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_') {
		return ""
	}
	p.NextChar()
	for p.HasNext() {
		char = p.Peek()
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
			p.NextChar()
		} else {
			break
		}
	}
	return string(p.source[start:p.pos])
}


// ParseSimpleInt parses a simple integer.
// Corresponds to C++ BasicParser::parseSimpleInt()
func (p *BasicParser) ParseSimpleInt(base int, maxLength uint) (int64, error) {
	start := p.pos
	// Handle sign
	sign := int64(1)
	if p.HasNext() && p.Peek() == '-' {
		sign = -1
		p.NextChar()
	}

	startDigit := p.pos
	for p.HasNext() {
		char := p.Peek()
		if (char >= '0' && char <= '9') {
			p.NextChar()
		} else {
			break
		}
	}
	numStr := string(p.source[startDigit:p.pos])
	if len(numStr) == 0 {
		return 0, p.error("expected integer")
	}

	val, err := strconv.ParseInt(numStr, base, 64)
	if err != nil {
		return 0, p.error(fmt.Sprintf("invalid integer: %%v", err))
	}
	return val * sign, nil
}

// ParseNumber parses a number (integer or float).
// Corresponds to C++ BasicParser::parseNumber() overloads.
func (p *BasicParser) ParseNumber() (interface{}, error) {
	p.SkipWhitespace(false)
	start := p.pos
	sign := 1
	if p.HasNext() && p.Peek() == '-' {
		sign = -1
		p.NextChar()
	}

	// Check for digits
	hasDigits := false
	for p.HasNext() && p.Peek() >= '0' && p.Peek() <= '9' {
		p.NextChar()
		hasDigits = true
	}

	isFloat := false
	if p.HasNext() && p.Peek() == '.' {
		isFloat = true
		p.NextChar()
		// Digits after decimal
		for p.HasNext() && p.Peek() >= '0' && p.Peek() <= '9' {
			p.NextChar()
			hasDigits = true
		}
	}

	if !hasDigits {
		p.pos = start // Rewind
		return nil, p.error("expected number")
	}

	numStr := string(p.source[start:p.pos])
	if isFloat {
		f, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, p.error(fmt.Sprintf("invalid float: %%v", err))
		}
		return f, nil
	} else {
		i, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return nil, p.error(fmt.Sprintf("invalid integer: %%v", err))
		}
		return i, nil
	}
}

// ParseString parses a string enclosed in 'delimiter' (e.g., quotes).
// Corresponds to C++ BasicParser::parseString(CharT chr, bool closeRequired)
func (p *BasicParser) ParseString(delimiter rune, closeRequired bool) (string, error) {
	if p.Peek() != delimiter {
		if closeRequired {
			return "", p.error(fmt.Sprintf("'%%c' expected", delimiter))
		}
		return "", nil // Not an error if not required and not found
	}
	p.NextChar() // Skip opening delimiter

	start := p.pos
	for p.HasNext() && p.Peek() != delimiter {
		if p.Peek() == '\\'' { // Handle escaped characters
			p.NextChar()
			if p.HasNext() {
				p.NextChar()
			}
		} else {
				p.NextChar()
		}
	}
	ss := string(p.source[start:p.pos])

	if p.Peek() == delimiter {
		p.NextChar() // Skip closing delimiter
	} else if closeRequired {
		return "", p.error(fmt.Sprintf("'%%c' expected to close string", delimiter))
	}
	return ss, nil
}

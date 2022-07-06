package sln

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type SolutionParser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewSolutionParser returns a new instance of Parser.
func NewSolutionParser(r io.Reader) *SolutionParser {
	return &SolutionParser{s: NewScanner(r)}
}

func (sp *SolutionParser) ParseString() (string, error) {
	tok, lit := sp.scanIgnoreWhitespace()
	if tok != QUOTE {
		return lit, nil
	} else {
		var s string
		for {
			tok, lit := sp.scan()
			if tok != QUOTE {
				s = s + lit
			} else {
				break
			}
		}
		return s, nil
	}
}

func (sp *SolutionParser) ParseProject() (Project, error) {
	var proj Project
	if ok, err := sp.expect(OPEN_PAREN); !ok {
		return proj, err
	}
	proj.TypeGUID, _ = sp.ParseString()
	if ok, err := sp.expect(CLOSE_PAREN, EQUAL); !ok {
		return proj, err
	}
	proj.Name, _ = sp.ParseString()
	if ok, err := sp.expect(COMMA); !ok {
		return proj, err
	}
	s, _ := sp.ParseString()
	proj.ProjectFile = strings.Replace(s, `\`, string(filepath.Separator), -1)
	if ok, err := sp.expect(COMMA); !ok {
		return proj, err
	}
	proj.ID, _ = sp.ParseString()
	if ok, err := sp.expect(ENDPROJECT); !ok {
		return proj, err
	}
	return proj, nil
}

func (sp *SolutionParser) expect(expected ...Token) (bool, error) {
	for _, exp := range expected {
		if tok, lit := sp.scanIgnoreWhitespace(); tok != exp {
			return false, fmt.Errorf("unexpected token %q", lit)
		}
	}
	return true, nil
}

// Parse parses a SQL SELECT statement.
func (sp *SolutionParser) Parse() (Solution, error) {
	var sln Solution
	for {
		tok, _ := sp.scanIgnoreWhitespace()
		switch tok {
		case EOF:
			sp.unscan()
		case PROJECT:
			proj, _ := sp.ParseProject()
			sln.Projects = append(sln.Projects, proj)
		}
		if tok == EOF {
			break
		}
	}
	return sln, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (sp *SolutionParser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if sp.buf.n != 0 {
		sp.buf.n = 0
		return sp.buf.tok, sp.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = sp.s.Scan()

	// Save it to the buffer in case we unscan later.
	sp.buf.tok, sp.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (sp *SolutionParser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = sp.scan()
	if tok == WS {
		tok, lit = sp.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (sp *SolutionParser) unscan() { sp.buf.n = 1 }

package parser

import (
	"bufio"
	"io"
	"strings"
)

// Comment types.
const (
	LCommand = "L"
	CCommand = "C"
	ACommand = "A"
)

// Parser is Hack assembly parser.
type Parser struct {
	currentCommand string
	s              *bufio.Scanner
	hasMoreCommand bool
}

// New creates a new Hack assembly parser.
func New(r io.Reader) *Parser {
	p := &Parser{
		s:              bufio.NewScanner(r),
		hasMoreCommand: true,
	}
	return p
}

// Advance reads until next command or eof.
func (p *Parser) Advance(ln *int) error {
	for {
		if !p.s.Scan() {
			p.currentCommand = ""
			p.hasMoreCommand = false
			return nil
		}
		l := p.s.Text()
		err := p.s.Err()
		if err != nil {
			return nil
		}
		*ln++
		l = strings.TrimSpace(l)
		commentPos := strings.Index(l, "//")
		if commentPos != -1 {
			l = l[:commentPos]
		}
		if len(l) == 0 {
			continue
		}
		p.currentCommand = strings.TrimSpace(l)
		break
	}
	return nil
}

// HasMoreCommand ...
func (p *Parser) HasMoreCommand() bool {
	return p.hasMoreCommand
}

// CommandType returns the current comment's type.
func (p *Parser) CommandType() string {
	if len(p.currentCommand) == 0 {
		return ""
	}
	if strings.HasPrefix(p.currentCommand, "@") {
		return ACommand
	}
	if strings.HasPrefix(p.currentCommand, "(") {
		return LCommand
	}
	return CCommand
}

func (p *Parser) Symbol() string {
	if p.CommandType() == ACommand {
		return strings.TrimPrefix(p.currentCommand, "@")
	}
	if p.CommandType() == LCommand {
		si := strings.Index(p.currentCommand, "(")
		ei := strings.Index(p.currentCommand, ")")
		return p.currentCommand[si+1 : ei]
	}
	return ""
}

func (p *Parser) Dest() string {
	if p.CommandType() != CCommand {
		return ""
	}
	i := strings.Index(p.currentCommand, "=")
	if i != -1 {
		return p.currentCommand[:i]
	}
	return ""
}

func (p *Parser) Comp() string {
	equalPos := strings.Index(p.currentCommand, "=")
	semicolonPos := strings.Index(p.currentCommand, ";")
	// dest = comp;jump
	if equalPos != -1 && semicolonPos != -1 {
		return p.currentCommand[equalPos+1 : semicolonPos]
	}
	// dest  = comp
	if equalPos != -1 && semicolonPos == -1 {
		return p.currentCommand[equalPos+1:]
	}
	// comp;jump
	if equalPos == -1 && semicolonPos != -1 {
		return p.currentCommand[0:semicolonPos]
	}
	return ""
}

func (p *Parser) Jump() string {
	if p.CommandType() != CCommand {
		return ""
	}
	semicolonPos := strings.Index(p.currentCommand, ";")
	if semicolonPos != -1 {
		return p.currentCommand[semicolonPos+1:]
	}
	return ""
}

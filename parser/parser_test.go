package parser

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	input := "//hello\n D=M;JMP//hello\n@111 // jfajfjfj j"
	p := New(strings.NewReader(input))
	line := 0
	p.Advance(&line)
	t.Logf("p %+v line %d", p, line)
	p.Advance(&line)
	t.Logf("p %+v line %d", p, line)
	p.Advance(&line)
	t.Logf("p %+v line %d", p, line)
}

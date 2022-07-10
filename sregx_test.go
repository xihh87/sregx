package sregx_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/zyedidia/sregx"
)

type Test struct {
	name  string
	input string
	want  string
}

func check(cmd sregx.Command, tests []Test, t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := cmd.Evaluate([]byte(tt.input))
			if !bytes.Equal([]byte(tt.want), out) {
				t.Errorf("got %q, want %q", out, tt.want)
			}
		})
	}
}

func TestS(t *testing.T) {
	cmd := sregx.S{
		Patt:    regexp.MustCompile("([A-Za-z]+) ([A-Za-z]+)"),
		Replace: []byte("$2 $1"),
	}

	tests := []Test{
		{"s1", "hello world", "world hello"},
	}

	check(cmd, tests, t)
}

func TestD(t *testing.T) {
	cmd := sregx.X{
		Patt: regexp.MustCompile("string"),
		Cmd:  sregx.D{},
	}

	tests := []Test{
		{"d1", "string", ""},
		{"d2", "hello string hi string test", "hello  hi  test"},
	}

	check(cmd, tests, t)
}

func TestCVar(t *testing.T) {
	// Renames c variables called 'n' to 'num'. Omits matches in strings.
	// expression: y/".*"/y/'.*'/x/[a-zA-Z0-9]+/g/n/v/../c/num/
	cmd := sregx.Y{
		Patt: regexp.MustCompile(`".*"`),
		Cmd: sregx.Y{
			Patt: regexp.MustCompile(`'.*'`),
			Cmd: sregx.X{
				Patt: regexp.MustCompile(`[a-zA-z0-9]+`),
				Cmd: sregx.G{
					Patt: regexp.MustCompile(`n`),
					Cmd: sregx.V{
						Patt: regexp.MustCompile(`..`),
						Cmd: sregx.C{
							Change: []byte("num"),
						},
					},
				},
			},
		},
	}

	cin := `#include <stdio.h>
int main() {
	char* n = "hello n \n";
	printf("%s\n", n);
}
	`
	cout := `#include <stdio.h>
int main() {
	char* num = "hello n \n";
	printf("%s\n", num);
}
	`

	tests := []Test{
		{"cvar1", "n", "num"},
		{"cvar2", cin, cout},
	}

	check(cmd, tests, t)
}

func TestICapitalize(t *testing.T) {
	// Program to capitalize 'i's
	// x/[A-Za-z]+/ g/i/ v/../ c/I/
	cmd := sregx.X{
		Patt: regexp.MustCompile("[A-Za-z]+"),
		Cmd: sregx.G{
			Patt: regexp.MustCompile("i"),
			Cmd: sregx.V{
				Patt: regexp.MustCompile(".."),
				Cmd: sregx.C{
					Change: []byte("I"),
				},
			},
		},
	}

	tests := []Test{
		{"i1", "i am making tests", "I am making tests"},
		{"i2", "ii i i i iii", "ii I I I iii"},
	}

	check(cmd, tests, t)
}

func TestICapitalizeAlternate(t *testing.T) {
	// Alternate program to capitalize 'i's
	// x/[A-Za-z]+/ g/^i$/ c/I/
	cmd := sregx.X{
		Patt: regexp.MustCompile("[A-Za-z]+"),
		Cmd: sregx.G{
			Patt: regexp.MustCompile("^i$"),
			Cmd: sregx.C{
				Change: []byte("I"),
			},
		},
	}

	tests := []Test{
		{"ialt1", "i am making tests", "I am making tests"},
		{"ialt2", "ii i i i iii", "ii I I I iii"},
	}

	check(cmd, tests, t)
}

func TestUnicodeInRegx(t *testing.T) {
    cmd := sregx.X {
        Patt: regexp.MustCompile("- +\\[ \\] "),
        Cmd: sregx.C{
            Change: []byte("- ☐ "),
        },
    }

    tests := []Test{
        {"pending1", "- [ ] ", "- ☐ "},
        {"pending2", "-     [ ] ", "- ☐ "},
    }

    check(cmd, tests, t)
}

package errors

import (
	goerrors "errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		New("error"),
		"%s",
		"error",
	}, {
		New("error"),
		"%v",
		"error",
	}, {
		New("error"),
		"%+v",
		"error\n" +
			"github.com/objenious/errors.TestFormatNew\n" +
			"\t.+/github.com/objenious/errors/format_test.go:26",
	}, {
		New("error"),
		"%q",
		`"error"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatErrorf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		Errorf("%s", "error"),
		"%v",
		"error",
	}, {
		Errorf("%s", "error"),
		"%+v",
		"error\n" +
			"github.com/objenious/errors.TestFormatErrorf\n" +
			"\t.+/github.com/objenious/errors/format_test.go:56",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrap(New("error"), "error2"),
		"%s",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%v",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%+v",
		"error\n" +
			"github.com/objenious/errors.TestFormatWrap\n" +
			"\t.+/github.com/objenious/errors/format_test.go:82",
	}, {
		Wrap(io.EOF, "error"),
		"%s",
		"error: EOF",
	}, {
		Wrap(io.EOF, "error"),
		"%v",
		"error: EOF",
	}, {
		Wrap(io.EOF, "error"),
		"%+v",
		"EOF\n" +
			"error\n" +
			"github.com/objenious/errors.TestFormatWrap\n" +
			"\t.+/github.com/objenious/errors/format_test.go:96",
	}, {
		Wrap(Wrap(io.EOF, "error1"), "error2"),
		"%+v",
		"EOF\n" +
			"error1\n" +
			"github.com/objenious/errors.TestFormatWrap\n" +
			"\t.+/github.com/objenious/errors/format_test.go:103\n",
	}, {
		Wrap(New("error with space"), "context"),
		"%q",
		`"context: error with space"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(io.EOF, "error%d", 2),
		"%s",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%v",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"EOF\n" +
			"error2\n" +
			"github.com/objenious/errors.TestFormatWrapf\n" +
			"\t.+/github.com/objenious/errors/format_test.go:134",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"github.com/objenious/errors.TestFormatWrapf\n" +
			"\t.+/github.com/objenious/errors/format_test.go:149",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWithStack(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		WithStack(io.EOF),
		"%s",
		[]string{"EOF"},
	}, {
		WithStack(io.EOF),
		"%v",
		[]string{"EOF"},
	}, {
		WithStack(io.EOF),
		"%+v",
		[]string{"EOF",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:175"},
	}, {
		WithStack(New("error")),
		"%s",
		[]string{"error"},
	}, {
		WithStack(New("error")),
		"%v",
		[]string{"error"},
	}, {
		WithStack(New("error")),
		"%+v",
		[]string{"error",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:189",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:189"},
	}, {
		WithStack(WithStack(io.EOF)),
		"%+v",
		[]string{"EOF",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:197",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:197"},
	}, {
		WithStack(WithStack(Wrapf(io.EOF, "message"))),
		"%+v",
		[]string{"EOF",
			"message",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:205",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:205",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:205"},
	}, {
		WithStack(Errorf("error%d", 1)),
		"%+v",
		[]string{"error1",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:216",
			"github.com/objenious/errors.TestFormatWithStack\n" +
				"\t.+/github.com/objenious/errors/format_test.go:216"},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func wrappedNew(message string) error { // This function will be mid-stack inlined in go 1.12+
	return New(message)
}

func TestFormatWrappedNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		wrappedNew("error"),
		"%+v",
		"error\n" +
			"github.com/objenious/errors.wrappedNew\n" +
			"\t.+/github.com/objenious/errors/format_test.go:231\n" +
			"github.com/objenious/errors.TestFormatWrappedNew\n" +
			"\t.+/github.com/objenious/errors/format_test.go:240",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	t.Helper()
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}

var stackLineR = regexp.MustCompile(`\.`)

// parseBlocks parses input into a slice, where:
//  - incase entry contains a newline, its a stacktrace
//  - incase entry contains no newline, its a solo line.
//
// Detecting stack boundaries only works incase the WithStack-calls are
// to be found on the same line, thats why it is optionally here.
//
// Example use:
//
// for _, e := range blocks {
//   if strings.ContainsAny(e, "\n") {
//     // Match as stack
//   } else {
//     // Match as line
//   }
// }
//
func parseBlocks(input string, detectStackboundaries bool) ([]string, error) {
	var blocks []string

	stack := ""
	wasStack := false
	lines := map[string]bool{} // already found lines

	for _, l := range strings.Split(input, "\n") {
		isStackLine := stackLineR.MatchString(l)

		switch {
		case !isStackLine && wasStack:
			blocks = append(blocks, stack, l)
			stack = ""
			lines = map[string]bool{}
		case isStackLine:
			if wasStack {
				// Detecting two stacks after another, possible cause lines match in
				// our tests due to WithStack(WithStack(io.EOF)) on same line.
				if detectStackboundaries {
					if lines[l] {
						if len(stack) == 0 {
							return nil, goerrors.New("len of block must not be zero here")
						}

						blocks = append(blocks, stack)
						stack = l
						lines = map[string]bool{l: true}
						continue
					}
				}

				stack = stack + "\n" + l
			} else {
				stack = l
			}
			lines[l] = true
		case !isStackLine && !wasStack:
			blocks = append(blocks, l)
		default:
			return nil, goerrors.New("must not happen")
		}

		wasStack = isStackLine
	}

	// Use up stack
	if stack != "" {
		blocks = append(blocks, stack)
	}
	return blocks, nil
}

func testFormatCompleteCompare(t *testing.T, n int, arg interface{}, format string, want []string, detectStackBoundaries bool) {
	gotStr := fmt.Sprintf(format, arg)

	got, err := parseBlocks(gotStr, detectStackBoundaries)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf("test %d: fmt.Sprintf(%s, err) -> wrong number of blocks: got(%d) want(%d)\n got: %s\nwant: %s\ngotStr: %q",
			n+1, format, len(got), len(want), prettyBlocks(got), prettyBlocks(want), gotStr)
	}

	for i := range got {
		if strings.ContainsAny(want[i], "\n") {
			// Match as stack
			match, err := regexp.MatchString(want[i], got[i])
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("test %d: block %d: fmt.Sprintf(%q, err):\ngot:\n%q\nwant:\n%q\nall-got:\n%s\nall-want:\n%s\n",
					n+1, i+1, format, got[i], want[i], prettyBlocks(got), prettyBlocks(want))
			}
		} else {
			// Match as message
			if got[i] != want[i] {
				t.Fatalf("test %d: fmt.Sprintf(%s, err) at block %d got != want:\n got: %q\nwant: %q", n+1, format, i+1, got[i], want[i])
			}
		}
	}
}

func prettyBlocks(blocks []string) string {
	var out []string

	for _, b := range blocks {
		out = append(out, fmt.Sprintf("%v", b))
	}

	return "   " + strings.Join(out, "\n   ")
}

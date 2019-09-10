package errors

import (
	goerrors "errors"
	"runtime"
	"testing"
)

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() Frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		initpc,
		"%s",
		"stack_test.go",
	}, {
		initpc,
		"%+s",
		"github.com/objenious/errors.init\n" +
			"\t.+/github.com/objenious/errors/stack_test.go",
	}, {
		0,
		"%s",
		"unknown",
	}, {
		0,
		"%+s",
		"unknown",
	}, {
		initpc,
		"%d",
		"9",
	}, {
		0,
		"%d",
		"0",
	}, {
		initpc,
		"%n",
		"init",
	}, {
		func() Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		`\(\*X\).ptr`,
	}, {
		func() Frame {
			var x X
			return x.val()
		}(),
		"%n",
		"X.val",
	}, {
		0,
		"%n",
		"",
	}, {
		initpc,
		"%v",
		"stack_test.go:9",
	}, {
		initpc,
		"%+v",
		"github.com/objenious/errors.init\n" +
			"\t.+/github.com/objenious/errors/stack_test.go:9",
	}, {
		0,
		"%v",
		"unknown:0",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.Frame, tt.format, tt.want)
	}
}

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/objenious/errors.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcname(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcname(%q): want: %q, got %q", tt.name, want, got)
		}
	}
}

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want []string
	}{{
		New("ooh"), []string{
			"github.com/objenious/errors.TestStackTrace\n" +
				"\t.+/github.com/objenious/errors/stack_test.go:121",
		},
	}, {
		Wrap(New("ooh"), "ahh"), []string{
			"github.com/objenious/errors.TestStackTrace\n" +
				"\t.+/github.com/objenious/errors/stack_test.go:126", // this is the stack of Wrap, not New
		},
	}, {
		goerrors.Unwrap(Wrap(New("ooh"), "ahh")), []string{
			"github.com/objenious/errors.TestStackTrace\n" +
				"\t.+/github.com/objenious/errors/stack_test.go:131", // this is the stack of New
		},
	}, {
		func() error { return New("ooh") }(), []string{
			`github.com/objenious/errors.TestStackTrace.func1` +
				"\n\t.+/github.com/objenious/errors/stack_test.go:136", // this is the stack of New
			"github.com/objenious/errors.TestStackTrace\n" +
				"\t.+/github.com/objenious/errors/stack_test.go:136", // this is the stack of New's caller
		},
	/*}, {
		goerrors.Unwrap(func() error {
			return func() error {
				return Errorf("hello %s", "world")
			}()
		}()), []string{
			`github.com/pkg/errors.TestStackTrace.func2.1` +
				"\n\t.+/github.com/pkg/errors/stack_test.go:145", // this is the stack of Errorf
			`github.com/pkg/errors.TestStackTrace.func2` +
				"\n\t.+/github.com/pkg/errors/stack_test.go:146", // this is the stack of Errorf's caller
			"github.com/pkg/errors.TestStackTrace\n" +
				"\t.+/github.com/pkg/errors/stack_test.go:147", // this is the stack of Errorf's caller's caller
		},*/
	}}
	for i, tt := range tests {
		x := GetStackTrace(tt.err)
		if x == nil {
			t.Fatalf("no stacktrace for test %d %+v", i+1, tt.err)
		}
		st := x.StackTrace()
		for j, want := range tt.want {
			testFormatRegexp(t, i, st[j], "%+v", want)
		}
	}
}

func stackTrace() StackTrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st stack = pcs[0:n]
	return st.StackTrace()
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		StackTrace
		format string
		want   string
	}{{
		nil,
		"%s",
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		nil,
		"%#v",
		`\[\]errors.Frame\(nil\)`,
	}, {
		make(StackTrace, 0),
		"%s",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%v",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%+v",
		"",
	}, {
		make(StackTrace, 0),
		"%#v",
		`\[\]errors.Frame{}`,
	}, {
		stackTrace()[:2],
		"%s",
		`\[stack_test.go stack_test.go\]`,
	}, {
		stackTrace()[:2],
		"%v",
		`\[stack_test.go:171 stack_test.go:218\]`,
	}, {
		stackTrace()[:2],
		"%+v",
		"\n" +
			"github.com/objenious/errors.stackTrace\n" +
			"\t.+/github.com/objenious/errors/stack_test.go:171\n" +
			"github.com/objenious/errors.TestStackTraceFormat\n" +
			"\t.+/github.com/objenious/errors/stack_test.go:222",
	}, {
		stackTrace()[:2],
		"%#v",
		`\[\]errors.Frame{stack_test.go:171, stack_test.go:230}`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.want)
	}
}

// a version of runtime.Caller that returns a Frame, not a uintptr.
func caller() Frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return Frame(frame.PC)
}

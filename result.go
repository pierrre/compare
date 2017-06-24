package compare

import (
	"bytes"
	"fmt"
)

// Result is a list of Difference.
type Result []Difference

// Merge merges 2 Result.
func (r Result) Merge(rm Result) Result {
	if len(rm) == 0 {
		return r
	}
	if len(r) == 0 {
		return rm
	}
	return append(r, rm...)
}

// Format implements fmt.Formatter.
//
// See Difference.Format() for supported verb and flag.
func (r Result) Format(s fmt.State, verb rune) {
	if verb != 'v' {
		_, _ = fmt.Fprintf(s, "%%!%c(%T)", verb, r)
		return
	}
	if len(r) == 0 {
		_, _ = s.Write(resultNoneBytes)
		return
	}
	for i, d := range r {
		if i > 0 {
			_, _ = s.Write(resultNewLineBytes)
		}
		d.Format(s, verb)
	}
}

var (
	resultNoneBytes    = []byte("<none>")
	resultNewLineBytes = []byte("\n")
)

// Difference represents a difference between 2 values.
type Difference struct {
	Path    Path
	Message string
	V1, V2  interface{}
}

// Format implements fmt.Formatter.
//
// It only supports the 'v' verb.
// By default, it show the path and message.
// The '+' flag shows values V1 and V2.
func (d Difference) Format(s fmt.State, verb rune) {
	// We use a buffer in order to reduce memory allocation.
	// fmt.State (and its real type) doesn't (yet?) implement WriteString().
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	if verb == 'v' {
		_, _ = buf.WriteString(PathString(d.Path) + ": " + d.Message)
		if s.Flag('+') {
			_, _ = fmt.Fprintf(buf, "\n\tv1="+d.getValueFormat(d.V1)+"\n\tv2="+d.getValueFormat(d.V2), d.V1, d.V2)
		}
	} else {
		_, _ = fmt.Fprintf(buf, "%%!%c(%T)", verb, d)
	}
	_, _ = buf.WriteTo(s)
	bufPool.Put(buf)
}

func (d Difference) getValueFormat(v interface{}) string {
	switch v.(type) {
	case string:
		return "%q"
	default:
		return "%v"
	}
}

const (
	msgOnlyOneIsValid        = "only one is valid"
	msgOnlyOneIsNil          = "only one is nil"
	msgTypeNotEqual          = "type not equal"
	msgCapacityNotEqual      = "capacity not equal"
	msgLengthNotEqual        = "length not equal"
	msgBoolNotEqual          = "bool not equal"
	msgIntNotEqual           = "int not equal"
	msgUintNotEqual          = "uint not equal"
	msgFloatNotEqual         = "float not equal"
	msgComplexNotEqual       = "complex not equal"
	msgStringNotEqual        = "string not equal"
	msgMapKeyNotDefined      = "map key not defined"
	msgUnsafePointerNotEqual = "unsafe pointer not equal"
	msgFuncPointerNotEqual   = "func pointer not equal"
	msgMethodNotEqual        = "method .%s() returned false"
	msgMethodCmpNotEqual     = "method .Cmp() returned %d"
)

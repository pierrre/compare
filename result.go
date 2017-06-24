package compare

import (
	"fmt"
	"io"
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
		_, _ = io.WriteString(s, "<none>")
		return
	}
	for i, d := range r {
		if i > 0 {
			_, _ = io.WriteString(s, "\n")
		}
		d.Format(s, verb)
	}
}

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
	if verb != 'v' {
		_, _ = fmt.Fprintf(s, "%%!%c(%T)", verb, d)
		return
	}
	_, _ = fmt.Fprintf(s, "%s: %s", PathString(d.Path), d.Message)
	if s.Flag('+') {
		_, _ = fmt.Fprintf(s, "\n\tv1="+d.getValueFormat(d.V1)+"\n\tv2="+d.getValueFormat(d.V2), d.V1, d.V2)
	}
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

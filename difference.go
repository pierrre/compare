package compare

import (
	"fmt"

	"github.com/pierrre/go-libs/unsafeio"
)

// Difference represents a difference between 2 values.
type Difference struct {
	Path    Path   `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
	V1      string `json:"v1,omitempty"`
	V2      string `json:"v2,omitempty"`
}

// Format implements [fmt.Formatter].
//
// It only supports the 'v' verb.
// By default, it show the path and message.
// The '+' flag shows values V1 and V2.
func (d Difference) Format(s fmt.State, verb rune) {
	if verb != 'v' {
		_, _ = fmt.Fprintf(s, "%%!%c(%T)", verb, d)
		return
	}
	d.Path.Format(s, verb)
	_, _ = unsafeio.WriteString(s, ": ")
	_, _ = unsafeio.WriteString(s, d.Message)
	if s.Flag('+') {
		if d.V1 != "" || d.V2 != "" {
			_, _ = unsafeio.WriteString(s, "\n\tv1=")
			_, _ = unsafeio.WriteString(s, d.V1)
			_, _ = unsafeio.WriteString(s, "\n\tv2=")
			_, _ = unsafeio.WriteString(s, d.V2)
		}
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
	msgMethodEqualFalse      = "method .Equal() returned false"
	msgMethodCmpNotEqual     = "method .Cmp() returned %d"
)

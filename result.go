package compare

import (
	"bytes"
	"fmt"
	"strconv"
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

func (r Result) String() string {
	switch len(r) {
	case 0:
		return "<none>"
	case 1:
		return r[0].String()
	}
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	for i, d := range r {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(d.String())
	}
	s := buf.String()
	bufPool.Put(buf)
	return s
}

// Difference represents a difference between 2 values.
type Difference struct {
	Path    Path
	Message string
	V1, V2  interface{}
}

func (d Difference) String() string {
	return fmt.Sprintf("%s: %s: v1=%v v2=%v", PathString(d.Path), d.Message, d.formatValue(d.V1), d.formatValue(d.V2))
}

func (d Difference) formatValue(v interface{}) string {
	switch v := v.(type) {
	case string:
		return strconv.Quote(v)
	}
	return fmt.Sprint(v)
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

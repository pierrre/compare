package compare

import (
	"reflect"
)

type ValueComparator func(st *State, v1, v2 reflect.Value) bool

// Func represents a comparison function.
// TODO remove
// It is guaranteed that both values are valid and of the same type.
// If the returned value "stop" is true, the comparison will stop.
type Func func(c *Comparator, st *State, v1, v2 reflect.Value) (r Result, stop bool)

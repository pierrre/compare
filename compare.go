// Package compare provide comparison utilities.
package compare

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"sync"

	"github.com/pierrre/compare/internal"
	"github.com/pierrre/go-libs/strconvio"
)

// Compare compares 2 values with [DefaultComparator].
func Compare(v1, v2 any) Result {
	return DefaultComparator.Compare(v1, v2)
}

// DefaultComparator is the default [Comparator].
var DefaultComparator = NewComparator()

// Comparator compares 2 values.
//
// It should be created with [NewComparator].
type Comparator struct {
	// MaxDepth is the maximum depth of the comparison.
	// If the value is reached, the comparison is stopped.
	// Default: 0 (no limit).
	MaxDepth int
	// SliceMaxDifferences is the maximum number of different items for a slice.
	// If the value is reached, the comparison is stopped for the current slice.
	// It is also used for array.
	// Set to 0 disables it.
	// Default: 10.
	SliceMaxDifferences int
	// MapMaxDifferences is the maximum number of different items for a map.
	// If the value is reached, the comparison is stopped for the current map.
	// Set to 0 disables it.
	// Default: 10.
	MapMaxDifferences int
	// Funcs is the list of custom comparison functions.
	// Default: []byte, reflect.Value, .Equal().
	Funcs []Func
}

// NewComparator returns a new [Comparator] initialized with default values.
func NewComparator() *Comparator {
	return &Comparator{
		SliceMaxDifferences: 10,
		MapMaxDifferences:   10,
		Funcs: []Func{
			NewBytesEqualFunc(),
			NewReflectValueFunc(),
			NewMethodEqualFunc(),
			NewMethodCmpFunc(),
		},
	}
}

// Compare compares 2 values.
func (c *Comparator) Compare(v1, v2 any) Result {
	stItf := statePool.Get()
	defer statePool.Put(stItf)
	st := stItf.(*State) //nolint:forcetypeassert // The pool only contains *State.
	st.reset()
	return c.compare(
		st,
		reflect.ValueOf(v1),
		reflect.ValueOf(v2),
	)
}

func (c *Comparator) compare(st *State, v1, v2 reflect.Value) Result {
	if c.MaxDepth > 0 && st.Depth >= c.MaxDepth {
		return nil
	}
	st.Depth++
	defer func() {
		st.Depth--
	}()
	if r, stop := c.compareValid(v1, v2); stop {
		return r
	}
	if r, stop := c.compareType(v1, v2); stop {
		return r
	}
	if r, stop := c.compareFuncs(st, v1, v2); stop {
		return r
	}
	return c.compareKind(st, v1, v2)
}

func (c *Comparator) checkRecursion(st *State, v1, v2 reflect.Value) bool {
	vp := Visited{
		V1: v1.Pointer(),
		V2: v2.Pointer(),
	}
	if slices.Contains(st.Visited, vp) {
		return true
	}
	st.Visited = append(st.Visited, vp)
	return false
}

func (c *Comparator) endRecursion(st *State) {
	st.Visited = st.Visited[:len(st.Visited)-1]
}

func (c *Comparator) compareValid(v1, v2 reflect.Value) (Result, bool) {
	vl1 := v1.IsValid()
	vl2 := v2.IsValid()
	if vl1 && vl2 {
		return nil, false
	}
	if vl1 == vl2 {
		return nil, true
	}
	return Result{Difference{
		Message: msgOnlyOneIsValid,
		V1:      strconv.FormatBool(vl1),
		V2:      strconv.FormatBool(vl2),
	}}, true
}

func (c *Comparator) compareType(v1, v2 reflect.Value) (Result, bool) {
	t1 := v1.Type()
	t2 := v2.Type()
	if t1 == t2 {
		return nil, false
	}
	return Result{Difference{
		Message: msgTypeNotEqual,
		V1:      t1.String(),
		V2:      t2.String(),
	}}, true
}

//nolint:gocyclo // Large switch/case is OK.
func (c *Comparator) compareKind(st *State, v1, v2 reflect.Value) Result {
	switch v1.Kind() { //nolint:exhaustive // All kinds are handled, Invalid should not happen.
	case reflect.Bool:
		return c.compareBool(v1, v2)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return c.compareInt(v1, v2)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return c.compareUint(v1, v2)
	case reflect.Float32, reflect.Float64:
		return c.compareFloat(v1, v2)
	case reflect.Complex64, reflect.Complex128:
		return c.compareComplex(v1, v2)
	case reflect.String:
		return c.compareString(v1, v2)
	case reflect.Array:
		return c.compareArray(st, v1, v2)
	case reflect.Slice:
		return c.compareSlice(st, v1, v2)
	case reflect.Interface:
		return c.compareInterface(st, v1, v2)
	case reflect.Pointer:
		return c.comparePointer(st, v1, v2)
	case reflect.Struct:
		return c.compareStruct(st, v1, v2)
	case reflect.Map:
		return c.compareMap(st, v1, v2)
	case reflect.UnsafePointer:
		return c.compareUnsafePointer(v1, v2)
	case reflect.Chan:
		return c.compareChan(v1, v2)
	case reflect.Func:
		return c.compareFunc(v1, v2)
	}
	return nil
}

func (c *Comparator) compareBool(v1, v2 reflect.Value) Result {
	b1 := v1.Bool()
	b2 := v2.Bool()
	if b1 == b2 {
		return nil
	}
	return Result{Difference{
		Message: msgBoolNotEqual,
		V1:      strconv.FormatBool(b1),
		V2:      strconv.FormatBool(b2),
	}}
}

func (c *Comparator) compareInt(v1, v2 reflect.Value) Result {
	i1 := v1.Int()
	i2 := v2.Int()
	if i1 == i2 {
		return nil
	}
	return Result{Difference{
		Message: msgIntNotEqual,
		V1:      strconv.FormatInt(i1, 10),
		V2:      strconv.FormatInt(i2, 10),
	}}
}

func (c *Comparator) compareUint(v1, v2 reflect.Value) Result {
	u1 := v1.Uint()
	u2 := v2.Uint()
	if u1 == u2 {
		return nil
	}
	return Result{Difference{
		Message: msgUintNotEqual,
		V1:      strconv.FormatUint(u1, 10),
		V2:      strconv.FormatUint(u2, 10),
	}}
}

func (c *Comparator) compareFloat(v1, v2 reflect.Value) Result {
	f1 := v1.Float()
	f2 := v2.Float()
	if f1 == f2 {
		return nil
	}
	bitSize := v1.Type().Bits()
	return Result{Difference{
		Message: msgFloatNotEqual,
		V1:      strconv.FormatFloat(f1, 'g', -1, bitSize),
		V2:      strconv.FormatFloat(f2, 'g', -1, bitSize),
	}}
}

func (c *Comparator) compareComplex(v1, v2 reflect.Value) Result {
	c1 := v1.Complex()
	c2 := v2.Complex()
	if c1 == c2 {
		return nil
	}
	bitSize := v1.Type().Bits()
	return Result{Difference{
		Message: msgComplexNotEqual,
		V1:      strconv.FormatComplex(c1, 'g', -1, bitSize),
		V2:      strconv.FormatComplex(c2, 'g', -1, bitSize),
	}}
}

func (c *Comparator) compareString(v1, v2 reflect.Value) Result {
	s1 := v1.String()
	s2 := v2.String()
	if s1 == s2 {
		return nil
	}
	return Result{Difference{
		Message: msgStringNotEqual,
		V1:      strconv.Quote(s1),
		V2:      strconv.Quote(s2),
	}}
}

func (c *Comparator) compareArray(st *State, v1, v2 reflect.Value) Result {
	var r Result
	diffCount := 0
	for i, n := 0, v1.Len(); i < n; i++ {
		ri := c.compareArrayIndex(st, v1, v2, i)
		r = append(r, ri...)
		if len(ri) > 0 {
			diffCount++
			if diffCount >= c.SliceMaxDifferences && c.SliceMaxDifferences > 0 {
				break
			}
		}
	}
	return r
}

func (c *Comparator) compareArrayIndex(st *State, v1, v2 reflect.Value, i int) Result {
	r := c.compare(st, v1.Index(i), v2.Index(i))
	if len(r) == 0 {
		return nil
	}
	for j, d := range r {
		d.Path = append(d.Path, PathElem{
			Index: toPtr(i),
		})
		r[j] = d
	}
	return r
}

func (c *Comparator) compareSlice(st *State, v1, v2 reflect.Value) Result {
	if r, stop := c.compareNilLenPointer(v1, v2); stop {
		return r
	}
	if c.checkRecursion(st, v1, v2) {
		return nil
	}
	defer c.endRecursion(st)
	return c.compareArray(st, v1, v2)
}

func (c *Comparator) compareInterface(st *State, v1, v2 reflect.Value) Result {
	if r, stop := c.compareNil(v1, v2); stop {
		return r
	}
	return c.compare(st, v1.Elem(), v2.Elem())
}

func (c *Comparator) comparePointer(st *State, v1, v2 reflect.Value) Result {
	if v1.Pointer() == v2.Pointer() {
		return nil
	}
	if c.checkRecursion(st, v1, v2) {
		return nil
	}
	defer c.endRecursion(st)
	return c.compare(st, v1.Elem(), v2.Elem())
}

func (c *Comparator) compareStruct(st *State, v1, v2 reflect.Value) Result {
	var r Result
	t := v1.Type()
	for i, n := 0, t.NumField(); i < n; i++ {
		r = append(r, c.compareStructField(st, v1, v2, i)...)
	}
	return r
}

func (c *Comparator) compareStructField(st *State, v1, v2 reflect.Value, i int) Result {
	r := c.compare(st, v1.Field(i), v2.Field(i))
	if len(r) == 0 {
		return nil
	}
	f := v1.Type().Field(i).Name
	for j, d := range r {
		d.Path = append(d.Path, PathElem{
			Struct: toPtr(f),
		})
		r[j] = d
	}
	return r
}

func (c *Comparator) compareMap(st *State, v1, v2 reflect.Value) Result {
	if r, stop := c.compareNilLenPointer(v1, v2); stop {
		return r
	}
	if c.checkRecursion(st, v1, v2) {
		return nil
	}
	defer c.endRecursion(st)
	var r Result
	diffCount := 0
	for _, k := range getSortedMapsKeys(v1, v2) {
		ri := c.compareMapKey(st, v1, v2, k)
		r = append(r, ri...)
		if len(ri) > 0 {
			diffCount++
			if diffCount >= c.MapMaxDifferences && c.MapMaxDifferences > 0 {
				break
			}
		}
	}
	return r
}

func getSortedMapsKeys(v1, v2 reflect.Value) []reflect.Value {
	ks := getMapsKeys(v1, v2)
	internal.SortMapsKeys(v1.Type().Key(), ks)
	return ks
}

func getMapsKeys(v1, v2 reflect.Value) []reflect.Value {
	ks := v1.MapKeys()
	for _, k2 := range v2.MapKeys() {
		if !v1.MapIndex(k2).IsValid() {
			ks = append(ks, k2)
		}
	}
	return ks
}

func (c *Comparator) compareMapKey(st *State, v1, v2, k reflect.Value) Result {
	v1 = v1.MapIndex(k)
	v2 = v2.MapIndex(k)
	vl1 := v1.IsValid()
	vl2 := v2.IsValid()
	if !vl1 || !vl2 {
		return Result{Difference{
			Path: Path{{
				Map: toPtr(fmt.Sprint(k)),
			}},
			Message: msgMapKeyNotDefined,
			V1:      strconv.FormatBool(vl1),
			V2:      strconv.FormatBool(vl2),
		}}
	}
	r := c.compare(st, v1, v2)
	if len(r) == 0 {
		return nil
	}
	pe := PathElem{
		Map: toPtr(fmt.Sprint(k)),
	}
	for i, d := range r {
		d.Path = append(d.Path, pe)
		r[i] = d
	}
	return r
}

func (c *Comparator) compareUnsafePointer(v1, v2 reflect.Value) Result {
	p1 := uintptr(v1.UnsafePointer())
	p2 := uintptr(v2.UnsafePointer())
	if p1 == p2 {
		return nil
	}
	return Result{Difference{
		Message: msgUnsafePointerNotEqual,
		V1:      "0x" + strconv.FormatUint(uint64(p1), 16),
		V2:      "0x" + strconv.FormatUint(uint64(p2), 16),
	}}
}

func (c *Comparator) compareChan(v1, v2 reflect.Value) Result {
	if r, stop := c.compareNil(v1, v2); stop {
		return r
	}
	if v1.Pointer() == v2.Pointer() {
		return nil
	}
	cap1 := v1.Cap()
	cap2 := v2.Cap()
	if cap1 != cap2 {
		return Result{Difference{
			Message: msgCapacityNotEqual,
			V1:      strconv.Itoa(cap1),
			V2:      strconv.Itoa(cap2),
		}}
	}
	len1 := v1.Len()
	len2 := v2.Len()
	if len1 != len2 {
		return Result{Difference{
			Message: msgLengthNotEqual,
			V1:      strconv.Itoa(len1),
			V2:      strconv.Itoa(len2),
		}}
	}
	return nil
}

func (c *Comparator) compareFunc(v1, v2 reflect.Value) Result {
	if r, stop := c.compareNil(v1, v2); stop {
		return r
	}
	p1 := uintptr(v1.UnsafePointer())
	p2 := uintptr(v2.UnsafePointer())
	if p1 == p2 {
		return nil
	}
	return Result{Difference{
		Message: msgFuncPointerNotEqual,
		V1:      runtime.FuncForPC(p1).Name(),
		V2:      runtime.FuncForPC(p2).Name(),
	}}
}

func (c *Comparator) compareNil(v1, v2 reflect.Value) (Result, bool) {
	nil1 := v1.IsNil()
	nil2 := v2.IsNil()
	if nil1 && nil2 {
		return nil, true
	}
	if nil1 != nil2 {
		return Result{Difference{
			Message: msgOnlyOneIsNil,
			V1:      strconv.FormatBool(nil1),
			V2:      strconv.FormatBool(nil2),
		}}, true
	}
	return nil, false
}

func (c *Comparator) compareNilLenPointer(v1, v2 reflect.Value) (Result, bool) {
	if r, stop := c.compareNil(v1, v2); stop {
		return r, true
	}
	len1 := v1.Len()
	len2 := v2.Len()
	if len1 != len2 {
		return Result{Difference{
			Message: msgLengthNotEqual,
			V1:      strconv.Itoa(len1),
			V2:      strconv.Itoa(len2),
		}}, true
	}
	if len1 == 0 {
		return nil, true
	}
	if v1.Pointer() == v2.Pointer() {
		return nil, true
	}
	return nil, false
}

var statePool = &sync.Pool{
	New: func() any {
		return &State{}
	},
}

// State represents the state of a comparison.
//
// Functions must restore the original state when they return.
type State struct {
	Depth   int
	Visited []Visited
}

func (st *State) reset() {
	st.Depth = 0
	st.Visited = st.Visited[:0]
}

// Visited represents a visited pair of values.
type Visited struct {
	V1, V2 uintptr
}

// Func represents a comparison function.
// It is guaranteed that both values are valid and of the same type.
// If the returned value "stop" is true, the comparison will stop.
type Func func(c *Comparator, st *State, v1, v2 reflect.Value) (r Result, stop bool)

func (c *Comparator) compareFuncs(st *State, v1, v2 reflect.Value) (Result, bool) {
	for _, f := range c.Funcs {
		if r, stop := f(c, st, v1, v2); stop {
			return r, true
		}
	}
	return nil, false
}

var typeByteSlice = reflect.TypeOf([]byte(nil))

// NewBytesEqualFunc returns a [Func] that compares byte slices with bytes.Equal().
func NewBytesEqualFunc() Func {
	return compareBytesEqual
}

func compareBytesEqual(c *Comparator, st *State, v1, v2 reflect.Value) (Result, bool) {
	if v1.Type() != typeByteSlice {
		return nil, false
	}
	if bytes.Equal(v1.Bytes(), v2.Bytes()) {
		return nil, true
	}
	// If the []byte are not equal,
	// we want to continue the comparison,
	// so we will know which elements are not equal.
	return nil, false
}

var typeReflectValue = reflect.TypeOf(reflect.Value{})

// NewReflectValueFunc returns a [Func] that compares reflect.Value.
func NewReflectValueFunc() Func {
	return compareReflectValue
}

func compareReflectValue(c *Comparator, st *State, v1, v2 reflect.Value) (Result, bool) {
	if v1.Type() != typeReflectValue {
		return nil, false
	}
	if !v1.CanInterface() || !v2.CanInterface() {
		// Stop the comparison here.
		// We don't want to compare the structs.
		return nil, true
	}
	v1 = v1.Interface().(reflect.Value) //nolint:forcetypeassert // The type assertion is already checked above.
	v2 = v2.Interface().(reflect.Value) //nolint:forcetypeassert // The type assertion is already checked above.
	return c.compare(st, v1, v2), true
}

// NewMethodEqualFunc returns a [Func] that compares with the method .Equal().
func NewMethodEqualFunc() Func {
	return compareMethodEqual
}

func compareMethodEqual(c *Comparator, st *State, v1, v2 reflect.Value) (Result, bool) {
	f, ok := getMethodEqualFunc(v1.Type())
	if !ok {
		return nil, false
	}
	if !v1.CanInterface() || !v2.CanInterface() {
		return nil, false
	}
	eqRes := f.Call([]reflect.Value{v1, v2})[0].Interface().(bool) //nolint:forcetypeassert // The type of the returned value is already checked.
	if eqRes {
		return nil, true
	}
	return Result{Difference{
		Message: msgMethodEqualFalse,
	}}, true
}

var (
	equalMethodFuncsLock sync.Mutex
	equalMethodFuncs     = make(map[reflect.Type]*reflect.Value)
)

func getMethodEqualFunc(typ reflect.Type) (reflect.Value, bool) {
	equalMethodFuncsLock.Lock()
	defer equalMethodFuncsLock.Unlock()
	fp, ok := equalMethodFuncs[typ]
	if ok {
		if fp != nil {
			return *fp, true
		}
		return reflect.Value{}, false
	}
	equalMethodFuncs[typ] = nil
	met, ok := typ.MethodByName("Equal")
	if !ok {
		return reflect.Value{}, false
	}
	metTyp := met.Type
	if metTyp.NumIn() != 2 || metTyp.In(0) != typ || metTyp.In(1) != typ || metTyp.NumOut() != 1 || metTyp.Out(0) != reflect.TypeOf(true) {
		return reflect.Value{}, false
	}
	equalMethodFuncs[typ] = &met.Func
	return met.Func, true
}

// NewMethodCmpFunc returns a [Func] that compares with the method .Cmp().
func NewMethodCmpFunc() Func {
	return compareMethodCmp
}

func compareMethodCmp(c *Comparator, st *State, v1, v2 reflect.Value) (Result, bool) {
	f, ok := getMethodCmpFunc(v1.Type())
	if !ok {
		return nil, false
	}
	if !v1.CanInterface() || !v2.CanInterface() {
		return nil, false
	}
	cmpRes := f.Call([]reflect.Value{v1, v2})[0].Interface().(int) //nolint:forcetypeassert // The type of the returned value is already checked.
	if cmpRes == 0 {
		return nil, true
	}
	return Result{Difference{
		Message: fmt.Sprintf(msgMethodCmpNotEqual, cmpRes),
	}}, true
}

var (
	cmdMethodFuncsLock sync.Mutex
	cmdMethodFuncs     = make(map[reflect.Type]*reflect.Value)
)

func getMethodCmpFunc(typ reflect.Type) (reflect.Value, bool) {
	cmdMethodFuncsLock.Lock()
	defer cmdMethodFuncsLock.Unlock()
	fp, ok := cmdMethodFuncs[typ]
	if ok {
		if fp != nil {
			return *fp, true
		}
		return reflect.Value{}, false
	}
	cmdMethodFuncs[typ] = nil
	met, ok := typ.MethodByName("Cmp")
	if !ok {
		return reflect.Value{}, false
	}
	metTyp := met.Type
	if metTyp.NumIn() != 2 || metTyp.In(0) != typ || metTyp.In(1) != typ || metTyp.NumOut() != 1 || metTyp.Out(0) != reflect.TypeOf(int(1)) {
		return reflect.Value{}, false
	}
	cmdMethodFuncs[typ] = &met.Func
	return met.Func, true
}

// Result is a list of [Difference].
type Result []Difference

// Format implements [fmt.Formatter].
//
// See [Difference.Format] for supported verb and flag.
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
	_, _ = io.WriteString(s, ": ")
	_, _ = io.WriteString(s, d.Message)
	if s.Flag('+') {
		if d.V1 != "" || d.V2 != "" {
			_, _ = io.WriteString(s, "\n\tv1=")
			_, _ = io.WriteString(s, d.V1)
			_, _ = io.WriteString(s, "\n\tv2=")
			_, _ = io.WriteString(s, d.V2)
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

// Path represents a field path, which is a list of [PathElem].
//
// Elements are stored in reverse order, the first element is the deepest.
// It helps to prepend elements to the path efficiently.
type Path []PathElem

// Format implements [fmt.Formatter].
//
// It only supports the 'v' verb.
func (p Path) Format(s fmt.State, verb rune) {
	if len(p) == 0 {
		_, _ = io.WriteString(s, ".")
		return
	}
	for i := len(p) - 1; i >= 0; i-- {
		p[i].Format(s, verb)
	}
}

// PathElem is a single element in a [Path].
type PathElem struct {
	Struct *string `json:"struct,omitempty"`
	Map    *string `json:"map,omitempty"`
	Index  *int    `json:"index,omitempty"`
}

// Format implements [fmt.Formatter].
//
// It only supports the 'v' verb.
func (e PathElem) Format(s fmt.State, verb rune) {
	switch {
	case e.Struct != nil:
		_, _ = io.WriteString(s, ".")
		_, _ = io.WriteString(s, *e.Struct)
	case e.Map != nil:
		_, _ = io.WriteString(s, "[")
		_, _ = io.WriteString(s, *e.Map)
		_, _ = io.WriteString(s, "]")
	case e.Index != nil:
		_, _ = io.WriteString(s, "[")
		_, _ = strconvio.WriteInt(s, int64(*e.Index), 10)
		_, _ = io.WriteString(s, "]")
	}
}

func toPtr[V any](v V) *V {
	return &v
}

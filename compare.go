// Package compare provide comparison utilities.
package compare

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pierrre/go-libs/bufpool"
)

// MaxSliceDifferences is the maximum number of differences for a slice.
// If the value is reached, the comparison is stopped for the current slice.
// It is also used for array.
// Set to 0 disables it.
var MaxSliceDifferences = 10

// Compare compares 2 values.
func Compare(v1, v2 any) Result {
	return compare(
		reflect.ValueOf(v1),
		reflect.ValueOf(v2),
	)
}

func compare(v1, v2 reflect.Value) Result {
	if r, stop := compareValid(v1, v2); stop {
		return r
	}
	if r, stop := compareType(v1, v2); stop {
		return r
	}
	if r, stop := compareFuncs(v1, v2); stop {
		return r
	}
	return compareKind(v1, v2)
}

func compareValid(v1, v2 reflect.Value) (Result, bool) {
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
		V1:      vl1,
		V2:      vl2,
	}}, true
}

func compareType(v1, v2 reflect.Value) (Result, bool) {
	t1 := v1.Type()
	t2 := v2.Type()
	if t1 == t2 {
		return nil, false
	}
	return Result{Difference{
		Message: msgTypeNotEqual,
		V1:      t1,
		V2:      t2,
	}}, true
}

//nolint:gocyclo // Large switch/case is OK.
func compareKind(v1, v2 reflect.Value) Result {
	switch v1.Kind() {
	case reflect.Bool:
		return compareBool(v1, v2)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareInt(v1, v2)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compareUint(v1, v2)
	case reflect.Float32, reflect.Float64:
		return compareFloat(v1, v2)
	case reflect.Complex64, reflect.Complex128:
		return compareComplex(v1, v2)
	case reflect.String:
		return compareString(v1, v2)
	case reflect.Array:
		return compareArray(v1, v2)
	case reflect.Slice:
		return compareSlice(v1, v2)
	case reflect.Interface:
		return compareInterface(v1, v2)
	case reflect.Pointer:
		return comparePointer(v1, v2)
	case reflect.Struct:
		return compareStruct(v1, v2)
	case reflect.Map:
		return compareMap(v1, v2)
	case reflect.UnsafePointer:
		return compareUnsafePointer(v1, v2)
	case reflect.Chan:
		return compareChan(v1, v2)
	case reflect.Func:
		return compareFunc(v1, v2)
	case reflect.Invalid:
		return nil
	}
	return nil
}

func compareBool(v1, v2 reflect.Value) Result {
	b1 := v1.Bool()
	b2 := v2.Bool()
	if b1 == b2 {
		return nil
	}
	return Result{Difference{
		Message: msgBoolNotEqual,
		V1:      b1,
		V2:      b2,
	}}
}

func compareInt(v1, v2 reflect.Value) Result {
	i1 := v1.Int()
	i2 := v2.Int()
	if i1 == i2 {
		return nil
	}
	return Result{Difference{
		Message: msgIntNotEqual,
		V1:      i1,
		V2:      i2,
	}}
}

func compareUint(v1, v2 reflect.Value) Result {
	u1 := v1.Uint()
	u2 := v2.Uint()
	if u1 == u2 {
		return nil
	}
	return Result{Difference{
		Message: msgUintNotEqual,
		V1:      u1,
		V2:      u2,
	}}
}

func compareFloat(v1, v2 reflect.Value) Result {
	f1 := v1.Float()
	f2 := v2.Float()
	if f1 == f2 {
		return nil
	}
	return Result{Difference{
		Message: msgFloatNotEqual,
		V1:      f1,
		V2:      f2,
	}}
}

func compareComplex(v1, v2 reflect.Value) Result {
	c1 := v1.Complex()
	c2 := v2.Complex()
	if c1 == c2 {
		return nil
	}
	return Result{Difference{
		Message: msgComplexNotEqual,
		V1:      c1,
		V2:      c2,
	}}
}

func compareString(v1, v2 reflect.Value) Result {
	s1 := v1.String()
	s2 := v2.String()
	if s1 == s2 {
		return nil
	}
	return Result{Difference{
		Message: msgStringNotEqual,
		V1:      s1,
		V2:      s2,
	}}
}

func compareArray(v1, v2 reflect.Value) Result {
	var r Result
	diffCount := 0
	for i, n := 0, v1.Len(); i < n; i++ {
		ri := compareArrayIndex(v1, v2, i)
		r = append(r, ri...)
		if len(ri) > 0 {
			diffCount++
			if diffCount >= MaxSliceDifferences && MaxSliceDifferences > 0 {
				break
			}
		}
	}
	return r
}

func compareArrayIndex(v1, v2 reflect.Value, i int) Result {
	r := compare(v1.Index(i), v2.Index(i))
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

var typeByteSlice = reflect.TypeOf([]byte(nil))

func compareSlice(v1, v2 reflect.Value) Result {
	if r, stop := compareNilLenPointer(v1, v2); stop {
		return r
	}
	if v1.Type() == typeByteSlice && bytes.Equal(v1.Bytes(), v2.Bytes()) {
		return nil
	}
	return compareArray(v1, v2)
}

func compareInterface(v1, v2 reflect.Value) Result {
	if r, stop := compareNil(v1, v2); stop {
		return r
	}
	return compare(v1.Elem(), v2.Elem())
}

func comparePointer(v1, v2 reflect.Value) Result {
	if v1.Pointer() == v2.Pointer() {
		return nil
	}
	return compare(v1.Elem(), v2.Elem())
}

func compareStruct(v1, v2 reflect.Value) Result {
	var r Result
	t := v1.Type()
	for i, n := 0, t.NumField(); i < n; i++ {
		r = append(r, compareStructField(v1, v2, i)...)
	}
	return r
}

func compareStructField(v1, v2 reflect.Value, i int) Result {
	r := compare(v1.Field(i), v2.Field(i))
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

func compareMap(v1, v2 reflect.Value) Result {
	if r, stop := compareNilLenPointer(v1, v2); stop {
		return r
	}
	var r Result
	for _, k := range getMapsKeys(v1, v2) {
		r = append(r, compareMapKey(v1, v2, k)...)
	}
	return r
}

func compareMapKey(v1, v2, k reflect.Value) Result {
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
			V1:      vl1,
			V2:      vl2,
		}}
	}
	r := compare(v1, v2)
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

func compareUnsafePointer(v1, v2 reflect.Value) Result {
	p1 := v1.Pointer()
	p2 := v2.Pointer()
	if p1 == p2 {
		return nil
	}
	return Result{Difference{
		Message: msgUnsafePointerNotEqual,
		V1:      p1,
		V2:      p2,
	}}
}

func compareChan(v1, v2 reflect.Value) Result {
	if r, stop := compareNil(v1, v2); stop {
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
			V1:      cap1,
			V2:      cap2,
		}}
	}
	len1 := v1.Len()
	len2 := v2.Len()
	if len1 != len2 {
		return Result{Difference{
			Message: msgLengthNotEqual,
			V1:      len1,
			V2:      len2,
		}}
	}
	return nil
}

func compareFunc(v1, v2 reflect.Value) Result {
	p1 := v1.Pointer()
	p2 := v2.Pointer()
	if p1 == p2 {
		return nil
	}
	return Result{Difference{
		Message: msgFuncPointerNotEqual,
		V1:      p1,
		V2:      p2,
	}}
}

func compareNil(v1, v2 reflect.Value) (Result, bool) {
	nil1 := v1.IsNil()
	nil2 := v2.IsNil()
	if nil1 && nil2 {
		return nil, true
	}
	if nil1 != nil2 {
		return Result{Difference{
			Message: msgOnlyOneIsNil,
			V1:      nil1,
			V2:      nil2,
		}}, true
	}
	return nil, false
}

func compareNilLenPointer(v1, v2 reflect.Value) (Result, bool) {
	if r, stop := compareNil(v1, v2); stop {
		return r, true
	}
	len1 := v1.Len()
	len2 := v2.Len()
	if len1 != len2 {
		return Result{Difference{
			Message: msgLengthNotEqual,
			V1:      len1,
			V2:      len2,
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

// Result is a list of Difference.
type Result []Difference

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
	V1, V2  any
}

// Format implements fmt.Formatter.
//
// It only supports the 'v' verb.
// By default, it show the path and message.
// The '+' flag shows values V1 and V2.
func (d Difference) Format(s fmt.State, verb rune) {
	// We use a buffer in order to reduce memory allocation.
	// fmt.State (and its real type) doesn't (yet?) implement WriteString().
	// TODO remove ?
	buf := bufPool.Get()
	if verb == 'v' {
		_, _ = buf.WriteString(d.Path.String() + ": " + d.Message)
		if s.Flag('+') {
			_, _ = fmt.Fprintf(buf, "\n\tv1="+d.getValueFormat(d.V1)+"\n\tv2="+d.getValueFormat(d.V2), d.V1, d.V2)
		}
	} else {
		_, _ = fmt.Fprintf(buf, "%%!%c(%T)", verb, d)
	}
	_, _ = buf.WriteTo(s)
	bufPool.Put(buf)
}

func (d Difference) getValueFormat(v any) string {
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

// Path represents a field path, which is a list of PathElem.
//
// Elements are stored in reverse order, the first element is the deepest.
// It helps to prepend elements to the path efficiently.
type Path []PathElem

// String returns the string value for a Path.
func (p Path) String() string {
	if len(p) == 0 {
		return "."
	}
	ss := make([]string, len(p))
	for i, e := range p {
		ss[len(ss)-i-1] = e.String()
	}
	return strings.Join(ss, "")
}

// PathElem is a single element in a Path.
type PathElem struct {
	Struct *string
	Map    *string
	Index  *int
}

func (e PathElem) String() string {
	if e.Struct != nil {
		return "." + *e.Struct
	}
	if e.Map != nil {
		return "[" + *e.Map + "]"
	}
	if e.Index != nil {
		return "[" + strconv.Itoa(*e.Index) + "]"
	}
	return ""
}

// Func represents a comparison function.
// It is guaranteed that both values: are valid, of the same type, and can be converted to any.
// If the returned value "stop" is true, the comparison will stop.
type Func func(v1, v2 reflect.Value) (r Result, stop bool)

var funcs []Func

// RegisterFunc registers a Func.
// It allows to handle manually the comparison for certain values.
func RegisterFunc(f Func) {
	funcs = append(funcs, f)
}

func compareFuncs(v1, v2 reflect.Value) (Result, bool) {
	if !v1.CanInterface() || !v2.CanInterface() {
		return nil, false
	}
	for _, f := range funcs {
		if r, stop := f(v1, v2); stop {
			return r, true
		}
	}
	return nil, false
}

func init() {
	RegisterFunc(compareValue)
	RegisterFunc(compareMethodEqual)
	RegisterFunc(compareMethodCmp)
}

var methodEqualNames []string

// RegisterMethodEqual registers an equal method.
// This methods must be callable as "v1.METHOD(v2) bool".
func RegisterMethodEqual(name string) {
	methodEqualNames = append(methodEqualNames, name)
}

func init() {
	RegisterMethodEqual("Equal")
	RegisterMethodEqual("Eq")
}

func compareMethodEqual(v1, v2 reflect.Value) (Result, bool) {
	for _, name := range methodEqualNames {
		r, stop := compareMethodEqualName(v1, v2, name)
		if stop {
			return r, true
		}
	}
	return nil, false
}

func compareMethodEqualName(v1, v2 reflect.Value, name string) (Result, bool) {
	m := v1.MethodByName(name)
	if !m.IsValid() {
		return nil, false
	}
	t := m.Type()
	if t.NumIn() != 1 || t.In(0) != v2.Type() || t.NumOut() != 1 || t.Out(0) != reflect.TypeOf(true) {
		return nil, false
	}
	if m.Call([]reflect.Value{v2})[0].Interface().(bool) { //nolint:forcetypeassert // The type of the returned value is already checked above.
		return nil, true
	}
	return Result{Difference{
		Message: fmt.Sprintf(msgMethodNotEqual, name),
		V1:      v1.Interface(),
		V2:      v2.Interface(),
	}}, true
}

func compareMethodCmp(v1, v2 reflect.Value) (Result, bool) {
	m := v1.MethodByName("Cmp")
	if !m.IsValid() {
		return nil, false
	}
	t := m.Type()
	if t.NumIn() != 1 || t.In(0) != v2.Type() || t.NumOut() != 1 || t.Out(0) != reflect.TypeOf(int(1)) {
		return nil, false
	}
	c := m.Call([]reflect.Value{v2})[0].Interface().(int) //nolint:forcetypeassert // The type of the returned value is already checked above.
	if c == 0 {
		return nil, true
	}
	return Result{Difference{
		Message: fmt.Sprintf(msgMethodCmpNotEqual, c),
		V1:      v1.Interface(),
		V2:      v2.Interface(),
	}}, true
}

var typeReflectValue = reflect.TypeOf(reflect.Value{})

func compareValue(v1, v2 reflect.Value) (Result, bool) {
	if v1.Type() != typeReflectValue {
		return nil, false
	}
	v1 = v1.Interface().(reflect.Value) //nolint:forcetypeassert // The type assertion is already checked above.
	v2 = v2.Interface().(reflect.Value) //nolint:forcetypeassert // The type assertion is already checked above.
	return compare(v1, v2), true
}

func getMapsKeys(v1, v2 reflect.Value) []reflect.Value {
	ks := v1.MapKeys()
	for _, k2 := range v2.MapKeys() {
		found := false
		for _, k := range ks {
			if len(compare(k2, k)) == 0 {
				found = true
				break
			}
		}
		if !found {
			ks = append(ks, k2)
		}
	}
	sortValues(ks, v1.Type().Key())
	return ks
}

func sortValues(s []reflect.Value, t reflect.Type) {
	sort.Slice(s, newSortLess(s, t))
}

func newSortLess(s []reflect.Value, t reflect.Type) func(i, j int) bool {
	switch t.Kind() { //nolint:exhaustive // We have a default case.
	case reflect.Bool:
		return newSortLessBool(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newSortLessInt(s)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return newSortLessUint(s)
	case reflect.Float32, reflect.Float64:
		return newSortLessFloat(s)
	case reflect.Complex64, reflect.Complex128:
		return newSortLessComplex(s)
	case reflect.String:
		return newSortLessString(s)
	default:
		return newSortLessGeneric(s)
	}
}

func newSortLessBool(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return !s[i].Bool() && s[j].Bool()
	}
}

func newSortLessInt(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].Int() < s[j].Int()
	}
}

func newSortLessUint(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].Uint() < s[j].Uint()
	}
}

func newSortLessFloat(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].Float() < s[j].Float()
	}
}

func newSortLessComplex(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		ci := s[i].Complex()
		cj := s[j].Complex()
		if real(ci) < real(cj) {
			return true
		}
		if real(ci) > real(cj) {
			return false
		}
		return imag(ci) < imag(cj)
	}
}

func newSortLessString(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].String() < s[j].String()
	}
}

func newSortLessGeneric(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return fmt.Sprint(s[i]) < fmt.Sprint(s[j])
	}
}

func toPtr[V any](v V) *V {
	return &v
}

var bufPool = bufpool.Pool{}

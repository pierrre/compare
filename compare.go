// Package compare provide comparison utilities.
package compare

import (
	"bytes"
	"fmt"
	"reflect"
)

// Compare compares 2 values.
func Compare(v1, v2 interface{}) Result {
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

// nolint: gocyclo
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
	case reflect.Ptr:
		return comparePtr(v1, v2)
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
		r = r.Merge(ri)
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
		d.Path = IndexedPath{
			Index: i,
			Next:  d.Path,
		}
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

func comparePtr(v1, v2 reflect.Value) Result {
	if v1.Pointer() == v2.Pointer() {
		return nil
	}
	return compare(v1.Elem(), v2.Elem())
}

func compareStruct(v1, v2 reflect.Value) Result {
	var r Result
	t := v1.Type()
	for i, n := 0, t.NumField(); i < n; i++ {
		r = r.Merge(compareStructField(v1, v2, i))
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
		d.Path = StructPath{
			Field: f,
			Next:  d.Path,
		}
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
		r = r.Merge(compareMapKey(v1, v2, k))
	}
	return r
}

func compareMapKey(v1, v2 reflect.Value, k reflect.Value) Result {
	v1 = v1.MapIndex(k)
	v2 = v2.MapIndex(k)
	vl1 := v1.IsValid()
	vl2 := v2.IsValid()
	if !vl1 || !vl2 {
		return Result{Difference{
			Path: MapPath{
				Key: fmt.Sprint(k),
			},
			Message: msgMapKeyNotDefined,
			V1:      vl1,
			V2:      vl2,
		}}
	}
	r := compare(v1, v2)
	if len(r) == 0 {
		return nil
	}
	ks := fmt.Sprint(k)
	for i, d := range r {
		d.Path = MapPath{
			Key:  ks,
			Next: d.Path,
		}
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

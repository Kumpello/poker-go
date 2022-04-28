package pointers

import "reflect"

// Pointer returns a pointer to v
func Pointer[V any](v V) *V {
	return &v
}

// NilPointer returns a pointer to v, but for default values it returns nil
func NilPointer[V any](v V) *V {
	// declare empty value for cmp
	var emptyV V
	if reflect.DeepEqual(v, emptyV) {
		return nil
	}
	return &v
}

package pointers

import "testing"

func Test_NilPointer_StringType_DefaultValue(t *testing.T) {
	var str string
	res := NilPointer(str)
	if res != nil {
		t.Fatalf("res should be nil")
	}
}

func Test_NilPointer_Int64_Type_DefaultValue(t *testing.T) {
	var i int64
	res := NilPointer(i)
	if res != nil {
		t.Fatalf("res should be nil")
	}
}

func Test_NilPointer_StringType_NotDefaultValue(t *testing.T) {
	str := "str"
	res := NilPointer(str)
	if res == nil {
		t.Fatalf("i should point to value")
	}
}

func Test_NilPointer_Int64_Type_NotDefaultValue(t *testing.T) {
	i := int64(64)
	res := NilPointer(i)
	if res == nil {
		t.Fatalf("i should point to value")
	}
}

func Test_Pointer_ComplexType_DefaultValue(t *testing.T) {
	type cpx struct {
		dummy1 int64
		dummy2 string
	}

	c := cpx{
		dummy1: 0,
		dummy2: "",
	}
	res := NilPointer(c)
	if res != nil {
		t.Fatalf("res should be a nil (all default values)")
	}
}

func Test_Pointer_ComplexType_NotDefaultValue(t *testing.T) {
	type cpx struct {
		dummy1 int64
		dummy2 string
	}

	c := cpx{
		dummy1: 10,
		dummy2: "",
	}
	res := NilPointer(c)
	if res == nil {
		t.Fatalf("res should point to value")
	}
}

func Test_Pointer_SliceType_DefaultValue(t *testing.T) {
	var slc []any
	res := NilPointer(slc)
	if res != nil {
		t.Fatalf("slc should be a nil (all default values)")
	}
}

func Test_Pointer_SliceType_NotDefaultValue(t *testing.T) {
	slc := []any{"a"}
	res := NilPointer(slc)
	if res == nil {
		t.Fatalf("slc should point to value")
	}
}

package core

// FieldType represent D value type
type FieldType int32

// DType enum
const (
	UnknownType FieldType = iota
	StringType
	IntTpye
	Int64Type
	UintType
	Uint64Type
	Float32Type
	Float64Type
	DurationType
	BoolType
)

// Field is for encoder
type Field struct {
	Key       string
	Value     interface{}
	Type      FieldType
	StringVal string
	Int64Val  int64
}

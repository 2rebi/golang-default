package def

import "reflect"

type (
	ErrorJustInit struct {
		Errors []*ErrorJustInitField
	}

	ErrorJustInitField struct {
		StructName string
		FieldName string
		TryValue string
		Field reflect.Value
	}
)

func (e *ErrorJustInit) Error() string {
	return ""
}


func (j *ErrorJustInitField) Error() string {
	return ""
}
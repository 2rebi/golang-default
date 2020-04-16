package def

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	ErrorJustInit struct {
		Errors []*ErrorJustInitField
	}

	ErrorJustInitField struct {
		StructName string
		FieldName string
		FieldType string
		TryValue string
		Cause error
		Target reflect.Value
	}
)

func (e *ErrorJustInit) Error() string {
	builder := strings.Builder{}
	builder.WriteString("Struct Errors\n")
	for i := range e.Errors {
		builder.WriteString("\t"+e.Errors[i].Error())
	}
	return builder.String()
}


func (j *ErrorJustInitField) Error() string {
	return fmt.Sprintf("struct:%s,field:(%s %s),try:%s,error:%s",
		j.StructName, j.FieldName, j.FieldType, j.TryValue, strings.ReplaceAll(j.Cause.Error(), "\n", "\n\t"))
}
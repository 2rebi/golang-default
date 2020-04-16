package main

import (
	"fmt"
	def "github.com/rebirthlee/golang-default"
)

type (
	ErrorSample struct {
		Number int `def:"string"`
	}

	NestedErrorSample struct {
		Number int `def:"string"`
		Nest ErrorSample `def:"dive"`
	}
)

func main() {
	sample := ErrorSample{}
	err := def.Init(&sample)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = def.JustInit(&sample)
	if err != nil {
		justErr := err.(*def.ErrorJustInit)
		fmt.Println(justErr.Error())

		for i := range justErr.Errors {
			fieldErr := justErr.Errors[i]
			fmt.Println(fieldErr.Error())
			//or do something
		}
	}

	println()
	println()

	nested := NestedErrorSample{}
	err = def.JustInit(&nested)
	if err != nil {
		justErr := err.(*def.ErrorJustInit)
		for i := range justErr.Errors {
			fieldErr := justErr.Errors[i]
			//do something

			nestedErr, ok := fieldErr.Cause.(*def.ErrorJustInit)
			if ok {
				fmt.Println(nestedErr.Error())
				// more do something
			}
		}
	}
}

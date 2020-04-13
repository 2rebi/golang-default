package main

import (
	"fmt"
	"github.com/rebirthlee/golang-default"
)

type (
	Sample struct {
		Arr [3]NestedSample `def:"dive,dive"`
		Numbers [3]int `def:"dive,32000"`
	}

	NestedSample struct {
		Number int `def:"777"`
		Name string `def:"Nested Sample"`
	}
)


func main() {
	sample := def.MustNew(Sample{}).(*Sample)
	fmt.Println(sample)
}

package main

import (
	"fmt"
	def "github.com/rebirthlee/golang-default"
)

type (
	Sample struct {
		Sli []NestedSample `def:"dive(3,7),dive"`
		Numbers []int `def:"dive(5),-32000"`
	}

	NestedSample struct {
		Number int `def:"333"`
		Name string `def:"Sample Nested"`
	}
)


func main() {
	sample := def.MustNew(Sample{}).(*Sample)
	fmt.Println(sample)
	fmt.Printf("sample.Sli len = %d, cap = %d\n", len(sample.Sli), cap(sample.Sli))
	fmt.Printf("sample.Numbers len = %d, cap = %d\n", len(sample.Numbers), cap(sample.Numbers))
}


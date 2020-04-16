package main

import (
	"fmt"
	def "github.com/rebirthlee/golang-default"
)

type (
	Sample struct {
		Interface interface{} `def:"{\"name\":\"rebirth lee\",\"age\":25}"`
		MapJson map[string]interface{} `def:"{\"items\":[\"item - 1\",\"item - 2\"]}"`

		NestedStruct Struct `def:"{\"title\":\"Struct Convert Json Test\",\"subtitle\":\"This is Struct\",\"num\":123}"`
		PtrNestedStruct *Struct `def:"{\"title\":\"Struct Pointer Convert Json Test\",\"subtitle\":\"This is Struct Pointer\",\"num\":321}"`
	}

	Struct struct {
		Title string `json:"title"`
		SubTitle string `json:"subtitle"`
		Number int `json:"num"`
	}
)

func main() {
	s := def.MustNew(Sample{}).(*Sample)
	fmt.Println(s)
	fmt.Println(s.Interface)
	fmt.Println(s.MapJson)
	fmt.Println(s.NestedStruct)
	fmt.Println(s.PtrNestedStruct)
}

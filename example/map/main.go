package main

import (
	"fmt"
	def "github.com/rebirthlee/golang-default"
)

type (
	NestedMap struct {
		PtrStructVal map[string]*Struct `def:"dive{\"john doe\":dive,\"some one\":dive,\"key\":dive}"`
		PtrStructKey map[*Struct]bool   `def:"dive{dive:true,dive:true,dive:true}"`
		MapVal map[string]map[*Struct]bool `def:"dive{\"key1\":dive{dive:true,dive:false},\"key2\":dive{dive:false,dive:false}}"`
	}

	Struct struct {
		Name string `def:"who?"`
	}
)

func main() {
	n := def.MustNew(NestedMap{}).(*NestedMap)
	fmt.Println(n)

	fmt.Println("---------PtrStructVal---------")
	for key, val := range n.PtrStructVal {
		fmt.Printf("%+v : %+v\n", key, val)
	}

	fmt.Println("---------PtrStructKey---------")
	for key, val := range n.PtrStructKey {
		fmt.Printf("%+v : %+v\n", key, val)
	}

	fmt.Println("---------MapVal---------")
	for key, val := range n.MapVal {
		fmt.Printf("MapVal[%s]\n", key)
		for inKey, inVal := range val {
			fmt.Printf("\t%+v : %+v\n", inKey, inVal)

		}
	}
}

package main

import (
	"fmt"
	"github.com/rebirthlee/golang-default"
)

type (
	Person struct {
		Age int `def:"-"`
		Name string `def:"rich guy"`
		Wallet *Wallet `def:"dive"`
	}

	Wallet struct {
		BrandName string `def:"Louis Vuitton"`
		Name *string `def:"unknown"`
		Money *int `def:"-"`
		Weight *int `def:"-"`
	}
)

func main() {
	{
		// Nested - New
		i, err := def.New(Person{})
		if err == nil {
			p := i.(*Person)
			fmt.Println(p)
			fmt.Println(p.Wallet)
		}
	}

	{
		// Nested - MustNew
		p := def.MustNew(Person{}).(*Person)
		fmt.Println(p)
		fmt.Println(p.Wallet)
	}

	{
		// Nested - JustNew
		//TODO
	}

	{
		// Nested - Init
		p := Person{}
		if err := def.Init(&p); err != nil {
			// ...err
			fmt.Println("Init, Handle Error")
		} else {
			fmt.Println(p)
			fmt.Println(p.Wallet)
		}
	}

	{
		// Nested - MustInit
		p := Person{}
		def.MustInit(&p)
		fmt.Println(p)
		fmt.Println(p.Wallet)
	}

	{
		// Nested - JustInit
		//TODO
	}
}

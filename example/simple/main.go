package main

import (
	"fmt"
	"github.com/rebirthlee/golang-default"
	"time"
)

type Person struct {
	Age int `def:"20"`
	Name string `def:"rebirth lee"`
	PocketName *string `def:"bitcoin"`
	AliveTime time.Duration `def:"175200h"`
}

func main() {

	{
		// Simple - New
		i, err := def.New(Person{})
		if err == nil {
			p := i.(*Person)
			fmt.Println(p)
		}
	}

	{
		// Simple - MustNew
		p := def.MustNew(Person{}).(*Person)
		fmt.Println(p)
	}

	{
		// Simple - JustNew
		i, err := def.JustNew(Person{})
		if err == nil {
			p := i.(*Person)
			fmt.Println(p)
		}
	}

	{
		// Simple - Init
		p := Person{}
		if err := def.Init(&p); err != nil {
			// ...err
			fmt.Println("Init, Handle Error")
		} else {
			fmt.Println(p)
		}
	}

	{
		// Simple - MustInit
		p := Person{}
		def.MustInit(&p)
		fmt.Println(p)
	}

	{
		// Simple - JustInit
		p := Person{}
		if err := def.JustInit(&p); err != nil {
			// ...err
			fmt.Println("JustInit, Handle Error")
		} else {
			fmt.Println(p)
		}
	}

}

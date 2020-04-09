package main

import (
	"fmt"
	def "github.com/rebirthlee/golang-default"
)

type Person struct {
	Age int `def:"20"`
	Name string `def:"hellp"`
	Say func() string `def:"person_default_say"`
}

func init() {
	if err := def.SetFunc("person_default_say", func(self interface{}) interface{} {
		p := self.(*Person)
		return func() string {
			return fmt.Sprintf("My name is %s, %d years old", p.Name, p.Age)
		}
	}); err != nil {
		panic(err)
	}
}

func main() {
	p := def.MustNew(Person{}).(*Person)
	fmt.Println(p)
	fmt.Println(p.Say())
}

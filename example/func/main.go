package main

import (
	"fmt"
	"github.com/rebirthlee/golang-default"
)

type Person struct {
	Age int `def:"20"`
	Name string `def:"hellp"`
	Do func() string `def:"person_default_do"`
}

func init() {
	if err := def.SetFunc("person_default_do", func(self *Person) interface{} {
		return self.IntroducingMySelf
	}); err != nil {
		panic(err)
	}
}

func main() {
	p := def.MustNew(Person{}).(*Person)
	fmt.Println(p)
	fmt.Println(p.Do())
	p.Name = "rebirth lee"
	p.Age = 25
	fmt.Println(p.Do())
}


func (p *Person) IntroducingMySelf() string {
	return fmt.Sprintf("My name is %s, %d years old", p.Name, p.Age)
}
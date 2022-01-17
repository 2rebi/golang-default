package main

import (
	"fmt"
	def "github.com/rebirthlee/golang-default"
)

type InitString string

type Sample struct {
	ExportField string `def:"export field"`
	notExportField string
	StringField *InitString
}

func (i *InitString) Init() {
	fmt.Printf("InitString(%p) call Init\n", i)
	*i = "Hello String Field"
}

func (s *Sample) Init() {
	fmt.Printf("Sample(%p) call Init\n", s)
	s.notExportField = "not export field"
}

func main() {

	{
		// New
		i, err := def.New(Sample{})
		if err == nil {
			s := i.(*Sample)
			showFields(s)
		}
	}

	{
		// MustNew
		s := def.MustNew(Sample{}).(*Sample)
		showFields(s)
	}

	{
		// JustNew
		i, err := def.JustNew(Sample{})
		if err == nil {
			s := i.(*Sample)
			showFields(s)
		}
	}

	{
		// Init
		s := Sample{}
		if err := def.Init(&s); err != nil {
			// ...err
			fmt.Println("Init, Handle Error")
		} else {
			showFields(&s)
		}
	}

	{
		// MustInit
		s := Sample{}
		def.MustInit(&s)
		showFields(&s)
	}

	{
		// JustInit
		s := Sample{}
		if err := def.JustInit(&s); err != nil {
			// ...err
			fmt.Println("JustInit, Handle Error")
		} else {
			showFields(&s)
		}
	}
}

func showFields(s *Sample) {
	fmt.Printf("Struct Address : %p\n", s)
	fmt.Println("p.ExportField :", s.ExportField)
	fmt.Println("p.notExportField :", s.notExportField)
	fmt.Printf("p.StringField : %p\n", s.StringField)
	fmt.Println("*p.StringField :", *s.StringField)
	println()
	println()
}
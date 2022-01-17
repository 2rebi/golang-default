# golang-default
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2FRebirthLee%2Fgolang-default%2Fbadge%3Fref%3Dmaster&style=flat)](https://actions-badge.atrox.dev/RebirthLee/golang-default/goto?ref=master)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/rebirthlee/golang-default)](https://github.com/RebirthLee/golang-default/releases)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/rebirthlee/golang-default)](https://golang.org/doc/go1.13)
[![license](https://img.shields.io/badge/license-BEER--WARE-green)](/LICENSE.md)

Initialize or New Struct with default value

## Setup
```
go get github.com/rebirthlee/golang-default
```

## Support type
- `int/8/16/32/64`, `uint/8/16/32/64`
```go
`def:"100"` 
`def:"-100"` // support only signed integer
// support hexadecimal, octal, binary
`def:"0xff"` // hexadecimal
`def:"0xFF"` // hexadecimal
`def:"0o77"` // octal
`def:"0b11001111"` // binary
```
- `float32/64`
```go
`def:"3.141592653589"`
`def:"-3.141592"`
```
- `complex64/128`
```go
 // `def:"{real part},{imaginary part}"`
`def:"3.14,-10"`
`def:"-3.14,3"`
```

- `time.Duration`
```go
// calling time.ParseDuration
`def:"1h"` // 1 * time.Hour
```

- `time.Time`
```go
`def:"now"` // time.Now()
`def:"+1h"` // time.Now().Add(1 * time.Hour)
`def:"-1h"` // time.Now().Add(-1 * time.Hour)
```

- Nested Struct
```go
type Parent struct {
	Name string `def:"Parent Struct"`
	OneChild Child `def:"dive"`
	TwoChild Child `def:"dive"`
}

type Child struct {
	Key string `def:"unknown"`
	Number int `def:"-1"`
}
```

- Pointer of type
```go
type Parent struct {
	Name *string `def:"Parent Struct"`
	OneChild *Child `def:"dive"`
	TwoChild *Child `def:"dive"`
}

type Child struct {
	Key *string `def:"unknown"`
	Number *int `def:"-1"`
}
```

- Array
```go
type Sample struct {
	Arr [3]int `def:"dive,-1"` // [-1, -1, -1]
}
```

- Slice
```go
type Sample struct {
	// `def:"dive({length},{capacity(optional)}),{value}"`
	Sli []int `def:"dive(5),-1"` // [-1, -1, -1, -1, -1]
  
	//cap(SliWithCap) == 7
	SliWithCap []NestedSample `def:"dive(3,7),dive"` // [{nested struct},{nested struct},{nested struct}] 
}

type NestedSample struct {
	Name string `def:"nested struct"`
}
```

- Map
```go
type Sample struct {
	DiveMap map[string]*Struct `def:"dive{\"john doe\":dive,\"some one\":dive,\"key\":dive}"`
/*
	{ 
		"john doe": &{who?},
		"some one": &{who?},
		"key": &{who?}
	}
*/ 
	StructKeyMap map[*Struct]bool `def:"dive{dive:true,dive:false,dive:true}"`
/*
	{ 
		&{who?}: true,
		&{who?}: false,
		&{who?}: true
	}
*/
	DiveNestedMap map[string]map[*Struct]bool `def:"dive{\"key1\":dive{dive:true,dive:false},\"key2\":dive{dive:false,dive:false}}"`
/*
	{ 
		"key1": {
			&{who?}: true,
			&{who?}: false
		},
		"key2": {
			&{who?}: false,
			&{who?}: false
		}
	}
*/
}

Struct struct {
	Name string `def:"who?"`
}
```

- Json
```go
type Sample struct {
	Arr [3]int `def:"[1,2,3]"` // [1,2,3] 
	Sli []string `def:"[\"slice 1\",\"slice 2\"]"` // [slice 1,slice 2]
	Map map[string]interface{} `def:"{\"key1\":123,\"key2\":\"value\",\"nested map\":{\"key\":\"val\"}}"` 
/*
	{ 
		"key1":123,
		"key2":"value",
		"nested map":{
			"key":"val"
		}
	}
*/

	Nested NestedSample `def:"{\"displayName\":\"nested struct type\"}"` // {nested struct type}
	PtrNested *NestedSample `def:"{\"displayName\":\"nested struct pointer type\"}"` // &{nested struct pointer type}
}

type NestedSample struct {
	Name string `json:"displayName"`
}
```

- Function

[Example](/example/func/main.go)

## Usage
### Simple

```go
import (
	"fmt"
	"github.com/rebirthlee/golang-default"
)

type Person struct {
	Age int `def:"20"`
	Name string `def:"hellp"`
}

...
var p Person
if err := def.Init(&p); err != nil {
	// error handle
}
fmt.Println(p) //out: {20 hellp}
```

### Init
If you got error, the next field of struct will not be initialized.

```go
if err := def.Init(&p); err != nil {
	// error handle
}
```

### JustInit
Even though it has an error, It will try to initialize all fields.
And you can know that error field of struct.

```go
if err := def.JustInit(&p); err != nil {
	justErr := err.(*def.ErrorJustInit)
	fmt.Println(justErr.Error())
	// error handle
}
```

### MustInit
It isn't return error. but it will be panic when you has an error.

```go
def.MustInit(&p) // hasn't return error.
```

### New
If you got error, it will be return nil.

```go
i, err := def.New(Person{})
if err != nil {
	// error handle
} else {
  p := i.(*Person)
  fmt.Println(p) //out: &{20 hellp}
}
```

### JustNew
Even though it has an error, It must return pointer of struct with error.

```go
i, err := def.JustNew(Person{})
if err != nil {
	justErr := err.(*def.ErrorJustInit)
	fmt.Println(justErr.Error())
	// error handle
} 
p := i.(*Person)
fmt.Println(p) //out: &{20 hellp}
```

### MustNew
It isn't return error. but it will be panic when you has an error.

```go
p := def.MustNew(Person{}).(*Person) // hasn't return error.
fmt.Println(p) //out: &{20 hellp}
```

License
---
[`THE BEER-WARE LICENSE (Revision 42)`](http://en.wikipedia.org/wiki/Beerware)

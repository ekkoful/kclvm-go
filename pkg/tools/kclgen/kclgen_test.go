package kclgen

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string `kcl:"name:name,type:str"`
	Age  int    `kcl:"name:age,type:int"`
}

func Test_parseGoStruct(t *testing.T) {
	p := &Person{}
	s := GenKclSchemaCode(p)
	fmt.Println(s)
}

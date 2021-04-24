package main

import (
	"testing"
)

func TestCreateProduct(t *testing.T) {
	p := &Product{
		Code:  "GO_WEST",
		Price: 20,
	}
	err := CreateProduct(p)
	if err != nil {
		panic(err)
	}
}

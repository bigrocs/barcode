package main

import (
	"fmt"
	"testing"
)

func TestAddGoods(t *testing.T) {
	data, err := GetURL("06923450605288")
	fmt.Println(data, err)
	t.Log(t)
}

package main

import (
	"fmt"
	"testing"

	"github.com/bigrocs/barcode/drives"
)

func TestAddGoods(t *testing.T) {
	chinatrace := &drives.Chinatrace{}
	data, err := chinatrace.Get("6917878036526")
	fmt.Println(data, err)
	t.Log(t)
}

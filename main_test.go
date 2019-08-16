package main

import (
	"fmt"
	"testing"

	"github.com/bigrocs/barcode/drives"
)

func TestAddGoods(t *testing.T) {
	chinatrace := &drives.Chinatrace{
		BaseHost: "http://webapi.chinatrace.org",
		// BaseHost: "https://www.chinatrace.org",
		Key: "V7N3Xpm4jpRon/WsZ8X/63G8oMeGdUkA8Luxs1CenTY=",
	}
	data, err := chinatrace.Get("6923450605289")
	fmt.Println(data, err)
	t.Log(t)
}

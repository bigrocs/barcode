package main

import (
	"fmt"
	"testing"

	"github.com/bigrocs/barcode/drives"
)

func TestAddGoods(t *testing.T) {
	chinatrace := &drives.Chinatrace{
		BaseHost: "http://webapi.chinatrace.org",
		Key:      "V7N3Xpm4jpRon/WsZ8X/63G8oMeGdUkA8Luxs1CenTY=",
	}
	data, err := chinatrace.Get("6917878036526")
	fmt.Println(data, err)
	t.Log(t)
}

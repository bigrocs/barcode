package main

import (
	"github.com/bigrocs/barcode/data"
)

// Barcode  商品仓库接口
type Barcode interface {
	Get(code string) (*data.Goods, error)
}

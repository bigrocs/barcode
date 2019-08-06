package chinatrace

import (
	"fmt"
	"testing"
)

func TestAddGoods(t *testing.T) {
	url, err := GetURL("06923450605288")
	fmt.Println(url, err)
	t.Log(t)
}

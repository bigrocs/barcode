package chinatrace

import (
	"fmt"
	"testing"
)

func TestAddGoods(t *testing.T) {
	url, err := GetURL("6938166920785")
	fmt.Println(url, err)
	t.Log(t)
}

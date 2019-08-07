package drives

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/bigrocs/barcode/data"
)

var (
	BASE_HOST string = "http://webapi.chinatrace.org"
	KEY       string = "V7N3Xpm4jpRon/WsZ8X/63G8oMeGdUkA8Luxs1CenTY="
)

// Chinatrace 国家食品(产品)安全追溯平台
type Chinatrace struct {
}

// Get 获取条码商品信息
func (srv *Chinatrace) Get(code string) (goods data.Goods, err error) {
	url, err := srv.getURL(code)
	fmt.Println(url)
	return goods, err
}

// getURL 获取请网址
func (srv *Chinatrace) getURL(code string) (url string, err error) {
	url = "/api/getProductData?productCode=" + code
	mac, err := srv.getMac(url)
	if err != nil {
		return url, err
	}
	url = BASE_HOST + url + "&mac=" + mac
	return url, err
}

// getMac 计算 mac
func (srv *Chinatrace) getMac(url string) (mac string, err error) {
	key, err := base64.StdEncoding.DecodeString(KEY)
	if err != nil {
		return mac, err
	}
	mac = srv.hmacSha256(url, string(key))
	fmt.Println(mac)
	return mac, err
}

// hmacSha256 加密
func (srv *Chinatrace) hmacSha256(src string, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(src))
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

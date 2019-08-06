package chinatrace

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	BASE_HOST string = "http://webapi.chinatrace.org"
	KEY       string = "V7N3Xpm4jpRon/WsZ8X/63G8oMeGdUkA8Luxs1CenTY="
)

// GetURL 获取请网址
func GetURL(code string) (url string, err error) {
	url = "/api/getProductData?productCode=" + code
	mac, err := GetMac(url)
	if err != nil {
		return url, err
	}
	url = BASE_HOST + url + "&mac=" + mac
	return url, err
}

// GetMac 计算 mac
func GetMac(url string) (mac string, err error) {
	key, err := base64.StdEncoding.DecodeString(KEY)
	if err != nil {
		return mac, err
	}
	mac = HmacSha256(url, string(key))
	fmt.Println(mac)
	return mac, err
}

// HmacSha256 加密
func HmacSha256(src string, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(src))
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

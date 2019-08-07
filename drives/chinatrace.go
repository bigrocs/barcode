package drives

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bigrocs/barcode/data"
)

var (
	IMAGE_HOST string = ""
	BASE_HOST  string = "http://webapi.chinatrace.org"
	KEY        string = "V7N3Xpm4jpRon/WsZ8X/63G8oMeGdUkA8Luxs1CenTY="
)

// Chinatrace 国家食品(产品)安全追溯平台
type Chinatrace struct {
}

// Get 获取条码商品信息
func (srv *Chinatrace) Get(code string) (goods *data.Goods, err error) {
	url, err := srv.getURL(code)
	headers := srv.Headers()
	httpContent, err := srv.request(url, headers)
	if err != nil {
		return goods, err
	}
	content, err := srv.response(httpContent)
	if err != nil {
		return goods, err
	}
	goods, err = srv.handerGoods(content)
	if err != nil {
		return goods, err
	}
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

// response 处理返回数据
func (srv *Chinatrace) handerGoods(content map[string]interface{}) (goods *data.Goods, err error) {
	goods = &data.Goods{}
	content = content["d"].(map[string]interface{})

	goods.Barcode = content["productCode"].(string)
	goods.Name = content["ItemName"].(string)
	goods.Image, _ = srv.handerImages(content["Image"])
	goods.BrandName = content["BrandName"].(string)
	goods.Specification = content["ItemSpecification"].(string)
	goods.Name = content["ItemName"].(string)
	goods.Unspsc, _ = strconv.ParseInt(content["ItemClassCode"].(string), 10, 64)
	goods.UnspscName = content["ItemClassCode"].(string)
	goods.Source = content["codeSource"].(string)
	goods.FirmName = content["FirmName"].(string)
	goods.FirmAddress = content["FirmAddress"].(string)
	goods.FirmStatus = content["FirmStatus"].(string)

	return goods, err
}

// handerImage 处理返回图片
func (srv *Chinatrace) handerImages(items interface{}) (images []string, err error) {
	for _, item := range items.([]interface{}) {
		// 转为 map 然后读取 Imageurl 然后转为 string
		img := item.(map[string]interface{})["Imageurl"].(string)
		if IMAGE_HOST != "" {
			img = strings.Replace(img, "http://www.anccnet.com", IMAGE_HOST, -1)
		}
		images = append(images, img)
	}
	return images, err
}

// response 处理返回数据
func (srv *Chinatrace) response(httpContent []byte) (content map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(httpContent), &content)
	if err != nil {
		return content, err
	}
	return content, err
}

// request 请求获取数据
func (srv *Chinatrace) request(url string, headers map[string]string) (httpContentString []byte, err error) {
	// http-Client
	client := &http.Client{}
	// request
	request, _ := http.NewRequest("GET", url, strings.NewReader(""))

	for k, v := range headers {
		request.Header.Set(k, v)
	}
	// post-request
	resp, err := client.Do(request)
	if err != nil {
		return httpContentString, err
	}
	defer resp.Body.Close()

	httpContent, err := ioutil.ReadAll(resp.Body)
	return httpContent, err

}

// Headers Headers 构建
func (srv *Chinatrace) Headers() (headers map[string]string) {
	return map[string]string{"Accept": "application/json", "Content-Type": "application/json;charset=utf-8"}
}

// getMac 计算 mac
func (srv *Chinatrace) getMac(url string) (mac string, err error) {
	key, err := base64.StdEncoding.DecodeString(KEY)
	if err != nil {
		return mac, err
	}
	mac = srv.hmacSha256(url, string(key))
	return mac, err
}

// hmacSha256 加密
func (srv *Chinatrace) hmacSha256(src string, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(src))
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

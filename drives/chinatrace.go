package drives

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bigrocs/barcode/data"
	"github.com/gomsa/tools/uitl"
)

// Chinatrace 国家食品(产品)安全追溯平台
type Chinatrace struct {
	// "http://xxx.xxx.com"
	ImageHost string `json:"image_host,omitempty"`
	// "http://webapi.chinatrace.org"
	BaseHost string `json:"base_host,omitempty"`
	// "V7N3Xpm4jpRon/WsZ8X/63G8oMeGdUkA8Luxs1CenTY="
	Key string `json:"key,omitempty"`
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
	err = srv.getSubGoods(content, goods)
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
	url = srv.BaseHost + url + "&mac=" + mac
	return url, err
}

// getSubGoods 获取商品子属性
func (srv *Chinatrace) getSubGoods(content map[string]interface{}, goods *data.Goods) (err error) {
	content = content["d"].(map[string]interface{})
	url := content["ItemDescription"].(string)
	// 请求商品详情页面
	headers := srv.Headers()
	httpContent, err := srv.request(url, headers)
	if err != nil {
		return err
	}
	// 截取商品信息
	html := uitl.ConvertToString(string(httpContent), "gbk", "utf8")
	regex, err := regexp.Compile(`SetValue\('(.*?)','(.*?)'\)`)
	atts := regex.FindAllStringSubmatch(html, -1)
	// 处理商品详情数据
	srv.handerSubGoods(atts, goods)
	return err
}

// handerSubGoods 处理商品详情数据
func (srv *Chinatrace) handerSubGoods(atts [][]string, goods *data.Goods) (err error) {
	for _, varr := range atts {
		switch varr[1] {
		case `Att_Sys_zh-cn_141_G`:
			if goods.Name == "无" {
				goods.Name = varr[2]
			}
		case `Att_Sys_en_us_141_G`:
			goods.EnName = varr[2]
		case `Att_Sys_en-us_141_G`:
			goods.EnName = varr[2]
		case `Att_Sys_zh-cn_101_G`:
			width, _ := strconv.ParseFloat(varr[2], 64)
			goods.Width = int64(width)
		case `Att_Sys_zh-cn_104_G`:
			switch varr[2] {
			case `厘米`:
				goods.Width = goods.Width * 10
			case `米`:
				goods.Width = goods.Width * 100
			}
		case `Att_Sys_zh-cn_106_G`:
			height, _ := strconv.ParseFloat(varr[2], 64)
			goods.Height = int64(height)
		case `Att_Sys_zh-cn_326_G`:
			switch varr[2] {
			case `厘米`:
				goods.Height = goods.Height * 10
			case `米`:
				goods.Height = goods.Height * 100
			}
		case `Att_Sys_zh-cn_118_G`:
			depth, _ := strconv.ParseFloat(varr[2], 64)
			goods.Depth = int64(depth)
		case `Att_Sys_zh-cn_331_G`:
			switch varr[2] {
			case `厘米`:
				goods.Depth = goods.Depth * 10
			case `米`:
				goods.Depth = goods.Depth * 100
			}

		case `Att_Sys_zh-cn_10_G`:
			netWeight, _ := strconv.ParseFloat(varr[2], 64)
			goods.NetWeight = int64(netWeight)
		case `Att_Sys_zh-cn_189_G`:
			switch varr[2] {
			case `千克`:
				goods.NetWeight = goods.NetWeight * 1000
			}
		case `Att_Sys_zh-cn_54_G`:
			grossWeight, _ := strconv.ParseFloat(varr[2], 64)
			goods.GrossWeight = int64(grossWeight)
		case `Att_Sys_zh-cn_84_G`:
			switch varr[2] {
			case `千克`:
				goods.GrossWeight = goods.GrossWeight * 1000
			}
		case `Att_Sys_zh-cn_22_G`:
			regex, _ := regexp.Compile(`.*\((.*?)\)`)
			atts := regex.FindAllStringSubmatch(varr[2], -1)
			goods.UnspscName = atts[0][1]
		case `Att_Sys_zh-cn_35_G`:
			goods.Unit = varr[2]
		case `Att_Sys_zh-cn_74_G`:
			goods.Country = varr[2]
		case `Att_Sys_zh-cn_405_G`:
			goods.Place = varr[2]
		}

	}
	return err
}

// handerGoods 处理商品数据
func (srv *Chinatrace) handerGoods(content map[string]interface{}) (goods *data.Goods, err error) {
	goods = &data.Goods{}
	content = content["d"].(map[string]interface{})

	goods.Barcode = content["productCode"].(string)
	goods.Name = content["ItemName"].(string)
	goods.Image, _ = srv.handerImages(content["Image"])
	goods.BrandName = content["BrandName"].(string)
	goods.Specification = strings.Replace(content["ItemSpecification"].(string), "×", "x", -1)
	goods.Name = content["ItemName"].(string)
	goods.Unspsc, _ = strconv.ParseInt(content["ItemClassCode"].(string), 10, 64)
	goods.Source = content["codeSource"].(string)
	goods.FirmName = content["FirmName"].(string)
	goods.FirmAddress = content["FirmAddress"].(string)
	goods.FirmStatus = content["FirmStatus"].(string)

	return goods, err
}

// handerImage 处理返回图片
func (srv *Chinatrace) handerImages(items interface{}) (images []string, err error) {
	// 防止没有图片报错
	switch items.(type) {
	case string:
	default:
		for _, item := range items.([]interface{}) {
			// 转为 map 然后读取 Imageurl 然后转为 string
			img := item.(map[string]interface{})["Imageurl"].(string)
			if srv.ImageHost != "" {
				img = strings.Replace(img, "http://www.anccnet.com", srv.ImageHost, -1)
			}
			images = append(images, img)
		}
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
	// 	Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
	// Accept-Encoding: gzip, deflate
	// Accept-Language: zh-CN,zh;q=0.9
	// Cache-Control: max-age=0
	// Connection: keep-alive
	// Cookie: ASP.NET_SessionId=qex3rbiyw5ft5diz1ltcql45
	// Host: v1.gds.org.cn
	// Referer: http://v1.gds.org.cn/goods.aspx?base_id=F25F56A9F703ED747848039802026E2BEF7B34F6AD959CFADBDABC7FD77685C0872F36D0507DA774
	// Upgrade-Insecure-Requests: 1
	// User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36
	return map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json;charset=utf-8",
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
	}
}

// getMac 计算 mac
func (srv *Chinatrace) getMac(url string) (mac string, err error) {
	key, err := base64.StdEncoding.DecodeString(srv.Key)
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

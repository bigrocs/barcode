package data

// Goods 返回商品结构
type Goods struct {
	Barcode       string   // 条码
	Name          string   // 产品名称
	EnName        string   // 产品英文名称
	Images         []string // 产品图片
	BrandName     string   // 品牌
	Specification string   // 规格
	Unit          string   // 单位
	Width         int64    // 宽(毫米)
	Height        int64    // 高(毫米)
	Depth         int64    // 深(毫米)
	NetWeight     int64    // 净重(克)
	GrossWeight   int64    // 总重(克)
	Unspsc        int64    // 商品及服务编码
	UnspscName    string   // 商品及服务编码分类名称
	Source        string   // 产地代码
	Place         string   // 产地
	Country       string   // 国家
	FirmName      string   // 公司名称
	FirmAddress   string   // 公司地址
	FirmStatus    string   // 公司状态
}

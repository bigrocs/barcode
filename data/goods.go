package data

// Goods 返回商品结构
type Goods struct {
	Barcode       string   // 条码
	Name          string   // 产品名称
	Image         []string // 产品图片
	BrandName     string   // 品牌
	Specification string   // 规格
	Width         string   // 宽(毫米)
	Height        string   // 高(毫米)
	Depth         string   // 深(毫米)
	NetWeight     string   // 净重(克)
	GrossWeight   string   // 总重(克)
	Unspsc        int64    // 商品及服务编码
	UnspscName    string   // 商品及服务编码分类名称
	Source        string   // 产地
	FirmName      string   // 公司名称
	FirmAddress   string   // 公司地址
	FirmStatus    string   // 公司状态
}

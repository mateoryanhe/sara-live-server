package ipgeo

// CountryInfo IP 所属国家信息
type CountryInfo struct {
	// ISO 3166-1 alpha-2,如 CN、US
	Code string
	// 国家名称(优先中文,否则英文)
	Name string
}

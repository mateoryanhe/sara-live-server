package ipgeo

// CountryInfo IP 所属国家信息
type CountryInfo struct {
	// ISO 3166-1 alpha-2,如 CN、US
	Code string
	// 展示名(优先 constants/country.NameZh)
	Name string
}

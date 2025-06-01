package IpQuery

type IPInfoResp struct {
	Status     string  `json:"status"`     // success or fail
	Message    string  `json:"message"`    // included only when status is fail, Can be one of the following: private range, reserved range, invalid query
	Continent  string  `json:"continent"`  // 洲
	Country    string  `json:"country"`    // 国家
	Region     string  `json:"region"`     // 省份
	RegionName string  `json:"regionName"` // 省份名称
	City       string  `json:"city"`       // 城市
	District   string  `json:"district"`   // 区
	Zip        string  `json:"zip"`        // 邮编
	Lat        float64 `json:"lat"`        // 纬度
	Lon        float64 `json:"lon"`        // 经度
	TimeZone   string  `json:"timezone"`   // 时区
	Query      string  `json:"query"`      // 查询的ip
}

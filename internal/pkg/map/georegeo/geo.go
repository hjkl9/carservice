package georegeo

import (
	"carservice/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 匹配级别
const (
// todo
)

type Geocode struct {
	Country  string `json:"country"`  // 地址所在的省份名, 例如：北京市。此处需要注意的是，中国的四大直辖市也算作省级单位。
	City     string `json:"city"`     // 地址所在的城市名, 例如：北京市
	Citycode string `json:"citycode"` // 城市编码, 例如：010
	District string `json:"district"` // 地址所在的区, 例如：朝阳区
	Street   string `json:"street"`   // 街道, 例如：阜通东大街
	Number   string `json:"number"`   // 门牌, 例如：6号
	Adcode   string `json:"adcode"`   // 区域编码, 例如：110101
	Location string `json:"location"` // 坐标点, 经度，纬度
	Level    string `json:"level"`    // 匹配级别
}

type GeoResp struct {
	Status   uint8     `json:"status"`   // 返回结果状态值 <返回值为 0 或 1，0 表示请求失败；1 表示请求成功>
	Count    uint      `json:"count"`    // 返回结果数目 <返回结果的个数>
	Info     string    `json:"info"`     // 返回状态说明 <当 status 为 0 时，info 会返回具体错误原因，否则返回“OK”。详情可以参阅info状态表 `https://lbs.amap.com/api/webservice/guide/tools/info`>
	Geocodes []Geocode `json:"geocodes"` // 地理编码信息列表 <结果对象列表>
}

type Geo struct {
	cfg config.AMapConf

	geocodes *[]Geocode
}

func NewGeo(cfg config.AMapConf) *Geo {
	s := make([]Geocode, 0)
	return &Geo{
		cfg:      cfg,
		geocodes: &s,
	}
}

func (g *Geo) ByAddress(address string) (*Geo, error) {
	url := "https://restapi.amap.com/v3/geocode/geo?key=%s&address=%s"
	resp, err := http.Get(fmt.Sprintf(url, g.cfg.Key, address))
	if err != nil {
		return g, err
	}
	body := resp.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return g, err
	}
	var geoResp GeoResp
	json.Unmarshal(bodyBytes, &geoResp)
	g.geocodes = &geoResp.Geocodes
	return g, nil
}

func (g *Geo) GetFirstGeoCode() (Geocode, bool) {
	if len(*(g.geocodes)) > 0 {
		return (*g.geocodes)[0], true
	}
	return Geocode{}, false
}

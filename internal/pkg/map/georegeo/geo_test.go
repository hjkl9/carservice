package georegeo_test

import (
	"carservice/internal/config"
	"carservice/internal/pkg/map/georegeo"
	"testing"
)

// TestGeoByAddress
// test ok.
func TestGeoByAddress(t *testing.T) {
	config := config.AMapConf{
		Key: "9c1395d0dfcba6e065e50c27634d05ea",
	}
	geo := georegeo.NewGeo(config)
	// rs, err := geo.ByAddress("广东省广州市白云区新市街棠安路31号") // 113.256671,23.190103
	rs, err := geo.ByAddress("广东省东莞市万江岳潭村南向十号") // 113.726964,23.039920
	if err != nil {
		t.Errorf("failed to search location, err: %s\n", err.Error())
		return
	}
	geoCode, ok := rs.GetFirstGeoCode()
	if ok {
		t.Logf("GeoCode: %#v\n", geoCode)
	}

	// for k, v := range (*rs).Geocodes {
	// 	t.Logf("-- %d --: mem(%p)%#v\n", k, &v, v)
	// }
}

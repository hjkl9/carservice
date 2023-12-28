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
	rs, err := geo.ByAddress("广东省东莞市万江岳潭村南向十号")
	if err != nil {
		t.Errorf("failed to search location, err: %s\n", err.Error())
		return
	}
	geoCode := rs.GetFirstGeoCode()
	t.Logf("GeoCode: %#v\n", geoCode)

	// for k, v := range (*rs).Geocodes {
	// 	t.Logf("-- %d --: mem(%p)%#v\n", k, &v, v)
	// }
}

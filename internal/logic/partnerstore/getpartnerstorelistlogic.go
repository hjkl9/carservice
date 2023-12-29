package partnerstore

import (
	"context"
	"strconv"
	"strings"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/map/georegeo"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPartnerStoreListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPartnerStoreListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPartnerStoreListLogic {
	return &GetPartnerStoreListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type StoreListItem struct {
	Id          uint    `db:"id"`
	Title       string  `db:"title"`
	FullAddress string  `db:"fullAddress"`
	Longitude   float32 `db:"longitude"` // 经度
	Latitude    float32 `db:"latitude"`  // 纬度
	Distance    float32 `db:"distance"`
}

func (l *GetPartnerStoreListLogic) GetPartnerStoreList(req *types.GetPartnerStoreListReq) (resp []types.PartnerStoreListItem, err error) {
	limit := 0
	if req.LimitGap <= 0 {
		limit = 15
	} else {
		limit = int(req.LimitGap)
	}
	// 实例化地理编码
	geo, err := georegeo.NewGeo(l.svcCtx.Config.AMapConf).ByAddress(req.Address)
	if err != nil {
		if serviceErr := err.(*georegeo.AMapError); serviceErr != nil {
			return nil, errcode.InternalServerError.Lazy(serviceErr.GetCode() + ": " + serviceErr.GetMsg())
		}
		return nil, errcode.InternalServerError.Lazy("第三方服务发生错误")
	}
	// 获取第一个
	geocode, ok := geo.GetFirstGeoCode()
	if !ok {
		return []types.PartnerStoreListItem{}, nil
	}
	// 分割经纬度
	var location struct {
		Longitude float64
		Latitude  float64
	}
	// 解析经纬度 //
	locationSlice := strings.Split(geocode.Location, ",")
	// 解析经度
	location.Longitude = func() float64 {
		v, _ := strconv.ParseFloat(locationSlice[0], 64)
		return v
	}()
	// 解析纬度
	location.Latitude = func() float64 {
		v, _ := strconv.ParseFloat(locationSlice[1], 64)
		return v
	}()
	var stores []*StoreListItem
	query := "SELECT `id`, `title`, `full_address` AS `fullAddress`, `longitude`, `latitude`, (ST_DISTANCE_SPHERE(POINT(?, ?), POINT(longitude, latitude))) / 1000 AS `distance` FROM `partner_stores` WHERE `status` = ? HAVING `distance` <= ?"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return nil, errcode.NewDatabaseErrorx().GetError(err)
	}
	if err = stmt.SelectContext(l.ctx, &stores, location.Longitude, location.Latitude, 1, limit); err != nil {
		return nil, errcode.NewDatabaseErrorx().GetError(err)
	}
	var interfaceData []types.PartnerStoreListItem
	for _, v := range stores {
		interfaceData = append(interfaceData, types.PartnerStoreListItem{
			Id:          (*v).Id,
			Title:       (*v).Title,
			FullAddress: (*v).FullAddress,
			Gap:         (*v).Distance,
			Unit:        "千米",
		})
	}
	// ! use fake list data.
	// l.fakeListData(&interfaceData)
	return interfaceData, nil
}

func (l *GetPartnerStoreListLogic) calculateGap(list *[]*StoreListItem) {
	// todo()
}

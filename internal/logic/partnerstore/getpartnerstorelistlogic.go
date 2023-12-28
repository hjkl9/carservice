package partnerstore

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/map/georegeo"
	"carservice/internal/svc"
	"carservice/internal/types"

	gofakeit "github.com/brianvoe/gofakeit/v6"
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
	Longitude   float64 `db:"longitude"` // 经度
	Latitude    float64 `db:"latitude"`  // 纬度
}

func (l *GetPartnerStoreListLogic) GetPartnerStoreList(req *types.GetPartnerStoreListReq) (resp []types.PartnerStoreListItem, err error) {
	// 实例化地理编码
	geo, err := georegeo.NewGeo(l.svcCtx.Config.AMapConf).ByAddress(req.Address)
	if err != nil {
		return nil, errcode.InternalServerError.Lazy("第三方服务获取时发生错误", err.Error())
	}
	// 获取第一个
	geocode, ok := geo.GetFirstGeoCode()
	if !ok {
		return nil, errcode.New(http.StatusNoContent, "-", "找不到该地址")
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
	fmt.Printf("经纬度: %#v\n", location)
	// todo: 数据库查询附近的门店
	// var stores []*StoreListItem
	// query := "SELECT * FROM `partner_stores` WHERE `longitude` <= ? AND `latitude` <= ?"
	// stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	// if err != nil {
	// 	return nil, errcode.NewDatabaseErrorx().GetError(err)
	// }
	// if err = stmt.SelectContext(l.ctx, &stores, 100.00, 100.00); err != nil {
	// 	return nil, errcode.NewDatabaseErrorx().GetError(err)
	// }
	var interfaceData []types.PartnerStoreListItem
	// for _, v := range stores {
	// 	interfaceData = append(interfaceData, types.PartnerStoreListItem{
	// 		Id:          (*v).Id,
	// 		Title:       (*v).Title,
	// 		FullAddress: (*v).FullAddress,
	// 	})
	// }
	// ! use fake list data.
	l.fakeListData(&interfaceData)
	return interfaceData, nil
}

func (l *GetPartnerStoreListLogic) calculateGap(list *[]*StoreListItem) {
	// todo()
}

func (l *GetPartnerStoreListLogic) fakeListData(dest *[]types.PartnerStoreListItem) {
	var i uint = 0
	for i = 0; i < 15; i++ {
		(*dest) = append((*dest), types.PartnerStoreListItem{
			Id:          i + 1,
			Title:       gofakeit.School(),
			FullAddress: gofakeit.Address().Street,
			Gap:         gofakeit.UintRange(0, 5000),
			Unit:        "米",
		})
	}
}

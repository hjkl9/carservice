package partnerstore

import (
	"context"

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
	Longitude   float64 `db:"longitude"`
	Latitude    float64 `db:"latitude"`
}

func (l *GetPartnerStoreListLogic) GetPartnerStoreList(req *types.GetPartnerStoreListReq) (resp []types.PartnerStoreListItem, err error) {
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

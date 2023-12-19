package carbrand

import (
	"context"
	"fmt"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BrandOptionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBrandOptionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BrandOptionListLogic {
	return &BrandOptionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BrandOptionListLogic) BrandOptionList() (resp []types.BrandOptionListItem, err error) {
	sql := "SELECT `brand_id` AS `id`, `brand_name` AS `label`, `pinyin` AS `pinyin` FROM `car_brands` WHERE `business_status` = ?"
	var list []types.BrandOptionListItem
	err = l.svcCtx.DBC.Select(&list, sql, "1")
	fmt.Println(list)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	}
	return list, nil
}

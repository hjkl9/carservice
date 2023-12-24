package carbrandseries

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OptionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOptionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OptionListLogic {
	return &OptionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OptionListLogic) OptionList(req *types.GetCarBrandSeriesOptionListReq) (resp []types.CarBrandSeriesOptionListItem, err error) {
	query := "SELECT `series_id` AS `id`, `series_name` AS `label`, '-' AS `pinyin` FROM `car_brand_series` WHERE `brand_id` = ? AND `business_status` = ?"
	var list []types.CarBrandSeriesOptionListItem
	if err = l.svcCtx.DBC.Select(&list, query, req.BrandId, 1); err != nil {
		return nil, errcode.InternalServerError.
			SetMsg("查询数据时出现错误").
			SetDetails(err.Error())
	}
	return list, nil
}

package carbrand

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

func (l *OptionListLogic) OptionList() (resp []types.CarBrandOptionListItem, err error) {
	sql := "SELECT `brand_id` AS `id`, `brand_name` AS `label`, `pinyin` AS `pinyin` FROM `car_brands` WHERE `business_status` = ? ORDER BY `pinyin` ASC"
	var list []types.CarBrandOptionListItem
	// 预处理查询语句
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, sql)
	stmt.SelectContext(l.ctx, &list, "1")
	// err = l.svcCtx.DBC.Select(&list, sql, "1")
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	}
	return list, nil
}

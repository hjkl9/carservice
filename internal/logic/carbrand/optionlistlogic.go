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
	// rs, err := l.svcCtx.RDBC.Ping(l.ctx).Result()
	// logc.Debugf(l.ctx, "connect result: %s\n", rs)
	// if err != nil {
	// 	return nil, api.NewApiCode(http.StatusInternalServerError, "1", "操作缓存时发生错误")
	// }
	// if rs == "PONG" {
	// 	k := "brands"
	// 	n, err := l.svcCtx.RDBC.Exists(l.ctx, k).Result()
	// 	if err != nil {
	// 		return nil, errcode.RedisGetErr
	// 	}
	// 	if !(n > 0) {
	// 		// 预处理查询语句
	// 		stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, sql)
	// 		if err != nil {
	// 			return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	// 		}
	// 		err = stmt.SelectContext(l.ctx, &list, "1")
	// 		if err != nil {
	// 			return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	// 		}
	// 		jsonBytes, err := json.Marshal(&list)
	// 		if err != nil {
	// 			return nil, errcode.JSONMarshalErr
	// 		}
	// 		if status := l.svcCtx.RDBC.Set(l.ctx, k, string(jsonBytes), 0); status.Err() != nil {
	// 			return nil, errcode.RedisSetErr
	// 		}
	// 		return list, nil
	// 	} else {
	// 		val, err := l.svcCtx.RDBC.Get(l.ctx, k).Result()
	// 		if err != nil {
	// 			return nil, errcode.RedisGetErr
	// 		}
	// 		if err = json.Unmarshal([]byte(val), &list); err != nil {
	// 			return nil, errcode.JSONUnmarshalErr
	// 		}
	// 		return list, nil
	// 	}
	// }
	// 预处理查询语句
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, sql)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	}
	err = stmt.SelectContext(l.ctx, &list, "1")
	// err = l.svcCtx.DBC.Select(&list, sql, "1")
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	}
	return list, nil
}

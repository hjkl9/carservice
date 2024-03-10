package carreplacement

import (
	"context"
	"fmt"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReplacementListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplacementListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplacementListLogic {
	return &ReplacementListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type replacement struct {
	Id          uint    `db:"id"`
	ParentId    uint    `db:"parentId"`
	Title       string  `db:"title"`
	EstF32Price float32 `db:"estF32Price"`
	EstU64Price uint64  `db:"estU64Price"`
	Counter     uint8   `db:"counter"`
}

type carSeries struct {
	officialPriceUp float64 `db:"official_price_up"`
}

type OfficialPrice struct {
	OfficialPriceUp   float64 `db:"officialPriceUp"`
	OfficialPriceDown float64 `db:"officialPriceDown"`
}

func (l *ReplacementListLogic) ReplacementList(req *types.CarReplacementReq) (resp []types.CarReplacement, err error) {
	// todo: 查询车型价格区间
	if req.CarSeriesId == 0 {
		return nil, errcode.InvalidParametersErr.SetMessage("无效的参数 carSeriesId")
	}

	var count int8
	query := "SELECT COUNT(1) AS `count` FROM `car_brand_series` WHERE `series_id` = ? LIMIT 1;"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "查询车型是否存在[预处理时]发生错误")
		return nil, errcode.DatabasePrepareErr
	}
	if err = stmt.GetContext(l.ctx, &count, req.CarSeriesId); err != nil {
		logc.Error(l.ctx, "查询车型是否存在[获取数据时]发生错误")
		return nil, errcode.DatabaseGetErr
	}

	if count == 0 {
		return nil, errcode.InvalidParametersErr.SetMessage("无效的参数 carSeriesId")
	}

	var officialPrice OfficialPrice
	query = "SELECT `official_price_up` AS `officialPriceUp`, `official_price_down` AS `officialPriceDown` FROM `car_brand_series` WHERE `series_id` = ? LIMIT 1;"
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "查询官方售价[预处理时]发生错误", err)
		return nil, errcode.DatabasePrepareErr
	}
	if err = stmt.GetContext(l.ctx, &officialPrice, req.CarSeriesId); err != nil {
		logc.Error(l.ctx, "查询官方售价[获取数据时]发生错误", err)
		return nil, errcode.DatabaseGetErr
	}
	// True: 高端
	// False: 低端
	var grade bool = func() bool {
		const P float64 = 30.00
		if officialPrice.OfficialPriceDown > P || officialPrice.OfficialPriceUp > P {
			return true
		}
		// 其他条件或规则
		// todo: 处理过滤规则
		return false
	}()
	var addQuery = func() string {
		if grade {
			return "`hm_est_f32_price` AS `estF32Price`, `hm_est_u64_price` AS `estU64Price`"
		} else {
			return "`lm_est_f32_price` AS `estF32Price`, `lm_est_u64_price` AS `estU64Price`"
		}
	}()

	// todo: 处理暂无报价的车型

	query = "SELECT `id` AS `id`, `parent_id` AS `parentId`, `title` AS `title`, %s , `counter` AS `counter` FROM car_replacements ORDER BY `sort`;"
	query = fmt.Sprintf(query, addQuery)

	var data []*replacement

	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "查询配件列表[预处理时]发生错误", err)
		return nil, errcode.DatabasePrepareErr
	}
	err = stmt.SelectContext(l.ctx, &data)
	if err != nil {
		logc.Error(l.ctx, "查询配件列表[获取数据时]发生错误", err)
		return nil, errcode.DatabaseGetErr
	}

	// list to tree
	var result = list2tree(data, 0)

	return result, nil
}

func list2tree(data []*replacement, pid uint) []types.CarReplacement {
	newTree := make([]types.CarReplacement, 0)
	for _, elem := range data {
		if elem.ParentId == pid {
			t := types.CarReplacement{
				Id: elem.Id,
				// ParentId:    elem.ParentId,
				Title:       elem.Title,
				EstF32Price: elem.EstF32Price,
				EstU64Price: elem.EstU64Price,
				Counter:     uint(elem.Counter),
			}
			t.Children = list2tree(data, elem.Id)
			newTree = append(newTree, t)
		}
	}
	return newTree
}

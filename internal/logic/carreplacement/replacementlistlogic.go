package carreplacement

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

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

type tree struct {
	replacement
	children []*tree `json:"children"`
}

func (l *ReplacementListLogic) ReplacementList() (resp []types.CarReplacement, err error) {
	query := "SELECT `id` AS `id`, `parent_id` AS `parentId`, `title` AS `title`, `est_f32_price` AS `estF32Price`, `est_u64_price` AS `estU64Price` , `counter` AS `counter` FROM car_replacements ORDER BY `sort`;"

	var data []*replacement

	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	}
	err = stmt.SelectContext(l.ctx, &data)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
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

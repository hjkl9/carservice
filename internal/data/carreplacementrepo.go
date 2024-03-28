package data

import (
	"carservice/internal/datatypes/carreplacement"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logc"
)

type CarReplacementRepo interface {
	GetEstPriceListByIdSet(
		ctx context.Context,
		f func() bool,
		idSet []uint,
	) ([]carreplacement.Replacement, error)
}

type carReplacement struct {
	db *sqlx.DB
}

func newCarReplacement(db *sqlx.DB) *carReplacement {
	return &carReplacement{db}
}

func (cr *carReplacement) GetEstPriceListByIdSet(
	ctx context.Context,
	f func() bool,
	idSet []uint,

) ([]carreplacement.Replacement, error) {
	var dest []carreplacement.Replacement

	subQuery1 := ""
	if f() {
		subQuery1 = "`hm_est_f32_price` AS `estF32Price`, `hm_est_u64_price` AS `estU64Price`"
	} else {
		subQuery1 = "`lm_est_f32_price` AS `estF32Price`, `lm_est_u64_price` AS `estU64Price`"
	}

	mainQuery := fmt.Sprintf("SELECT `id` AS `id`, `parent_id` AS `parentId`, `title` AS `title`, %s , `counter` AS `counter` FROM car_replacements WHERE `id` IN (?) ORDER BY `sort`;", subQuery1)
	query, args, _ := sqlx.In(mainQuery, idSet)
	stmtx, err := cr.db.PreparexContext(ctx, query)
	if err != nil {
		logc.Error(ctx, "查询配件列表[预处理时]发生错误", err)
		return nil, err
	}
	err = stmtx.SelectContext(ctx, &dest, args...)
	if err != nil {
		logc.Error(ctx, "查询配件列表[获取数据时]发生错误", err)
		return nil, err
	}

	return dest, nil
}

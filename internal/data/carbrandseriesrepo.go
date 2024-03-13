package data

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logc"
)

const SplitPrice float64 = 30.00

type CarBrandSeriesRepo interface {
	CheckIfSeriesExists(ctx context.Context, brandId, seriesId uint) (bool, error)
	GetOfficialPrice(ctx context.Context, seriesId interface{}) (float64, float64, error)
	CheckGradeByCarSeries(down, up float64) bool
}

type carBrandSeries struct {
	db *sqlx.DB
}

func newCarBrandSeries(db *sqlx.DB) *carBrandSeries {
	return &carBrandSeries{db}
}

func (cbs *carBrandSeries) CheckIfSeriesExists(ctx context.Context, brandId, seriesId uint) (bool, error) {
	var count uint8
	query := "SELECT COUNT(1) AS `count` FROM `car_brands` `cb` JOIN `car_brand_series` `cbs` ON `cb`.`brand_id` = `cbs`.`brand_id` WHERE `cbs`.`brand_id` = ? AND `cbs`.`series_id` = ? LIMIT 1;"
	stmt, err := cbs.db.PreparexContext(ctx, query)
	if err != nil {
		return false, err
	}
	if err = stmt.GetContext(ctx, &count, brandId, seriesId); err != nil {
		return false, err
	}
	return count == 1, nil
}

func (cbs *carBrandSeries) GetOfficialPrice(ctx context.Context, seriesId interface{}) (float64, float64, error) {
	query := "SELECT `official_price_up` AS `officialPriceUp`, `official_price_down` AS `officialPriceDown` FROM `car_brand_series` WHERE `series_id` = ? LIMIT 1;"

	var officialPrice struct {
		OfficialPriceUp   float64 `db:"officialPriceUp"`
		OfficialPriceDown float64 `db:"officialPriceDown"`
	}

	stmtx, err := cbs.db.PreparexContext(ctx, query)
	if err != nil {
		return 0.0, 0.0, err
	}
	if err = stmtx.GetContext(ctx, &officialPrice, seriesId); err != nil {
		logc.Error(ctx, "查询官方售价[获取数据时]发生错误", err)
		return 0.0, 0.0, err
	}

	return officialPrice.OfficialPriceDown, officialPrice.OfficialPriceUp, nil
}

func (cbs *carBrandSeries) CheckGradeByCarSeries(down, up float64) bool {
	if down > SplitPrice || up > SplitPrice {
		return true
	}
	// 其他条件或规则
	// todo: 处理过滤规则

	return false
}

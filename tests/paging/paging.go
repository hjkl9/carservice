package paging

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Paging struct {
	List     interface{}
	Previous int
	Next     int
	PrevUrl  string
	NextUrl  string
}

type SeriesListItem struct {
	SeriesId   uint   `db:"series_id"`
	SeriesName string `db:"series_name"`
}

func GetData(db *sqlx.DB) error {
	ctx := context.Background()
	query := "SELECT `series_id`, `series_name` FROM `car_brand_series` WHERE `series_id` > 0 LIMIT 20"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var sli SeriesListItem
		err = rows.Scan(&sli)
		if err != nil {
			return err
		}

	}
	return nil
}

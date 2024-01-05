package data

import (
	"carservice/internal/datatypes/partnerstore"
	"context"

	"github.com/jmoiron/sqlx"
)

type PartnerStoreRepo interface {
	GetListByLatLng(context.Context, GetListByLatLngArgs) (*[]*partnerstore.StoreListItem, error)
}

type partnerStore struct {
	db *sqlx.DB
}

func newPartnerStore(db *sqlx.DB) *partnerStore {
	return &partnerStore{db}
}

type GetListByLatLngArgs struct {
	Lat      float32
	Lng      float32
	LimitGap uint8
}

func (ps *partnerStore) GetListByLatLng(ctx context.Context, args GetListByLatLngArgs) (*[]*partnerstore.StoreListItem, error) {
	var items []*partnerstore.StoreListItem
	query := "SELECT `id`, `title`, `full_address` AS `fullAddress`, `longitude`, `latitude`, (ST_DISTANCE_SPHERE(POINT(?, ?), POINT(longitude, latitude))) / 1000 AS `distance` FROM `partner_stores` WHERE `status` = ? HAVING `distance` <= ?"
	stmt, err := ps.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = stmt.SelectContext(
		ctx,
		&items,
		args.Lng,
		args.Lat,
		1,
		args.LimitGap,
	); err != nil {
		return nil, err
	}
	return &items, nil
}

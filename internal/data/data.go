package data

import "github.com/jmoiron/sqlx"

type RepoFactory interface {
	UserRelated() UserRepo
	UserOrder() UserOrderRepo
	PartnerStoreRelated() PartnerStoreRepo
	PingRelated() PingRepo
	CarBrandSeriesRepo() CarBrandSeriesRepo
	CarReplacementRepo() CarReplacementRepo
}

type datastore struct {
	dbc *sqlx.DB
}

func NewDatastore(dbc *sqlx.DB) RepoFactory {
	return &datastore{dbc}
}

func (ds *datastore) UserRelated() UserRepo {
	return newUser(ds.dbc)
}

func (ds *datastore) PingRelated() PingRepo {
	return newPing(ds.dbc)
}

func (ds *datastore) UserOrder() UserOrderRepo {
	return newUserOrder(ds.dbc)
}

func (ds *datastore) PartnerStoreRelated() PartnerStoreRepo {
	return newPartnerStore(ds.dbc)
}

func (ds *datastore) CarBrandSeriesRepo() CarBrandSeriesRepo {
	return newCarBrandSeries(ds.dbc)
}

func (ds *datastore) CarReplacementRepo() CarReplacementRepo {
	return newCarReplacement(ds.dbc)
}

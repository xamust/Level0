package cashdata

import (
	"context"
	"service/internal/app/model"
	"service/internal/app/store"
)

type CashData struct {
	config   *Config
	ctx      context.Context
	CashMass []model.Order
	store    *store.Store
}

//init new cash...
func NewCash(config *Config, ctx context.Context, cashMap []model.Order, store *store.Store) *CashData {
	return &CashData{
		config:   config,
		ctx:      ctx,
		CashMass: cashMap,
		store:    store,
	}
}

//append to cash slice...
func (c *CashData) SetCash(order model.Order) {
	c.CashMass = append(c.CashMass, order)
}

//restore cash...
func (c *CashData) RestoreCash() error {
	//getLastId...
	id, err := c.store.GetLastDataId(c.ctx)
	if err != nil {
		return err
	}
	for i, j := id, 1; i > 0; i, j = i-1, j+1 {
		//get order by id's from db...
		mod, err := c.store.GetDataById(c.ctx, i)
		if err != nil {
			return err
		}
		//append in cash slice...
		c.SetCash(*mod)
		//check count of cash data from configs...
		if j == c.config.CountOfCash {
			break
		}
	}
	return nil
}

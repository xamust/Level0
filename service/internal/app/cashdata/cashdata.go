package cashdata

import (
	"context"
	"fmt"
	"service/internal/app/model"
	"service/internal/app/store"
)

type CashData struct {
	ctx     context.Context
	CashMap map[int]model.Order
	store   store.Repository
}

//init new cash...
func NewCash(ctx context.Context, cashMap map[int]model.Order, store store.Repository) *CashData {
	return &CashData{
		ctx:     ctx,
		CashMap: cashMap,
		store:   store,
	}
}

//append to cash slice...
func (c *CashData) SetCash(id int, order model.Order) {
	c.CashMap[id] = order
}

//restore cash...
func (c *CashData) RestoreCash(id int) error {
	if id > 0 {
		//get order by id's from db...
		mod, err := c.store.GetDataById(c.ctx, id)
		if err != nil {
			return err
		}
		//append in cash slice...
		c.SetCash(id, *mod)
		go func() error {
			if err := c.RestoreCash(id - 1); err != nil {
				return err
			}
			return nil
		}()
	}
	return nil
}

//get by id, from cash map...
func (c *CashData) GetById(id int) (*model.Order, error) {
	//if id exist, return order...
	if _, ok := c.CashMap[id]; ok {
		mod := c.CashMap[id]
		return &mod, nil
	}
	//if not, return err...
	return nil, fmt.Errorf("Order number %d, not in the database", id)
}

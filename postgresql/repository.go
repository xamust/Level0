package postgresql

import (
	"L0/model"
	"context"
)

type Repository interface {
	Close()
	InsertData(ctx context.Context, order model.Order) error
	InsertDataOrder(ctx context.Context, order model.Order) (int, error)
	InsertDataDelivery(ctx context.Context, order model.Order) (int, error)
	InsertDataOrdersDelivery(ctx context.Context, order model.Order, id int) error
	InsertDataItem(ctx context.Context, order model.Order) error
	InsertDataPayments(ctx context.Context, order model.Order) error
}

var myImpl Repository

func SetRepo(repo Repository) {
	myImpl = repo
}

func Close() {
	myImpl.Close()
}

func InsertData(ctx context.Context, order model.Order) error {
	return myImpl.InsertData(ctx, order)
}

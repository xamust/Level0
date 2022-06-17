package store

import (
	"context"
	"service/internal/app/model"
)

type Repository interface {
	Open() error
	Close()
	InsertData(ctx context.Context, order model.Order) (int, error)
	InsertDataOrder(ctx context.Context, order model.Order) (int, error)
	InsertDataDelivery(ctx context.Context, order model.Order) (int, error)
	InsertDataOrdersDelivery(ctx context.Context, order model.Order, id int) error
	InsertDataItem(ctx context.Context, order model.Order) error
	InsertDataPayments(ctx context.Context, order model.Order) error
	GetDataById(ctx context.Context, id int) (*model.Order, error)
	GetDataOrder(ctx context.Context, id int) (*model.Order, error)
	GetDataDelivery(ctx context.Context, orderUid string) (delivery *model.Delivery, err error)
	GetDataOrdersDelivery(ctx context.Context, orderUid string) (idDelivery int, err error)
	GetDataItems(ctx context.Context, trackNumber string) (items []model.Item, err error)
	GetDataPayment(ctx context.Context, transaction string) (payment *model.Payment, err error)
	GetLastDataId(ctx context.Context) (int, error)
	InsertIncorrectData(ctx context.Context, message interface{}) error
}

var impl Repository

func SetRepository(repo Repository) {
	impl = repo
}
func Open() {
	impl.Open()
}
func Close() {
	impl.Close()
}
func InsertData(ctx context.Context, order model.Order) (int, error) {
	return impl.InsertData(ctx, order)
}
func InsertDataOrder(ctx context.Context, order model.Order) (int, error) {
	return InsertDataOrder(ctx, order)
}
func InsertDataDelivery(ctx context.Context, order model.Order) (int, error) {
	return InsertDataDelivery(ctx, order)
}
func InsertDataOrdersDelivery(ctx context.Context, order model.Order, id int) error {
	return InsertDataOrdersDelivery(ctx, order, id)
}
func InsertDataItem(ctx context.Context, order model.Order) error {
	return InsertDataItem(ctx, order)
}
func InsertDataPayments(ctx context.Context, order model.Order) error {
	return InsertDataPayments(ctx, order)
}
func GetDataById(ctx context.Context, id int) (*model.Order, error) {
	return GetDataById(ctx, id)
}
func GetDataOrder(ctx context.Context, id int) (*model.Order, error) {
	return GetDataOrder(ctx, id)
}
func GetDataDelivery(ctx context.Context, orderUid string) (delivery *model.Delivery, err error) {
	return GetDataDelivery(ctx, orderUid)
}
func GetDataOrdersDelivery(ctx context.Context, orderUid string) (idDelivery int, err error) {
	return GetDataOrdersDelivery(ctx, orderUid)
}
func GetDataItems(ctx context.Context, trackNumber string) (items []model.Item, err error) {
	return GetDataItems(ctx, trackNumber)
}
func GetDataPayment(ctx context.Context, transaction string) (payment *model.Payment, err error) {
	return GetDataPayment(ctx, transaction)
}
func GetLastDataId(ctx context.Context) (int, error) {
	return GetLastDataId(ctx)
}
func InsertIncorrectData(ctx context.Context, message interface{}) error {
	return InsertIncorrectData(ctx, message)
}

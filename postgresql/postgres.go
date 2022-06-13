package postgresql

import (
	"L0/model"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

//new postgres...
func New(url string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepo{db: db}, nil
}

//close db...
func (r *PostgresRepo) Close() {
	r.db.Close()
}

//insert data...
func (r *PostgresRepo) InsertData(ctx context.Context, order model.Order) error {

	_, err := r.InsertDataOrder(ctx, order)
	if err != nil {
		return err
	}

	idDelivery, err := r.InsertDataDelivery(ctx, order)
	if err != nil {
		return err
	}

	if err = r.InsertDataOrdersDelivery(ctx, order, idDelivery); err != nil {
		return err
	}

	if err = r.InsertDataItem(ctx, order); err != nil {
		return err
	}

	if err = r.InsertDataPayments(ctx, order); err != nil {
		return err
	}

	return nil
}

//insert into order...
func (r *PostgresRepo) InsertDataOrder(ctx context.Context, order model.Order) (int, error) {
	//insert into order...
	var id int
	if err := r.db.QueryRowContext(ctx, "INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id",
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("InsertDataOrder: %v", err)
	}
	return id, nil
}

//insert into delivery...
func (r *PostgresRepo) InsertDataDelivery(ctx context.Context, order model.Order) (int, error) {
	var id int
	if err := r.db.QueryRowContext(ctx, "INSERT INTO delivery(name, phone, zip, city, address, region, email) "+
		"VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("InsertDataDelivery: %v", err)
	}
	return id, nil
}

//insert into orders_delivery...
func (r *PostgresRepo) InsertDataOrdersDelivery(ctx context.Context, order model.Order, id int) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO orders_delivery(order_id,delivery_id) "+
		"VALUES ($1,$2)",
		order.OrderUid,
		id); err != nil {
		return fmt.Errorf("InsertDataOrderDelivery: %v", err)
	}
	return nil
}

//insert into items...
func (r *PostgresRepo) InsertDataItem(ctx context.Context, order model.Order) error {
	for _, v := range order.Items {
		if _, err := r.db.ExecContext(ctx, "INSERT INTO item(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)"+
			"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
			v.ChrtId,
			v.TrackNumber,
			v.Price,
			v.Rid,
			v.Name,
			v.Sale,
			v.Size,
			v.TotalPrice,
			v.NmId,
			v.Brand,
			v.Status,
		); err != nil {
			return fmt.Errorf("InsertDataItem: %v", err)
		}
	}
	return nil
}

//insert into payments...
func (r *PostgresRepo) InsertDataPayments(ctx context.Context, order model.Order) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO payment(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_free)"+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)",
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFree,
	); err != nil {
		return fmt.Errorf("InsertDataPayments: %v", err)
	}
	return nil
}

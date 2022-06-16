package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"service/internal/app/model"
)

type Store struct {
	config *Config
	db     *sql.DB
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

//new postgres ...
func (r *Store) Open() error {
	//open db...
	db, err := sql.Open("postgres", r.config.DBUrl)
	if err != nil {
		return fmt.Errorf("Open postgres: %v", err)
	}
	//ping db...
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Ping postgres: %v", err)
	}
	//set to struct...
	r.db = db
	return nil
}

//close db...
func (r *Store) Close() {
	r.db.Close()
}

//insert data...
func (r *Store) InsertData(ctx context.Context, order model.Order) (int, error) {

	id, err := r.InsertDataOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	idDelivery, err := r.InsertDataDelivery(ctx, order)
	if err != nil {
		return 0, err
	}

	if err = r.InsertDataOrdersDelivery(ctx, order, idDelivery); err != nil {
		return 0, err
	}

	if err = r.InsertDataItem(ctx, order); err != nil {
		return 0, err
	}

	if err = r.InsertDataPayments(ctx, order); err != nil {
		return 0, err
	}

	return id, nil
}

//insert into order...
func (r *Store) InsertDataOrder(ctx context.Context, order model.Order) (int, error) {
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
func (r *Store) InsertDataDelivery(ctx context.Context, order model.Order) (int, error) {
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
func (r *Store) InsertDataOrdersDelivery(ctx context.Context, order model.Order, id int) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO orders_delivery(order_id,delivery_id) "+
		"VALUES ($1,$2)",
		order.OrderUid,
		id); err != nil {
		return fmt.Errorf("InsertDataOrderDelivery: %v", err)
	}
	return nil
}

//insert into items...
func (r *Store) InsertDataItem(ctx context.Context, order model.Order) error {
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
func (r *Store) InsertDataPayments(ctx context.Context, order model.Order) error {
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

//get data...
func (r *Store) GetDataById(ctx context.Context, id int) (*model.Order, error) {
	var nM model.Order

	//get order data...
	dataOrder, err := r.GetDataOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	//insert data to model...
	nM = *dataOrder

	//get order data...
	deliveryData, err := r.GetDataDelivery(ctx, nM.OrderUid)
	if err != nil {
		return nil, err
	}
	//insert data to model...
	nM.Delivery = *deliveryData

	//get payments data...
	paymentData, err := r.GetDataPayment(ctx, nM.OrderUid)
	if err != nil {
		return nil, err
	}
	//insert data to model...
	nM.Payment = *paymentData

	//get items data...
	itemsData, err := r.GetDataItems(ctx, nM.TrackNumber)
	if err != nil {
		return nil, err
	}
	//insert data to model...
	nM.Items = itemsData

	return &nM, nil
}

//get data items...
func (r *Store) GetDataOrder(ctx context.Context, id int) (*model.Order, error) {
	var order model.Order
	if err := r.db.QueryRowContext(ctx, "SELECT o.order_uid,o.track_number,o.entry,o.locale, o.internal_signature,o.customer_id,o.delivery_service,o.shardkey,o.sm_id,o.date_created,o.oof_shard "+
		"FROM orders o WHERE id = $1", id).Scan(
		&order.OrderUid,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.DateCreated,
		&order.OofShard,
	); err != nil {
		return nil, fmt.Errorf("GetDataOrder: %v", err)
	}
	return &order, nil
}

//get delivery data...
func (r *Store) GetDataDelivery(ctx context.Context, orderUid string) (delivery *model.Delivery, err error) {
	delivery = &model.Delivery{}
	idDelivery, err := r.GetDataOrdersDelivery(ctx, orderUid)
	if err != nil {
		return nil, err
	}
	if err = r.db.QueryRowContext(ctx, "SELECT d.name,d.phone,d.zip,d.city,d.address,d.region,d.email "+
		"FROM delivery d WHERE id = $1", idDelivery).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	); err != nil {
		return nil, fmt.Errorf("GetDataDelivery: %v", err)
	}
	return delivery, nil
}

//get data orders_delivery...
func (r *Store) GetDataOrdersDelivery(ctx context.Context, orderUid string) (idDelivery int, err error) {
	if err = r.db.QueryRowContext(ctx, "SELECT od.delivery_id FROM orders_delivery od WHERE order_id = $1",
		orderUid).Scan(&idDelivery); err != nil {
		return 0, fmt.Errorf("GetDataOrdersDelivery: %v", err)
	}
	return
}

//get data items...
func (r *Store) GetDataItems(ctx context.Context, trackNumber string) (items []model.Item, err error) {
	rows, err := r.db.QueryContext(ctx, "SELECT chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status FROM item i WHERE track_number = $1",
		trackNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.Item
		if err = rows.Scan(
			&temp.ChrtId,
			&temp.TrackNumber,
			&temp.Price,
			&temp.Rid,
			&temp.Name,
			&temp.Sale,
			&temp.Size,
			&temp.TotalPrice,
			&temp.NmId,
			&temp.Brand,
			&temp.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, temp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return
}

//get data payment...
func (r *Store) GetDataPayment(ctx context.Context, transaction string) (payment *model.Payment, err error) {
	payment = &model.Payment{}
	if err = r.db.QueryRowContext(ctx, "SELECT p.transaction,p.request_id,p.currency,p.provider,p.amount,p.payment_dt,p.bank,p.delivery_cost,p.goods_total,p.custom_free "+
		"FROM payment p WHERE p.transaction = $1", transaction,
	).Scan(
		&payment.Transaction,
		&payment.RequestId,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFree,
	); err != nil {
		return nil, fmt.Errorf("GetDataPayment: %v", err)
	}
	return
}

//get last id order from db, for cash restore...
func (r *Store) GetLastDataId(ctx context.Context) (int, error) {
	//SELECT id FROM orders ORDER BY id DESC LIMIT 1
	//SELECT MAX(id) FROM orders
	var id int
	if err := r.db.QueryRowContext(ctx, "SELECT id FROM orders ORDER BY id DESC LIMIT 1").Scan(&id); err != nil {
		return 0, fmt.Errorf("GetLastDataId: %v", err)
	}
	return id, nil
}

//insert incorrect data to table...
func (r *Store) InsertIncorrectData(ctx context.Context, message interface{}) error {
	if err := r.db.QueryRowContext(ctx, "INSERT INTO incorrect_messages(message)"+
		"VALUES ($1)",
		message); err != nil {
		return fmt.Errorf("InsertIncorrectData: %v", err)
	}
	return nil
	return nil
}

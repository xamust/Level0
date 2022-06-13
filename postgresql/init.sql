DROP TABLE if exists orders,delivery,item,payment,orders_delivery;

CREATE TABLE orders(
                       id  serial primary key,
                       order_uid varchar(32) unique,
                       track_number varchar(32) unique,
                       entry varchar(32),
                       locale  varchar(32),
                       internal_signature varchar(32),
                       customer_id varchar(32),
                       delivery_service varchar(32),
                       shardkey integer,
                       sm_id integer,
                       date_created varchar(32),
                       oof_shard integer
);

CREATE TABLE delivery (
                        id  serial primary key,
                        name varchar(255),
                        phone varchar(32),
                        zip integer,
                        city varchar(255),
                        address varchar(255),
                        region varchar(255),
                        email varchar(255)
);

CREATE TABLE payment (
                         id  serial primary key,
                         transaction varchar(32) references orders(order_uid),
                         request_id varchar(32),
                         currency varchar(3),
                         provider varchar(32),
                         amount integer,
                         payment_dt integer,
                         bank varchar(32),
                         delivery_cost integer,
                         goods_total integer,
                         custom_free integer
);

CREATE TABLE item (
                      id serial primary key,
                      chrt_id integer,
                      track_number varchar(32) references orders(track_number),
                      price integer,
                      rid varchar(255),
                      name varchar(32),
                      sale integer,
                      size varchar(32),
                      total_price integer,
                      nm_id integer,
                      brand varchar(255),
                      status integer
);

CREATE TABLE orders_delivery (
                    order_id varchar(32) references orders(order_uid),
                    delivery_id integer references delivery(id)
);
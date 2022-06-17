package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

var testModel = &Order{
	OrderUid:    "a563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACKA",
	Entry:       "WBIL",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "a563feb7b2b84b6test",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFree:   0,
	},
	Items: []Item{{
		ChrtId:      9934930,
		TrackNumber: "WBILMTESTTRACKA",
		Price:       453,
		Rid:         "ab4219087a764ae0btest",
		Name:        "Mascaras",
		Sale:        30,
		Size:        "0",
		TotalPrice:  317,
		NmId:        2389212,
		Brand:       "Vivienne Sabo",
		Status:      202,
	},
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerId:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmId:              99,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
}
var testModel2 = &Order{
	OrderUid:    "b563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACKB",
	Entry:       "WBIL",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "b563feb7b2b84b6test",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFree:   0,
	},
	Items: []Item{{
		ChrtId:      9934930,
		TrackNumber: "WBILMTESTTRACKB",
		Price:       453,
		Rid:         "ab4219087a764ae0btest",
		Name:        "Mascaras",
		Sale:        30,
		Size:        "0",
		TotalPrice:  317,
		NmId:        2389212,
		Brand:       "Vivienne Sabo",
		Status:      202,
	}, {
		ChrtId:      9934930,
		TrackNumber: "WBILMTESTTRACKB",
		Price:       453,
		Rid:         "ab4219087a764ae0btest",
		Name:        "Mascaras",
		Sale:        25,
		Size:        "13",
		TotalPrice:  317,
		NmId:        2389212,
		Brand:       "Красный скороход",
		Status:      202,
	},
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerId:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmId:              99,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
}
var testModel3 = &Order{
	OrderUid:    "c563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACKC",
	Entry:       "WBIL",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "c563feb7b2b84b6test",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFree:   0,
	},
	Items: []Item{{
		ChrtId:      9934930,
		TrackNumber: "WBILMTESTTRACKC",
		Price:       453,
		Rid:         "ab4219087a764ae0btest",
		Name:        "Mascaras",
		Sale:        30,
		Size:        "0",
		TotalPrice:  317,
		NmId:        2389212,
		Brand:       "Vivienne Sabo",
		Status:      202,
	},
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerId:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmId:              99,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
}
var testModel4 = &Order{
	OrderUid:    "d563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACKD",
	Entry:       "WBIL",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "d563feb7b2b84b6test",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFree:   0,
	},
	Items: []Item{{
		ChrtId:      9934930,
		TrackNumber: "WBILMTESTTRACKD",
		Price:       453,
		Rid:         "ab4219087a764ae0btest",
		Name:        "Mascaras",
		Sale:        30,
		Size:        "0",
		TotalPrice:  317,
		NmId:        2389212,
		Brand:       "Vivienne Sabo",
		Status:      202,
	},
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerId:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmId:              99,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
}

func main() {

	nc, err := nats.Connect("nats://127.0.0.1:4223")
	if err != nil {
		log.Fatal(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	ec.Publish("test", testModel)
	ec.Publish("test", testModel2)
	ec.Publish("test", "testetstetstetstetst")
	ec.Publish("test", testModel3)
	ec.Publish("test", "{{{test}}}")
	ec.Publish("test", testModel4)
	defer nc.Close()
	defer ec.Close()
}

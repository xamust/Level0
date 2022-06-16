package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
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

func main() {

	//pg, err := postgresql.New("dbname=test_db user=testUser password=password host=localhost port=5432 sslmode=disable")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println("Success conn to DB")
	//ctx := context.Background()
	//
	//nM, err := pg.GetDataById(ctx, 1)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println(testModel2)
	//fmt.Println(*nM)
	//
	//if err = pg.InsertData(ctx, testModel2); err != nil {
	//	log.Fatalln(err)
	//}

	nc, err := nats.Connect("nats://127.0.0.1:4223")
	if err != nil {
		log.Fatal(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	ec.Publish("test", testModel2)
	//defer nc.Close()
	//defer ec.Close()
	//fmt.Println(nc.IsConnected())
	//sendCh := make(chan *Order)
	//ec.BindSendChan("test", sendCh)
	//time.Sleep(time.Second * 5)
	//sendCh <- testModel
	//time.Sleep(time.Second * 5)
	//sendCh <- testModel2

	// Simple Async Subscriber
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	fmt.Println(nc.NumSubscriptions())

	//client()
	// Simple Publisher

	//nc.Publish("test", []byte("Hello World"))

	time.Sleep(time.Second * 1)

	/*
		// Connect to a server
		//nc, _ := natsapp.Connect(natsapp.DefaultURL)
		nc, err := natsapp.Connect("natsapp://127.0.0.1:4223")
		if err != nil {
			log.Fatalln(err)
		}
		// Simple Publisher
		nc.Publish("foo", []byte("Hello World"))
		/*
			// Simple Async Subscriber
			nc.Subscribe("foo", func(m *natsapp.Msg) {
				fmt.Printf("Received a message: %s\n", string(m.Data))
			})

			// Responding to a request message
			nc.Subscribe("request", func(m *natsapp.Msg) {
				m.Respond([]byte("answer is 42"))
			})

			// Simple Sync Subscriber
			sub, err := nc.SubscribeSync("foo")
			m, err := sub.NextMsg(natsapp.DefaultDrainTimeout)

			// Channel Subscriber
			ch := make(chan *natsapp.Msg, 64)
			sub, err := nc.ChanSubscribe("foo", ch)
			msg := <- ch

			// Unsubscribe
			sub.Unsubscribe()

			// Drain
			sub.Drain()

			// Requests
			msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)

			// Replies
			nc.Subscribe("help", func(m *natsapp.Msg) {
				nc.Publish(m.Reply, []byte("I can help!"))
			})

			// Drain connection (Preferred for responders)
			// Close() not needed if this is called.
			nc.Drain()

			// Close connection
			nc.Close()
	*/
}

func client() {
	nc, err := nats.Connect("natsapp://127.0.0.1:8223", nats.Name("API PublishBytes Example"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	if err := nc.Publish("updates", []byte("All is Well")); err != nil {
		log.Fatal(err)
	}
}

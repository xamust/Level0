package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	//client()

	// Connect to a server
	//nc, _ := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect("nats://127.0.0.1:4223")
	if err != nil {
		log.Fatalln(err)
	}
	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))

	// Simple Async Subscriber
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	// Responding to a request message
	nc.Subscribe("request", func(m *nats.Msg) {
		m.Respond([]byte("answer is 42"))
	})

	// Simple Sync Subscriber
	sub, err := nc.SubscribeSync("foo")
	m, err := sub.NextMsg(nats.DefaultDrainTimeout)
	fmt.Println(&m.Data)
	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	sub, err = nc.ChanSubscribe("foo", ch)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sub.IsValid())

	msg := <-ch
	fmt.Println(msg)
	// Unsubscribe
	sub.Unsubscribe()

	// Drain
	sub.Drain()

	// Requests
	msg, err = nc.Request("help", []byte("help me"), 10*time.Millisecond)
	fmt.Println(msg)
	// Replies
	nc.Subscribe("help", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("I can help!"))
	})

	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	nc.Drain()

	// Close connection
	nc.Close()
}

func client() {
	nc, err := nats.Connect("nats://127.0.0.1:8223", nats.Name("API PublishBytes Example"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	if err := nc.Publish("updates", []byte("All is Well")); err != nil {
		log.Fatal(err)
	}
}

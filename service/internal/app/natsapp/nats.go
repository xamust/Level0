package natsapp

import (
	"github.com/nats-io/nats.go"
)

type NatsService struct {
	config       *Config
	NatsConn     *nats.Conn
	subscription *nats.Subscription
}

func New(config *Config) *NatsService {
	return &NatsService{
		config: config,
	}
}

func (n *NatsService) Connect() error {
	nc, err := nats.Connect(n.config.NatsAddr)
	if err != nil {
		return err
	}
	n.NatsConn = nc
	return nil
}

func (n *NatsService) Close() {
	//close...
	if n.NatsConn != nil {
		n.NatsConn.Close()
	}
	//unsubs...
	if n.subscription != nil {
		n.Unsubscribe()
	}
}

//unsubs...
func (n *NatsService) Unsubscribe() error {
	return n.subscription.Unsubscribe()
}

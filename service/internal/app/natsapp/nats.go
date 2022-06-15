package natsapp

import (
	"github.com/nats-io/nats.go"
	"service/internal/app/model"
)

type NatsService struct {
	config       *Config
	nats         *nats.Conn
	subscription *nats.Subscription
	receivedMsg  chan model.Order
}

func New(config *Config) *NatsService {
	return &NatsService{
		config: config,
	}
}

func (n *NatsService) Connect() (*NatsService, error) {
	nc, err := nats.Connect(n.config.NatsAddr)
	if err != nil {
		return nil, err
	}
	n.nats = nc
	return &NatsService{nats: nc}, nil
}

func (n *NatsService) Close() {
	//close...
	if n.nats != nil {
		n.nats.Close()
	}
	//unsubs...
	if n.subscription != nil {
		n.Unsubscribe()
	}
	//close ch...
	close(n.receivedMsg)
}

func (n *NatsService) ChannelSubscribe() (chan *nats.Msg, error) {
	ch := make(chan *nats.Msg, 64)
	sub, err := n.nats.ChanSubscribe(n.config.NatsSubs, ch)
	if err != nil {
		return nil, err
	}
	n.subscription = sub
	return ch, nil
}

func (n *NatsService) JSONEncodedConn() (chan *model.Order, error) {
	ec, err := nats.NewEncodedConn(n.nats, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	defer ec.Close()
	recvCh := make(chan *model.Order)
	ec.BindSendChan(n.config.NatsSubs, recvCh)
	return recvCh, nil
}

func (n *NatsService) Unsubscribe() error {
	return n.subscription.Unsubscribe()
}

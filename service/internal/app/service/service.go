package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"service/internal/app/model"
	"service/internal/app/natsapp"
	"service/internal/app/store"
	"strconv"
)

type Service struct {
	config *Config
	logger *logrus.Logger
	store  *store.Store
	nats   *natsapp.NatsService
	mux    *mux.Router
}

func NewService(config *Config) *Service {
	return &Service{
		config: config,
		logger: logrus.New(),
		mux:    mux.NewRouter(),
	}
}

func (s *Service) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("Logger ready...")

	s.configureRouter()
	s.logger.Info("Router ready...")

	if err := s.configureNats(); err != nil {
		s.logger.Fatalln(err)
	}
	s.logger.Info("Nats ready...")

	_, err := s.nats.ChannelSubscribe()
	if err != nil {
		s.logger.Fatalln(err)
	}
	s.logger.Infof("Nats subs on %v...", s.config.NatsApp.NatsSubs)

	recvCh, err := s.nats.JSONEncodedConn()
	if err != nil {
		s.logger.Fatalln(err)
	}

	go func(chan *model.Order) {
		msg := <-recvCh
		s.logger.Info(msg)
		//	msg := <-ch
		//	s.logger.Info(msg.Data)
	}(recvCh)

	if err = s.configureStore(); err != nil {
		s.logger.Fatalln(err)
	}
	s.logger.Info("Store ready...")

	s.logger.Info(fmt.Sprintf("Starting server (bind on %v)...", s.config.BindAddr))
	return http.ListenAndServe(s.config.BindAddr, s.mux)
}

func (s *Service) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *Service) configureStore() error {
	newStore := store.New(s.config.Store)
	if err := newStore.Open(); err != nil {
		return err
	}
	s.store = newStore
	return nil
}

func (s *Service) configureNats() error {
	newNats := natsapp.New(s.config.NatsApp)
	if err := newNats.Connect(); err != nil {
		return err
	}
	s.nats = newNats

	return nil
}

func (s *Service) configureRouter() {
	s.mux.HandleFunc("/{id:[0-9]+}", s.GetById)
}

func (s *Service) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	s.logger.Print(id)
}

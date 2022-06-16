package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
	"service/internal/app/cashdata"
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
	cash   *cashdata.CashData
	mux    *mux.Router
	ctx    context.Context
}

func NewService(config *Config) *Service {
	return &Service{
		config: config,
		logger: logrus.New(),
		mux:    mux.NewRouter(),
	}
}

//start service...
func (s *Service) Start() error {
	var correct, incorrect int
	//init context...
	s.ctx = context.Background()

	//config logger...
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("Logger ready...")

	//config router (gorilla/mux)...
	s.configureRouter()
	s.logger.Info("Router ready...")

	//config db...
	if err := s.configureStore(); err != nil {
		s.logger.Fatalln(err)
	}
	s.logger.Info("Store ready...")

	//config and restore cash...
	s.configureCash()
	if err := s.cash.RestoreCash(); err != nil {
		s.logger.Error("Cash not restored! ", err)
	} else {
		s.logger.Infof("Cash ready. %d element restore from db", len(s.cash.CashMass))
	}

	//config and connect to Nats...
	if err := s.configureNats(); err != nil {
		s.logger.Fatalln(err)
	}
	s.logger.Info("Conn to Nats ready...")

	//subscribe to channel Nats, from config and read msg...
	s.nats.NatsConn.Subscribe(s.config.NatsApp.NatsSubs, func(m *nats.Msg) {
		var mod model.Order

		//unmarshall data from msg...
		if err := json.Unmarshal(m.Data, &mod); err != nil {
			//counter incorrect msg, for logger info...
			incorrect++
			s.logger.Errorf("Unmarshal: %v", err)
			//if data can't unmarshall, write to special table, not to lose, for debug...
			s.store.InsertIncorrectData(s.ctx, string(m.Data))
		} else {
			//if data correctly unmarshalled, writes to tables in db...
			_, err = s.store.InsertData(s.ctx, mod)
			if err != nil {
				//not error, because service always restart....
				s.logger.Errorf("Insert data: %v", err)
			} else {
				//insert data to cash...
				//counter correct msg, for logger info...
				correct++
				//append data to cash slice...
				s.cash.SetCash(mod)
				s.logger.Infof("Correct messages: %d, incorrect: %d", correct, incorrect)
			}
		}
	})
	s.logger.Infof("Nats subs on %v...", s.config.NatsApp.NatsSubs)

	s.logger.Info(fmt.Sprintf("Starting server (bind on %v)...", s.config.BindAddr))
	//start web server...
	return http.ListenAndServe(s.config.BindAddr, s.mux)
}

//config logger...
func (s *Service) configureLogger() error {
	//get level for logrus from configs...
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	//set level for logrus...
	s.logger.SetLevel(level)
	return nil
}

//config db...
func (s *Service) configureStore() error {
	//get config for db from configs...
	newStore := store.New(s.config.Store)
	//open new db...
	if err := newStore.Open(); err != nil {
		return err
	}
	//set store to struct...
	s.store = newStore
	return nil
}

func (s *Service) configureNats() error {
	//get config for nats from configs...
	newNats := natsapp.New(s.config.NatsApp)
	//new connect to nats...
	if err := newNats.Connect(); err != nil {
		return err
	}
	//set nats to struct...
	s.nats = newNats
	return nil
}

//config cash...
func (s *Service) configureCash() {
	//make new cash slice
	cashSlice := make([]model.Order, 0)
	//init new cash ...
	cash := cashdata.NewCash(s.config.CashData, s.ctx, cashSlice, s.store)
	//set cash to struct..
	s.cash = cash
}

func (s *Service) configureRouter() {
	//register handle on router... (:8080/[0-9]+)
	//with regexp...
	s.mux.HandleFunc("/{id:[0-9]+}", s.GetById)
}

//handle...
func (s *Service) GetById(w http.ResponseWriter, r *http.Request) {
	//check correct http method
	if r.Method == http.MethodGet {
		//read regexp id, from handle [0-9]+...
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			//return err page...
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		//get order from db, by id...
		order, err := s.store.GetDataById(s.ctx, id)
		if err != nil {
			s.logger.Error(order)
		}
		//marshalling order data by id for page...
		result, err := json.Marshal(order)
		//return page with data by id...
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

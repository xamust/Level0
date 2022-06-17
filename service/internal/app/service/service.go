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

	//config db and cash...
	go func() {
		//config db...
		if err := s.configureStore(); err != nil {
			s.logger.Fatalln(err)
		}
		s.logger.Info("Store ready...")

		//config and restore cash...
		s.configureCash()
		//getLastId from db, for restore cash...
		id, err := s.store.GetLastDataId(s.ctx)
		if err != nil {
			s.logger.Error("Cash not restored! ", err)
		}
		//restore cash...
		if err = s.cash.RestoreCash(id); err != nil {
			//error handling...
			s.logger.Error("Cash not restored! ", err)
		}
	}()

	//config Nats...
	go func() {
		if err := s.configureNats(); err != nil {
			s.logger.Fatalln(err)
		}
		s.logger.Info("Conn to Nats ready...")

		//subscribe to channel Nats, from config and read msg...
		go s.nats.NatsConn.Subscribe(s.config.NatsApp.NatsSubs, func(m *nats.Msg) {
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
				id, err := s.store.InsertData(s.ctx, mod)
				if err != nil {
					//not error, because service always restart....
					s.logger.Errorf("Insert data: %v", err)
				} else {
					//insert data to cash...
					//counter correct msg, for logger info...
					correct++
					//append data to cash slice...
					s.cash.SetCash(id, mod)
					s.logger.Infof("Correct messages: %d, incorrect: %d", correct, incorrect)
				}
			}
		})
		s.logger.Infof("Nats subs on %v...", s.config.NatsApp.NatsSubs)
	}()

	s.logger.Info(fmt.Sprintf("Starting service (bind on %v)...", s.config.BindAddr))
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

//config Nats...
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
	cashSlice := make(map[int]model.Order)
	//init new cash ...
	cash := cashdata.NewCash(s.ctx, cashSlice, s.store)
	//set cash to struct..
	s.cash = cash
}

//config router...
func (s *Service) configureRouter() {
	//register handle on router... (:8080/[0-9]+)
	//with regexp...
	s.mux.HandleFunc("/{id:[0-9]+}", s.GetById).Methods(http.MethodGet, http.MethodOptions)
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
		//order, err := s.store.GetDataById(s.ctx, id)
		//if err != nil {
		//	s.logger.Error(order)
		//}

		//get order from cash, by id...
		order, err := s.cash.GetById(id)
		if err != nil {
			//return err page...
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		//marshalling order data by id for page...
		result, err := json.Marshal(order)
		if err != nil {
			//return err page...
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		//return page with data by id...
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

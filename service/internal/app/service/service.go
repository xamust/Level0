package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"service/internal/app/store"
	"strconv"
)

type Service struct {
	config *Config
	logger *logrus.Logger
	store  *store.Store
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

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

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

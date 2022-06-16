package service

import "service/internal/app/model"

type mainCashMap map[int]*model.Order

func SetCash(order model.Order) {

}

func GetCash() mainCashMap {
	return mainCashMap{}
}

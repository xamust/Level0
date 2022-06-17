package cashdata

import (
	"service/internal/app/model"
)

type Repository interface {
	SetCash(id int, order model.Order)
	RestoreCash(id int) error
	GetById(id int) (*model.Order, error)
}

var impl Repository

func SetCash(id int, order model.Order) {

}
func RestoreCash(id int) error {
	return impl.RestoreCash(id)
}
func GetById(id int) (*model.Order, error) {
	return impl.GetById(id)
}

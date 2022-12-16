package storage

import (
	"car_rental/genprotos/rental"
)

type StorageI interface {
	CreateRental(id string, req *rental.CreateRentalRequest) error
	GetRentalByID(id string) (*rental.GetRentalByIDResponse, error)
	GetRentalList(offset, limit int, search string) (*rental.GetRentalListResponse, error)
	UpdateRental(entity *rental.UpdateRentalRequest) error
	DeleteRental(id string) error
	GetRentalsByUserId(offset, limit int, search, customer_id string) (*rental.GetRentalsByUserIdResponse, error)
}

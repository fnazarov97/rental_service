package rental

import (
	"car_rental/genprotos/rental"
	"car_rental/storage"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RentalService is a struct that implements the server interface
type RentalService struct {
	Stg storage.StorageI
	rental.UnimplementedRentalServiceServer
}

// CreateRental ...
func (a *RentalService) CreateRental(c context.Context, req *rental.CreateRentalRequest) (*rental.Rental, error) {
	id := uuid.New()
	err := a.Stg.CreateRental(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.CreateRental: %s", err.Error())
	}
	res, err := a.Stg.GetRentalByID(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.Stg.GetRentalByID: %s", err.Error())
	}
	return &rental.Rental{
		RentalId:   res.RentalId,
		CarId:      res.Car.CarId,
		CustomerId: res.Customer.Id,
		StartDate:  res.StartDate,
		EndDate:    res.EndDate,
		Payment:    res.Payment,
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}, nil
}

// GetRentalByID ...
func (a *RentalService) GetRentalByID(c context.Context, req *rental.GetRentalByIDRequest) (*rental.GetRentalByIDResponse, error) {
	res, err := a.Stg.GetRentalByID(req.RentalId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.GetRentalByID: %s", err.Error())
	}
	return res, nil
}

// GetRentalList ...
func (a *RentalService) GetRentalList(c context.Context, req *rental.GetRentalListRequest) (*rental.GetRentalListResponse, error) {
	res, err := a.Stg.GetRentalList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.GetRentalList: %s", err.Error())
	}
	return res, nil
}

// UpdateRental ...
func (a *RentalService) UpdateRental(c context.Context, req *rental.UpdateRentalRequest) (*rental.Rental, error) {
	err := a.Stg.UpdateRental(req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.UpdateRental: %s", err.Error())
	}
	res, e := a.Stg.GetRentalByID(req.RentalId)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.UpdateRental: %s", e.Error())
	}
	return &rental.Rental{
		RentalId:   res.RentalId,
		CarId:      res.Car.CarId,
		CustomerId: res.Customer.Id,
		StartDate:  res.StartDate,
		EndDate:    res.EndDate,
		Payment:    res.Payment,
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}, nil
}

// DeleteRental ...
func (a *RentalService) DeleteRental(c context.Context, req *rental.DeleteRentalRequest) (*rental.Rental, error) {
	res, e := a.Stg.GetRentalByID(req.RentalId)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.UpdateRental: %s", e.Error())
	}
	err := a.Stg.DeleteRental(req.RentalId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.DeleteRental: %s", err.Error())
	}

	return &rental.Rental{
		RentalId:   res.RentalId,
		CarId:      res.Car.CarId,
		CustomerId: res.Customer.Id,
		StartDate:  res.StartDate,
		EndDate:    res.EndDate,
		Payment:    res.Payment,
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}, nil
}

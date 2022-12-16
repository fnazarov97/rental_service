package rental

import (
	"car_rental/clients"

	// "car_rental/genprotos/authorization"
	"car_rental/genprotos/authorization"
	"car_rental/genprotos/car"
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
	Services clients.ServiceManageI
	// grpcClients *clients.GrpcClients
}

func NewRentalService(s clients.GrpcClients, stg storage.StorageI) *RentalService {
	return &RentalService{
		Stg:      stg,
		Services: &s,
	}
}

// CreateRental ...
func (a *RentalService) CreateRental(c context.Context, req *rental.CreateRentalRequest) (*rental.Rental, error) {
	id := uuid.New()

	car, e := a.Services.CarService().GetCarByID(c, &car.GetCarByIDRequest{
		Id: req.CarId,
	})
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.Services.CarService().GetCarByID: %s", e.Error())
	}
	customer, e := a.Services.AuthService().GetUserByID(c, &authorization.GetUserByIDRequest{
		Id: req.CustomerId,
	})

	if e != nil {
		return nil, status.Errorf(codes.Unauthenticated, "a.Services.AuthService().GetUserByID: %s", e.Error())
	}

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
		CarId:      car.CarId,
		CustomerId: customer.Id,
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
	car, e := a.Services.CarService().GetCarByID(c, &car.GetCarByIDRequest{
		Id: res.CarId,
	})
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.Services.CarService().GetCarByID: %s", e.Error())
	}
	res.Car.CarId = car.CarId
	res.Car.Model = car.Model
	res.Car.Color = car.Color
	res.Car.Year = car.Year
	res.Car.Mileage = car.Mileage
	res.Car.BrandId = car.Brand.BrandId
	res.CreatedAt = car.CreatedAt
	res.UpdatedAt = car.UpdatedAt

	customer, err := a.Services.AuthService().GetUserByID(c, &authorization.GetUserByIDRequest{
		Id: res.CustomerId,
	})
	res.Customer.Id = customer.Id
	res.Customer.Fname = customer.Fname
	res.Customer.Lname = customer.Lname
	res.Customer.Username = customer.Username
	res.Customer.Password = customer.Password
	res.Customer.UserType = customer.UserType
	res.Customer.Address = customer.Address
	res.Customer.Phone = customer.Phone
	res.CreatedAt = customer.CreatedAt
	res.UpdatedAt = customer.UpdatedAt
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.grpcClients.Authorization.GetUserByID: %s", err.Error())
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
	_, err := a.Services.AuthService().GetUserByID(c, &authorization.GetUserByIDRequest{
		Id: req.CustomerId,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Services.AuthService().GetUserByID: %s", err.Error())
	}

	_, err = a.Services.CarService().GetCarByID(c, &car.GetCarByIDRequest{
		Id: req.CarId,
	})

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Services.CarService().GetCarByID: %s", err.Error())
	}

	err = a.Stg.UpdateRental(req)
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

// GetRentalsByUserId ...
func (a *RentalService) GetRentalsByUserId(c context.Context, req *rental.GetRentalsByUserIdRequest) (*rental.GetRentalsByUserIdResponse, error) {

	res, err := a.Stg.GetRentalsByUserId(int(req.Offset), int(req.Limit), req.Search, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.GetRentalByID: %s", err.Error())
	}

	return res, nil
}

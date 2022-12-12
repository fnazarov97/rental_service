package clients

import (
	"car_rental/config"
	"car_rental/genprotos/authorization"
	"car_rental/genprotos/car"
	"car_rental/genprotos/rental"

	"google.golang.org/grpc"
)

type ServiceManageI interface {
	CarService() car.CarServiceClient
	AuthService() authorization.AuthServiceClient
}

type GrpcClients struct {
	Car           car.CarServiceClient
	Rental        rental.RentalServiceClient
	Authorization authorization.AuthServiceClient
	conns         []*grpc.ClientConn
}

func NewGrpcClients(cfg config.Config) (*GrpcClients, error) {
	connCar, err := grpc.Dial(cfg.CarServiceGrpcHost+cfg.CarServiceGrpcPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	car := car.NewCarServiceClient(connCar)

	connRental, err := grpc.Dial(cfg.RentalServiceGrpcHost+cfg.RentalServiceGrpcPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	rental := rental.NewRentalServiceClient(connRental)

	connAuthorization, err := grpc.Dial(cfg.AuthorizationServiceGrpcHost+cfg.AuthorizationServiceGrpcPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	authorization := authorization.NewAuthServiceClient(connAuthorization)
	conns := make([]*grpc.ClientConn, 0)
	return &GrpcClients{
		Car:           car,
		Rental:        rental,
		Authorization: authorization,
		conns:         append(conns, connCar, connRental, connAuthorization),
	}, nil
}

func (g *GrpcClients) CarService() car.CarServiceClient {
	return g.Car
}

func (g *GrpcClients) AuthService() authorization.AuthServiceClient {
	return g.Authorization
}

func (c *GrpcClients) Close() {
	for _, v := range c.conns {
		v.Close()
	}
}

package postgres

import (
	"car_rental/genprotos/rental"
	"errors"
	"time"
)

//"rental_id", "car_id", "customer_id", "start_date", "end_date",
//"payment", "created_at", "updated_at", "deleted_at"
// CreateRental ...
func (p Postgres) CreateRental(id string, req *rental.CreateRentalRequest) error {

	_, err := p.DB.Exec(`
	INSERT INTO
	 "rentals"("rental_id", "car_id", "customer_id", "start_date", "end_date", "payment", "created_at") 
	VALUES($1, $2, $3, $4, $5, $6, now())
	`, id, req.CarId, req.CustomerId, req.StartDate, req.EndDate, req.Payment)
	if err != nil {
		return err
	}
	return nil
}

// GetRentalByID ...
func (p Postgres) GetRentalByID(id string) (*rental.GetRentalByIDResponse, error) {
	//---bu yerga car va user malumotlari kerak!--- will bw update soon
	res := &rental.GetRentalByIDResponse{
		Car:      &rental.GetRentalByIDResponse_Car{},
		Customer: &rental.GetRentalByIDResponse_User{},
	}
	var deletedAt *time.Time
	var updatedAt *string
	err := p.DB.QueryRow(`SELECT 
	"rental_id", "car_id", "customer_id", "start_date", "end_date",
	"payment", "created_at", "updated_at"
    FROM "rentals" WHERE "rental_id" = $1`, id).Scan(
		&res.RentalId,
		&res.CarId,
		&res.CustomerId,
		&res.StartDate,
		&res.EndDate,
		&res.Payment,
		&res.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return res, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if deletedAt != nil {
		return res, errors.New("rental not found")
	}

	return res, err
}

// GetRentalList ...
func (p Postgres) GetRentalList(offset, limit int, search string) (*rental.GetRentalListResponse, error) {
	resp := &rental.GetRentalListResponse{
		Rentals: make([]*rental.Rental, 0),
	}

	rows, err := p.DB.Queryx(`
	SELECT
	"rental_id", "car_id", "customer_id", "start_date", "end_date", "payment", "created_at", "updated_at"
	FROM 
		"rentals" WHERE deleted_at IS NULL AND (start_date ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		a := &rental.Rental{}
		var updatedAt *string
		err := rows.Scan(
			&a.RentalId, &a.CarId, &a.CustomerId, &a.StartDate, &a.EndDate, &a.Payment, &a.CreatedAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}

		resp.Rentals = append(resp.Rentals, a)
	}

	return resp, err
}

// UpdateRental ...
func (p Postgres) UpdateRental(entity *rental.UpdateRentalRequest) error {

	res, err := p.DB.NamedExec(`
	UPDATE 
		"rentals"  
	SET 
		"car_id"=:ca, "customer_id"=:cu, "start_date"=:s, "end_date"=:e, "payment"=:p, updated_at=now() WHERE deleted_at IS NULL AND rental_id=:id`, map[string]interface{}{
		"id": entity.RentalId,
		"ca": entity.CarId,
		"cu": entity.CustomerId,
		"s":  entity.StartDate,
		"e":  entity.EndDate,
		"p":  entity.Payment,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("rental not found")
}

// DeleteRental ...
func (p Postgres) DeleteRental(id string) error {
	res, err := p.DB.Exec("UPDATE rentals SET deleted_at=now() WHERE rental_id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("rental not found")
}

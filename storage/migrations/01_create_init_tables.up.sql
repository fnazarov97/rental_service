BEGIN;

CREATE TABLE IF NOT EXISTS "rentals" (
	"rental_id" CHAR(36) NOT NULL PRIMARY KEY,
	"car_id" CHAR(36) NOT NULL,
	"customer_id" CHAR(36) NOT NULL,
	"start_date" VARCHAR(50) NOT NULL,
	"end_date" VARCHAR(50) NOT NULL,
	"payment" VARCHAR(50) NOT NULL,
	"created_at" TIMESTAMP DEFAULT now(),
	"updated_at" TIMESTAMP,
	"deleted_at" TIMESTAMP
);

COMMIT;
package car

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRepository struct {
	db *pgxpool.Pool
}

// NewPgxRepository creates a new postgresql repository
func NewPgxRepository(db *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{
		db: db,
	}
}

// /////////////////////////////////////// start migration (create table of cars if it doesn't exist)
func (r *PgxRepository) Migrate(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS cars(
        id SERIAL PRIMARY KEY,
        brand TEXT NOT NULL,
        model TEXT NOT NULL,
		color TEXT NOT NULL,
        price FLOAT NOT NULL
    );
    `
	_, err := r.db.Exec(ctx, query)
	return err
}

////////////////////////////////////////////// create a car data

func (r *PgxRepository) CreateCar(ctx context.Context, c Car) (*Car, error) {
	var id int64
	err := r.db.QueryRow(ctx, "INSERT INTO cars(brand, model, color, price) values($1, $2, $3, $4) RETURNING id", c.Brand, c.Model, c.Color, c.Price).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	c.ID = id

	return &c, nil
}

// //////////////////////////////////////////// get all cars data
func (r *PgxRepository) AllCars(ctx context.Context) ([]Car, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var allCars []Car

	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Color, &car.Price); err != nil {
			return nil, err
		}
		allCars = append(allCars, car)
	}
	return allCars, nil
}

// //////////////////////////////////////////// get a single car data by ID
func (r *PgxRepository) GetCarById(ctx context.Context, id int64) (*Car, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM cars WHERE id = $1", id)

	var car Car
	if err := row.Scan(&car.ID, &car.Brand, &car.Model, &car.Color, &car.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &car, nil
}

// //////////////////////////////////////////// update a single car data
func (r *PgxRepository) UpdateCar(ctx context.Context, id int64, updated Car) (*Car, error) {
	res, err := r.db.Exec(ctx, "UPDATE cars SET brand=$1, model=$2, color=$3, price=$4 WHERE id = $5", updated.Brand, updated.Model, updated.Color, updated.Price, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

// //////////////////////////////////////////// Delete a single car data
func (r *PgxRepository) DeleteCar(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, "DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}

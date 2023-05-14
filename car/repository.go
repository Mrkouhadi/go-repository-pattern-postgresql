package car

import (
	"context"
	"errors"
)

/*
**
Of course, when using the repository, a whole bunch of other errors can also occur, e.g., network errors, connection failures, timeouts, etc. However, these are unexpected errors,
and not every application needs to handle them in a special way. But user errors like the above should always be handled so that the user can react and correct bad input.
**
*/
var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate(ctx context.Context) error
	CreateCar(ctx context.Context, car Car) (*Car, error)
	AllCars(ctx context.Context) ([]Car, error)
	GetCarById(ctx context.Context, id int64) (*Car, error)
	UpdateCar(ctx context.Context, id int64, updated Car) (*Car, error)
	DeleteCar(ctx context.Context, id int64) error
}

/******
- Migrate(ctx context.Context) error - The method responsible for migrating the repository, i.e., adjusting the PostgreSQL table to the domain object and importing the initial data. In our case, this function will be responsible for creating a new Cars table. So there is no need to log into the GUI database client and manually create the table. However, it is important to remember that this function should be executed first, before reading or writing to the database.
- CreateCar: Method that creates a new Car record in the repository. It returns the Car saved to the database with the generated ID or an error in case of problems.
- AllCars: It extracts all records from the repository and returns as a slice or returns an error if there are problems.
- GetCarById: Gets a single Car record with the specified ID or returns an error if there are problems. The ID of each Car must be unique, so there is no risk of having two records with the same ID.
- UpdateCar: Updates the Car record with the id identifier with the values found in the updated struct. It returns the updated record or an error in case of problems.
- DeleteCar: Deletes a record with the id identifier and returns an error if there are problems.
*******/

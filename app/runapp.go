package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/mrkouhadi/go-repository-pattern-postgresql/car"
)

func RunRepository(ctx context.Context, carRepo car.Repository) {
	///////////////////////////////////////// start migration (create table of cars if it doesn't exist)
	fmt.Println("1. MIGRATE REPOSITORY")
	if err := carRepo.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	////////////////////////////////////////////// create a car data
	fmt.Println("2. CREATE RECORDS OF REPOSITORY")
	audi := car.Car{
		Brand: "Audi",
		Model: "A8",
		Color: "Black",
		Price: 33457.8,
	}
	bmw := car.Car{
		Brand: "BMW",
		Model: "X6",
		Color: "Red",
		Price: 49875.3,
	}
	mercedes := car.Car{
		Brand: "Mercedes",
		Model: "Cl250",
		Color: "White",
		Price: 70875.3,
	}
	// audi
	createAudi, err := carRepo.CreateCar(ctx, audi)
	if errors.Is(err, car.ErrDuplicate) {
		fmt.Printf("record: %+v already exists\n", audi)
	} else if err != nil {
		log.Fatal(err)
	}
	// bmw
	createBmw, err := carRepo.CreateCar(ctx, bmw)
	if errors.Is(err, car.ErrDuplicate) {
		fmt.Printf("record: %+v already exists\n", bmw)
	} else if err != nil {
		log.Fatal(err)
	}
	// bmw
	createMercedes, err := carRepo.CreateCar(ctx, mercedes)
	if errors.Is(err, car.ErrDuplicate) {
		fmt.Printf("record: %+v already exists\n", mercedes)
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n%+v\n%+v\n", createAudi, createBmw, createMercedes)

	////////////////////////////////////////////// get a single car data by ID
	fmt.Println("3. GET RECORD BY ID")
	c, err := carRepo.GetCarById(ctx, 1)
	if errors.Is(err, car.ErrNotExist) {
		log.Println("record: It does not exist in the repository")
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", c)
	////////////////////////////////////////////// get all cars data
	fmt.Println("4. GET ALL cars")
	all, err := carRepo.AllCars(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, car := range all {
		fmt.Printf("%+v\n", car)
	}
	////////////////////////////////////////////// update a single car data
	fmt.Println("5. UPDATE RECORD")
	createAudi.Price = 91699.777 // the price was  33457.8 in the besginning
	if _, err := carRepo.UpdateCar(ctx, createAudi.ID, *createAudi); err != nil {
		if errors.Is(err, car.ErrDuplicate) {
			fmt.Printf("record: %+v already exists\n", createAudi)
		} else if errors.Is(err, car.ErrUpdateFailed) {
			fmt.Printf("update of record: %+v failed", createBmw)
		} else {
			log.Fatal(err)
		}
	}
	////////////////////////////////////////////// get all cars data
	fmt.Println("6. GET ALL cars for the second time after updating the price of 'audi'")
	all, err = carRepo.AllCars(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, car := range all {
		fmt.Printf("%+v\n", car)
	}
	////////////////////////////////////////////// Delete a single car data
	fmt.Println("6. DELETE RECORD")
	if err := carRepo.DeleteCar(ctx, createBmw.ID); err != nil {
		if errors.Is(err, car.ErrDeleteFailed) {
			fmt.Printf("delete of record: %d failed", createBmw.ID)
		} else {
			log.Fatal(err)
		}
	}
	////////////////////////////////////////////// get all cars data
	fmt.Println("7. GET ALL cars for the second time after deleting 'BMW'")
	all, err = carRepo.AllCars(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, car := range all {
		fmt.Printf("%+v\n", car)
	}
}

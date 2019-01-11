package models_test

import (
	"fmt"
	"go-api/models"
	"testing"
)

func TestFruitCreate(t *testing.T) {
	f := &models.Fruit{
		Code: "123",
	}
	affectedRow, err := f.Create(ctx)
	fmt.Println(affectedRow, err, f)
}

func TestFruitUpdate(t *testing.T) {
	f := &models.Fruit{
		Code: "222",
	}
	affectedRow, err := f.Update(ctx, 1)
	fmt.Println(affectedRow, err)
}

func TestFruitDelete(t *testing.T) {
	affectedRow, err := models.Fruit{}.Delete(ctx, 2)
	fmt.Println(affectedRow, err)
}

func TestFruitGetAll(t *testing.T) {
	total, items, err := models.Fruit{}.GetAll(ctx, nil, nil, 0, 2, nil)
	fmt.Println(total, items, err)
}
func TestFruitGetById(t *testing.T) {
	has, v, err := models.Fruit{}.GetById(ctx, 1)
	fmt.Println(has, v, err)
}

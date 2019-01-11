package models

import (
	"context"

	"github.com/go-xorm/xorm"

	"go-api/factory"
)

type Fruit struct {
	Id        int64  `json:"id" xorm:"int64 notnull autoincr pk 'id'"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Price     int64  `json:"price"`
	StoreCode string `json:"storeCode"`
	CreatedAt string `json:"createdAt" xorm:"created"`
	UpdatedAt string `json:"updatedAt" xorm:"updated"`
	DeletedAt string `json:"deletedAt" xorm:"deleted"`
}

type Store struct {
	Id   int64  `json:"id" xorm:"pk autoincr 'id'"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (d *Fruit) Create(ctx context.Context) (affectedRow int64, err error) {
	affectedRow, err = factory.DB(ctx).Insert(d)
	return
}
func (Fruit) GetFullById(ctx context.Context, id int64) (has bool, result interface{}, err error) {

	var fruitDto struct {
		Id        int64  `json:"id"`
		Name      string `json:"name"`
		Color     string `json:"color"`
		Price     int64  `json:"price"`
		StoreName string `json:"storeName" xorm:"store_name"` // note: xorm:"store_name" ==== b.name as store_name
	}
	has, err = factory.DB(ctx).Table("fruit").Alias("a").
		Join("inner", []string{"store", "b"}, "a.store_code = b.code").
		Select(`a.id,a.name,a.color,a.price,b.name as store_name`).
		Where("a.id=?", id).Get(&fruitDto)
	result = fruitDto
	return
}
func (Fruit) GetById(ctx context.Context, id int64) (has bool, fruit *Fruit, err error) {
	fruit = &Fruit{}
	has, err = factory.DB(ctx).Where("id=?", id).Get(fruit)
	return
}
func (Fruit) GetAll(ctx context.Context, sortby, order []string, offset, limit int, options *FruitSearchOption) (totalCount int64, items []*Fruit, err error) {
	q := func() xorm.Session {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		if options != nil {
			if len(options.Name) != 0 {
				q.Where("name like ?", options.Name+"%")
			}
		}
		return *q
	}

	errc, totalCountc, fruitc := make(chan error), make(chan int64, 1), make(chan []*Fruit, 1)
	go func(qNew xorm.Session) {
		v, err := qNew.Count(&Fruit{})
		totalCountc <- v
		errc <- err
	}(q())

	go func(qNew xorm.Session) {
		var v []*Fruit
		err := qNew.Limit(limit, offset).Find(&v)
		fruitc <- v
		errc <- err
	}(q())
	for i := 0; i < 2; i++ {
		if err := <-errc; err != nil {
			return 0, nil, err
		}
	}
	totalCount = <-totalCountc
	items = <-fruitc
	return
}
func (d *Fruit) Update(ctx context.Context, id int64) (affectedRow int64, err error) {
	affectedRow, err = factory.DB(ctx).Cols("name,color,price").Where("id=?", id).Update(d)
	return
}

func (Fruit) Delete(ctx context.Context, id int64) (affectedRow int64, err error) {
	affectedRow, err = factory.DB(ctx).Where("id=?", id).Delete(&Fruit{})
	return
}

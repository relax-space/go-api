package models

import (
	"context"
	"time"

	"github.com/go-xorm/xorm"

	"go-api/factory"
)

type Fruit struct {
	Id        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Price     int64     `json:"price"`
	StoreCode string    `json:"storeCode"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}

type Store struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type FruitStoreDto struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Price     int64  `json:"price"`
	StoreName string `json:"storeName" xorm:"store_name"` // note: xorm:"store_name" ==== b.name as store_name
}

func (d *Fruit) Create(ctx context.Context) (affectedRow int64, err error) {
	affectedRow, err = factory.DB(ctx).Insert(d)
	return
}

func (Fruit) GetById(ctx context.Context, id int64) (has bool, fruit Fruit, err error) {
	has, err = factory.DB(ctx).Where("id=?", id).Get(&fruit)
	return
}

func (Fruit) GetByCode(ctx context.Context, code string) (has bool, fruit Fruit, err error) {
	has, err = factory.DB(ctx).Where("code=?", code).Get(&fruit)
	return
}

func (Fruit) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []*Fruit, err error) {
	queryBuilder := func() xorm.Interface {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	totalCount, err = queryBuilder().Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		return
	}
	return
}

func (d *Fruit) Update(ctx context.Context, id int64) (affectedRow int64, err error) {
	affectedRow, err = factory.DB(ctx).Where("id=?", id).Update(d)
	return
}

func (Fruit) Delete(ctx context.Context, id int64) (affectedRow int64, err error) {
	affectedRow, err = factory.DB(ctx).Where("id=?", id).Delete(&Fruit{})
	return
}

func (Fruit) GetWithStoreById(ctx context.Context, id int64) (has bool, dto FruitStoreDto, err error) {
	has, err = factory.DB(ctx).Table("fruit").Alias("a").
		Join("inner", []string{"store", "b"}, "a.store_code = b.code").
		Select(`a.id,a.name,a.color,a.price,b.name as store_name`).
		Where("a.id=?", id).Get(&dto)
	return
}

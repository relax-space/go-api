package models

import (
	"context"
	"time"

	"github.com/relax-space/go-api/factory"
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

func (d *Fruit) Create(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(d)
}

func (Fruit) GetById(ctx context.Context, id int64) (bool, Fruit, error) {
	fruit := Fruit{}
	has, err := factory.DB(ctx).Where("id=?", id).Get(&fruit)
	return has, fruit, err
}

func (Fruit) GetByCode(ctx context.Context, code string) (bool, Fruit, error) {
	fruit := Fruit{}
	has, err := factory.DB(ctx).Where("code=?", code).Get(&fruit)
	return has, fruit, err
}

func (Fruit) GetAll(ctx context.Context, sortby, order []string, offset, limit int, withHasMore bool) (bool, int64, []Fruit, error) {
	query := factory.DB(ctx)
	if err := setSortOrder(query, sortby, order); err != nil {
		return false, 0, nil, err
	}

	var (
		items      []Fruit
		hasMore    bool
		totalCount int64
		err        error
	)
	if withHasMore {
		err = query.Limit(limit+1, offset).Find(&items)
		if len(items) == limit+1 {
			items = items[:limit]
			hasMore = true
		}
	} else {
		totalCount, err = query.Limit(limit, offset).FindAndCount(&items)
	}
	return hasMore, totalCount, items, err
}

func (d *Fruit) Update(ctx context.Context, id int64) (int64, error) {
	return factory.DB(ctx).Where("id=?", id).Update(d)

}

func (Fruit) Delete(ctx context.Context, id int64) (int64, error) {
	return factory.DB(ctx).Where("id=?", id).Delete(&Fruit{})
}

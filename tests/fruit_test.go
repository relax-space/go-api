package tests_test

import (
	"testing"
	
	"github.com/relax-space/go-api/models"
	"github.com/pangpanglabs/goutils/test"
)

func TestFruitCRUD(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		_,count, fruits, err := models.Fruit{}.GetAll(ctx,[]string{"id"},[]string{"asc"},0, 1,false)
		test.Ok(t, err)
		test.Assert(t,count>=1,"result must be greater than 1")
		v :=fruits[0]
		test.Equals(t, int64(1), v.Id)
		test.Equals(t, "apple", v.Code)
		test.Equals(t, "apple", v.Name)
		test.Equals(t, int64(11), v.Price)
		test.Equals(t, "10001", v.StoreCode)
	})
	t.Run("GetAll_setSortOrder_err1", func(t *testing.T) {
		_,_, _, err := models.Fruit{}.GetAll(ctx,[]string{"id"},[]string{"asc1"},0, 1,false)
		test.Assert(t,err != nil,"When the array lengths of sortby and order are inconsistent, you should throw an error")
	})
	t.Run("GetAll_setSortOrder_err2", func(t *testing.T) {
		_,_, _, err := models.Fruit{}.GetAll(ctx,[]string{"id","code"},[]string{"asc1"},0, 1,false)
		test.Assert(t,err != nil,"When the array lengths of sortby and order are inconsistent, you should throw an error")
	})
	t.Run("GetAll_setSortOrder_err3", func(t *testing.T) {
		_,_, _, err := models.Fruit{}.GetAll(ctx,[]string{"id"},[]string{"asc","desc"},0, 1,false)
		test.Assert(t,err != nil,"When the array lengths of sortby and order are inconsistent, you should throw an error")
	})
	t.Run("GetAll_setSortOrder_err4", func(t *testing.T) {
		_,_, _, err := models.Fruit{}.GetAll(ctx,nil,[]string{"asc"},0, 1,false)
		test.Assert(t,err != nil,"When the array lengths of sortby and order are inconsistent, you should throw an error")
	})
	t.Run("GetAll_setSortOrder_desc", func(t *testing.T) {
		_,count, _, err := models.Fruit{}.GetAll(ctx,[]string{"id"},[]string{"desc"},0, 1,false)
		test.Ok(t, err)
		test.Assert(t,count>=1,"result must be greater than 1")
	})
	t.Run("GetAll_setSortOrder_one_order", func(t *testing.T) {
		_,count, fruits, err := models.Fruit{}.GetAll(ctx,[]string{"id","code"},[]string{"asc"},0, 1,false)
		test.Ok(t, err)
		test.Assert(t,count>=1,"result must be greater than 1")
		v :=fruits[0]
		test.Equals(t, int64(1), v.Id)
		test.Equals(t, "apple", v.Code)
		test.Equals(t, "apple", v.Name)
		test.Equals(t, int64(11), v.Price)
		test.Equals(t, "10001", v.StoreCode)
	})
	t.Run("GetAll_setSortOrder_one_order_desc", func(t *testing.T) {
		_,count, _, err := models.Fruit{}.GetAll(ctx,[]string{"id","code"},[]string{"desc"},0, 1,false)
		test.Ok(t, err)
		test.Assert(t,count>=1,"result must be greater than 1")
	})
	t.Run("GetAll_hasMore", func(t *testing.T) {
		hasMore,_, fruits, err := models.Fruit{}.GetAll(ctx,[]string{"id"},[]string{"asc"},0, 1,true)
		test.Ok(t, err)
		v :=fruits[0]
		test.Equals(t, true,hasMore)
		test.Equals(t, int64(1), v.Id)
		test.Equals(t, "apple", v.Code)
		test.Equals(t, "apple", v.Name)
		test.Equals(t, int64(11), v.Price)
		test.Equals(t, "10001", v.StoreCode)
	})
	t.Run("Get", func(t *testing.T) {
		has, v, err := models.Fruit{}.GetById(ctx, 1)
		test.Ok(t, err)
		test.Equals(t, true, has)
		test.Equals(t, int64(1), v.Id)
		test.Equals(t, "apple", v.Code)
		test.Equals(t, "apple", v.Name)
		test.Equals(t, int64(11), v.Price)
		test.Equals(t, "10001", v.StoreCode)
	})
	t.Run("Create", func(t *testing.T) {
		code := "apple"
		f := &models.Fruit{
			Code: code,
		}
		affectedRow, err := f.Create(ctx)
		test.Ok(t, err)
		test.Equals(t, int64(1), affectedRow)
		test.Assert(t, f.Id != int64(0), "create failure")
	})

	t.Run("Update", func(t *testing.T) {
		var price int64 = 10
		f := &models.Fruit{
			Price: price,
		}
		affectedRow, err := f.Update(ctx, 1)
		test.Ok(t, err)
		test.Equals(t, int64(1), affectedRow)
	})
	t.Run("Delete", func(t *testing.T) {
		affectedRow, err := models.Fruit{}.Delete(ctx, 1)
		test.Ok(t, err)
		test.Equals(t, int64(1), affectedRow)
	})
	rollback()
}


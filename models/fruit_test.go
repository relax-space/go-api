package models_test

import (
	"go-api/models"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFruit(t *testing.T) {

	var ID int64 = 2
	Convey("Test Fruit", t, func() {

		Convey("Attempting to retrieve the fruit should return fruit", func() {
			has, v, err := models.Fruit{}.GetById(ctx, ID)
			So(err, ShouldBeNil)
			So(has, ShouldBeTrue)
			So(v.Id, ShouldEqual, 2)
		})

		Convey("Given a fruit in the database", func() {
			code := "apple"
			f := &models.Fruit{
				Code: code,
			}
			affectedRow, err := f.Create(ctx)
			So(err, ShouldBeNil)
			So(affectedRow, ShouldEqual, int64(1))
			ID = f.Id
			has, v, err := models.Fruit{}.GetById(ctx, ID)
			So(err, ShouldBeNil)
			So(has, ShouldBeTrue)
			So(v.Code, ShouldEqual, code)

		})
		Convey("Update fruit price in the database", func() {
			var price int64 = 10
			f := &models.Fruit{
				Price: price,
			}
			affectedRow, err := f.Update(ctx, ID)
			So(err, ShouldBeNil)
			So(affectedRow, ShouldEqual, int64(1))
			has, v, err := models.Fruit{}.GetById(ctx, ID)
			So(err, ShouldBeNil)
			So(has, ShouldBeTrue)
			So(v.Price, ShouldEqual, price)
		})
		Convey("Delete a fruit in the database", func() {
			affectedRow, err := models.Fruit{}.Delete(ctx, ID)
			So(err, ShouldBeNil)
			So(affectedRow, ShouldEqual, int64(1))
			has, v, err := models.Fruit{}.GetById(ctx, ID)
			So(err, ShouldBeNil)
			So(has, ShouldBeFalse)
			So(v.Id, ShouldEqual, int64(0))
		})

	})
}

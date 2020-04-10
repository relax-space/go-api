package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/relax-space/go-api/models"
	"github.com/hublabs/common/api"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type FruitApiController struct {
}

// localhost:8080/docs
func (d FruitApiController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")
	g.GET("", d.GetAll).
		AddParamQueryNested(SearchInput{})
	g.GET("/:id", d.GetOne).
		AddParamPath("", "id", "id").AddParamQuery("", "with_store", "with_store", false)
	g.PUT("/:id", d.Update).
		AddParamPath("", "id", "id").
		AddParamBody(models.Fruit{}, "fruit", "only can modify name,color,price", true)
	g.POST("", d.Create).
		AddParamBody(models.Fruit{}, "fruit", "new fruit", true)
	g.DELETE("/:id", d.Delete).
		AddParamPath("", "id", "id")
}

/*
localhost:8080/fruits
localhost:8080/fruits?name=apple
localhost:8080/fruits?skipCount=0&maxResultCount=2
localhost:8080/fruits?skipCount=0&maxResultCount=2&sortby=store_code&order=desc
*/
func (FruitApiController) GetAll(c echo.Context) error {
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}
	hasMore,totalCount, items, err := models.Fruit{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount,v.WithHasMore)
	if err != nil {
		return renderFail(c,api.ErrorDB.New(err))
	}
	return renderSuccArray(c, v.WithHasMore, hasMore, totalCount, items)
}

/*
localhost:8080/fruits/1?with_store=true
localhost:8080/fruits/1
*/
func (d FruitApiController) GetOne(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameterParsingFailed.New(err,fmt.Sprintf("id:%v",c.Param("id"))))
	}

	var withStore bool
	if len(c.QueryParam("with_store")) != 0 {
		withStore, err = strconv.ParseBool(c.QueryParam("with_store"))
		if err != nil {
			return renderFail(c, api.ErrorParameterParsingFailed.New(err,fmt.Sprintf("with_store:%v",c.Param("with_store"))))
		}
	}
	if withStore == true {
		_, fruit, err := models.Fruit{}.GetWithStoreById(c.Request().Context(), id)
		if err != nil {
			return renderFail(c, api.ErrorNotFound.New(err))
		}
		return renderSucc(c, http.StatusOK, fruit)
	}

	_, fruit, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if fruit.Id == 0 {
		return renderFail(c, api.ErrorNotFound.New(nil))
	}
	return renderSucc(c, http.StatusOK, fruit)
}

/*
localhost:8080/fruits
 {
        "code": "AA01",
        "name": "Apple",
        "color": "",
        "price": 2,
        "store_code": ""
    }
*/
func (d FruitApiController) Create(c echo.Context) error {
	var v models.Fruit
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	has, _, err := models.Fruit{}.GetByCode(c.Request().Context(), v.Code)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if has {
		return renderFail(c, api.ErrorHasExisted.New(nil))
	}
	affectedRow, err := v.Create(c.Request().Context())
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if affectedRow == int64(0) {
		return renderFail(c, api.ErrorNotCreated.New(nil))
	}
	return renderSucc(c, http.StatusCreated, v)
}

/*
localhost:8080/fruits
 {
        "price": 21,
    }
*/
func (d FruitApiController) Update(c echo.Context) error {
	var v models.Fruit
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameterParsingFailed.New(err,fmt.Sprintf("id:%v",c.Param("id"))))
	}
	has, _, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))

	}
	if has == false {
		return renderFail(c, api.ErrorNotFound.New(nil))
	}
	affectedRow, err := v.Update(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if affectedRow == int64(0) {
		return renderFail(c, api.ErrorNotUpdated.New(nil))
	}
	return renderSucc(c, http.StatusOK, v)
}

/*
localhost:8080/fruits/45
*/
func (d FruitApiController) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameterParsingFailed.New(err,fmt.Sprintf("id:%v",c.Param("id"))))
	}
	has, v, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if has == false {
		return renderFail(c, api.ErrorNotFound.New(nil))
	}
	affectedRow, err := models.Fruit{}.Delete(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if affectedRow == int64(0) {
		return renderFail(c, api.ErrorNotDeleted.New(nil))
	}
	return renderSucc(c, http.StatusOK, v)
}

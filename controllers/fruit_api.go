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
// FruitAPIController define a struct
type FruitAPIController struct {
}

// Init set router
func (d FruitAPIController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")
	g.GET("", d.GetAll).
		AddParamQueryNested(SearchInput{})
	g.GET("/:id", d.GetOne).
		AddParamPath("", "id", "id")
	g.PUT("/:id", d.Update).
		AddParamPath("", "id", "id").
		AddParamBody(struct{
			Code 	string 	`json:"code"`
			Name 	string 	`json:"name"`
			Color 	string 	`json:"color"`
			Price 	int64 	`json:"price"`
		}{
			Code:"apple",
			Name:"apple",
			Color:"red",
			Price:12,
		}, "fruit", "only can modify name,color,price", true)
	g.POST("", d.Create).
		AddParamBody(struct{
			Code 	string 	`json:"code"`
			Name 	string 	`json:"name"`
			Color 	string 	`json:"color"`
			Price 	int64 	`json:"price"`
		}{
			Code:"banana",
			Name:"banana",
			Color:"yellow",
			Price:16,
		}, "fruit", "new fruit", true)
	g.DELETE("/:id", d.Delete).
		AddParamPath("", "id", "id")
}
// GetAll search multi fruits
func (FruitAPIController) GetAll(c echo.Context) error {
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

// GetOne query a fruit
func (d FruitAPIController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameterParsingFailed.New(err,fmt.Sprintf("id:%v",c.Param("id"))))
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

// Create a fruit
func (d FruitAPIController) Create(c echo.Context) error {
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

// Update a fruit
func (d FruitAPIController) Update(c echo.Context) error {
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

// Delete a fruit
func (d FruitAPIController) Delete(c echo.Context) error {
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

package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/relax-space/go-api/controllers"
	"github.com/relax-space/go-api/models"
	"github.com/hublabs/common/api"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/test"
)


func TestFruitAPICRUD(t *testing.T) {

	t.Run("GetAll", func(t *testing.T) {
		expFruits := []models.Fruit{
			models.Fruit{
				Id:        int64(1),
				Code:      "apple",
				Name:      "apple",
				Color:     "red",
				Price:     int64(11),
				StoreCode: "10001",
			},

			models.Fruit{
				Id:        int64(2),
				Code:      "banana",
				Name:      "banana",
				Color:     "yellow",
				Price:     int64(14),
				StoreCode: "10002",
			},
		}
		req := httptest.NewRequest(echo.GET, "/v1/fruits?maxResultCount=2", nil)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetAll, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				TotalCount int            `json:"totalCount"`
				Items      []models.Fruit `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Assert(t,v.Result.TotalCount>=2,"result must be greater than 2")
		test.Equals(t, expFruits[0].Id, v.Result.Items[0].Id)
		test.Equals(t, expFruits[0].Code, v.Result.Items[0].Code)
		test.Equals(t, expFruits[0].Name, v.Result.Items[0].Name)
		test.Equals(t, expFruits[0].Color, v.Result.Items[0].Color)
		test.Equals(t, expFruits[0].Price, v.Result.Items[0].Price)
		test.Equals(t, expFruits[0].StoreCode, v.Result.Items[0].StoreCode)

		test.Equals(t, expFruits[1].Id, v.Result.Items[1].Id)
		test.Equals(t, expFruits[1].Code, v.Result.Items[1].Code)
		test.Equals(t, expFruits[1].Name, v.Result.Items[1].Name)
		test.Equals(t, expFruits[1].Color, v.Result.Items[1].Color)
		test.Equals(t, expFruits[1].Price, v.Result.Items[1].Price)
		test.Equals(t, expFruits[1].StoreCode, v.Result.Items[1].StoreCode)

	})

	t.Run("GetAll_Default", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/fruits", nil)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetAll, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				TotalCount int            `json:"totalCount"`
				HasMore  	bool          `json:"hasMore"`
				Items      []models.Fruit `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, true, v.Success)
		test.Equals(t, false, v.Result.HasMore)
		test.Assert(t,v.Result.TotalCount>=2,"result must be greater than 2")

	})

	t.Run("GetAll_hasMore", func(t *testing.T) {
		expFruits := []models.Fruit{
			models.Fruit{
				Id:        int64(1),
				Code:      "apple",
				Name:      "apple",
				Color:     "red",
				Price:     int64(11),
				StoreCode: "10001",
			},
		}
		req := httptest.NewRequest(echo.GET, "/v1/fruits?maxResultCount=1&withHasMore=true", nil)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetAll, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				HasMore  	bool          `json:"hasMore"`
				Items      []models.Fruit `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, true, v.Success)
		test.Equals(t, true, v.Result.HasMore)
		test.Equals(t, expFruits[0].Id, v.Result.Items[0].Id)
		test.Equals(t, expFruits[0].Code, v.Result.Items[0].Code)
		test.Equals(t, expFruits[0].Name, v.Result.Items[0].Name)
		test.Equals(t, expFruits[0].Color, v.Result.Items[0].Color)
		test.Equals(t, expFruits[0].Price, v.Result.Items[0].Price)
		test.Equals(t, expFruits[0].StoreCode, v.Result.Items[0].StoreCode)

	})

	t.Run("GetOne", func(t *testing.T) {
		expFruit := models.Fruit{
			Id:        int64(1),
			Code:      "apple",
			Name:      "apple",
			Color:     "red",
			Price:     int64(11),
			StoreCode: "10001",
		}
		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetOne, c))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, true, v.Success)

		test.Equals(t, expFruit.Id, v.Result.Id)
		test.Equals(t, expFruit.Code, v.Result.Code)
		test.Equals(t, expFruit.Name, v.Result.Name)
		test.Equals(t, expFruit.Color, v.Result.Color)
		test.Equals(t, expFruit.Price, v.Result.Price)
		test.Equals(t, expFruit.StoreCode, v.Result.StoreCode)
	})

	createFruits := []models.Fruit{
		models.Fruit{
			Code:  "fruit#1",
			Color: "red",
		},
		models.Fruit{
			Code:  "fruit#2",
			Color: "green",
		},
	}
	for i, p := range createFruits {
		pb, _ := json.Marshal(p)
		t.Run(fmt.Sprint("Create#", i+1), func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/v1/fruits", bytes.NewReader(pb))
			rec := httptest.NewRecorder()
			test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Create, echo.New().NewContext(req, rec)))
			test.Equals(t, http.StatusCreated, rec.Code)

			var v struct {
				Result  models.Fruit `json:"result"`
				Success bool         `json:"success"`
			}
			test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
			test.Equals(t, true, v.Success)
		})
	}

	t.Run("Update", func(t *testing.T) {
		expFruit := models.Fruit{
			Id:        int64(1),
			Code:      "apple",
			Name:      "iphone",
			Color:     "red",
			Price:     int64(11),
			StoreCode: "10001",
		}
		pb, _ := json.Marshal(expFruit)
		req := httptest.NewRequest(echo.PUT, "/v1/fruits/:id", bytes.NewReader(pb))
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", expFruit.Id))
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Update, c))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, true, v.Success)

		test.Equals(t, expFruit.Id, v.Result.Id)
		test.Equals(t, expFruit.Code, v.Result.Code)
		test.Equals(t, expFruit.Name, v.Result.Name)
		test.Equals(t, expFruit.Color, v.Result.Color)
		test.Equals(t, expFruit.Price, v.Result.Price)
		test.Equals(t, expFruit.StoreCode, v.Result.StoreCode)
	})

	t.Run("Delete", func(t *testing.T) {
		expFruit := models.Fruit{
			Id:        int64(1),
			Code:      "apple",
			Name:      "apple",
			Color:     "red",
			Price:     int64(11),
			StoreCode: "10001",
		}
		req := httptest.NewRequest(echo.DELETE, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", 1))
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Delete, c))
		test.Equals(t, http.StatusOK, rec.Code)
		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, true, v.Success)

		test.Equals(t, expFruit.Id, v.Result.Id)
		test.Equals(t, expFruit.Code, v.Result.Code)
		test.Equals(t, expFruit.Name, v.Result.Name)
		test.Equals(t, expFruit.Color, v.Result.Color)
		test.Equals(t, expFruit.Price, v.Result.Price)
		test.Equals(t, expFruit.StoreCode, v.Result.StoreCode)
	})
}


func TestFruitAPICRUDFail(t *testing.T) {

	t.Run("GetAll_bind", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/fruits?maxResultCount=A", nil)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetAll, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusBadRequest, rec.Code)

		var v struct {
			Result struct {
				TotalCount int            `json:"totalCount"`
				Items      []models.Fruit `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t,false,v.Success)
		test.Equals(t, api.ErrorParameter.Code,v.Error.Code)
	})

	t.Run("GetAll_DB", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/fruits?maxResultCount=1&sortby=id1", nil)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetAll, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Result struct {
				HasMore  	bool          `json:"hasMore"`
				Items      []models.Fruit `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
	})

	t.Run("GetOne_id", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("A")
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetOne, c))
		test.Equals(t, http.StatusBadRequest, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorParameterParsingFailed.Code,v.Error.Code)
	})

	t.Run("GetOne_NotFound", func(t *testing.T) {
		
		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1010")
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.GetOne, c))
		test.Equals(t, http.StatusNotFound, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorNotFound.Code,v.Error.Code)
	})

	t.Run("Create_bind", func(t *testing.T) {
		
		b :=[]byte("lalala")
		req := httptest.NewRequest(echo.POST, "/v1/fruits", bytes.NewReader(b))
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Create, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusBadRequest, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Equals(t, false, v.Success)
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, api.ErrorParameter.Code,v.Error.Code)
	})

	t.Run("Create_exist", func(t *testing.T) {
		fruit := models.Fruit{
			Code:  "apple",
			Color: "red",
		}
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.POST, "/v1/fruits", bytes.NewReader(pb))
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Create, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusBadRequest, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Equals(t, false, v.Success)
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, api.ErrorHasExisted.Code,v.Error.Code)
	})

	t.Run("Create_Long_string", func(t *testing.T) {
		fruit := models.Fruit{
			Code:  "apple1010",
			Color: "redredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredred",
		}
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.POST, "/v1/fruits", bytes.NewReader(pb))
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Create, echo.New().NewContext(req, rec)))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
	})

	t.Run("Update_bind", func(t *testing.T) {
		
		b :=[]byte("lalala")
		req := httptest.NewRequest(echo.PUT, "/v1/fruits/:id", bytes.NewReader(b))
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("A")
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Update, c))
		test.Equals(t, http.StatusBadRequest, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorParameter.Code,v.Error.Code)
	})

	t.Run("Update_Id", func(t *testing.T) {
		fruit := models.Fruit{
			Code:      "apple2",
		}
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.PUT, "/v1/fruits/:id", bytes.NewReader(pb))
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("A")
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Update, c))
		test.Equals(t, http.StatusBadRequest, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorParameterParsingFailed.Code,v.Error.Code)
	})


	t.Run("Update_notFound", func(t *testing.T) {
		fruit := models.Fruit{
			Id:        int64(1010),
			Code:      "apple2",
		}
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.PUT, "/v1/fruits/:id", bytes.NewReader(pb))
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", fruit.Id))
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Update, c))
		test.Equals(t, http.StatusNotFound, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorNotFound.Code,v.Error.Code)
	})

	t.Run("Update_long_string", func(t *testing.T) {
		fruit := models.Fruit{
			Id:        int64(1),
			Color: "redredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredredred",
		}
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.PUT, "/v1/fruits/:id", bytes.NewReader(pb))
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", fruit.Id))
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Update, c))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
	})

	t.Run("Delete_Bind", func(t *testing.T) {
		req := httptest.NewRequest(echo.DELETE, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", "A"))
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Delete, c))
		test.Equals(t, http.StatusBadRequest, rec.Code)
		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorParameterParsingFailed.Code,v.Error.Code)
	})

	t.Run("Delete_notFount", func(t *testing.T) {
		req := httptest.NewRequest(echo.DELETE, "/v1/fruits/:id", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", 1010))
		test.Ok(t, handleWithFilter(controllers.FruitAPIController{}.Delete, c))
		test.Equals(t, http.StatusNotFound, rec.Code)
		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorNotFound.Code,v.Error.Code)

	})

}


func TestFruitAPICRUDFailDB(t *testing.T) {

	t.Run("GetAll", func(t *testing.T) {

		req := httptest.NewRequest(echo.GET, "/v1/fruits?maxResultCount=1", nil)
		c,rec := SetContextWithDBClose(req)
		test.Ok(t,controllers.FruitAPIController{}.GetAll(c))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
		
	})

	t.Run("GetOne", func(t *testing.T) {

		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id", nil)
		c,rec := SetContextWithDBClose(req)
		c.SetParamNames("id")
		c.SetParamValues("1")
		test.Ok(t,controllers.FruitAPIController{}.GetOne(c))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
		
	})

	t.Run("Create", func(t *testing.T) {
		fruit := models.Fruit{
			Code:  "apple123",
			Color: "red",
		}
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.GET, "/v1/fruits", bytes.NewReader(pb))
		c,rec := SetContextWithDBClose(req)
		test.Ok(t,controllers.FruitAPIController{}.Create(c))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
	})

	t.Run("Update", func(t *testing.T) {
		fruit := models.Fruit{
			Id:        int64(1),
			Code:      "apple234",
		}
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id",  bytes.NewReader(pb))
		c,rec := SetContextWithDBClose(req)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", fruit.Id))
		test.Ok(t,controllers.FruitAPIController{}.Update(c))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/fruits/:id", nil)
		c,rec := SetContextWithDBClose(req)
		c.SetParamNames("id")
		c.SetParamValues("1")
		test.Ok(t,controllers.FruitAPIController{}.Delete(c))
		test.Equals(t, http.StatusInternalServerError, rec.Code)

		var v struct {
			Success bool         `json:"success"`
			Error 	api.Error 	`json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, false, v.Success)
		test.Equals(t, api.ErrorDB.Code,v.Error.Code)
	})


}



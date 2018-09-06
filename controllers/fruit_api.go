package controllers

import (
	"errors"
	"go-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type FruitApiController struct {
}

func (d FruitApiController) Init(g *echo.Group) {
	g.GET("", d.GetAll)
	g.GET("/:id", d.GetOne)
	g.PUT("/:id", d.Update)
	g.POST("", d.Create)
	g.DELETE("/:id", d.Delete)
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
		return ReturnApiListFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}
	name := c.QueryParam("name")
	totalCount, items, err := models.Fruit{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount,
		&models.FruitSearchOption{Name: name})
	if err != nil {
		return ReturnApiListFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if len(items) == 0 {
		return ReturnApiListFail(c, http.StatusNotFound, ApiErrorNotFound, err)
	}
	return ReturnApiListSucc(c, http.StatusOK, totalCount, items)
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
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	return d.create(c, &v)
}

func (FruitApiController) create(c echo.Context, fruit *models.Fruit) (err error) {
	has, _, err := models.Fruit{}.GetById(c.Request().Context(), fruit.Id)
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if has {
		ReturnApiFail(c, http.StatusOK, ApiErrorHasExist, nil)
		err = errors.New(ApiErrorHasExist.Message)
		return
	}
	affectedRow, err := fruit.Create(c.Request().Context())
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if affectedRow == int64(0) {
		ReturnApiFail(c, http.StatusOK, ApiErrorNotChanged, nil)
		err = errors.New(ApiErrorNotChanged.Message)
		return
	}
	return ReturnApiSucc(c, http.StatusCreated, fruit)
}

/*
localhost:8080/fruits/1?full=1
localhost:8080/fruits/1
*/
func (d FruitApiController) GetOne(c echo.Context) error {
	full := c.QueryParam("full")
	if len(full) != 0 {
		return d.GetOneFull(c)
	} else {
		return d.GetOnePart(c)
	}
}

func (FruitApiController) GetOnePart(c echo.Context) (err error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err, map[string]interface{}{"id": c.Param("id")})
		return
	}
	has, v, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if !has {
		ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
		err = errors.New(ApiErrorNotFound.Message)
		return
	}
	err = ReturnApiSucc(c, http.StatusOK, v)
	return
}

func (FruitApiController) GetOneFull(c echo.Context) (err error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err, map[string]interface{}{"id": c.Param("id")})
		return
	}
	has, v, err := models.Fruit{}.GetFullById(c.Request().Context(), id)
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if !has {
		ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
		err = errors.New(ApiErrorNotFound.Message)
		return
	}
	err = ReturnApiSucc(c, http.StatusOK, v)
	return
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
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err, map[string]interface{}{"id": c.Param("id")})
	}
	v.Id = id
	return d.update(c, &v)
}

func (FruitApiController) update(c echo.Context, v *models.Fruit) (err error) {
	has, _, err := models.Fruit{}.GetById(c.Request().Context(), v.Id)
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if !has {
		ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
		err = errors.New(ApiErrorNotFound.Message)
		return
	}
	affectedRow, err := v.Update(c.Request().Context(), v.Id)
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if affectedRow == int64(0) {
		ReturnApiFail(c, http.StatusOK, ApiErrorNotChanged, err)
		err = errors.New(ApiErrorNotChanged.Message)
		return
	}
	err = ReturnApiSucc(c, http.StatusNoContent, nil)
	return
}

/*
localhost:8080/fruits/45
*/
func (d FruitApiController) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err, map[string]interface{}{"id": c.Param("id")})
	}
	return d.delete(c, id)
}

func (FruitApiController) delete(c echo.Context, id int64) (err error) {
	has, _, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if !has {
		ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
		err = errors.New(ApiErrorNotFound.Message)
		return
	}
	affectedRow, err := models.Fruit{}.Delete(c.Request().Context(), id)
	if err != nil {
		ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
		return
	}
	if affectedRow == int64(0) {
		ReturnApiFail(c, http.StatusOK, ApiErrorNotChanged, err)
		err = errors.New(ApiErrorNotChanged.Message)
		return
	}
	err = ReturnApiSucc(c, http.StatusNoContent, nil)
	return
}

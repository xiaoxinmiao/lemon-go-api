package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-api/models"
	"net/http/httptest"
	"testing"

	"go-api/controllers"

	"github.com/relax-space/go-kit/test"

	"github.com/labstack/echo"
)

func Test_fruit_GetAll(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/?name=2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	test.Ok(t, handleWithFilter(controllers.FruitApiController{}.GetAll, echoApp.NewContext(req, rec)))
	fmt.Println(string(rec.Body.Bytes()))
	fmt.Printf("http status:%v", rec.Result().StatusCode)
}

func Test_fruit_GetFullById(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echoApp.NewContext(req, rec)
	c.SetPath("fruits/:id?full=true")
	c.SetParamNames("id")
	c.SetParamValues("2")
	test.Ok(t, handleWithFilter(controllers.FruitApiController{}.GetOneFull, c))
	fmt.Println(string(rec.Body.Bytes()))
	fmt.Printf("http status:%v", rec.Result().StatusCode)
}

func Test_fruit_Create(t *testing.T) {
	fruit := &models.Fruit{
		Code: "123",
	}
	b, _ := json.Marshal(fruit)
	req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	test.Ok(t, handleWithFilter(controllers.FruitApiController{}.Create, echoApp.NewContext(req, rec)))
	fmt.Println(string(rec.Body.Bytes()))
	fmt.Printf("http status:%v", rec.Result().StatusCode)
}

func Test_fruit_Update(t *testing.T) {
	fruit := &models.Fruit{
		Code: "2222",
	}
	b, _ := json.Marshal(fruit)
	req := httptest.NewRequest(echo.POST, "/1", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	test.Ok(t, handleWithFilter(controllers.FruitApiController{}.Update, echoApp.NewContext(req, rec)))
	fmt.Println(string(rec.Body.Bytes()))
	fmt.Printf("http status:%v", rec.Result().StatusCode)
}

func Test_fruit_Delete(t *testing.T) {
	id := "12"
	req := httptest.NewRequest(echo.POST, "/"+id, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	test.Ok(t, handleWithFilter(controllers.FruitApiController{}.Delete, echoApp.NewContext(req, rec)))
	fmt.Println(string(rec.Body.Bytes()))
	fmt.Printf("http status:%v", rec.Result().StatusCode)
}

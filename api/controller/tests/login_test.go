package tests

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
	"main/api/controller"
	"main/bootstrap"
	"main/domain"
	"main/repository"
	"main/usecase"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	app := bootstrap.App()
	env := app.Env
	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()
	timeout := time.Duration(env.ContextTimeout) * time.Second
	ur := repository.NewUserRepository(db, domain.UserCollection)
	lc := &controller.LoginController{
		LoginUseCase: usecase.NewLoginUseCase(ur, timeout),
		Env:          env,
	}
	r := chi.NewRouter()

	Convey("given the 1st POST http request for /public/login", t, func() {
		jsonByte, _ := json.Marshal(domain.LoginRequest{ // email and password of mocked user
			"testadmin@mail.ru",
			"qwerty123",
		})
		req := httptest.NewRequest("POST", "/public/login", bytes.NewReader(jsonByte))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/public/login", lc.Login)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
			Convey("and the response should contain a message", func() {
				So(responseBody, ShouldContainSubstring, "access_token")
			})
		})
	})

	Convey("given the 2nd POST http request for /public/login with wrong password", t, func() {
		jsonByte, _ := json.Marshal(domain.LoginRequest{
			"testadmin@mail.ru",
			"WrongPassword",
		})
		req := httptest.NewRequest("POST", "/public/login", bytes.NewReader(jsonByte))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/public/login", lc.Login)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
			Convey("and the response should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "Invalid credentials")
			})
		})
	})

	Convey("given the 3rd POST http request for /public/login with wrong email", t, func() {
		jsonByte, _ := json.Marshal(domain.LoginRequest{
			"WrongEmail@mail.ru",
			"qwerty123",
		})
		req := httptest.NewRequest("POST", "/public/login", bytes.NewReader(jsonByte))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/public/login", lc.Login)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
			Convey("and the response should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "User with this email not found")
			})
		})
	})
}

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

func TestSignUp(t *testing.T) {
	app := bootstrap.App()
	env := app.Env
	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()
	timeout := time.Duration(env.ContextTimeout) * time.Second
	ur := repository.NewUserRepository(db, domain.UserCollection)
	sc := &controller.SignupController{
		SignupUseCase: usecase.NewSignupUseCase(ur, timeout),
		Env:           env,
	}
	r := chi.NewRouter()

	Convey("given the 1st http request for /public/signup", t, func() {
		jsonByte, _ := json.Marshal(domain.SignupRequest{
			"dobryak",
			"dobryachok123@mail.ru",
			"123123",
		})
		req := httptest.NewRequest("POST", "/public/signup", bytes.NewReader(jsonByte))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/public/signup", sc.Signup)
			r.ServeHTTP(resp, req)
			Convey("then the response should be 201", func() {
				So(resp.Code, ShouldEqual, 201)
			})
		})
	})

	Convey("given the 2nd http request for /public/signup with same email", t, func() {
		jsonByte, _ := json.Marshal(domain.SignupRequest{
			"ZlodeySPochtoyDobryaka",
			"dobryachok123@mail.ru",
			"321321",
		})
		req := httptest.NewRequest("POST", "/public/signup", bytes.NewReader(jsonByte))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/public/signup", sc.Signup)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "User with this email already exists")
			})
		})
	})

	Convey("given the 3rd http request for /public/signup without field", t, func() {
		jsonByte, _ := json.Marshal(domain.SignupRequest{
			"Pustyak",
			"",
			"321321",
		})
		req := httptest.NewRequest("POST", "/public/signup", bytes.NewReader(jsonByte))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/public/signup", sc.Signup)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "Invalid data")
			})
		})
	})

}

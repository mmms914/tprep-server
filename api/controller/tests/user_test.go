package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
	"main/api/controller"
	"main/bootstrap"
	"main/domain"
	"main/repository"
	"main/storage"
	"main/usecase"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	app := bootstrap.App()
	env := app.Env
	s := app.Storage
	us := storage.NewUserStorage(s, domain.UserBucket)
	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()
	timeout := time.Duration(env.ContextTimeout) * time.Second
	ur := repository.NewUserRepository(db, domain.UserCollection)
	uc := &controller.UserController{
		UserUseCase: usecase.NewUserUseCase(ur, us, timeout),
	}
	r := chi.NewRouter()
	userID := "69a9f624-41e9-4818-bbd4-6af79644dbd1" // mocked user ID
	Convey("given the 1st GET http request for /user", t, func() {
		req := httptest.NewRequest("GET", "/user", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Get("/user", uc.Get)
			r.ServeHTTP(resp, req)
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("given the 2nd GET http request for /user with wrong user ID", t, func() {
		req := httptest.NewRequest("GET", "/user", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID+"fake"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Get("/user", uc.Get)
			r.ServeHTTP(resp, req)
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
		})
	})

	Convey("given the 3rd PUT http request for /user", t, func() {
		jsonByte, _ := json.Marshal(domain.User{
			Username: "MyUsernameWasUpdated",
			Email:    "testadmin@mail.ru",
		})
		req := httptest.NewRequest("PUT", "/user", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/user", uc.Update)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
			Convey("and then the response should contain a message", func() {
				So(responseBody, ShouldContainSubstring, "User updated")
			})
		})
	})

	Convey("given a 4th PUT http request for /user with empty field", t, func() {
		jsonByte, _ := json.Marshal(domain.User{
			Username: "",
			Email:    "testadmin@mail.ru",
		})
		req := httptest.NewRequest("PUT", "/user", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/user", uc.Update)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
			Convey("and then the response should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "Invalid data")
			})
		})
	})

}

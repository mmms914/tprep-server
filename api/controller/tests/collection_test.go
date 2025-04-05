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
	"main/usecase"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestCollection(t *testing.T) {
	app := bootstrap.App()
	env := app.Env
	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()
	timeout := time.Duration(env.ContextTimeout) * time.Second

	ur := repository.NewUserRepository(db, domain.UserCollection)
	uuc := usecase.NewUserUseCase(ur, timeout)
	cr := repository.NewCollectionRepository(db, domain.CollectionCollection)

	cc := &controller.CollectionController{
		CollectionUseCase: usecase.NewCollectionUseCase(cr, timeout),
		UserUseCase:       uuc,
	}
	var createdID string
	var localID int
	collID := "47034bef-d4a0-4bd8-aae1-703b2bc079a8" // mocked collection ID
	r := chi.NewRouter()
	Convey("given the 1st POST http request for /collection", t, func() {
		jsonByte, _ := json.Marshal(domain.Collection{
			Name: "MyCollection",
		})
		req := httptest.NewRequest("POST", "/collection", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mocked"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/collection", cc.Create)
			r.ServeHTTP(resp, req)
			responsebody := resp.Body.String()
			var collectionResponse domain.Collection
			err := json.Unmarshal([]byte(responsebody), &collectionResponse)
			So(err, ShouldBeNil)
			createdID = collectionResponse.ID
			Convey("then the response should be 201", func() {
				So(resp.Code, ShouldEqual, 201)
			})
		})
	})

	Convey("given the 2nd POST http request for /collection without name", t, func() {
		jsonByte, _ := json.Marshal(domain.Collection{
			Name: "",
		})
		req := httptest.NewRequest("POST", "/collection", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mocked"))
		resp := httptest.NewRecorder()

		Convey("when the request handled by router", func() {
			r.Post("/collection", cc.Create)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "Invalid name")
			})
		})
	})

	Convey("given the 3rd GET http request for /collection/{id}/", t, func() {
		req := httptest.NewRequest("GET", "/collection/47034bef-d4a0-4bd8-aae1-703b2bc079a8/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Get("/collection/{id}/", cc.Get)
			r.ServeHTTP(resp, req)
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("given a 4th GET http request for /collection/{id}/ by not owner", t, func() {
		req := httptest.NewRequest("GET", "/collection/"+collID+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockery"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Get("/collection/{id}/", cc.Get)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "You are not the owner of this collection")
			})
		})
	})

	Convey("given a 5th GET http request for /collection/{id}/ to a non-existent collection", t, func() {
		req := httptest.NewRequest("GET", "/collection/fakeID/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockery"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Get("/collection/{id}/", cc.Get)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "There is no collection with this ID")
			})
		})
	})

	Convey("given a 6th PUT http request for /collection/{id}/", t, func() {
		jsonByte, _ := json.Marshal(domain.Collection{
			Name:     "MyCollectionWasUpdated",
			IsPublic: true,
		})
		req := httptest.NewRequest("PUT", "/collection/"+createdID+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mocked"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/", cc.Update)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
			Convey("and the response body should contain a message", func() {
				So(responseBody, ShouldContainSubstring, "Collection updated")
			})
		})
	})

	Convey("given a 7th PUT http request for /collection/{id}/ with empty name field", t, func() {
		jsonByte, _ := json.Marshal(domain.Collection{
			Name:     "",
			IsPublic: true,
		})
		req := httptest.NewRequest("PUT", "/collection/"+createdID+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mocked"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/", cc.Update)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "Invalid name")
			})
		})
	})

	Convey("given a 8th PUT http request for /collection/{id}/ with fakeID", t, func() {
		jsonByte, _ := json.Marshal(domain.Collection{
			Name:     "Matanaliz",
			IsPublic: false,
		})
		req := httptest.NewRequest("PUT", "/collection/"+createdID+"kakyazdesokazalsa"+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mocked"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/", cc.Update)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "There is no collection with this ID")
			})
		})
	})

	Convey("given a 9th PUT http request for /collection/{id}/ with wrong userID", t, func() {
		jsonByte, _ := json.Marshal(domain.Collection{
			Name:     "Matanaliz",
			IsPublic: false,
		})
		req := httptest.NewRequest("PUT", "/collection/"+createdID+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "fakeID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/", cc.Update)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "You are not the owner of this collection")
			})
		})
	})

	Convey("given a 10th DELETE http request for /collection/{id}/ with wrong user", t, func() {
		req := httptest.NewRequest("DELETE", "/collection/"+createdID+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedFake"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Delete("/collection/{id}/", cc.Delete)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "You are not the owner of this collection")
			})
		})
	})

	Convey("given a 11th DELETE http request for /collection/{id}/ with wrong id", t, func() {
		req := httptest.NewRequest("DELETE", "/collection/"+createdID+"abbabab"+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mocked"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Delete("/collection/{id}/", cc.Delete)
			r.ServeHTTP(resp, req)
			responsebody := resp.Body.String()
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
			Convey("and the response body should contain an error message", func() {
				So(responsebody, ShouldContainSubstring, "There is no collection with this ID")
			})
		})
	})

	Convey("given a 12th DELETE http request for /collection/{id}/", t, func() {
		req := httptest.NewRequest("DELETE", "/collection/"+createdID+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mocked"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Delete("/collection/{id}/", cc.Delete)
			r.ServeHTTP(resp, req)
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("given a 13th Card-POST http request for /collection/{id}/card/", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "How are you?",
			Answer:   "Im fine",
		})
		req := httptest.NewRequest("POST", "/collection/"+collID+"/card/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/collection/{id}/card/", cc.CreateCard)
			r.ServeHTTP(resp, req)
			responsebody := resp.Body.String()
			var cardResponse domain.Card
			err := json.Unmarshal([]byte(responsebody), &cardResponse)
			So(err, ShouldBeNil)
			localID = cardResponse.LocalID
			Convey("then the response should be 201", func() {
				So(resp.Code, ShouldEqual, 201)
			})
		})
	})

	Convey("given a 14th Card-POST http request for /collection/{id}/card/ with empty field", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "",
			Answer:   "Im fine??",
		})
		req := httptest.NewRequest("POST", "/collection/"+collID+"/card/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/collection/{id}/card/", cc.CreateCard)
			r.ServeHTTP(resp, req)
			responsebody := resp.Body.String()
			Convey("then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
			Convey("and the response should contain an error message", func() {
				So(responsebody, ShouldContainSubstring, "Invalid body data")
			})
		})
	})

	Convey("given a 15th Card-POST http request for /collection/{id}/card/ with invalid collection ID", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "Are you OK?",
			Answer:   "Im fine!",
		})
		req := httptest.NewRequest("POST", "/collection/fakeID/card/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/collection/{id}/card/", cc.CreateCard)
			r.ServeHTTP(resp, req)
			responsebody := resp.Body.String()
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
			Convey("and the response should contain an error message", func() {
				So(responsebody, ShouldContainSubstring, "There is no collection with this ID")
			})
		})
	})

	Convey("given a 16th Card-POST http request for /collection/{id}/card/ with invalid user ID", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "Are you OK?",
			Answer:   "Im fine!",
		})
		req := httptest.NewRequest("POST", "/collection/"+collID+"/card/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "fakeID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Post("/collection/{id}/card/", cc.CreateCard)
			r.ServeHTTP(resp, req)
			responsebody := resp.Body.String()
			Convey("then the response should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
			Convey("and the response should contain an error message", func() {
				So(responsebody, ShouldContainSubstring, "You are not the owner of this collection")
			})
		})
	})

	Convey("given a 17th Card-PUT http request for /collection/{id}/card/{cardID}", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "Are you OK?",
			Answer:   "No, Im coding now!",
		})
		req := httptest.NewRequest("PUT", "/collection/"+collID+"/card/"+strconv.Itoa(localID)+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/card/{cardID}/", cc.UpdateCard)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
			Convey("and the response body should contain a message", func() {
				So(responseBody, ShouldContainSubstring, "Card updated")
			})
		})
	})

	Convey("given a 18th Card-PUT http request for /collection/{id}/card/{cardID} with empty field", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "",
			Answer:   "where is my question?",
		})
		req := httptest.NewRequest("PUT", "/collection/"+collID+"/card/"+strconv.Itoa(localID)+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/card/{cardID}/", cc.UpdateCard)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "Invalid body data")
			})
		})
	})

	Convey("given a 19th Card-PUT http request for /collection/{id}/card/{cardID} with wrong card ID", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "Im here!",
			Answer:   "where is my question?",
		})
		req := httptest.NewRequest("PUT", "/collection/"+collID+"/card/"+strconv.Itoa(localID+123)+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/card/{cardID}/", cc.UpdateCard)
			r.ServeHTTP(resp, req)
			Convey("then the response should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("given a 20th Card-PUT http request for /collection/{id}/card/{cardID} with wrong collection ID", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "Im here",
			Answer:   "where is my question?",
		})
		req := httptest.NewRequest("PUT", "/collection/fakeID/card/"+strconv.Itoa(localID)+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/card/{cardID}/", cc.UpdateCard)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "There is no collection with this ID")
			})
		})
	})

	Convey("given a 21th Card-PUT http request for /collection/{id}/card/{cardID} with wrong user ID", t, func() {
		jsonByte, _ := json.Marshal(domain.Card{
			Question: "Im here!",
			Answer:   "where is my question?",
		})
		req := httptest.NewRequest("PUT", "/collection/"+collID+"/card/"+strconv.Itoa(localID)+"/", bytes.NewReader(jsonByte))
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "fakeID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Put("/collection/{id}/card/{cardID}/", cc.UpdateCard)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "You are not the owner of this collection")
			})
		})
	})

	Convey("given a 22th Card-DELETE http request for /collection/{id}/card/{cardID} with wrong card ID", t, func() {
		req := httptest.NewRequest("DELETE", "/collection/"+collID+"/card/"+strconv.Itoa(localID+123)+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Delete("/collection/{id}/card/{cardID}/", cc.DeleteCard)
			r.ServeHTTP(resp, req)
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("given a 23th Card-DELETE http request for /collection/{id}/card/{cardID} with wrong collection ID", t, func() {
		req := httptest.NewRequest("DELETE", "/collection/fakeID/card/"+strconv.Itoa(localID)+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Delete("/collection/{id}/card/{cardID}/", cc.DeleteCard)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "There is no collection with this ID")
			})
		})
	})

	Convey("given a 24th Card-DELETE http request for /collection/{id}/card/{cardID} with wrong user ID", t, func() {
		req := httptest.NewRequest("DELETE", "/collection/"+collID+"/card/"+strconv.Itoa(localID)+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "fakeID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Delete("/collection/{id}/card/{cardID}/", cc.DeleteCard)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
			Convey("and the response body should contain an error message", func() {
				So(responseBody, ShouldContainSubstring, "You are not the owner of this collection")
			})
		})
	})

	Convey("given a 25th Card-DELETE http request for /collection/{id}/card/{cardID}", t, func() {
		req := httptest.NewRequest("DELETE", "/collection/"+collID+"/card/"+strconv.Itoa(localID)+"/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "mockedID"))
		resp := httptest.NewRecorder()
		Convey("when the request handled by router", func() {
			r.Delete("/collection/{id}/card/{cardID}/", cc.DeleteCard)
			r.ServeHTTP(resp, req)
			responseBody := resp.Body.String()
			Convey("then the response should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
			Convey("and the response body should contain a message", func() {
				So(responseBody, ShouldContainSubstring, "Card deleted successfully")
			})
		})
	})
}

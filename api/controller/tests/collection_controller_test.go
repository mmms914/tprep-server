package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"main/api/controller"
	"main/domain"
	mocks "main/mocks/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCollectionController_Create_Success(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	mockUserUseCase := new(mocks.UserUseCase)
	controller := &controller.CollectionController{
		CollectionUseCase: mockCollUseCase,
		UserUseCase:       mockUserUseCase,
	}
	userID := "user-id"
	collection := domain.Collection{
		Name:   "New Collection",
		Author: userID,
	}
	createdID := "created-id"
	fullCollection := domain.Collection{
		ID:       createdID,
		Name:     collection.Name,
		IsPublic: false,
		Author:   userID,
	}

	mockCollUseCase.On("Create", mock.Anything, &collection, userID).
		Return(createdID, nil)
	mockCollUseCase.On("GetByID", mock.Anything, createdID).
		Return(fullCollection, nil)

	bodyJSON, _ := json.Marshal(collection)
	req := httptest.NewRequest(http.MethodPost, "/collection", strings.NewReader(string(bodyJSON)))
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
	rr := httptest.NewRecorder()

	controller.Create(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	var response domain.CollectionInfo
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, createdID, response.ID)
	assert.Equal(t, collection.Name, response.Name)
	mockCollUseCase.AssertExpectations(t)
	mockUserUseCase.AssertExpectations(t)
}

func TestCollectionController_Create_EmptyName(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	controller := &controller.CollectionController{
		CollectionUseCase: mockCollUseCase,
	}
	invalidCollection := domain.Collection{
		Name: "",
	}
	bodyJSON, _ := json.Marshal(invalidCollection)

	req := httptest.NewRequest(http.MethodPost, "/collection", strings.NewReader(string(bodyJSON)))
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "user-id"))
	rr := httptest.NewRecorder()
	controller.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	assert.Contains(t, rr.Body.String(), "Invalid name")
}

func TestCollectionController_Update_Success(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	controller := &controller.CollectionController{
		CollectionUseCase: mockCollUseCase,
	}
	userID := "user-id"
	collID := "collection-id"
	updCollection := domain.Collection{Name: "UPDname"}
	existingCollection := domain.Collection{
		ID:     collID,
		Name:   "Old Name",
		Author: userID,
	}

	mockCollUseCase.On("GetByID", mock.Anything, collID).Return(existingCollection, nil)
	mockCollUseCase.On("PutByID", mock.Anything, collID, &updCollection).Return(nil)

	bodyJSON, _ := json.Marshal(updCollection)
	req := httptest.NewRequest(http.MethodPut, "/collection/{id}", strings.NewReader(string(bodyJSON)))
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", collID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	rr := httptest.NewRecorder()
	controller.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	var resp domain.SuccessResponse
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, "Collection updated", resp.Message)

	mockCollUseCase.AssertExpectations(t)
}

func TestCollectionController_Update_EmptyName(t *testing.T) {
	controller := &controller.CollectionController{}
	collection := domain.Collection{Name: ""}
	bodyJSON, _ := json.Marshal(collection)
	req := httptest.NewRequest(http.MethodPut, "/collection/{id}", strings.NewReader(string(bodyJSON)))
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "user-id"))
	rr := httptest.NewRecorder()
	controller.Update(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	var resp domain.SuccessResponse
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Contains(t, "Invalid name", resp.Message)
}

func TestCollectionController_Update_NotFound(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	controller := &controller.CollectionController{CollectionUseCase: mockCollUseCase}

	collection := domain.Collection{Name: "Updated"}
	bodyJSON, _ := json.Marshal(collection)
	collID := "coll-id"
	req := httptest.NewRequest(http.MethodPut, "/collections/{id}", bytes.NewReader(bodyJSON))
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "user-id"))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", collID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	rr := httptest.NewRecorder()

	mockCollUseCase.On("GetByID", mock.Anything, collID).Return(domain.Collection{}, errors.New("Not found"))
	controller.Update(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
	var resp domain.SuccessResponse
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Contains(t, "There is no collection with this ID", resp.Message)

}

func TestCollectionController_Update_NotAuthor(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	controller := &controller.CollectionController{
		CollectionUseCase: mockCollUseCase,
	}
	collection := domain.Collection{Name: "NewName", Author: "other-author"}
	bodyJSON, _ := json.Marshal(collection)
	collID := "collID"

	mockCollUseCase.On("GetByID", mock.Anything, collID).Return(collection, nil)

	req := httptest.NewRequest(http.MethodPut, "/collection/{id}", strings.NewReader(string(bodyJSON)))
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", "userID"))
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", collID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	rr := httptest.NewRecorder()

	controller.Update(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Result().StatusCode)
	var resp domain.SuccessResponse
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Contains(t, "You are not the owner of this collection", resp.Message)
}

func TestCollectionController_Delete_Success(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	mockUserUseCase := new(mocks.UserUseCase)
	controller := &controller.CollectionController{
		CollectionUseCase: mockCollUseCase,
		UserUseCase:       mockUserUseCase,
	}
	userID := "user-id"
	collID := "coll-id"
	collection := domain.Collection{
		ID:     collID,
		Author: userID,
		Cards: []domain.Card{
			{LocalID: 1, Attachment: "pic.jpg"},
		},
	}
	mockCollUseCase.On("GetByID", mock.Anything, collID).Return(collection, nil)
	mockCollUseCase.On("DeleteByID", mock.Anything, collID, userID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/collection/{id}", nil)
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", collID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	rr := httptest.NewRecorder()
	controller.Delete(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resp domain.SuccessResponse
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&resp))
	assert.Contains(t, resp.Message, "Collection deleted successfully")

	mockCollUseCase.AssertExpectations(t)
	mockUserUseCase.AssertExpectations(t)
}

func TestCollectionController_Delete_NotFound(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	controller := &controller.CollectionController{
		CollectionUseCase: mockCollUseCase,
	}

	userID := "user-id"
	collID := "wrong-id"

	mockCollUseCase.On("GetByID", mock.Anything, collID).Return(domain.Collection{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/collections/{id}", nil)
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", collID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	rr := httptest.NewRecorder()
	controller.Delete(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	var resp domain.SuccessResponse
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&resp))
	assert.Contains(t, resp.Message, "There is no collection with this ID")
}

func TestCollectionController_Delete_NotAuthor(t *testing.T) {
	mockCollUseCase := new(mocks.CollectionUseCase)
	controller := &controller.CollectionController{
		CollectionUseCase: mockCollUseCase,
	}

	userID := "attacker-id"
	collID := "collection-id"
	collection := domain.Collection{
		ID:     collID,
		Author: "real-owner-id",
	}

	mockCollUseCase.On("GetByID", mock.Anything, collID).Return(collection, nil)

	req := httptest.NewRequest(http.MethodDelete, "/collections/{id}", nil)
	req = req.WithContext(context.WithValue(req.Context(), "x-user-id", userID))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", collID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	rr := httptest.NewRecorder()
	controller.Delete(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusForbidden, res.StatusCode)
	var resp domain.SuccessResponse
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&resp))
	assert.Contains(t, resp.Message, "You are not the owner of this collection")
}

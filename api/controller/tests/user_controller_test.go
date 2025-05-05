package tests

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"io"
	"main/api/controller"
	"main/bootstrap"
	"main/domain"
	mocks "main/mocks/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUserController_Get_OwnProfile(t *testing.T) {
	mockUserUseCase := new(mocks.UserUseCase)
	mockCollectionUseCase := new(mocks.CollectionUseCase)
	mockHistoryUseCase := new(mocks.HistoryUseCase)

	userID := "test-user-id"

	mockUser := domain.User{
		ID:          userID,
		Username:    "johndoe",
		Email:       "john@example.com",
		HasPicture:  true,
		Collections: []string{"col1", "col2"},
		Statistics:  domain.UserStatistics{},
		Favourite:   []string{"fav1"},
	}

	// настройка ожиданий
	mockUserUseCase.On("GetByID", mock.Anything, userID).Return(mockUser, nil)

	uc := &controller.UserController{
		UserUseCase:       mockUserUseCase,
		CollectionUseCase: mockCollectionUseCase,
		HistoryUseCase:    mockHistoryUseCase,
	}

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	ctx := context.WithValue(req.Context(), "x-user-id", userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	uc.Get(rr, req)

	res := rr.Result()

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)

	var userInfo domain.UserInfo
	err := easyjson.Unmarshal(body, &userInfo)
	assert.NoError(t, err)
	assert.Equal(t, userID, userInfo.ID)
	assert.Equal(t, "johndoe", userInfo.Username)
	assert.Equal(t, "john@example.com", userInfo.Email)

	mockUserUseCase.AssertExpectations(t)
}

func TestUserController_Get_InvalidUserID(t *testing.T) {
	mockUserUseCase := new(mocks.UserUseCase)
	mockCollectionUseCase := new(mocks.CollectionUseCase)
	mockHistoryUseCase := new(mocks.HistoryUseCase)

	invalidUserID := "invalid-user-id"
	expectedErr := errors.New("User not found")

	mockUserUseCase.On("GetByID", mock.Anything, invalidUserID).Return(domain.User{}, expectedErr)

	uc := &controller.UserController{
		UserUseCase:       mockUserUseCase,
		CollectionUseCase: mockCollectionUseCase,
		HistoryUseCase:    mockHistoryUseCase,
	}
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	ctx := context.WithValue(req.Context(), "x-user-id", invalidUserID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	uc.Get(rr, req)

	res := rr.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockUserUseCase.AssertExpectations(t)
}

func TestUserController_SignUp_Success(t *testing.T) {
	mockUseCase := new(mocks.SignupUseCase)

	controller := &controller.SignupController{
		SignupUseCase: mockUseCase,
		Env: &bootstrap.Env{
			AccessTokenSecret:      "access-secret",
			RefreshTokenSecret:     "refresh-secret",
			AccessTokenExpiryHour:  1,
			RefreshTokenExpiryHour: 24,
		},
	}

	reqBody := domain.SignupRequest{
		Username: "NewUser",
		Email:    "NewUser@example.com",
		Password: "qwerty123",
	}
	bodyJSON, _ := easyjson.Marshal(reqBody)

	mockUseCase.On("GetUserByEmail", mock.Anything, reqBody.Email).
		Return(domain.User{}, errors.New("User with this email already exists"))

	mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
		Return("new-user-id", nil)

	mockUseCase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), "access-secret", 1).
		Return("access-token", time.Now().Add(time.Hour), nil)

	mockUseCase.On("CreateRefreshToken", mock.AnythingOfType("*domain.User"), "refresh-secret", 24).
		Return("refresh-token", time.Now().Add(24*time.Hour), nil)

	req := httptest.NewRequest(http.MethodPost, "/public/signup", strings.NewReader(string(bodyJSON)))
	rr := httptest.NewRecorder()

	controller.Signup(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.NotEmpty(t, res.Header.Get("X-Access-Expires-After"))
	assert.NotEmpty(t, res.Header.Get("X-Refresh-Expires-After"))

	var resp domain.SignupResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "new-user-id", resp.UserID)
	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)

	mockUseCase.AssertExpectations(t)
}

func TestUserController_SignUp_InvalidJSON(t *testing.T) {
	controller := &controller.SignupController{}

	req := httptest.NewRequest(http.MethodPost, "/public/signup", strings.NewReader("ivalid-json"))
	rr := httptest.NewRecorder()
	controller.Signup(rr, req)

	res := rr.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserController_SignUp_MissingField(t *testing.T) {
	controller := &controller.SignupController{}
	reqBody := domain.SignupRequest{
		Username: "",
		Email:    "",
		Password: "qwerty123",
	}
	bodyJSON, _ := easyjson.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/public/signup", strings.NewReader(string(bodyJSON)))
	rr := httptest.NewRecorder()
	controller.Signup(rr, req)

	res := rr.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserController_SignUp_UserAlreadyExists(t *testing.T) {
	mockUseCase := new(mocks.SignupUseCase)
	controller := &controller.SignupController{
		SignupUseCase: mockUseCase,
	}
	reqBody := domain.SignupRequest{
		Username: "existing",
		Email:    "existing@example.com",
		Password: "qwerty123",
	}
	bodyJSON, _ := easyjson.Marshal(reqBody)
	mockUseCase.On("GetUserByEmail", mock.Anything, reqBody.Email).Return(domain.User{}, nil)
	req := httptest.NewRequest(http.MethodPost, "/public/signup", strings.NewReader(string(bodyJSON)))
	rr := httptest.NewRecorder()
	controller.Signup(rr, req)
	res := rr.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Contains(t, rr.Body.String(), "User with this email already exists")
}

func TestUserController_Login_Success(t *testing.T) {
	mockUseCase := new(mocks.LoginUseCase)
	controller := &controller.LoginController{
		LoginUseCase: mockUseCase,
		Env: &bootstrap.Env{
			AccessTokenSecret:      "access-secret",
			RefreshTokenSecret:     "refresh-secret",
			AccessTokenExpiryHour:  1,
			RefreshTokenExpiryHour: 24,
		},
	}

	password := "qwerty123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := domain.User{
		ID:       "user-id",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	reqBody := domain.LoginRequest{
		Email:    user.Email,
		Password: password,
	}
	bodyJSON, _ := easyjson.Marshal(reqBody)

	mockUseCase.On("GetUserByEmail", mock.Anything, reqBody.Email).
		Return(user, nil)
	mockUseCase.On("CreateAccessToken", &user, "access-secret", 1).
		Return("access-token", time.Now().Add(time.Hour), nil)
	mockUseCase.On("CreateRefreshToken", &user, "refresh-secret", 24).
		Return("refresh-token", time.Now().Add(24*time.Hour), nil)

	req := httptest.NewRequest(http.MethodPost, "/public/login", strings.NewReader(string(bodyJSON)))
	rr := httptest.NewRecorder()
	controller.Login(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.NotEmpty(t, res.Header.Get("X-Access-Expires-After"))
	assert.NotEmpty(t, res.Header.Get("X-Refresh-Expires-After"))

	var resp domain.LoginResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)
	mockUseCase.AssertExpectations(t)
}

func TestUserController_Login_InvalidJSON(t *testing.T) {
	controller := &controller.LoginController{}

	req := httptest.NewRequest(http.MethodPost, "/public/login", strings.NewReader("invalid-json"))
	rr := httptest.NewRecorder()
	controller.Login(rr, req)

	res := rr.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserController_Login_UserNotFound(t *testing.T) {
	mockUseCase := new(mocks.LoginUseCase)
	controller := &controller.LoginController{
		LoginUseCase: mockUseCase,
	}
	reqBody := domain.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password",
	}
	bodyJSON, _ := easyjson.Marshal(reqBody)
	mockUseCase.On("GetUserByEmail", mock.Anything, "notfound@example.com").
		Return(domain.User{}, errors.New("User with this email not found"))
	req := httptest.NewRequest(http.MethodPost, "/public/login", strings.NewReader(string(bodyJSON)))
	rr := httptest.NewRecorder()
	controller.Login(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestUserController_Login_WrongPassword(t *testing.T) {
	mockUseCase := new(mocks.LoginUseCase)
	controller := &controller.LoginController{
		LoginUseCase: mockUseCase,
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	user := domain.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	reqBody := domain.LoginRequest{
		Email:    user.Email,
		Password: "wrong-password",
	}
	bodyJSON, _ := easyjson.Marshal(reqBody)
	mockUseCase.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil)
	req := httptest.NewRequest(http.MethodPost, "/public/login", strings.NewReader(string(bodyJSON)))
	rr := httptest.NewRecorder()
	controller.Login(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserController_Update_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	controller := &controller.UserController{
		UserUseCase: mockUseCase,
	}
	userID := "user-id"
	updatedUser := domain.User{
		Username: "UPDusername",
		Email:    "UPDemail@example.com",
	}
	mockUseCase.On("PutByID", mock.Anything, userID, &updatedUser).
		Return(nil)
	bodyJSON, _ := easyjson.Marshal(updatedUser)
	req := httptest.NewRequest(http.MethodPut, "/user/{id}", strings.NewReader(string(bodyJSON)))
	ctx := context.WithValue(req.Context(), "x-user-id", userID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	controller.Update(rr, req)
	res := rr.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	var response domain.SuccessResponse
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "User updated", response.Message)
	mockUseCase.AssertExpectations(t)
}

func TestUserController_Update_EmptyFields(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	controller := &controller.UserController{
		UserUseCase: mockUseCase,
	}

	testCases := []struct {
		name string
		user domain.User
	}{
		{
			name: "Empty Username",
			user: domain.User{
				Username: "",
				Email:    "valid@example.com",
			},
		},
		{
			name: "Empty Email",
			user: domain.User{
				Username: "valid_username",
				Email:    "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bodyJSON, _ := easyjson.Marshal(tc.user)
			req := httptest.NewRequest(http.MethodPut, "/user/{id}", strings.NewReader(string(bodyJSON)))
			rr := httptest.NewRecorder()
			controller.Update(rr, req)
			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			assert.Contains(t, rr.Body.String(), "Invalid data")
		})
	}
}

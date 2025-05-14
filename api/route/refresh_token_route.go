package route

import (
	"main/api/controller"
	"main/bootstrap"
	"main/database"
	"main/domain"
	"main/repository"
	"main/usecase"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewRefreshTokenRouter(env *bootstrap.Env, timeout time.Duration, db database.Database, r chi.Router) {
	ur := repository.NewUserRepository(db, domain.UserCollection)
	rtc := &controller.RefreshTokenController{
		RefreshTokenUseCase: usecase.NewRefreshTokenUseCase(ur, timeout),
		Env:                 env,
	}

	r.Post("/public/refreshToken", rtc.RefreshToken)
}

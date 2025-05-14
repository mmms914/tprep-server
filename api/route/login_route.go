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

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db database.Database, r chi.Router) {
	ur := repository.NewUserRepository(db, domain.UserCollection)
	lc := &controller.LoginController{
		LoginUseCase: usecase.NewLoginUseCase(ur, timeout),
		Env:          env,
	}

	r.Post("/public/login", lc.Login)
}

package route

import (
	"github.com/go-chi/chi/v5"
	"main/api/controller"
	"main/bootstrap"
	"main/database"
	"main/domain"
	"main/repository"
	"main/usecase"
	"time"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db database.Database, r chi.Router) {
	ur := repository.NewUserRepository(db, domain.UserCollection)
	uc := &controller.UserController{
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
	}
	r.Get("/user", uc.Get)
	r.Put("/user", uc.Update)
}

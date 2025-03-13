package route

import (
	"github.com/go-chi/chi/v5"
	"main/api/controller"
	"main/bootstrap"
	"main/database"
	"main/domain"
	"main/repository"
	"main/storage"
	"main/usecase"
	"time"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db database.Database, s storage.Client, r chi.Router) {
	ur := repository.NewUserRepository(db, domain.UserCollection)
	us := storage.NewUserStorage(s, domain.UserBucket)

	cr := repository.NewCollectionRepository(db, domain.CollectionCollection)

	uc := &controller.UserController{
		UserUseCase:       usecase.NewUserUseCase(ur, us, timeout),
		CollectionUseCase: usecase.NewCollectionUseCase(cr, timeout),
	}
	r.Route("/user", func(r chi.Router) {
		r.Get("/", uc.Get)
		r.Put("/", uc.Update)
		r.Route("/picture", func(r chi.Router) {
			r.Get("/", uc.GetProfilePicture)
			r.Put("/", uc.UploadProfilePicture)
			r.Delete("/", uc.RemoveProfilePicture)
		})
	})
}

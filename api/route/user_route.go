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
	cs := storage.NewCollectionStorage(s, domain.CollectionBucket)

	uhr := repository.NewUserHistoryRepository(db, domain.UserHistoryCollection)
	chr := repository.NewCollectionHistoryRepository(db, domain.CollectionHistoryCollection)

	cr := repository.NewCollectionRepository(db, domain.CollectionCollection)

	uc := &controller.UserController{
		UserUseCase:       usecase.NewUserUseCase(ur, us, timeout),
		CollectionUseCase: usecase.NewCollectionUseCase(cr, cs, ur, timeout),
		HistoryUseCase:    usecase.NewHistoryUseCase(uhr, chr, cr, ur, timeout),
	}
	r.Route("/user", func(r chi.Router) {
		r.Get("/", uc.Get)
		r.Put("/", uc.Update)
		r.Route("/picture", func(r chi.Router) {
			r.Get("/", uc.GetProfilePicture)
			r.Put("/", uc.UploadProfilePicture)
			r.Delete("/", uc.RemoveProfilePicture)
		})
		r.Route("/history", func(r chi.Router) {
			r.Get("/", uc.GetHistory)
		})
	})
}

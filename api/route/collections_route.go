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

func NewCollectionRouter(env *bootstrap.Env, timeout time.Duration, db database.Database, s storage.Client, r chi.Router) {
	ur := repository.NewUserRepository(db, domain.UserCollection)
	us := storage.NewUserStorage(s, domain.UserBucket)

	uhr := repository.NewUserHistoryRepository(db, domain.UserHistoryCollection)
	chr := repository.NewCollectionHistoryRepository(db, domain.CollectionHistoryCollection)

	cr := repository.NewCollectionRepository(db, domain.CollectionCollection)
	cc := &controller.CollectionController{
		CollectionUseCase: usecase.NewCollectionUseCase(cr, timeout),
		UserUseCase:       usecase.NewUserUseCase(ur, us, timeout),
		HistoryUseCase:    usecase.NewHistoryUseCase(uhr, chr, cr, ur, timeout),
	}
	r.Route("/collection", func(r chi.Router) {
		r.Post("/", cc.Create)
		r.Get("/search", cc.Search)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", cc.Get)
			r.Put("/", cc.Update)
			r.Delete("/", cc.Delete)
			r.Put("/like", cc.AddLike)
			r.Put("/unlike", cc.RemoveLike)
			r.Route("/card", func(r chi.Router) {
				r.Post("/", cc.CreateCard)
				r.Route("/{cardID}", func(r chi.Router) {
					r.Put("/", cc.UpdateCard)
					r.Delete("/", cc.DeleteCard)
				})
			})
		})
		r.Route("/training", func(r chi.Router) {
			r.Post("/", cc.AddTraining)
		})
	})
}

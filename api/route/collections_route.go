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

func NewCollectionRouter(env *bootstrap.Env, timeout time.Duration, db database.Database, r chi.Router) {
	ur := repository.NewUserRepository(db, domain.UserCollection)
	uuc := usecase.NewUserUseCase(ur, timeout)

	cr := repository.NewCollectionRepository(db, domain.CollectionCollection)
	cc := &controller.CollectionController{
		CollectionUseCase: usecase.NewCollectionUseCase(cr, timeout),
		UserUseCase:       uuc,
	}
	r.Route("/collection", func(r chi.Router) {
		r.Post("/", cc.Create)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", cc.Get)
			r.Put("/", cc.Update)
			r.Delete("/", cc.Delete)
			r.Route("/card", func(r chi.Router) {
				r.Post("/", cc.CreateCard)
				r.Route("/{cardID}", func(r chi.Router) {
					r.Put("/", cc.UpdateCard)
					r.Delete("/", cc.DeleteCard)
				})
			})
		})
	})
}

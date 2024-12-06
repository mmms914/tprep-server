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
	cr := repository.NewCollectionRepository(db, domain.CollectionCollection)
	cc := &controller.CollectionController{
		CollectionUseCase: usecase.NewCollectionUseCase(cr, timeout),
	}

	r.Get("/collection/{id}", cc.Get)
	r.Post("/collection", cc.Create)
	r.Put("/collection/{id}", cc.Update)
	r.Delete("/collection/{id}", cc.Delete)

	r.Post("/collection/{id}/card", cc.CreateCard)
	r.Put("/collection/{id}/card/{cardID}", cc.UpdateCard)
	r.Delete("/collection/{id}/card/{cardID}", cc.DeleteCard)
}

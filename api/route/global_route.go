package route

import (
	"github.com/go-chi/chi/v5"
	"main/api/controller"
	"main/bootstrap"
	"main/usecase"
	"time"
)

func NewGlobalRouter(env *bootstrap.Env, timeout time.Duration, r chi.Router) {
	gc := &controller.GlobalController{
		GlobalUseCase: usecase.NewGlobalUseCase(timeout),
		Env:           env,
	}

	r.Get("/global/getTrainingPlan", gc.GetTrainingPlan)
}

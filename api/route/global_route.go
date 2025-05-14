package route

import (
	"main/api/controller"
	"main/bootstrap"
	"main/usecase"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewGlobalRouter(env *bootstrap.Env, timeout time.Duration, r chi.Router) {
	gc := &controller.GlobalController{
		GlobalUseCase: usecase.NewGlobalUseCase(timeout),
		Env:           env,
	}

	r.Get("/global/getTrainingPlan", gc.GetTrainingPlan)
	r.Post("/global/addMetrics", gc.AddMetrics)
}

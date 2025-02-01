package usecase

import (
	"main/domain"
	"main/internal"
	"time"
)

type globalUseCase struct {
	contextTimeout time.Duration
}

func NewGlobalUseCase(timeout time.Duration) domain.GlobalUseCase {
	return &globalUseCase{
		contextTimeout: timeout,
	}
}

func (gu *globalUseCase) GetTrainingPlan(start int, end int, preferredTime int) []int {
	return internal.CalculateTrainingPlan(start, end, preferredTime)
}

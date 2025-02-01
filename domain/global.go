package domain

type PlanResponse struct {
	Count int
	Items []int
}

type GlobalUseCase interface {
	GetTrainingPlan(start int, end int, preferredTime int) []int
}

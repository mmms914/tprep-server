package domain

type PlanResponse struct {
	Count int   `bson:"count" json:"count"`
	Items []int `bson:"items" json:"items"`
}

type GlobalUseCase interface {
	GetTrainingPlan(start int, end int, preferredTime int) []int
}

package domain

type PlanResponse struct {
	Count int   `bson:"count" json:"count"`
	Items []int `bson:"items" json:"items"`
}

type MetricsRequest struct {
	FilterClicks     int `json:"favourite_filter_button"`
	ProfileClicks    int `json:"favourite_profile_button"`
	LastInAppTime    int `json:"last_in_app_time"`
	SumTrainingsTime int `json:"sum_trainings_time"`
	TrainingsCount   int `json:"trainings_count"`
}

type GlobalUseCase interface {
	GetTrainingPlan(start int, end int, preferredTime int) []int
}

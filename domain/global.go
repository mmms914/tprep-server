package domain

type PlanResponse struct {
	Count int   `bson:"count" json:"count"`
	Items []int `bson:"items" json:"items"`
}

type FavouriteButtonRequest struct {
	FilterClicks  int `json:"favourite_filter_button"`
	ProfileClicks int `json:"favourite_profile_button"`
}

type GlobalUseCase interface {
	GetTrainingPlan(start int, end int, preferredTime int) []int
	TrackFavouriteButtons(filterClicks int, profileClicks int) error
}

package controller

import (
	"encoding/json"
	"main/api/middleware"
	"main/bootstrap"
	"main/domain"
	"net/http"
	"strconv"
)

type GlobalController struct {
	GlobalUseCase domain.GlobalUseCase
	Env           *bootstrap.Env
}

func (gc *GlobalController) GetTrainingPlan(w http.ResponseWriter, r *http.Request) {
	startDate, err := strconv.Atoi(r.URL.Query().Get("start_date"))
	if err != nil {
		http.Error(w, jsonError("invalid start date"), http.StatusBadRequest)
		return
	}

	endDate, err := strconv.Atoi(r.URL.Query().Get("end_date"))
	if err != nil {
		http.Error(w, jsonError("invalid end date"), http.StatusBadRequest)
		return
	}

	preferredTime, err := strconv.Atoi(r.URL.Query().Get("preferred_time"))
	if err != nil || preferredTime < 0 || preferredTime > 86399 {
		http.Error(w, jsonError("invalid preferred time"), http.StatusBadRequest)
		return
	}

	if endDate-startDate < 24*3600 {
		http.Error(w, jsonError("difference between the beginning and the end should be at least a day"), http.StatusBadRequest)
		return
	}

	planItems := gc.GlobalUseCase.GetTrainingPlan(startDate, endDate, preferredTime)
	planResponse := domain.PlanResponse{
		Count: len(planItems),
		Items: planItems,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(planResponse)
}
func (gc *GlobalController) TrackFavouriteButtons(w http.ResponseWriter, r *http.Request) {
	var req domain.FavouriteButtonRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, jsonError("Invalid data"), http.StatusBadRequest)
		return
	}
	if req.FilterClicks < 0 || req.ProfileClicks < 0 {
		http.Error(w, jsonError("The quantity cannot be negative"), http.StatusBadRequest)
		return
	}
	middleware.FavouriteButtonClicks.WithLabelValues("filter_favourite_button").Add(float64(req.FilterClicks))
	middleware.FavouriteButtonClicks.WithLabelValues("profile_favourite_button").Add(float64(req.ProfileClicks))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

package internal

func CalculateTrainingPlan(start int, finish int, preferredTime int) []int { // time in unix timestamp
	// ax^2 + bx + c = y
	//nolint:mnd // business logic
	b := (32*3600 - 24*3600) / 2 // magic numbers
	//nolint:mnd // business logic
	a := 8*3600 - b
	c := start

	var trainings []int
	x := 0
	lastDate := 0
	for {
		newDate := a*x*x + b*x + c
		if newDate > finish {
			break
		}
		//nolint:mnd // business logic
		if x > 2 {
			newDate = newDate/86400*86400 + preferredTime
			if newDate/86400 == lastDate/86400 {
				x++
				continue
			}
		}
		trainings = append(trainings, newDate)
		lastDate = newDate
		x++
	}

	return trainings
}

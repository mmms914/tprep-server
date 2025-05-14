package internal

import (
	"main/domain"
)

func CalcStatistics(items []domain.HistoryItem) domain.UserStatistics {
	count := 0
	sumPercentage := 0.0
	//nolint:mnd // percent conversion
	for _, item := range items {
		count++
		if item.AllCardsCount != 0 {
			sumPercentage += float64(len(item.CorrectCards)) / float64(item.AllCardsCount) * 100.0
		}
	}
	return domain.UserStatistics{
		TotalTrainings:   count,
		MediumPercentage: int(sumPercentage) / count,
	}
}

package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"main/domain"
	"main/internal"
	"time"
)

type historyUseCase struct {
	userHistoryRepository       domain.UserHistoryRepository
	collectionHistoryRepository domain.CollectionHistoryRepository
	collectionRepository        domain.CollectionRepository
	userRepository              domain.UserRepository
	contextTimeout              time.Duration
}

func NewHistoryUseCase(userHistoryRepository domain.UserHistoryRepository, collectionHistoryRepository domain.CollectionHistoryRepository, collectionRepository domain.CollectionRepository, userRepository domain.UserRepository, timeout time.Duration) domain.HistoryUseCase {
	return &historyUseCase{
		userHistoryRepository:       userHistoryRepository,
		collectionRepository:        collectionRepository,
		collectionHistoryRepository: collectionHistoryRepository,
		userRepository:              userRepository,
		contextTimeout:              timeout,
	}
}

func (hu *historyUseCase) AddTraining(c context.Context, userID string, historyItem domain.HistoryItem) error {
	ctx, cancel := context.WithTimeout(c, hu.contextTimeout)
	defer cancel()

	err := hu.userHistoryRepository.UpdateByID(ctx, userID, historyItem)
	if err != nil {
		return err
	}

	userHistory, err := hu.GetUserHistoryFromTime(ctx, userID, 0)
	if err != nil {
		return err
	}

	newStatistics := internal.CalcStatistics(userHistory)

	update := bson.D{
		{"$set", bson.D{
			{"statistics", newStatistics},
		}},
	}
	err = hu.userRepository.UpdateByID(ctx, userID, update)
	if err != nil {
		return err
	}

	smallHistoryItem := domain.SmallHistoryItem{
		CollectionName: historyItem.CollectionName,
		Time:           historyItem.Time,
		CorrectCards:   historyItem.CorrectCards,
		IncorrectCards: historyItem.IncorrectCards,
		AllCardsCount:  historyItem.AllCardsCount,
	}

	update = bson.D{
		{"$inc", bson.D{{"trainings", 1}}},
	}
	_, err = hu.collectionRepository.UpdateByID(ctx, historyItem.CollectionID, update)
	if err != nil {
		return err
	}

	return hu.collectionHistoryRepository.UpdateByID(ctx, historyItem.CollectionID, smallHistoryItem)
}

func (hu *historyUseCase) GetUserHistoryFromTime(c context.Context, userID string, fromTime int) ([]domain.HistoryItem, error) {
	ctx, cancel := context.WithTimeout(c, hu.contextTimeout)
	defer cancel()

	userHistoryUpdate := make([]domain.HistoryItem, 0)

	allUserHistory, err := hu.userHistoryRepository.GetByID(ctx, userID)
	if err != nil {
		return userHistoryUpdate, err
	}

	for _, item := range allUserHistory.Items {
		if item.Time >= fromTime {
			userHistoryUpdate = append(userHistoryUpdate, item)
		}
	}
	return userHistoryUpdate, nil
}

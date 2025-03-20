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
	userRepository              domain.UserRepository
	contextTimeout              time.Duration
}

func NewHistoryUseCase(userHistoryRepository domain.UserHistoryRepository, collectionHistoryRepository domain.CollectionHistoryRepository, userRepository domain.UserRepository, timeout time.Duration) domain.HistoryUseCase {
	return &historyUseCase{
		userHistoryRepository:       userHistoryRepository,
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

	userHistory, err := hu.GetUserHistory(ctx, userID)
	if err != nil {
		return err
	}

	newStatistics := internal.CalcStatistics(userHistory.Items)

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
	}
	return hu.collectionHistoryRepository.UpdateByID(ctx, historyItem.CollectionID, smallHistoryItem)
}

func (hu *historyUseCase) GetUserHistory(c context.Context, userID string) (domain.UserHistory, error) {
	ctx, cancel := context.WithTimeout(c, hu.contextTimeout)
	defer cancel()

	return hu.userHistoryRepository.GetByID(ctx, userID)
}

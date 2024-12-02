package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"main/domain"
	"time"
)

type collectionUseCase struct {
	collectionRepository domain.CollectionRepository
	contextTimeout       time.Duration
}

func NewCollectionUseCase(collectionRepository domain.CollectionRepository, timeout time.Duration) domain.CollectionUseCase {
	return &collectionUseCase{
		collectionRepository: collectionRepository,
		contextTimeout:       timeout,
	}
}

func (cu *collectionUseCase) Create(c context.Context, collection *domain.Collection) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	return cu.collectionRepository.Create(ctx, collection)
}

func (cu *collectionUseCase) UpdateNameByID(c context.Context, collectionID string, collectionName string) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	update := bson.D{
		{"$set", bson.D{
			{"name", collectionName},
		}},
	}
	return cu.collectionRepository.UpdateByID(ctx, collectionID, update)
}

func (cu *collectionUseCase) UpdateTypeByID(c context.Context, collectionID string, collectionType bool) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	update := bson.D{
		{"$set", bson.D{
			{"is_public", collectionType},
		}},
	}
	return cu.collectionRepository.UpdateByID(ctx, collectionID, update)
}

func (cu *collectionUseCase) DeleteByID(c context.Context, collectionID string) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	return cu.collectionRepository.DeleteByID(ctx, collectionID)
}

func (cu *collectionUseCase) GetByID(c context.Context, collectionID string) (*domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	return cu.collectionRepository.GetByID(ctx, collectionID)
}

func (cu *collectionUseCase) AddCard(c context.Context, collectionID string, card *domain.Card) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$push", bson.D{
			{"cards", card},
		}},
		{"$inc", bson.D{
			{"max_id", 1},
		}},
	}
	return cu.collectionRepository.UpdateByID(ctx, collectionID, update)
}

func (cu *collectionUseCase) DeleteCard(c context.Context, collectionID string, cardLocalID int) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$pull", bson.D{
			{"cards", bson.D{
				{"id", cardLocalID},
			}},
		}},
	}
	return cu.collectionRepository.UpdateByID(ctx, collectionID, update)
}

func (cu *collectionUseCase) UpdateCard(c context.Context, collectionID string, card *domain.Card) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	filter := bson.D{
		{"id", collectionID},
		{"cards.local_id", card.LocalID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"cards.$.question", card.Question},
			{"cards.$.answer", card.Answer},
		}},
	}
	return cu.collectionRepository.Update(ctx, filter, update)
}

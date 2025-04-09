package usecase

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"io"
	"main/database"
	"main/domain"
	"slices"
	"strconv"
	"time"
)

type collectionUseCase struct {
	collectionRepository domain.CollectionRepository
	collectionStorage    domain.CollectionStorage
	userRepository       domain.UserRepository
	contextTimeout       time.Duration
}

func NewCollectionUseCase(collectionRepository domain.CollectionRepository, collectionStorage domain.CollectionStorage, userRepository domain.UserRepository, timeout time.Duration) domain.CollectionUseCase {
	return &collectionUseCase{
		collectionRepository: collectionRepository,
		collectionStorage:    collectionStorage,
		userRepository:       userRepository,
		contextTimeout:       timeout,
	}
}

func (cu *collectionUseCase) Create(c context.Context, collection *domain.Collection) (string, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	collection.Cards = make([]domain.Card, 0)
	return cu.collectionRepository.Create(ctx, collection)
}

func (cu *collectionUseCase) PutByID(c context.Context, collectionID string, collection *domain.Collection) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$set", bson.D{
			{"name", collection.Name},
			{"is_public", collection.IsPublic},
		}},
	}
	res, err := cu.collectionRepository.UpdateByID(ctx, collectionID, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("collection not exists")
	}
	return nil
}

func (cu *collectionUseCase) AddLike(c context.Context, collectionID string) (*domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$inc", bson.D{{"likes", 1}}},
	}
	res, err := cu.collectionRepository.UpdateByID(ctx, collectionID, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, errors.New("collection not exists")
	}
	updated, err := cu.collectionRepository.GetByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (cu *collectionUseCase) RemoveLike(c context.Context, collectionID string) (*domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	current, err := cu.collectionRepository.GetByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	if current.Likes <= 0 {
		return &current, nil
	}

	update := bson.D{
		{"$inc", bson.D{{"likes", -1}}},
	}

	res, err := cu.collectionRepository.UpdateByID(ctx, collectionID, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, errors.New("collection not exists")
	}

	updated, err := cu.collectionRepository.GetByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (cu *collectionUseCase) DeleteByID(c context.Context, collectionID string) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	return cu.collectionRepository.DeleteByID(ctx, collectionID)
}

func (cu *collectionUseCase) GetByID(c context.Context, collectionID string) (domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	return cu.collectionRepository.GetByID(ctx, collectionID)
}

func (cu *collectionUseCase) SearchPublic(c context.Context, text string, count int, offset int, sortBy string, category string, userID string) ([]domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	var filter interface{}
	if text == "" {
		filter = bson.M{
			"is_public": true,
		}
	} else {
		filter = bson.M{
			"is_public": true,
			"$text": bson.M{
				"$search": text,
			},
		}
	}

	if category == "favourite" {
		user, err := cu.userRepository.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if user.Favourite == nil {
			user.Favourite = make([]string, 0)
		}

		if text == "" {
			filter = bson.M{
				"_id": bson.M{
					"$in": user.Favourite,
				},
			}
		} else {
			filter = bson.M{
				"_id": bson.M{
					"$in": user.Favourite,
				},
				"$text": bson.M{
					"$search": text,
				},
			}
		}

	}

	opts := database.FindOptions{
		Limit:  int64(count),
		Skip:   int64(offset),
		SortBy: sortBy,
	}

	collections, err := cu.collectionRepository.GetByFilter(ctx, filter, opts)
	if err == nil && len(collections) == 0 && category != "favourite" {
		filter = bson.M{
			"is_public": true,
			"name": bson.Regex{
				Pattern: ".*" + text + ".*",
				Options: "i",
			},
		}
		return cu.collectionRepository.GetByFilter(ctx, filter, opts)
	}

	return collections, err
}

func (cu *collectionUseCase) SearchPublicByAuthor(c context.Context, author string) ([]domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	var filter interface{}
	filter = bson.M{
		"author":    author,
		"is_public": true,
	}

	collections, err := cu.collectionRepository.GetByFilter(ctx, filter, database.FindOptions{})

	return collections, err
}

func (cu *collectionUseCase) AddCard(c context.Context, collectionID string, card *domain.Card) (domain.Card, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	collection, err := cu.collectionRepository.GetByID(c, collectionID)
	if err != nil {
		return domain.Card{}, err
	}

	card.LocalID = collection.MaxId
	if card.OtherAnswers.Items == nil {
		card.OtherAnswers.Items = make([]string, 0)
	}
	update := bson.D{
		{"$push", bson.D{
			{"cards", card},
		}},
		{"$inc", bson.D{
			{"max_id", 1},
		}},
	}

	answer := domain.Card{
		LocalID:      card.LocalID,
		Question:     card.Question,
		Answer:       card.Answer,
		Attachment:   card.Attachment,
		OtherAnswers: card.OtherAnswers,
	}

	_, err = cu.collectionRepository.UpdateByID(ctx, collectionID, update)
	return answer, err
}

func (cu *collectionUseCase) DeleteCard(c context.Context, collectionID string, cardLocalID int) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	update := bson.D{
		{"$pull", bson.D{
			{"cards", bson.D{
				{"local_id", cardLocalID},
			}},
		}},
	}
	res, err := cu.collectionRepository.UpdateByID(ctx, collectionID, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("Card not exists")
	}

	return nil
}

func (cu *collectionUseCase) UpdateCard(c context.Context, collectionID string, card *domain.Card) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	filter := bson.D{
		{"_id", collectionID},
		{"cards.local_id", card.LocalID},
	}

	if card.OtherAnswers.Items == nil {
		card.OtherAnswers.Items = make([]string, 0)
	}

	update := bson.D{
		{"$set", bson.D{
			{"cards.$.question", card.Question},
			{"cards.$.answer", card.Answer},
			{"cards.$.other_answers", card.OtherAnswers},
		}},
	}
	res, err := cu.collectionRepository.Update(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("card not exists")
	}
	return nil
}

func (cu *collectionUseCase) GetCardPhoto(c context.Context, objectName string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	return cu.collectionStorage.GetObject(ctx, objectName)
}

func (cu *collectionUseCase) UploadCardPhoto(c context.Context, userID string, collectionID string, cardID int, picture io.Reader, size int64) (string, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	user, err := cu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.Limits.TotalFileSize+int(size) > domain.MAX_TOTAL_FILE_SIZE {
		return "", errors.New("you reached the total file size limit")
	}

	if !slices.Contains(user.Collections, collectionID) {
		return "", errors.New("access denied")
	}

	user.Limits.TotalFileSize += int(size)

	update := bson.D{
		{"$set", bson.D{
			{"limits", user.Limits},
		}},
	}

	err = cu.userRepository.UpdateByID(ctx, userID, update)
	if err != nil {
		return "", err
	}

	objectName := collectionID + "_" + strconv.Itoa(cardID) + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	err = cu.collectionStorage.PutObject(ctx, objectName, picture, size)
	if err != nil {
		return "", err
	}

	filter := bson.D{
		{"_id", collectionID},
		{"cards.local_id", cardID},
	}

	update = bson.D{
		{"$set", bson.D{
			{"cards.$.attachment", objectName},
		}},
	}

	_, err = cu.collectionRepository.Update(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (cu *collectionUseCase) RemoveCardPicture(c context.Context, userID, collectionID string, cardID int, objectName string) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	obj, err := cu.collectionStorage.GetObject(ctx, objectName)
	if err != nil {
		return errors.New("Card picture not exists")
	}

	// update card in collection
	filter := bson.D{
		{"_id", collectionID},
		{"cards.local_id", cardID},
	}

	update := bson.D{
		{"$set", bson.D{
			{"cards.$.attachment", ""},
		}},
	}

	_, err = cu.collectionRepository.Update(ctx, filter, update)
	if err != nil {
		return err
	}
	//

	// update user limits
	update = bson.D{
		{"$inc", bson.D{
			{"limits.total_file_size", -len(obj)},
		}},
	}

	err = cu.userRepository.UpdateByID(ctx, userID, update)
	if err != nil {
		return err
	}
	//

	return cu.collectionStorage.RemoveObject(ctx, objectName)
}

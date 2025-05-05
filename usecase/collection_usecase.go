package usecase

import (
	"context"
	"errors"
	"io"
	"main/database"
	"main/domain"
	"main/repository"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type collectionUseCase struct {
	collectionRepository domain.CollectionRepository
	collectionStorage    domain.CollectionStorage
	userRepository       domain.UserRepository
	contextTimeout       time.Duration
}

func NewCollectionUseCase(
	collectionRepository domain.CollectionRepository, collectionStorage domain.CollectionStorage,
	userRepository domain.UserRepository, timeout time.Duration,
) domain.CollectionUseCase {
	return &collectionUseCase{
		collectionRepository: collectionRepository,
		collectionStorage:    collectionStorage,
		userRepository:       userRepository,
		contextTimeout:       timeout,
	}
}

func (cu *collectionUseCase) Create(c context.Context, collection *domain.Collection, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	collection.Cards = make([]domain.Card, 0)
	collection.NameLower = strings.ToLower(collection.Name)

	client := repository.GetClient()
	session, err := client.StartSession()
	if err != nil {
		return "", err
	}

	defer session.EndSession(ctx)

	id, err := session.WithTransaction(ctx, func(transactionCtx context.Context) (interface{}, error) {
		var id string
		id, err = cu.collectionRepository.Create(transactionCtx, collection)
		if err != nil {
			return nil, err
		}

		err = cu.userRepository.AddCollection(transactionCtx, userID, id, "collections")
		if err != nil {
			return nil, err
		}

		return id, nil
	})
	if err != nil {
		return "", err
	}

	return id.(string), nil
}

func (cu *collectionUseCase) PutByID(c context.Context, collectionID string, collection *domain.Collection) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: collection.Name},
			{Key: "is_public", Value: collection.IsPublic},
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

func (cu *collectionUseCase) AddLike(
	c context.Context, collectionID string, userID string) (*domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	client := repository.GetClient()
	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(transactionCtx context.Context) (interface{}, error) {
		err = cu.userRepository.AddCollection(transactionCtx, userID, collectionID, "favourite")
		if err != nil {
			return nil, err
		}

		update := bson.D{
			{Key: "$inc", Value: bson.D{{Key: "likes", Value: 1}}},
		}
		var res database.UpdateResult
		res, err = cu.collectionRepository.UpdateByID(transactionCtx, collectionID, update)
		if err != nil {
			return nil, err
		}
		if res.MatchedCount == 0 {
			return nil, errors.New("collection not exists")
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := cu.collectionRepository.GetByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (cu *collectionUseCase) RemoveLike(
	c context.Context, collectionID string, userID string) (*domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	current, err := cu.collectionRepository.GetByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	if current.Likes <= 0 {
		return &current, nil
	}

	client := repository.GetClient()
	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(transactionCtx context.Context) (interface{}, error) {
		err = cu.userRepository.DeleteCollection(transactionCtx, userID, collectionID, "favourite")
		if err != nil {
			return nil, err
		}

		update := bson.D{
			{Key: "$inc", Value: bson.D{{Key: "likes", Value: -1}}},
		}
		var res database.UpdateResult
		res, err = cu.collectionRepository.UpdateByID(transactionCtx, collectionID, update)
		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, errors.New("collection not exists")
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := cu.collectionRepository.GetByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (cu *collectionUseCase) DeleteByID(c context.Context, collectionID, userID string) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	client := repository.GetClient()
	session, err := client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(transactionCtx context.Context) (interface{}, error) {
		err = cu.userRepository.DeleteCollection(transactionCtx, userID, collectionID, "collections")
		if err != nil {
			return nil, err
		}
		var coll domain.Collection
		coll, err = cu.GetByID(transactionCtx, collectionID)
		if err != nil {
			return nil, err
		}

		for _, elem := range coll.Cards {
			if elem.Attachment != "" {
				err = cu.RemoveCardPicture(transactionCtx, userID, collectionID, elem.LocalID, elem.Attachment)
				if err != nil && err.Error() != "card picture not exists" {
					return nil, err
				}
			}
		}

		return nil, cu.collectionRepository.DeleteByID(transactionCtx, collectionID)
	})

	return err
}

func (cu *collectionUseCase) GetByID(c context.Context, collectionID string) (domain.Collection, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()
	return cu.collectionRepository.GetByID(ctx, collectionID)
}

func (cu *collectionUseCase) SearchPublic(
	c context.Context, text string, count int, offset int, sortBy string, category string, userID string,
) ([]domain.Collection, error) {
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

	var filter interface{} = bson.M{
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

	card.LocalID = collection.MaxID
	if card.OtherAnswers.Items == nil {
		card.OtherAnswers.Items = make([]string, 0)
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "cards", Value: card},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "max_id", Value: 1},
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
		{Key: "$pull", Value: bson.D{
			{Key: "cards", Value: bson.D{
				{Key: "local_id", Value: cardLocalID},
			}},
		}},
	}
	res, err := cu.collectionRepository.UpdateByID(ctx, collectionID, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("card not exists")
	}

	return nil
}

func (cu *collectionUseCase) UpdateCard(c context.Context, collectionID string, card *domain.Card) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	filter := bson.D{
		{Key: "_id", Value: collectionID},
		{Key: "cards.local_id", Value: card.LocalID},
	}

	if card.OtherAnswers.Items == nil {
		card.OtherAnswers.Items = make([]string, 0)
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "cards.$.question", Value: card.Question},
			{Key: "cards.$.answer", Value: card.Answer},
			{Key: "cards.$.other_answers", Value: card.OtherAnswers},
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

func (cu *collectionUseCase) UploadCardPhoto(
	c context.Context, userID string, collectionID string, cardID int, picture io.Reader, size int64,
) (string, error) {
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
		{Key: "$set", Value: bson.D{
			{Key: "limits", Value: user.Limits},
		}},
	}

	_, err = cu.userRepository.UpdateByID(ctx, userID, update)
	if err != nil {
		return "", err
	}

	objectName := collectionID + "_" + strconv.Itoa(cardID) + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	err = cu.collectionStorage.PutObject(ctx, objectName, picture, size)
	if err != nil {
		return "", err
	}

	filter := bson.D{
		{Key: "_id", Value: collectionID},
		{Key: "cards.local_id", Value: cardID},
	}

	update = bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "cards.$.attachment", Value: objectName},
		}},
	}

	_, err = cu.collectionRepository.Update(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (cu *collectionUseCase) RemoveCardPicture(
	c context.Context,
	userID, collectionID string,
	cardID int,
	objectName string,
) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	obj, err := cu.collectionStorage.GetObject(ctx, objectName)
	if err != nil {
		return errors.New("card picture not exists")
	}

	// update card in collection
	filter := bson.D{
		{Key: "_id", Value: collectionID},
		{Key: "cards.local_id", Value: cardID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "cards.$.attachment", Value: ""},
		}},
	}

	_, err = cu.collectionRepository.Update(ctx, filter, update)
	if err != nil {
		return err
	}
	//

	// update user limits
	update = bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "limits.total_file_size", Value: -len(obj)},
		}},
	}

	_, err = cu.userRepository.UpdateByID(ctx, userID, update)
	if err != nil {
		return err
	}
	//

	return cu.collectionStorage.RemoveObject(ctx, objectName)
}

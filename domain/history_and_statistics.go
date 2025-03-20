package domain

import "context"

const (
	UserHistoryCollection       = "user_history"
	CollectionHistoryCollection = "collection_history"
)

type UserStatistics struct {
	TotalTrainings   int `bson:"total_trainings" json:"total_trainings"`
	MediumPercentage int `bson:"medium_percentage" json:"medium_percentage"`
}

type HistoryItem struct {
	CollectionID   string      `bson:"collection_id" json:"collection_id"`
	CollectionName string      `bson:"collection_name" json:"collection_name"`
	Time           int         `bson:"time" json:"time"`
	CorrectCards   []int       `bson:"correct_cards" json:"correct_cards"`
	IncorrectCards []int       `bson:"incorrect_cards" json:"incorrect_cards"`
	AllCardsCount  int         `bson:"all_cards_count" json:"all_cards_count"`
	Errors         []ErrorItem `bson:"errors" json:"errors"`
}

type SmallHistoryItem struct {
	CollectionName string `bson:"collection_name" json:"collection_name"`
	Time           int    `bson:"time" json:"time"`
	CorrectCards   []int  `bson:"correct_cards" json:"correct_cards"`
	IncorrectCards []int  `bson:"incorrect_cards" json:"incorrect_cards"`
}

type ErrorItem struct {
	CardID      int    `bson:"card_id" json:"card_id"`
	Question    string `bson:"question" json:"question"`
	Answer      string `bson:"answer" json:"answer"`
	Type        string `bson:"type" json:"type"`
	UserAnswer  string `bson:"user_answer" json:"user_answer"`
	BlankAnswer string `bson:"blank_answer" json:"blank_answer"`
}

type UserHistory struct {
	UserID string        `bson:"_id" json:"user_id"`
	Items  []HistoryItem `bson:"items" json:"items"`
}

type CollectionHistory struct {
	CollectionID string             `bson:"_id" json:"collection_id"`
	Items        []SmallHistoryItem `bson:"items" json:"items"`
}

type UserHistoryRepository interface {
	CreateIfNotExists(c context.Context, userID string) error
	UpdateByID(c context.Context, userID string, item HistoryItem) error
	GetByID(c context.Context, userID string) (UserHistory, error)
}

type CollectionHistoryRepository interface {
	CreateIfNotExists(c context.Context, collectionID string) error
	UpdateByID(c context.Context, collectionID string, item SmallHistoryItem) error
	GetByID(c context.Context, collectionID string) (CollectionHistory, error)
}

type HistoryUseCase interface {
	GetUserHistory(c context.Context, userID string) (UserHistory, error)
	AddTraining(c context.Context, userID string, historyItem HistoryItem) error
}

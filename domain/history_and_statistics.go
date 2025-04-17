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
	CollectionID   string            `bson:"collection_id" json:"collection_id"`
	CollectionName string            `bson:"collection_name" json:"collection_name"`
	Time           int               `bson:"time" json:"time"`
	CorrectCards   []int             `bson:"correct_cards" json:"correct_cards"`
	IncorrectCards []int             `bson:"incorrect_cards" json:"incorrect_cards"`
	AllCardsCount  int               `bson:"all_cards_count" json:"all_cards_count"`
	Errors         []ErrorItem       `bson:"errors" json:"errors"`
	RightAnswers   []RightAnswerItem `bson:"right_answers" json:"right_answers"`
}

type SmallHistoryItem struct {
	CollectionName string `bson:"collection_name" json:"collection_name"`
	Time           int    `bson:"time" json:"time"`
	CorrectCards   []int  `bson:"correct_cards" json:"correct_cards"`
	IncorrectCards []int  `bson:"incorrect_cards" json:"incorrect_cards"`
	AllCardsCount  int    `bson:"all_cards_count" json:"all_cards_count"`
}

type ErrorItem struct {
	CardID      int    `bson:"card_id" json:"card_id"`
	Question    string `bson:"question" json:"question"`
	Answer      string `bson:"answer" json:"answer"`
	Type        string `bson:"type" json:"type"`
	UserAnswer  string `bson:"user_answer" json:"user_answer"`
	BlankAnswer string `bson:"blank_answer" json:"blank_answer"`
	Attachment  string `bson:"attachment" json:"attachment"`
}

type RightAnswerItem struct {
	CardID int    `bson:"card_id" json:"card_id"`
	Type   string `bson:"type" json:"type"`
}

type UserHistory struct {
	UserID string        `bson:"_id" json:"user_id"`
	Items  []HistoryItem `bson:"items" json:"items"`
}

type CollectionHistory struct {
	CollectionID string             `bson:"_id" json:"collection_id"`
	Items        []SmallHistoryItem `bson:"items" json:"items"`
}

type UserHistoryArray struct {
	Count int           `json:"count"`
	Items []HistoryItem `json:"items"`
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
	GetUserHistoryFromTime(c context.Context, userID string, fromTime int) ([]HistoryItem, error)
	AddTraining(c context.Context, userID string, historyItem HistoryItem) error
}

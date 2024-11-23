package models

type Card struct {
	LocalID  int    `bson:"local_id" json:"local_id"`
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer" json:"answer"`
}
type GlobalValues struct {
	MaxCollectionID int `bson:"max_collection_id"`
}

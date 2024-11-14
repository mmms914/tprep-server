package models

type Card struct {
	LocalID  int    `bson:"local_id" json:"local_id"`
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer" json:"answer"`
}

type Collection struct {
	ID       int    `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	IsPublic bool   `bson:"is_public" json:"is_public"`
	Cards    []Card `bson:"cards" json:"cards"`
}

package domain

type Card struct {
	LocalID  int    `bson:"local_id" json:"local_id"`
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer"   json:"answer"`

	Attachment   string       `bson:"attachment"    json:"attachment"`
	OtherAnswers OtherAnswers `bson:"other_answers" json:"other_answers"`
}

type OtherAnswers struct {
	Count int      `bson:"count" json:"count"`
	Items []string `bson:"items" json:"items"`
}

type UploadCardPhotoResult struct {
	ObjectName string `bson:"object_name" json:"object_name"`
}

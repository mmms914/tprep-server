package models

import "errors"

type Collection struct {
	ID       int    `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	IsPublic bool   `bson:"is_public" json:"is_public"`
	Cards    []Card `bson:"cards" json:"cards"`
}

type CollectionInfo struct {
	Id    int `bson:"id"`
	MaxId int `bson:"max_id"`
}

func (c *Collection) AddCard(card Card) {
	c.Cards = append(c.Cards, card)
}

func (c *Collection) DeleteCard(localId int) int {
	for index, value := range c.Cards {
		if value.LocalID == localId {
			c.Cards = append(c.Cards[:index], c.Cards[index+1:]...)
			return 1
		}
	}
	return 0
}

func (c *Collection) UpdateCard(localId int, cardPrototype Card) (int, error) {
	for index, value := range c.Cards {
		if value.LocalID == localId {
			if cardPrototype.Question != "" {
				c.Cards[index].Question = cardPrototype.Question
			}
			if cardPrototype.Answer != "" {
				c.Cards[index].Answer = cardPrototype.Answer
			}
			return index, nil
		}
	}
	return -1, errors.New("card not found")
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"company name"`
	Zip string `json:"zip" bson:"zip code"`
	Website string `json:"website"`
}
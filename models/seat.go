package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Seat struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
	Row  string             `bson:"row" json:"row"`
	Col  int                `bson:"col" json:"col"`
}

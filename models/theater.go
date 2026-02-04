package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Theater struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name   string             `bson:"name" json:"name"`
	AllRow int                `bson:"all_row" json:"all_row"`
	AllCol int                `bson:"all_col" json:"all_col"`
	Seats  []Seat             `bson:"seats" json:"seats"`
}

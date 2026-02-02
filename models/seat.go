package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Seat struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Row         int                `bson:"row" json:"row"`
	Col         int                `bson:"col" json:"col"`
	BookedAt    *time.Time         `bson:"booked_at,omitempty"`
	ExpiredAt   *time.Time         `bson:"expired_at,omitempty"`
	AlreadySold bool               `bson:"already_sold" json:"already_sold" default:"false"`
}

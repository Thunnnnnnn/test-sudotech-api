package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName      string             `bson:"full_name" json:"full_name"`
	Email         string             `bson:"email" json:"email"`
	Firstname     string             `bson:"first_name" json:"first_name"`
	Surname       string             `bson:"surname" json:"surname"`
	Picture       string             `bson:"picture" json:"picture"`
	VerifiedEmail bool               `bson:"verified_email" json:"verified_email"`
}

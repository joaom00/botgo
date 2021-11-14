package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type coin struct {
	Symbol   string  `bson:"symbol"`
	Quantity float64 `bson:"quantity"`
}

type wallet struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserID     string             `bson:"user_id"`
	Amount     float64            `bson:"amount"`
	Coins      []coin             `bson:"coins"`
	SalaryTime time.Time
}

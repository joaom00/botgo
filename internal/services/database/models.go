package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Coin struct {
	Symbol   string  `bson:"symbol"`
	Quantity float64 `bson:"quantity"`
}

type Wallet struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID string             `bson:"user_id"`
	Amount float64            `bson:"amount"`
	Coins  []Coin             `bson:"coins"`
}

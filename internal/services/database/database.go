package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	mongoOnce           sync.Once
)

func AddSalary(userID string) (Wallet, error) {
	wallet, err := GetWallet(userID)
	if err != nil {
		return wallet, err
	}

	client, err := getMongoClient()
	if err != nil {
		return wallet, err
	}

	walletCollection := client.Database("jbot").Collection("wallets")

	wallet.Amount += 30
	filter := bson.D{primitive.E{Key: "user_id", Value: userID}}
	updater := bson.D{primitive.E{Key: "$set", Value: wallet}}

	_, err = walletCollection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

func Invest(wallet Wallet, cryptoCoin string, quantity float64) error {
	client, err := getMongoClient()
	if err != nil {
		return err
	}

	cryptoPrice, err := getCryptoCoin(cryptoCoin)
	if err != nil {
		return err
	}
	total := quantity / cryptoPrice

	walletCollection := client.Database("jbot").Collection("wallets")

	i := indexOf(wallet.Coins, cryptoCoin)

	if i == -1 {
		wallet.Coins = append(wallet.Coins, Coin{Symbol: cryptoCoin, Quantity: total})
	} else {
		wallet.Coins[i].Quantity += total
	}
	wallet.Amount -= quantity

	filter := bson.D{primitive.E{Key: "user_id", Value: wallet.UserID}}
	updater := bson.D{primitive.E{Key: "$set", Value: wallet}}

	_, err = walletCollection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	return nil
}

func GetWallet(userID string) (Wallet, error) {
	result := Wallet{}

	client, err := getMongoClient()
	if err != nil {
		return result, err
	}

	walletCollection := client.Database("jbot").Collection("wallets")

	filter := bson.D{primitive.E{Key: "user_id", Value: userID}}

	err = walletCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateWallet(userID string) (Wallet, error) {
	client, err := getMongoClient()
	if err != nil {
		return Wallet{}, err
	}

	walletCollection := client.Database("jbot").Collection("wallets")

	result := Wallet{
		UserID: userID,
		Amount: 100,
		Coins:  []Coin{},
	}

	_, err = walletCollection.InsertOne(context.TODO(), result)
	if err != nil {
		return Wallet{}, err
	}

	return result, nil
}

func getMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_CREDENTIALS"))

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}

		clientInstance = client
	})

	return clientInstance, clientInstanceError
}

func indexOf(arr []Coin, cryptoCoin string) int {
	for i, c := range arr {
		if c.Symbol == cryptoCoin {
			return i
		}
	}
	return -1
}

type CryptoCoin struct {
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}

func getCryptoCoin(symbol string) (float64, error) {
	res, err := http.Get(fmt.Sprintf("https://brapi.ga/api/v2/crypto?coin=%s&currency=BRL", strings.ToUpper(symbol)))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var coins struct{ Coins []CryptoCoin }

	if err = json.Unmarshal(body, &coins); err != nil {
		return 0, err
	}

	return coins.Coins[0].RegularMarketPrice, nil
}

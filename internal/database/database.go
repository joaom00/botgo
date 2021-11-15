package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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

const (
	DB      = "jbot"
	WALLETS = "wallets"
)

func (w *wallet) AddSalary() error {
	client, err := getMongoClient()
	if err != nil {
		return err
	}

	t := w.SalaryTime.Add(24 * time.Hour)
	if t.After(time.Now()) {
		return errors.New("only one salary per day")
	}

	w.Amount += 30
	w.SalaryTime = time.Now()
	filter := bson.D{primitive.E{Key: "user_id", Value: w.UserID}}
	updater := bson.D{primitive.E{Key: "$set", Value: w}}

	walletCollection := client.Database("jbot").Collection("wallets")
	_, err = walletCollection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	return nil

}

func (w *wallet) Invest(cryptoCoin string, quantity float64) (float64, error) {
	cryptos, err := getCryptoCoin([]string{cryptoCoin})
	if err != nil {
		return 0, err
	}
	total := quantity / cryptos[0].RegularMarketPrice

	i := indexOf(w.Coins, cryptoCoin)

	var totalInvested float64
	if i == -1 {
		w.Coins = append(w.Coins, coin{Symbol: cryptoCoin, Quantity: total})
		totalInvested = total
	} else {
		w.Coins[i].Quantity += total
		totalInvested = w.Coins[i].Quantity
	}
	w.Amount -= quantity

	filter := bson.D{primitive.E{Key: "user_id", Value: w.UserID}}
	updater := bson.D{primitive.E{Key: "$set", Value: w}}

	client, err := getMongoClient()
	if err != nil {
		return 0, err
	}

	walletCollection := client.Database(DB).Collection(WALLETS)
	_, err = walletCollection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return 0, err
	}

	return totalInvested, nil
}

func GetWallet(userID string) (wallet, error) {
	result := wallet{}

	client, err := getMongoClient()
	if err != nil {
		return result, err
	}

	walletCollection := client.Database(DB).Collection(WALLETS)

	filter := bson.D{primitive.E{Key: "user_id", Value: userID}}

	err = walletCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateWallet(userID string) (wallet, error) {
	client, err := getMongoClient()
	if err != nil {
		return wallet{}, err
	}

	walletCollection := client.Database(DB).Collection(WALLETS)

	result := wallet{
		UserID: userID,
		Amount: 100,
		Coins:  []coin{},
	}

	_, err = walletCollection.InsertOne(context.TODO(), result)
	if err != nil {
		return wallet{}, err
	}

	return result, nil
}

func (w *wallet) Total() float64 {
	var total float64
	var symbols []string

	for _, coin := range w.Coins {
		symbols = append(symbols, coin.Symbol)
	}

	cryptos, _ := getCryptoCoin(symbols)

	for i, coin := range w.Coins {
		total += cryptos[i].RegularMarketPrice * coin.Quantity
	}

	return total + w.Amount
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

func indexOf(arr []coin, cryptoCoin string) int {
	for i, c := range arr {
		if c.Symbol == cryptoCoin {
			return i
		}
	}
	return -1
}

func getCryptoCoin(symbols []string) ([]cryptoCoin, error) {
	res, err := http.Get(fmt.Sprintf("https://brapi.ga/api/v2/crypto?coin=%s&currency=BRL", strings.Join(symbols, ",")))
	if err != nil {
		return []cryptoCoin{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []cryptoCoin{}, err
	}

	var coins struct{ Coins []cryptoCoin }

	if err = json.Unmarshal(body, &coins); err != nil {
		return []cryptoCoin{}, err
	}

	return coins.Coins, nil
}

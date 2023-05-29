package testhelpers

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Testdb struct {
	client *mongo.Client
	*db.Store
}

func (t *Testdb) Teardown() error {
	dbname := os.Getenv(db.MongoDBNmeEnvName)
	if err := t.client.Database(dbname).Drop(context.TODO()); err != nil {
		return err
	}

	return nil

}
func Setup(t *testing.T) *Testdb {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Error()
	}

	dburi := os.Getenv("MONGO_DB_URL_TEST")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &Testdb{
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
			Hotel:   hotelStore,
		},
		client: client,
	}
}

func InsertTestUser(firstname, lastname string, store *db.Store) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", firstname, lastname),
		FirstName: firstname,
		LastName:  lastname,
		Password:  fmt.Sprintf("%s_%s", firstname, lastname),
	})

	if err != nil {
		log.Fatal(err)
	}

	user, err = store.User.InsertUser(context.Background(), user)

	if err != nil {
		log.Fatal(err)

	}

	return user
}

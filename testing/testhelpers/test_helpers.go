package testhelpers

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDBName = "hotel-reservation-test"
	testDBUri  = "mongodb://localhost:27017"
)

type Testdb struct {
	client *mongo.Client
	*db.Store
}

func (t *Testdb) Teardown() error {
	if err := t.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		return err
	}

	return nil

}
func Setup(t *testing.T) *Testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDBUri))
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

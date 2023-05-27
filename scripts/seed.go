package main

import (
	"context"
	"fmt"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/db/fixtures"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	store  *db.Store
	ctx    = context.Background()
)

func main() {

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	fixtures.AddUser(store, "admin", "admin", true)

	hotel := fixtures.AddHotel(store, "Aubrale", "Pakistan", 3.6, nil)
	room := fixtures.AddRoom(store, "large", true, 55.55, hotel.ID)

	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println("booking ---->", booking.ID.Hex())

}

func init() {

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store = &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   db.NewMongoHotelStore(client),
	}

}

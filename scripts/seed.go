package main

import (
	"context"
	"fmt"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/db/fixtures"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	store  *db.Store
	ctx    = context.Background()
)

func main() {

	user := fixtures.AddUser(store, "james", "normal", false)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Printf("%s -> %s\n", admin.Email, api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "Auberine", "Asfall", 4, nil)
	room := fixtures.AddRoom(store, "large", true, 55.55, hotel.ID)

	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println("booking ---->", booking.ID.Hex())

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("random hotel location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}

}

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	DBURL := os.Getenv("MONGO_DB_URL")
	DBNAME := os.Getenv("MONGO_DB_NAME")
	fmt.Println(DBURL)
	fmt.Println(DBNAME)

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURL))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(DBNAME).Drop(ctx); err != nil {
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

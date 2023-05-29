package fixtures

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {

	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}

	resp, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	booking.ID = resp.ID
	return booking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)

	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomsIDS = rooms

	if rooms == nil {
		roomsIDS = []primitive.ObjectID{}
	}

	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomsIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{

		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})

	user.IsAdmin = admin
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

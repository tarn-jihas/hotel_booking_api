package main

import (
	"context"
	"hotel-reservation/api"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configutration
// 1. MongoDB Endpoint
// 2. ListenAddress of our HTTP server
// 3. JWT Secret
// 4. MongoDBName

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	var (
		// stores init
		userstore    = db.NewMongoUserStore(client)
		hotelstore   = db.NewMongoHotelStore(client)
		roomstore    = db.NewMongoRoomStore(client, hotelstore)
		bookingstore = db.NewMongoBookingStore(client)
		stores       = &db.Store{
			User:    userstore,
			Hotel:   hotelstore,
			Room:    roomstore,
			Booking: bookingstore,
		}
		// handlers init
		userHandler    = api.NewUserHandler(userstore)
		hotelHandler   = api.NewHotelHandler(stores)
		authHandler    = api.NewAuthHandler(userstore)
		roomHandler    = api.NewRoomHandler(stores)
		bookingHandler = api.NewBookingHandler(stores)

		app  = fiber.New(config)
		auth = app.Group("/api")

		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userstore))
		admin = apiv1.Group("/admin", middleware.AdminAuth)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetHotelRooms)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)

	apiv1.Get("/room", roomHandler.HandlerGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// booking handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBookings)

	// admin route
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)

}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

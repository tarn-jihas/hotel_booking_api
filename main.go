package main

import (
	"context"
	"flag"
	"hotel-reservation/api"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
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
	app.Listen(*listenAddr)

}

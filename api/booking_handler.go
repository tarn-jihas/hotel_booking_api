package api

import (
	"hotel-reservation/db"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (b *BookingHandler) HandleCancelBookings(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := b.store.Booking.GetBookingsByID(c.Context(), id)

	if err != nil {
		return err
	}

	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	if err := b.store.Booking.UpdateBooking(c.Context(), booking.ID.Hex(), bson.M{
		"$set": bson.M{"canceled": true},
	}); err != nil {
		return err
	}

	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})

}

func (b *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {

	bookings, err := b.store.Booking.GetBookings(c.Context(), bson.M{})

	if err != nil {
		return err
	}

	return c.JSON(bookings)

}

func (b *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {

	id := c.Params("id")
	booking, err := b.store.Booking.GetBookingsByID(c.Context(), id)

	if err != nil {
		return err
	}

	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	return c.JSON(booking)
}

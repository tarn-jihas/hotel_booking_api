package api

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}

	return nil
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
	store db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: *store,
	}
}

func (h *RoomHandler) HandlerGetRooms(c *fiber.Ctx) error {

	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})

	if err != nil {
		return err
	}

	// for _, room := range rooms {
	// 	h.IsRoomAvailableForBooking()
	// }
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}
	roomOID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	ok, err = h.IsRoomAvailableForBooking(c.Context(), roomOID, params)
	if err != nil {
		return err
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room with the %s id is already booked", roomOID.Hex()),
		})
	}

	booking := types.Booking{
		RoomID:     roomOID,
		UserID:     user.ID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) IsRoomAvailableForBooking(c context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(c, where)

	if err != nil {
		return false, err
	}
	fmt.Printf("%+v", bookings)
	ok := len(bookings) == 0
	return ok, nil

}

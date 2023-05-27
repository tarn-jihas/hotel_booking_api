package handletest

import (
	"encoding/json"
	"hotel-reservation/api"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db/fixtures"
	"hotel-reservation/testing/testhelpers"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := testhelpers.Setup(t)

	defer db.Teardown()
	var (
		user    = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel   = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room    = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from    = time.Now()
		till    = from.AddDate(0, 0, 5)
		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app     = fiber.New()
		route   = app.Group("/:id", middleware.JWTAuthentication(db.User))
		_       = booking

		bookinHandler = api.NewBookingHandler(db.Store)
	)

	token := api.CreateTokenFromUser(user)

	route.Get("/", bookinHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", "/"+booking.ID.Hex(), nil)
	req.Header.Add("x-api-token", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

}

func TestAdminGetBookings(t *testing.T) {
	db := testhelpers.Setup(t)

	defer db.Teardown()
	var (
		user    = fixtures.AddUser(db.Store, "james", "foo", false)
		admin   = fixtures.AddUser(db.Store, "james", "foo", true)
		hotel   = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room    = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from    = time.Now()
		till    = from.AddDate(0, 0, 5)
		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app     = fiber.New()
		jtwapp  = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		_       = booking

		bookinHandler = api.NewBookingHandler(db.Store)
	)

	token := api.CreateTokenFromUser(admin)

	jtwapp.Get("/", bookinHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("x-api-token", token)
	resp, err := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}
	if err != nil {
		t.Fatal(err)
	}

	var bookings []*types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	have := bookings[0]

	if have.UserID != booking.UserID {
		t.Fatalf("expected  %s got %s", booking.UserID, have.UserID)
	}

	// test non-admin cannot access the bookings
	token = api.CreateTokenFromUser(user)

	jtwapp.Get("/", bookinHandler.HandleGetBookings)

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("x-api-token", token)
	resp, err = app.Test(req)

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code but got this: %d", resp.StatusCode)
	}
	if err != nil {
		t.Fatal(err)
	}

}

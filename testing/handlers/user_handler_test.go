package handletest

import (
	"bytes"
	"encoding/json"
	"hotel-reservation/api"
	"hotel-reservation/testing/testhelpers"
	"hotel-reservation/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

var (
	tdb         *testhelpers.Testdb
	userHandler *api.UserHandler
	app         *fiber.App
	testUser    *types.User
)

func testMain(t *testing.T) {

	tdb = testhelpers.Setup(t)

	userHandler = api.NewUserHandler(tdb.User)

	app = fiber.New()

	testUser = testhelpers.InsertTestUser("strahinja", "keselj", tdb.Store)

}

func TestPostUser(t *testing.T) {
	testMain(t)

	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "postUser@post.com",
		FirstName: "Post",
		LastName:  "User",
		Password:  "1232132131",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var user types.User

	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("expecting id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("EncryptedPassword should not be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected username %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected username %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected username %s but got %s", params.Email, user.Email)
	}

	defer tdb.Teardown()

}
func TestGetUser(t *testing.T) {
	testMain(t)
	app.Get("/:id", userHandler.HandleGetUser)

	req := httptest.NewRequest("GET", "/"+testUser.ID.Hex(), nil)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var user types.User

	json.NewDecoder(resp.Body).Decode(&user)

	if user.ID.Hex() != testUser.ID.Hex() {
		t.Errorf("expected ID %s but got %s", testUser.ID.Hex(), user.ID.Hex())
	}

	defer tdb.Teardown()

}

func TestGetUsers(t *testing.T) {
	testMain(t)

	for i := 0; i < 5; i++ {
		testhelpers.InsertTestUser("test", "user", tdb.Store)
	}

	app.Get("/", userHandler.HandleGetUsers)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var users []types.User

	json.NewDecoder(resp.Body).Decode(&users)

	if len(users) <= 1 {
		t.Errorf("expected more than one user, but there was only %v returned.", len(users))
	}

	defer tdb.Teardown()

}

func TestPutUser(t *testing.T) {
	testMain(t)
	app.Put("/:id", userHandler.HandlePutUser)
	params := types.UpdateUserParams{
		FirstName: "TestUpdateFirst",
		LastName:  "TestUpdateLast",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("PUT", "/"+testUser.ID.Hex(), bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var userId map[string]any

	json.NewDecoder(resp.Body).Decode(&userId)

	if _, ok := userId["updated"]; !ok {

		t.Errorf("user was not updated succesfully.")

	}

	defer tdb.Teardown()

}

func TestDeleteUser(t *testing.T) {
	testMain(t)
	app.Delete("/:id", userHandler.HandleDeleteUser)
	req := httptest.NewRequest("DELETE", "/"+testUser.ID.Hex(), nil)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)
	var userId map[string]any

	json.NewDecoder(resp.Body).Decode(&userId)

	if _, ok := userId["deleted"]; !ok {

		t.Errorf("user was not deleted succesfully.")

	}

	defer tdb.Teardown()

}

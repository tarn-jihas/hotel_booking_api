package handletest

import (
	"bytes"
	"encoding/json"
	"hotel-reservation/api"
	"hotel-reservation/testing/testhelpers"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"reflect"

	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSucess(t *testing.T) {
	tdb := testhelpers.Setup(t)
	app := fiber.New()
	defer tdb.Teardown()
	authHandler := api.NewAuthHandler(tdb.Store.User)
	insertedUser := testhelpers.InsertTestUser("james", "foo", tdb.Store)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authparams := types.AuthParams{
		Email:    "james@foo.com",
		Password: "james_foo",
	}

	b, _ := json.Marshal(authparams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}
	var aresp api.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&aresp); err != nil {
		t.Fatal(err)
	}

	if aresp.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth response")
	}

	// Set the encrypted password to an empty string, because we do NOT return that in any
	// JSON response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, aresp.User) {
		t.Fatalf("Expected user needs to be the inserted user ")
	}

}

func TestAuthenticateWithWrongPasswordFailure(t *testing.T) {
	tdb := testhelpers.Setup(t)
	app := fiber.New()

	defer tdb.Teardown()

	authHandler := api.NewAuthHandler(tdb.Store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authparams := types.AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepassbutincorrect",
	}

	b, _ := json.Marshal(authparams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}
	var genResp types.GenericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response type to be msg to be <invalid credentials> %s", genResp.Msg)
	}
}

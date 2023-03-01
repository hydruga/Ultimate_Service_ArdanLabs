package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hydruga/ultimate_service/app/business/data/tests"
	"github.com/hydruga/ultimate_service/app/services/sales-api/handlers"
)

// UserTests holds methods for each user subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type UserTests struct {
	app        http.Handler
	userToken  string
	adminToken string
}

// TestUsers is the entry point for testing user management functions.
func TestUsers(t *testing.T) {
	test := tests.NewIntegration(
		t,
		tests.DBContainer{
			Image: "postgres:13-alpine",
			Port:  "5432",
			Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
		},
	)
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := UserTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken:  test.Token("user@example.com", "gophers"),
		adminToken: test.Token("admin@example.com", "gophers"),
	}

	// t.Run("getToken404", tests.getToken404)
	t.Run("getToken200", tests.getToken200)
	// t.Run("postUser400", tests.postUser400)
	// t.Run("postUser401", tests.postUser401)
	// t.Run("postUser403", tests.postUser403)
	// t.Run("getUser400", tests.getUser400)
	// t.Run("getUser403", tests.getUser403)
	// t.Run("getUser404", tests.getUser404)
	// t.Run("deleteUserNotFound", tests.deleteUserNotFound)
	// t.Run("putUser404", tests.putUser404)
	// t.Run("crudUsers", tests.crudUser)
}

func (ut *UserTests) getToken200(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/users/token", nil)
	w := httptest.NewRecorder()

	r.SetBasicAuth("admin@example.com", "gophers")
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to issues tokens to known users.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen fetching a token with valid credentials.", testID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the response.", tests.Success, testID)

			var got struct {
				Token string `json:"token"`
			}
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response.", tests.Success, testID)

			// TODO(jlw) Should we ensure the token is valid?
		}
	}
}

package handlers

import (
	"assignments-thebriando/servers/gateway/models/users"
	"assignments-thebriando/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestUsersHandlers(t *testing.T) {
	ctx := InitializeFake()
	requestBody, _ := json.Marshal(map[string]string{
		"Email":        "example@example.com",
		"Password":     "1234567",
		"PasswordConf": "1234567",
		"UserName":     "example",
		"FirstName":    "firstName",
		"LastName":     "lastName",
	})

	badRequestBody, _ := json.Marshal(map[string]string{
		"Password":     "1234567",
		"PasswordConf": "1234567",
		"UserName":     "example",
		"FirstName":    "firstName",
		"LastName":     "lastName",
	})

	passHash, _ := bcrypt.GenerateFromPassword([]byte("123456"), 13)

	user := &users.User{
		ID:        4,
		Email:     "example4@example.com",
		PassHash:  passHash,
		UserName:  "aUser4",
		FirstName: "firstName4",
		LastName:  "lastName4",
		PhotoURL:  "https://www.gravatar.com/avatar/user4@example.com",
	}

	cases := []struct {
		method              string
		contentType         string
		expectedStatusCode  int
		expectedContentType string
		id                  string
		requestBody         []byte
		expectedReturn      *users.User
	}{
		// test with wrong method
		{
			http.MethodGet,
			"application/json",
			http.StatusMethodNotAllowed,
			"application/json",
			"83901280938120",
			requestBody,
			user,
		},
		// test with wrong content type
		{
			http.MethodPost,
			"",
			http.StatusUnsupportedMediaType,
			"application/json",
			"83901280938120",
			requestBody,
			user,
		},
		// test valid method and content, bad request body
		{
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"application/json",
			"83901280938120",
			badRequestBody,
			user,
		},
		// test valid
		{
			http.MethodPost,
			"application/json",
			http.StatusCreated,
			"application/json",
			"83901280938120",
			requestBody,
			user,
		},
	}

	for _, c := range cases {
		req, err := http.NewRequest(c.method, "/v1/users", bytes.NewBuffer(c.requestBody))
		if err != nil {
			t.Errorf("Error with request: %v", err)
		}
		req.Header.Set("Content-Type", c.contentType)
		req.Header.Set("Authorization", "Bearer "+c.id)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ctx.UsersHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.expectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, c.expectedStatusCode)
		}
	}
}
func TestSpecificUserHandler(t *testing.T) {
	ctx := InitializeFake()
	user := users.User{
		ID:        1,
		Email:     "example@example.com",
		PassHash:  []byte("123456"),
		UserName:  "aUser",
		FirstName: "firstName",
		LastName:  "lastName",
		PhotoURL:  "https://www.gravatar.com/avatar/user@example.com",
	}
	update := &users.Updates{
		FirstName: "newFirstName",
		LastName:  "newLastName",
	}
	updateJSON, err := json.Marshal(update)
	if err != nil {
		t.Errorf("Error marshaling update")
	}
	testPatch := user
	testPatch.FirstName = update.FirstName
	testPatch.LastName = update.LastName
	cases := []struct {
		method              string
		contentType         string
		expectedStatusCode  int
		expectedContentType string
		id                  string
		requestBody         []byte
		expectedReturn      *users.User
	}{
		// test with invalid id
		{
			http.MethodGet,
			"application/json",
			http.StatusNotFound,
			"application/json",
			"83901280938120",
			[]byte(""),
			&user,
		},
		// test with wrong method
		{
			http.MethodPost,
			"application/json",
			http.StatusMethodNotAllowed,
			"application/json",
			"83901280938120",
			[]byte(""),
			&user,
		},
		// test valid id
		{
			http.MethodGet,
			"application/json",
			http.StatusOK,
			"application/json",
			"1",
			[]byte(""),
			&user,
		},
		// test get "me"
		{
			http.MethodGet,
			"application/json",
			http.StatusOK,
			"application/json",
			"me",
			[]byte(""),
			&user,
		},
		// test patch "me"
		{
			http.MethodPatch,
			"application/json",
			http.StatusOK,
			"application/json",
			"me",
			[]byte(updateJSON),
			&testPatch,
		},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.method, "/v1/users/"+c.id, bytes.NewBuffer(c.requestBody))
		if err != nil {
			t.Errorf("Error with request: %v", err)
		}
		sessState := SessionState{
			SessionTime: time.Now(),
			User:        user,
		}
		sid, err := sessions.NewSessionID(ctx.SigningKey)
		if err != nil {
			t.Errorf(err.Error())
		}
		ctx.SessionStore.Save(sid, sessState)
		req.Header.Set("Content-Type", c.contentType)
		req.Header.Set("Authorization", "Bearer "+string(sid))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ctx.SpecificUserHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.expectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, c.expectedStatusCode)
		}
		compareUser := new(users.User)
		json.Unmarshal([]byte(rr.Body.String()), &compareUser)
		if c.expectedStatusCode == http.StatusOK && c.method == http.MethodGet {
			if compareUser.UserName != c.expectedReturn.UserName {
				t.Errorf("Response username: %v, doesn't match expected: %v",
					compareUser.UserName, c.expectedReturn.UserName)
			}
			if compareUser.FirstName != c.expectedReturn.FirstName {
				t.Errorf("Response firstname: %v, doesn't match expected: %v",
					compareUser.FirstName, c.expectedReturn.FirstName)
			}
			if compareUser.LastName != c.expectedReturn.LastName {
				t.Errorf("Response lastname: %v, doesn't match expected: %v",
					compareUser.LastName, c.expectedReturn.LastName)
			}
			if compareUser.PhotoURL != c.expectedReturn.PhotoURL {
				t.Errorf("Response photourl: %v, doesn't match expected: %v",
					compareUser.PhotoURL, c.expectedReturn.PhotoURL)
			}
		}
		if c.expectedStatusCode == http.StatusOK && c.method == http.MethodPatch {
			if compareUser.FirstName != update.FirstName {
				t.Errorf("Response photourl: %v, doesn't match expected: %v",
					compareUser.FirstName, update.FirstName)
			}
			if compareUser.LastName != update.LastName {
				t.Errorf("Response photourl: %v, doesn't match expected: %v",
					compareUser.LastName, update.LastName)
			}
		}
		ctx.SessionStore.Delete(sid)
	}
}

func TestSessionsHandler(t *testing.T) {

	credentials := &users.Credentials{
		Email:    "example@example.com",
		Password: "123456",
	}
	emptyCredentials := &users.Credentials{
		Email:    "",
		Password: "",
	}
	invalidEmailCredentials := &users.Credentials{
		Email:    "randomemail@email.com",
		Password: "1234567",
	}
	invalidPasswordCredentials := &users.Credentials{
		Email:    "example@example.com",
		Password: "1234567",
	}

	cases := []struct {
		name               string
		method             string
		expectedStatusCode int
		contentType        string
		credential         *users.Credentials
	}{
		{
			"Post Request/ Most Valid",
			http.MethodPost,
			http.StatusCreated,
			"application/json",
			credentials,
		},
		{
			"Post Request/ Credential email is not found in User",
			http.MethodPost,
			http.StatusUnauthorized,
			"application/json",
			invalidEmailCredentials,
		},
		{
			"Post Request/ Invalid PasswordCredential",
			http.MethodPost,
			http.StatusUnauthorized,
			"application/json",
			invalidPasswordCredentials,
		},
		{
			"Post Request/ Invalid Content Type",
			http.MethodPost,
			http.StatusUnsupportedMediaType,
			"",
			emptyCredentials,
		},
		{
			"Get Request",
			http.MethodGet,
			http.StatusMethodNotAllowed,
			"",
			emptyCredentials,
		},
		{
			"Delete Request",
			http.MethodDelete,
			http.StatusMethodNotAllowed,
			"",
			emptyCredentials,
		},
	}
	context := InitializeFake()
	for _, c := range cases {
		marshal, _ := json.Marshal(c.credential)
		req, err := http.NewRequest(c.method, "/v1/users", bytes.NewBuffer(marshal))
		req.Header.Set("Content-Type", c.contentType)

		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(context.SessionsHandler)
		handler.ServeHTTP(rr, req)

		if c.method == http.MethodGet {
			if status := rr.Code; status != c.expectedStatusCode {
				t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
					status, c.expectedStatusCode)
			}
		}
		if c.method == http.MethodDelete {
			if status := rr.Code; status != c.expectedStatusCode {
				t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
					status, c.expectedStatusCode)
			}
		}

		//check the post -> wrong content type
		if c.method == http.MethodPost {
			if c.contentType != "application/json" {
				if status := rr.Code; status != c.expectedStatusCode {
					t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
						status, c.expectedStatusCode)
				}
			}
		}

		//if credential is not found in user
		if c.method == http.MethodPost {
			if c.contentType == "application/json" {
				user, _ := context.UserStore.GetByEmail(c.credential.Email)
				if user != nil {
					err = user.Authenticate(c.credential.Password)
					if err != nil { //fails to authenticate
						if status := rr.Code; status != c.expectedStatusCode {
							t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
								status, c.expectedStatusCode)
						}
					}
					if err == nil {
						sessID, _ := sessions.NewSessionID(context.SigningKey)
						context.SessionStore.Get(sessID, user)
						if status := rr.Code; status != c.expectedStatusCode {
							t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
								status, c.expectedStatusCode)
						}
					}
				} else {
					if status := rr.Code; status != c.expectedStatusCode {
						t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
							status, c.expectedStatusCode)
					}
				}

			}
		}
	}
}

func TestSpecificSessionsHandler(t *testing.T) {
	user, _ := json.Marshal(map[string]string{
		"Email":        "example@example.com",
		"Password":     "1234567",
		"PasswordConf": "1234567",
		"UserName":     "example",
		"FirstName":    "firstName",
		"LastName":     "lastName",
	})
	cases := []struct {
		name               string
		method             string
		expectedStatusCode int
		mine               string
	}{
		{
			"Post Request",
			http.MethodPost,
			http.StatusMethodNotAllowed,
			"",
		},
		{
			"Get Request",
			http.MethodGet,
			http.StatusMethodNotAllowed,
			"",
		},
		{
			"Delete Request",
			http.MethodDelete,
			http.StatusOK,
			"mine",
		},
		{
			"Delete Request",
			http.MethodDelete,
			http.StatusForbidden,
			"",
		},
	}
	context := InitializeFake()
	for _, c := range cases {
		req, err := http.NewRequest(c.method, "/v1/sessions/"+c.mine, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(context.SpecificSessionHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.expectedStatusCode {
			t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
				status, c.expectedStatusCode)
		}
		url, _ := ioutil.ReadAll(rr.Body)
		stringURL := string(url)
		if strings.Contains(stringURL, "mine") {
			sessID, _ := sessions.NewSessionID(context.SigningKey)
			context.SessionStore.Save(sessID, user)
			context.SessionStore.Delete(sessID)
		} else {
			if status := rr.Code; status != c.expectedStatusCode {
				t.Errorf("%s handler returned wrong status code: got %v want %v", c.name,
					status, c.expectedStatusCode)
			}
		}
	}
}

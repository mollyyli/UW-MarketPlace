package handlers

import (
	"UW-Marketplace/servers/gateway/models/users"
	"UW-Marketplace/servers/gateway/sessions"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		header := strings.Split(r.Header.Get("Content-Type"), ",")
		if header[0] != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Request body must be in JSON"))
		} else {
			body, _ := ioutil.ReadAll(r.Body)
			var newUser users.NewUser
			json.Unmarshal([]byte(body), &newUser)
			if newUser.LastName == "" || newUser.FirstName == "" || newUser.Email == "" || newUser.PasswordConf == "" || newUser.Password == "" || newUser.UserName == "" {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				userErr := newUser.Validate()
				if userErr != nil {
					w.WriteHeader(http.StatusBadRequest)
				} else {
					user, err := newUser.ToUser()
					sessionState := SessionState{
						SessionTime: time.Now(),
						User:        *user,
					}
					if err == nil {
						// sid, err := sessions.NewSessionID(ctx.SigningKey)
						// log.Println("sid", sid)
						// log.Println("newsession", err)
						w.Header().Set("Content-Type", "application/json")
						// w.Header().Set("Authorization", "Bearer "+sid.String())
						sid, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, &sessionState, w)
						log.Println("siderr", err)
						log.Println("user", user)
						insertedUser, err := ctx.UserStore.Insert(user)
						if err != nil {
							log.Println("Insert error", err)
						}
						sessionState.User = *insertedUser
						log.Println("sessionState", sessionState)
						log.Println("save", ctx.SessionStore.Save(sid, &sessionState))
						profile, _ := json.Marshal(&user)
						w.WriteHeader(http.StatusCreated)
						w.Write(profile)
					}
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/users/")

	// check if auth header exists
	// if not return http.unauth'ed error
	if authBear := r.Header.Get("Authorization"); len(authBear) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	} else {

		// The current user must be authenticated to call this handler regardless of HTTP method. If the user is not authenticated, respond immediately with an http.StatusUnauthorized (401) error status code

		if r.Method != http.MethodGet && r.Method != http.MethodPatch {
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			if id == "me" {
				sessionID, err := sessions.GetSessionID(r, ctx.SigningKey)
				if err != nil {
					http.Error(w, err.Error(), 400)
				}
				var currentUser users.User
				sessionState := SessionState{
					SessionTime: time.Now(),
					User:        currentUser,
				}
				ctx.SessionStore.Get(sessionID, &sessionState)
				currentUser = sessionState.User

				id = strconv.FormatInt(currentUser.ID, 10)

				if strconv.FormatInt(currentUser.ID, 10) != id {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Access forbidden"))
				}
			}
			intID, _ := strconv.ParseInt(id, 10, 64)
			if r.Method == http.MethodGet {
				user, err := ctx.UserStore.GetByID(intID)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("No user found with given ID"))
				} else {
					json, _ := json.Marshal(user)
					w.Header().Set("Content-Type", "application/json")
					w.Write(json)
					w.WriteHeader(http.StatusOK)
				}

			} else if r.Method == http.MethodPatch {
				if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
					w.WriteHeader(http.StatusUnsupportedMediaType)
					w.Write([]byte("Request body must be in JSON"))
				} else {
					var updates users.Updates
					bodyMarshal, err := ioutil.ReadAll(r.Body)
					if err == nil {
						json.Unmarshal([]byte(bodyMarshal), &updates)
					}
					updatedUser, err := ctx.UserStore.Update(intID, &updates)
					w.Header().Set("Content-Type", "application/json")

					userMarshal, err := json.Marshal(updatedUser)
					if err == nil {
						w.Write(userMarshal)
					}
					w.WriteHeader(http.StatusOK)
				}

			}
		}
	}
}

func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Request body must be in JSON"))
			return
		}
		var credentials users.Credentials

		marshal, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal([]byte(marshal), &credentials)
			user, err := ctx.UserStore.GetByEmail(credentials.Email)
			if user == nil {
				time.Sleep(600 * time.Millisecond)
			}
			if err != nil || user.Authenticate(credentials.Password) != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid credentials"))
			} else {
				sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, &SessionState{}, w)
				headerIP := r.Header.Get("X-Forwarded-For")
				currentIP := r.RemoteAddr
				if len(headerIP) != 0 {
					currentIP = headerIP
				}
				// strID := strconv.FormatInt(user.ID, 10)
				signIn := users.UserSignIn{
					ID:         int64(0),
					UserID:     user.ID,
					SignInTime: time.Now().String(),
					IP:         currentIP,
				}
				ctx.UserStore.InsertSignIn(&signIn)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(marshal)
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (ctx *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	mine := strings.TrimPrefix(r.URL.Path, "/v1/sessions/")
	if r.Method == http.MethodDelete {
		if mine != "mine" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access forbidden"))

		} else {
			sid, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
			if err != nil {
				log.Println(sid)
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Write([]byte("Signed out"))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

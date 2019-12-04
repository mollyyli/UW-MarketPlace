package handlers

import (
	"UW-Marketplace/servers/gateway/models/users"
	"UW-Marketplace/servers/gateway/sessions"
	"encoding/json"
	"io/ioutil"
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
			return
		} else {
			body, _ := ioutil.ReadAll(r.Body)
			var newUser users.NewUser
			json.Unmarshal([]byte(body), &newUser)
			_, err := ctx.UserStore.GetByUserName(newUser.UserName)

			if err == nil {

				http.Error(w, "Username already exists", http.StatusBadRequest)
				return
			}
			_, err = ctx.UserStore.GetByEmail(newUser.Email)

			if err == nil {
				http.Error(w, "Email already exists", http.StatusBadRequest)
				return
			}
			if newUser.LastName == "" || newUser.FirstName == "" || newUser.Email == "" || newUser.PasswordConf == "" || newUser.Password == "" || newUser.UserName == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				userErr := newUser.Validate()
				if userErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				} else {
					user, err := newUser.ToUser()
					sessionState := SessionState{
						SessionTime: time.Now(),
						User:        *user,
					}
					sid, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, &sessionState, w)
					insertedUser, err := ctx.UserStore.Insert(user)
					sessionState.User = *insertedUser
					ctx.SessionStore.Save(sid, &sessionState)
					if err == nil {
						w.Header().Set("Content-Type", "application/json")
						w.Header().Set("Authorization", "Bearer "+sid.String())
						profile, _ := json.Marshal(&user)
						w.WriteHeader(http.StatusCreated)
						w.Write(profile)
						return
					}
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// SpecificUserHandler : givn specific ID
func (ctx *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/users/")
	var currentUser *users.User

	if authBear := r.Header.Get("Authorization"); len(authBear) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	} else {
		if r.Method != http.MethodGet && r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		} else {
			if id == "me" {
				sessionState := &SessionState{}
				_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
				if err != nil {
					http.Error(w, "Can't find state", http.StatusNotFound)
					return
				}
				currentUser = &sessionState.User
				id = string(currentUser.ID)
			} else {
				id64, _ := strconv.ParseInt(id, 10, 64)
				var err2 error
				currentUser, err2 = ctx.UserStore.GetByID(id64)
				if err2 != nil {
					http.Error(w, "No user found with ID", http.StatusNotFound)
					return
				}
			}
			id = strconv.FormatInt(currentUser.ID, 10)
			intID := currentUser.ID
			if r.Method == http.MethodGet {
				user, err := ctx.UserStore.GetByID(intID)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("No user found with given ID"))
					return
				} else {
					json, _ := json.Marshal(user)
					w.Header().Set("Content-Type", "application/json")
					w.Write(json)
					w.WriteHeader(http.StatusOK)
					return
				}

			} else if r.Method == http.MethodPatch {
				if strconv.FormatInt(currentUser.ID, 10) != id {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Access forbidden"))
					return
				}
				if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
					w.WriteHeader(http.StatusUnsupportedMediaType)
					w.Write([]byte("Request body must be in JSON"))
					return
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
					return
				}

			}
		}
	}
}

// func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
// 			w.WriteHeader(http.StatusUnsupportedMediaType)
// 			w.Write([]byte("Request body must be in JSON"))
// 			return
// 		}
// 		var credentials users.Credentials

// 		marshal, err := ioutil.ReadAll(r.Body)

// 		if err == nil {
// 			json.Unmarshal([]byte(marshal), &credentials)
// 			user, err := ctx.UserStore.GetByEmail(credentials.Email)
// 			if user == nil {
// 				time.Sleep(600 * time.Millisecond)
// 			}
// 			if err != nil || user.Authenticate(credentials.Password) != nil {
// 				w.WriteHeader(http.StatusUnauthorized)
// 				w.Write([]byte("Invalid credentials"))
// 			} else {
// 				sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, &SessionState{}, w)
// 				headerIP := r.Header.Get("X-Forwarded-For")
// 				currentIP := r.RemoteAddr
// 				if len(headerIP) != 0 {
// 					currentIP = headerIP
// 				}
// 				// strID := strconv.FormatInt(user.ID, 10)
// 				signIn := users.UserSignIn{
// 					ID:         int64(0),
// 					UserID:     user.ID,
// 					SignInTime: time.Now().String(),
// 					IP:         currentIP,
// 				}
// 				ctx.UserStore.InsertSignIn(&signIn)
// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(http.StatusCreated)
// 				w.Write(marshal)
// 			}
// 		}
// 	} else {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }

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
				return
			} else {
				sessionState := SessionState{
					SessionTime: time.Now(),
					User:        *user,
				}
				sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessionState, w)
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
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
func (ctx *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	mine := strings.TrimPrefix(r.URL.Path, "/v1/sessions/")
	if r.Method == http.MethodDelete {
		if mine != "mine" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access forbidden"))

		} else {
			_, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Write([]byte("Signed out"))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

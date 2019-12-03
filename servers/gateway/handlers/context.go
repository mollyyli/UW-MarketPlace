package handlers

import (
	fakeusers "UW-Marketplace/servers/gateway/models/fake_users"
	"UW-Marketplace/servers/gateway/models/users"
	"UW-Marketplace/servers/gateway/sessions"
	"time"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

type Context struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	SockStore    SocketStore
}

func InitializeFake() *Context {
	sessStore := sessions.NewMemStore(time.Hour, time.Hour)
	usrStore, _ := fakeusers.ConnectToFakeDB()
	return &Context{
		SigningKey:   "test",
		SessionStore: sessStore,
		UserStore:    usrStore,
	}
}

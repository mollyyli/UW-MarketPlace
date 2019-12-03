package handlers

import (
	"UW-Marketplace/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

// SessionState contians the sessionTime for the user
type SessionState struct {
	SessionTime time.Time  `json:"session_ime"`
	User        users.User `json:"user"`
}

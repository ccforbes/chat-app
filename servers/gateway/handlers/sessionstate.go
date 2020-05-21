package handlers

import (
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/models/users"
)

//SessionState represents the state of a User's session
type SessionState struct {
	BeginTime time.Time   `json:"beginTime,omitempty"`
	User      *users.User `json:"user,omitempty"`
}

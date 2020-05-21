package handlers

import (
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/indexes"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/sessions"
)

//HandlerCtx allows for access to the SessionStore
//and the UserStore
type HandlerCtx struct {
	SigningKey   string
	SessionStore *sessions.RedisStore
	UserStore    *users.SQLStore
	SearchTree   *indexes.TrieNode
	Notifier     *Notifier
}

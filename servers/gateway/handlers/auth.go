package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/sessions"
)

//UsersHandler handles requests for the "users" resource
func (ctx *HandlerCtx) UsersHandler(w http.ResponseWriter, r *http.Request) {
	//If request method is POST
	if r.Method == http.MethodPost {
		// checks Content-Type value
		ctype := r.Header.Get("Content-Type")
		if !strings.Contains(ctype, "application/json") {
			http.Error(w, "The request body must be JSON", http.StatusUnsupportedMediaType)
			return
		}

		// Decodes JSON into a NewUser struct. Adds to database
		// after validation
		var nu users.NewUser
		err := json.NewDecoder(r.Body).Decode(&nu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := nu.ToUser()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		testEmail, _ := ctx.UserStore.GetByEmail(user.Email)
		if testEmail.Email == user.Email {
			http.Error(w, "Email taken", http.StatusBadRequest)
			return
		}
		testUserName, _ := ctx.UserStore.GetByUserName(user.UserName)
		if testUserName.UserName == user.UserName {
			http.Error(w, "UserName taken", http.StatusBadRequest)
			return
		}
		insertedUser, err := ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx.updateSearchTree(user.UserName, user.ID, "add")
		ctx.updateSearchTree(user.FirstName, user.ID, "add")
		ctx.updateSearchTree(user.LastName, user.ID, "add")
		ctx.updateSearchTree(strings.Join(strings.Split(user.FullName(), " "), ""), user.ID, "add")

		// Begin new session with the user
		newSession := &SessionState{
			BeginTime: time.Now(),
			User:      user,
		}
		sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, &newSession, w)

		// Response to client:
		// Status Code 201, Content-Type=applicaiton/json, encoded JSON object
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		testUser, err := ctx.UserStore.GetByID(user.ID)
		if err != nil {
			http.Error(w, "User was not correctly created", http.StatusInternalServerError)
			return
		}
		if !reflect.DeepEqual(insertedUser, testUser) {
			http.Error(w, "Retrieved incorrect information containing user", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "Unable to send user information", http.StatusBadGateway)
			return
		}
		return
	}
	if r.Method == http.MethodGet {
		ctx.SearchHandler(w, r)
		return
	}
	http.Error(w, "The HTTP request is not allowed.", http.StatusMethodNotAllowed)
	return

}

//SpecificUserHandler handles requests for a specific user.
func (ctx *HandlerCtx) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	// checks for authentication
	var sessionState SessionState
	sid, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &sessionState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	uid := sessionState.User.ID
	firstName := sessionState.User.FirstName
	lastName := sessionState.User.LastName
	fullName := sessionState.User.FullName()

	// If method is GET
	if r.Method == http.MethodGet {
		// get the user profile from the database
		paths := strings.Split(r.URL.Path, "/")

		// added me resource
		resource := paths[len(paths)-1]
		if resource != "me" {
			uid, _ = strconv.ParseInt(resource, 10, 64)
		}

		user, err := ctx.UserStore.GetByID(uid)
		// if user not found, 404
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// client response
		// Status Code 200, Content-Type=application/json, encoded JSON object
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Unable to send user information", http.StatusBadGateway)
			return
		}

		// if method is PATCH
	} else if r.Method == http.MethodPatch {
		// if URL resource is not "me" or authenticated user ID
		paths := strings.Split(r.URL.Path, "/")
		resource := paths[len(paths)-1]
		if resource != "me" {
			requestUID, _ := strconv.ParseInt(resource, 10, 64)
			if uid != requestUID {
				http.Error(w, "Unable to update user information", http.StatusForbidden)
				return
			}
		}

		// check if Content-Type is application/json
		ctype := r.Header.Get("Content-Type")
		if !strings.Contains(ctype, "application/json") {
			http.Error(w, "The request body must be JSON", http.StatusUnsupportedMediaType)
			return
		}

		// udpate user profile
		var updates users.Updates
		err = json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			http.Error(w, "Unable to update user", http.StatusInternalServerError)
			return
		}
		err = sessionState.User.ApplyUpdates(&updates)
		if err != nil {
			http.Error(w, "Unable to update user", http.StatusInternalServerError)
			return
		}
		ctx.SessionStore.Save(sid, sessionState)
		var updatedSession SessionState
		err := ctx.SessionStore.Get(sid, &updatedSession)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		updatedUser := updatedSession.User

		_, err = ctx.UserStore.Update(uid, &updates)
		if err != nil {
			http.Error(w, "Unable to update user", http.StatusBadRequest)
			return
		}

		// update search tree
		if updates.FirstName != "" {
			ctx.updateSearchTree(firstName, uid, "remove")
			ctx.updateSearchTree(updates.FirstName, uid, "add")
		}
		if updates.LastName != "" {
			ctx.updateSearchTree(lastName, uid, "remove")
			ctx.updateSearchTree(updates.LastName, uid, "add")
		}
		ctx.updateSearchTree(strings.Join(strings.Split(fullName, " "), ""), uid, "remove")
		ctx.updateSearchTree(strings.Join(strings.Split(updatedUser.FullName(), " "), ""), uid, "add")

		// response to client
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(sessionState.User); err != nil {
			http.Error(w, "Unable to send user information", http.StatusBadGateway)
			return
		}

		// if method is something else
	} else {
		http.Error(w, "The HTTP request is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

//SessionsHandler handles requests for the "sessions" resource and allows clients
// to begin a new session using an existing user's credentials
func (ctx *HandlerCtx) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	// If method is NOT POST
	if r.Method != http.MethodPost {
		http.Error(w, "The HTTP request was not allowed", http.StatusMethodNotAllowed)
		return
	}

	// checks value for "Content-Type"
	ctype := r.Header.Get("Content-Type")
	if !strings.Contains(ctype, "application/json") {
		http.Error(w, "The request body must be JSON", http.StatusUnsupportedMediaType)
		return
	}

	// decode body into users.Credentials
	var credentials users.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// find user in store
	user, err := ctx.UserStore.GetByEmail(credentials.Email)
	// if user not found
	if err != nil {
		time.Sleep(1 * time.Second)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// if user not authenticated
	err = user.Authenticate(credentials.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Begin new session with the user
	newSession := &SessionState{
		BeginTime: time.Now(),
		User:      user,
	}
	sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, &newSession, w)

	// response to client
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Unable to send user information", http.StatusBadGateway)
		return
	}
}

//SpecificSessionHandler handles requests related to a specific authenticated session
func (ctx *HandlerCtx) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	// if method is DELETE
	if r.Method == http.MethodDelete {
		if !strings.HasSuffix(r.URL.Path, "mine") {
			http.Error(w, "Unable to end session", http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "The HTTP request is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

//SearchHandler handles requests search queries from an authenticated user
func (ctx *HandlerCtx) SearchHandler(w http.ResponseWriter, r *http.Request) {
	// checks for authentication
	var sessionState SessionState
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &sessionState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// checks the query parameter
	query := r.URL.Query().Get("q")
	query = strings.ToLower(strings.Join(strings.Split(query, " "), ""))
	if query == "" {
		http.Error(w, "No users searched. Requires search token.", http.StatusBadRequest)
	}

	// find users in search tree based on query parameter
	resultsList := ctx.SearchTree.Find(query, 20)
	IDSet := make(map[int64]struct{})
	resultUsers := []*users.User{}
	for _, uid := range resultsList {
		user, err := ctx.UserStore.GetByID(uid)
		log.Printf("User ID to be added: %d", uid)
		_, added := IDSet[uid]
		IDSet[uid] = struct{}{}
		log.Printf("In the set: %t", added)
		if err != nil {
			http.Error(w, "Could not load user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !added {
			resultUsers = append(resultUsers, user)
			log.Print("added")
		}
	}

	// sort by username
	sort.Slice(resultUsers, func(i, j int) bool {
		return resultUsers[i].UserName < resultUsers[j].UserName
	})

	for _, user := range resultUsers {
		log.Print(user.UserName)
	}

	// write the results
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resultUsers); err != nil {
		http.Error(w, "Unable to load users: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctx *HandlerCtx) updateSearchTree(name string, uid int64, method string) {
	newName := strings.Split(strings.ToLower(name), " ")
	for _, n := range newName {
		if method == "add" {
			ctx.SearchTree.Add(n, uid)
		}
		if method == "remove" {
			ctx.SearchTree.Remove(n, uid)
		}
	}
}

package users

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/indexes"
)

// SQLStore represents users.Store
type SQLStore struct {
	db *sql.DB
}

//NewSQLStore constructs a new SQLStore
func NewSQLStore(db *sql.DB) *SQLStore {
	if db == nil {
		panic("nil database pointer passed to NewSqlStore")
	}
	return &SQLStore{
		db: db,
	}
}

//Constants for SQL queries and statements
const sqlColumnListNoID = "email,passHash,username,firstName,lastName,photoUrl"

//SQLSelectFromQuery is a select query with included id column
const SQLSelectFromQuery = "select id," + sqlColumnListNoID + " from Users where "
const sqlInsertUser = "insert into Users(" + sqlColumnListNoID + ") values(?,?,?,?,?,?)"

//SQLUpdateUser is an update statement with both names
const SQLUpdateUser = "update Users set firstName=?,lastName=? where id=?"

//SQLDeleteUser is a delete statement
const SQLDeleteUser = "delete from Users where id=?"

const sqlColumnListNamesAndID = "select id, username, firstName, lastName from Users"

//GetByID returns the User with the given ID
func (s *SQLStore) GetByID(id int64) (*User, error) {
	rows, err := s.db.Query(SQLSelectFromQuery+"id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user, err := getUser(rows)
	if user.ID != id {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//GetByEmail returns the User with the given email
func (s *SQLStore) GetByEmail(email string) (*User, error) {
	rows, err := s.db.Query(SQLSelectFromQuery+"email = ?", email)
	defer rows.Close()
	if err != nil {
		return nil, ErrUserNotFound
	}
	return getUser(rows)
}

//GetByUserName returns the User with the given Username
func (s *SQLStore) GetByUserName(username string) (*User, error) {
	rows, err := s.db.Query(SQLSelectFromQuery+"username = ?", username)
	defer rows.Close()
	if err != nil {
		return nil, ErrUserNotFound
	}
	return getUser(rows)
}

// Helper function to scan rows into a single User
// returns user and a potential error
func getUser(rows *sql.Rows) (*User, error) {
	user := &User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, ErrUserNotFound
		}
	}
	if err := rows.Err(); err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (s *SQLStore) Insert(user *User) (*User, error) {
	result, err := s.db.Exec(sqlInsertUser, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		return nil, errors.New("Unable to insert user into database")
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("Unable to unser user into database")
	}

	user.ID = newID
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (s *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	_, err := s.db.Exec(SQLUpdateUser, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, errors.New("Unable to update user")
	}
	return s.GetByID(id)
}

//Delete deletes the user with the given ID
func (s *SQLStore) Delete(id int64) error {
	_, err := s.db.Exec(SQLDeleteUser, id)
	if err != nil {
		return errors.New("Unable to delete user")
	}
	return nil
}

//LoadExistingUsers loads all existing users into the search tree
func (s *SQLStore) LoadExistingUsers() (*indexes.TrieNode, error) {
	users := &struct {
		ID        int64
		UserName  string
		FirstName string
		LastName  string
	}{}
	newTrie := indexes.NewTrieNode()

	// get all users in db
	rows, err := s.db.Query(sqlColumnListNamesAndID)
	if err != nil {
		return nil, errors.New("Unable to find users")
	}

	// scan each row
	for rows.Next() {
		if err := rows.Scan(&users.ID, &users.UserName,
			&users.FirstName, &users.LastName); err != nil {
			return nil, errors.New("Unable to load users")
		}

		// to lowercase
		users.UserName = strings.ToLower(users.UserName)
		users.FirstName = strings.ToLower(users.FirstName)
		users.LastName = strings.ToLower(users.LastName)

		// add each name. find fields with spaces
		newTrie.Add(users.UserName, users.ID)
		firstName := strings.Split(users.FirstName, " ")
		for _, name := range firstName {
			newTrie.Add(name, users.ID)
		}
		lastName := strings.Split(users.LastName, " ")
		for _, name := range lastName {
			newTrie.Add(name, users.ID)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, ErrUserNotFound
	}

	return newTrie, err
}

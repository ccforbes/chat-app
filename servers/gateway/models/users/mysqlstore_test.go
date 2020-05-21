package users

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID is a test function for the SQLStore's GetByID
func TestGetByID(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		// TODO: update based on the name of your type struct
		mainSQLStore := &SQLStore{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// TODO: update to match the query used in your Store implementation
		query := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where id=?"

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating sqlmock: %v", err)
	}

	mainSQLStore := NewSQLStore(db)

	validEmail := "test@test.com"

	r := mock.NewRows([]string{"id", "email", "passHash", "username", "firstName", "lastName", "photoUrl"}).
		AddRow(1, "test@test.com", nil, "testing", "test", "tester", "test")
	query := SQLSelectFromQuery + "email=?"
	mock.ExpectQuery(query).WithArgs(validEmail).WillReturnRows(r)

	user, err := mainSQLStore.GetByEmail(validEmail)
	if err != nil {
		t.Errorf("unexpected error during successful query: %v", err)
	}

	if user == nil {
		t.Errorf("nil user returned from GetByEmail")
	} else if user.Email != validEmail {
		t.Errorf("incorrect email: expected %s but got %s", validEmail, user.Email)
	}
}

func TestGetByUserName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating sqlmock: %v", err)
	}

	mainSQLStore := NewSQLStore(db)

	validUsername := "testing"

	r := mock.NewRows([]string{"id", "email", "passHash", "username", "firstName", "lastName", "photoUrl"}).
		AddRow(1, "test@test.com", nil, "testing", "test", "tester", "test")
	query := SQLSelectFromQuery + "username=?"
	mock.ExpectQuery(query).WithArgs(validUsername).WillReturnRows(r)

	user, err := mainSQLStore.GetByUserName(validUsername)
	if err != nil {
		t.Errorf("unexpected error during successful query: %v", err)
	}

	if user == nil {
		t.Error("nil user returned from GetByUserName")
	} else if user.UserName != validUsername {
		t.Errorf("incorrect username: expected %s but got %s", validUsername, user.UserName)
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating sqlmock: %v", err)
	}

	mainSQLStore := NewSQLStore(db)

	newUser := NewUser{
		Email:        "test@test.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "testing",
		FirstName:    "test",
		LastName:     "tester",
	}

	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("cannot convert to user: %v", err)
	}

	expectedSQL := regexp.QuoteMeta(sqlInsertUser)

	var newID int64 = 1

	mock.ExpectExec(expectedSQL).
		WithArgs(
			user.Email,
			user.PassHash,
			user.UserName,
			user.FirstName,
			user.LastName,
			user.PhotoURL,
		).
		WillReturnResult(sqlmock.NewResult(newID, 1))

	insertedUser, err := mainSQLStore.Insert(user)
	if err != nil {
		t.Errorf("unexpected error during successful insert: %v", err)
	}

	if insertedUser == nil {
		t.Error("nil user returned from insert")
	} else if insertedUser.ID != newID {
		t.Errorf("incorrect new ID: expected %d but got %d", newID, insertedUser.ID)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating sqlmock: %v", err)
	}

	mainSQLStore := NewSQLStore(db)

	testUpdate := &Updates{
		FirstName: "Stuart",
		LastName:  "Reges",
	}

	expectedSQL := regexp.QuoteMeta(SQLUpdateUser)

	var newID int64 = 1

	r := mock.NewRows([]string{"id", "email", "passHash", "username", "firstName", "lastName", "photoUrl"}).
		AddRow(1, "test@test.com", nil, "testing", "Stuart", "Reges", "test")

	mock.ExpectExec(expectedSQL).
		WithArgs(
			testUpdate.FirstName,
			testUpdate.LastName,
			newID,
		).
		WillReturnResult(sqlmock.NewResult(newID, 1))

	mock.ExpectQuery(SQLSelectFromQuery).
		WillReturnRows(r)

	updatedUser, err := mainSQLStore.Update(newID, testUpdate)
	if err != nil {
		t.Errorf("unexpected error during successful update: %v", err)
	}

	if updatedUser == nil {
		t.Error("nil user returned from insert")
	} else if updatedUser.ID != newID {
		t.Errorf("incorrect ID: expected %d but got %d", newID, updatedUser.ID)
	} else if updatedUser.FirstName != testUpdate.FirstName {
		t.Errorf("incorrect first name update: expected %s but got %s", testUpdate.FirstName, updatedUser.FirstName)
	} else if updatedUser.LastName != testUpdate.LastName {
		t.Errorf("incorrect last name update: expected %s but got %s", testUpdate.LastName, updatedUser.LastName)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating sqlmock: %v", err)
	}

	mainSQLStore := NewSQLStore(db)

	var id int64 = 1

	mock.NewRows([]string{"id", "email", "passHash", "username", "firstName", "lastName", "photoUrl"}).
		AddRow(1, "test@test.com", nil, "testing", "test", "tester", "test")

	mock.ExpectExec(SQLDeleteUser).WithArgs(id).WillReturnResult(sqlmock.NewResult(id, 1))

	err = mainSQLStore.Delete(1)
	if err != nil {
		t.Errorf("unexpected error occurred on a successful delete case")
	}
}

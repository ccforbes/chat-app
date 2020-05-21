package users

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
import (
	"crypto/md5"
	"encoding/base64"
	"net/mail"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name        string
		nu          NewUser
		hint        string
		expectError bool
	}{
		{
			"Invalid Email Address",
			NewUser{
				Email: "invalid email address",
			},
			"Check the email address if is a valid one",
			true,
		},
		{
			"Short Password",
			NewUser{
				Email:    "test@mail.com",
				Password: "fails",
			},
			"Check the length of the password",
			true,
		},
		{
			"Unmatching Passwords",
			NewUser{
				Email:        "test@mail.com",
				Password:     "password",
				PasswordConf: "doesnotmatch",
			},
			"Check to see if the passwords match",
			true,
		},
		{
			"Empty UserName",
			NewUser{
				Email:        "test@mail.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "",
			},
			"Make sure that the UserName is not empty",
			true,
		},
		{
			"Spaces in UserName",
			NewUser{
				Email:        "test@mail.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "has spaces",
			},
			"Make sure that the UserName does not contain spaces",
			true,
		},
		{
			"Valid New User",
			NewUser{
				Email:        "test@mail.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "testing123",
				FirstName:    "Valid",
				LastName:     "User",
			},
			"The user is valid",
			false,
		},
	}

	for _, c := range cases {
		err := c.nu.Validate()
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error validating NewUser: %v\nHINT: %s", c.name, err, c.hint)
		}
		if err == nil {
			if _, emailErr := mail.ParseAddress(c.nu.Email); emailErr != nil {
				t.Errorf("case %s: email is invalid\nHINT: %s", c.name, c.hint)
			}
			if len(c.nu.Password) < 6 {
				t.Errorf("case %s: password is less than 6 chars\nHINT: %s", c.name, c.hint)
			}
			if c.nu.Password != c.nu.PasswordConf {
				t.Errorf("case %s: password do not match\nHINT: %s", c.name, c.hint)
			}
			if len(c.nu.UserName) == 0 {
				t.Errorf("case %s: username is zero-length\nHINT: %s", c.name, c.hint)
			}
			if strings.Contains(c.nu.UserName, " ") {
				t.Errorf("case %s: username contains spaces\nHINT: %s", c.name, c.hint)
			}
		}
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		Name string
		Hint string
		nu   NewUser
	}{
		{
			"Regular Email",
			"The email may be incorrectly hashed",
			NewUser{
				Email:        "test@mail.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "testing123",
				FirstName:    "Valid",
				LastName:     "User",
			},
		},
		{
			"Email has uppercase",
			"Make sure that the email is all lowercase",
			NewUser{
				Email:        "TEST@mail.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "testing123",
				FirstName:    "Valid",
				LastName:     "User",
			},
		},
		{
			"Email has spaces",
			"Make sure that the email has no trailing or leading spaces",
			NewUser{
				Email:        " TEST@mail.com ",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "testing123",
				FirstName:    "Valid",
				LastName:     "User",
			},
		},
	}

	for _, c := range cases {
		u, err := c.nu.ToUser()
		if err != nil {
			t.Errorf("case %s: unexpeted error, the NewUser should be validated", c.Name)
		}

		h := md5.New()
		h.Write([]byte(c.nu.Email))
		photoHash := base64.URLEncoding.EncodeToString(h.Sum(nil))
		expectedPhotoURL := gravatarBasePhotoURL + photoHash

		if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(c.nu.Password)); err != nil {
			t.Errorf("case %s: password was not correctly hashed", c.Name)
		}

		if u.PhotoURL != expectedPhotoURL {
			t.Errorf("case %s: photo URL should be: %s\nHINT: %s", c.Name, expectedPhotoURL, c.Hint)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		Name           string
		Hint           string
		u              User
		expectedOutput string
	}{
		{
			"Both Fields Set",
			"The name should return as normal: <first name> <last name>",
			User{
				FirstName: "Bob",
				LastName:  "Smith",
			},
			"Bob Smith",
		},
		{
			"Neither Field Set",
			"First and last name is empty",
			User{},
			"",
		},
		{
			"No First Name",
			"The last name should be returned",
			User{
				LastName: "Smith",
			},
			"Smith",
		},
		{
			"No Last Name",
			"The first name should be returned",
			User{
				FirstName: "Bob",
			},
			"Bob",
		},
	}

	for _, c := range cases {
		if c.u.FullName() != c.expectedOutput {
			t.Errorf("case %s: full name is %s when it should be %s\nHINT: %s", c.Name, c.u.FullName(), c.expectedOutput, c.Hint)
		}
	}
}

func TestAuthenicate(t *testing.T) {
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password"), 13)
	cases := []struct {
		Name          string
		Hint          string
		u             User
		InputPassword string
		expectedError bool
	}{
		{
			"Correct Password",
			"The password is correct",
			User{
				PassHash: passHash,
			},
			"password",
			false,
		},
		{
			"Incorrect Password",
			"The password is incorrect",
			User{
				PassHash: passHash,
			},
			"incorrect",
			true,
		},
		{
			"Empty Password",
			"The password is empty",
			User{
				PassHash: passHash,
			},
			"",
			true,
		},
	}

	for _, c := range cases {
		err := c.u.Authenticate(c.InputPassword)
		if err != nil && !c.expectedError {
			t.Errorf("case %s: unexpected error authenticating user: %v\nHINT: %s", c.Name, err, c.Hint)
		}
		if err == nil && c.expectedError {
			t.Errorf("case %s: the user should not be authenticated\nHINT: %s", c.Name, c.Hint)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		Name           string
		Hint           string
		u              User
		updates        Updates
		expectedOutput User
		expectedError  bool
	}{
		{
			"Both Fields Set",
			"Both first and last names should be updated",
			User{
				FirstName: "Bob",
				LastName:  "Smith",
			},
			Updates{
				FirstName: "Stuart",
				LastName:  "Reges",
			},
			User{
				FirstName: "Stuart",
				LastName:  "Reges",
			},
			false,
		},
		{
			"First Name Set",
			"First Name should be updated",
			User{
				FirstName: "Bob",
				LastName:  "Smith",
			},
			Updates{
				FirstName: "Stuart",
			},
			User{
				FirstName: "Stuart",
				LastName:  "Smith",
			},
			false,
		},
		{
			"Last Name Set",
			"Last Name should be updated",
			User{
				FirstName: "Bob",
				LastName:  "Smith",
			},
			Updates{
				LastName: "Reges",
			},
			User{
				FirstName: "Bob",
				LastName:  "Reges",
			},
			false,
		},
		{
			"Fields Not Set",
			"No updates should take place",
			User{
				FirstName: "Bob",
				LastName:  "Smith",
			},
			Updates{},
			User{
				FirstName: "Bob",
				LastName:  "Smith",
			},
			true,
		},
	}

	for _, c := range cases {
		err := c.u.ApplyUpdates(&c.updates)
		if err != nil && !c.expectedError {
			t.Errorf("case %s: unexpected updating user: %v\nHINT: %s", c.Name, err, c.Hint)
		}
		if err == nil && c.expectedError {
			t.Errorf("case %s: the user was updated when the fields were empty: %v\nHINT: %s", c.Name, err, c.Hint)
		}
		if c.u.FullName() != c.expectedOutput.FullName() {
			t.Errorf("case %s: the user was not properly updated\nHINT: %s", c.Name, c.Hint)
		}
	}
}

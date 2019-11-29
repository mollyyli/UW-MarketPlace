package users

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"

	// "log"

	"golang.org/x/crypto/bcrypt"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
func TestValidate(t *testing.T) {
	users := []NewUser{
		{
			"Invalid email",
			"1234567",
			"1234567",
			"asdf",
			"firstname",
			"lastname",
		},

		{
			"user@example.com",
			"12345",
			"12345",
			"asdf",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123457",
			"asdf",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123456",
			"",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123456",
			"asd f",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123456",
			"asdf",
			"firstname",
			"lastname",
		},
	}
	cases := []struct {
		user        NewUser
		hint        string
		expectError bool
	}{
		{
			users[0],
			"Remember to return an error when the email is invalid",
			true,
		},
		{
			users[1],
			"Remember to return an error if password is not long enough",
			true,
		},
		{
			users[2],
			"Remember to return an error when the password and passwordconf don't match",
			true,
		},
		{
			users[3],
			"Remember to return an error when the username is empty",
			true,
		},
		{
			users[4],
			"Remember to return an error when the username contains any empty space",
			true,
		},
		{
			users[5],
			"Remember not to return an error when all the information is valid",
			false,
		},
	}

	for _, c := range cases {
		err := c.user.Validate()
		if err != nil && !c.expectError {
			t.Errorf("unexpected error validating new user: %v\nHINT: %s", err, c.hint)
		}
	}
}

func TestToUser(t *testing.T) {
	newuser := []NewUser{
		{
			"Invalid email",
			"1234567",
			"1234567",
			"asdf",
			"firstname",
			"lastname",
		},

		{
			"user@example.com",
			"12345",
			"12345",
			"asdf",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123457",
			"asdf",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123456",
			"",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123456",
			"asd f",
			"firstname",
			"lastname",
		},
		{
			"user@example.com",
			"123456",
			"123456",
			"asdf",
			"firstname",
			"lastname",
		},
	}

	cases := []struct {
		newUser     NewUser
		hint        string
		expectError bool
	}{
		{
			newuser[0],
			"Remember to return an error if NewUser cannot be validated",
			true,
		},
		{
			newuser[1],
			"Remember to return an error if NewUser cannot be validated",
			true,
		},
		{
			newuser[2],
			"Remember to return an error if NewUser cannot be validated",
			true,
		},
		{
			newuser[3],
			"Remember to return an error if NewUser cannot be validated",
			true,
		},
		{
			newuser[4],
			"Remember to return an error if NewUser cannot be validated",
			true,
		},
		{
			newuser[5],
			"Remember not to return an error when all the information is valid",
			false,
		},
	}
	for _, c := range cases {
		user, err := c.newUser.ToUser()
		if err != nil && !c.expectError {
			t.Errorf("unexpected error validating new user: %v\nHINT: %s", err, c.hint)
		}
		if err == nil {

			hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), 13)

			hash := md5.New()
			hash.Write([]byte(strings.Trim(strings.ToLower("user@example.com"), " ")))
			photoURL := gravatarBasePhotoURL + hex.EncodeToString(hash.Sum(nil))
			users := []*User{
				{
					0,
					"user@example.com",
					hashed,
					"asdf",
					"firstname",
					"lastname",
					photoURL,
				},
			}

			error2 := bcrypt.CompareHashAndPassword(user.PassHash, []byte("123456"))
			if error2 != nil {
				t.Errorf(strings.Split(error2.Error(), ": ")[1])
			}

			if user.ID != users[0].ID {
				t.Errorf("Invalid user returned")
			}

			if user.Email != users[0].Email {
				t.Errorf("Invalid user returned")
			}

			if user.UserName != users[0].UserName {
				t.Errorf("Invalid user returned")
			}

			if user.FirstName != users[0].FirstName {
				t.Errorf("Invalid user returned")
			}

			if user.LastName != users[0].LastName {
				t.Errorf("Invalid user returned")
			}

			if user.PhotoURL != users[0].PhotoURL {
				t.Errorf("Invalid user returned")
			}
		}
	}
}

func TestFullName(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), 13)

	users := []User{
		{
			0,
			"user@example.com",
			hashed,
			"asdf",
			"",
			"",
			"https://www.gravatar.com/avatar/user@example.com",
		},
		{
			0,
			"user@example.com",
			hashed,
			"asdf",
			"",
			"lastname",
			"https://www.gravatar.com/avatar/user@example.com",
		},
		{
			0,
			"user@example.com",
			hashed,
			"asdf",
			"firstname",
			"",
			"https://www.gravatar.com/avatar/user@example.com",
		},
		{
			0,
			"user@example.com",
			hashed,
			"asdf",
			"firstname",
			"lastname",
			"https://www.gravatar.com/avatar/user@example.com",
		},
	}
	cases := []struct {
		user        User
		hint        string
		fullname    string
		expectError bool
	}{
		{
			users[0],
			"Empty firstname and last name",
			"",
			true,
		},
		{
			users[1],
			"Empty firstname",
			users[1].LastName,
			true,
		},
		{
			users[2],
			"Empty lastname",
			users[2].FirstName,
			true,
		},
		{
			users[3],
			"Valid fullname",
			users[3].FirstName + " " + users[3].LastName,
			false,
		},
	}
	for _, c := range cases {
		name := c.user.FullName()

		if name != c.fullname {
			t.Errorf("Name is invalid")
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	users := []User{
		{
			0,
			"user@example.com",
			[]byte("123456"),
			"asdf",
			"firstname",
			"lastname",
			"https://www.gravatar.com/avatar/user@example.com",
		},
	}
	updates := []*Updates{
		{
			"newfirstname",
			"newlastname",
		},
	}
	cases := []struct {
		user        User
		update      *Updates
		expectError bool
	}{
		{
			users[0],
			updates[0],
			false,
		},
	}
	for _, c := range cases {
		c.user.ApplyUpdates(c.update)
		if c.user.FirstName != c.update.FirstName {
			t.Errorf("First name did not update")
		}

		if c.user.LastName != c.update.LastName {
			t.Errorf("Last name did not update")
		}
	}
}

func TestAuthenticate(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcryptCost)
	users := []User{
		{
			0,
			"user@example.com",
			hashed,
			"asdf",
			"firstname",
			"lastname",
			"https://www.gravatar.com/avatar/user@example.com",
		},
	}

	cases := []struct {
		user User
	}{
		{
			users[0],
		},
	}

	for _, c := range cases {
		err := c.user.Authenticate("123456")
		if err != nil {
			t.Errorf("Could not authenticate")
		}
		err = c.user.Authenticate("1234567")
		if err == nil {
			t.Errorf("Expected error")
		}
	}
}

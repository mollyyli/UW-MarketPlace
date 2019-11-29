package users

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	//TODO: validate the new user according to these rules:
	//- Email field must be a valid email address (hint: see mail.ParseAddress)
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return err
	}

	//- Password must be at least 6 characters
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password not long enough")
	}
	//- Password and PasswordConf must match

	if strings.Compare(nu.Password, nu.PasswordConf) != 0 {
		return fmt.Errorf("Password and PasswordConf don't match")
	}

	if len(nu.UserName) <= 0 {
		return fmt.Errorf("Username must be non-zero length")
	}

	if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("Username must not contain spaces")
	}
	//- UserName must be non-zero length and may not contain spaces
	//use fmt.Errorf() to generate appropriate error messages if
	//the new user doesn't pass one of the validation rules

	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//TODO: call Validate() to validate the NewUser and
	//return any validation errors that may occur.
	if nu.Validate() != nil {
		return nil, nu.Validate()
	}

	// url := md5.Sum([]byte(strings.Trim(strings.ToLower(nu.Email), " ")))
	hash := md5.New()
	hash.Write([]byte(strings.Trim(strings.ToLower(nu.Email), " ")))
	photoURL := gravatarBasePhotoURL + hex.EncodeToString(hash.Sum(nil))
	user := User{ID: 0, Email: nu.Email, UserName: nu.UserName, FirstName: nu.FirstName, LastName: nu.LastName, PhotoURL: photoURL}
	user.SetPassword(nu.Password)
	return &user, nil
	//if valid, create a new *User and set the fields
	//based on the field values in `nu`.
	//Leave the ID field as the zero-value; your Store
	//implementation will set that field to the DBMS-assigned
	//primary key value.
	//Set the PhotoURL field to the Gravatar PhotoURL
	//for the user's email address.
	//see https://en.gravatar.com/site/implement/hash/
	//and https://en.gravatar.com/site/implement/images/

	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password

}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	//TODO: implement according to comment above
	if len(u.FirstName) == 0 && len(u.LastName) == 0 {
		return ""
	} else if len(u.FirstName) == 0 {
		return u.LastName
	} else if len(u.LastName) == 0 {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName

}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//TODO: use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err == nil {
		u.PassHash = hashed
	}
	return err
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt

	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	return err

}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	//TODO: set the fields of `u` to the values of the related
	//field in the `updates` struct
	_, err := json.Marshal(updates)
	if err == nil {
		u.FirstName = updates.FirstName
		u.LastName = updates.LastName
	}

	return err
}

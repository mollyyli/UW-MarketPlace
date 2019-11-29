package fakeusers

import (
	"assignments-thebriando/servers/gateway/models/users"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// FakeMySQLConnection contains a fake "database"
type FakeMySQLConnection struct {
	db       []*users.User
	signInDB []*users.UserSignIn
}

// ConnectToFakeDB connects to the fake database
func ConnectToFakeDB() (*FakeMySQLConnection, error) {
	passHash, _ := bcrypt.GenerateFromPassword([]byte("123456"), 13)
	fakeUsers := []*users.User{
		&users.User{
			ID:        1,
			Email:     "example@example.com",
			PassHash:  passHash,
			UserName:  "aUser",
			FirstName: "firstName",
			LastName:  "lastName",
			PhotoURL:  "https://www.gravatar.com/avatar/user@example.com",
		},
		&users.User{
			ID:        2,
			Email:     "example1@example.com",
			PassHash:  passHash,
			UserName:  "aUser1",
			FirstName: "firstName1",
			LastName:  "lastName1",
			PhotoURL:  "https://www.gravatar.com/avatar/user1@example.com",
		},
		&users.User{
			ID:        3,
			Email:     "example2@example.com",
			PassHash:  passHash,
			UserName:  "aUser2",
			FirstName: "firstName2",
			LastName:  "lastName2",
			PhotoURL:  "https://www.gravatar.com/avatar/user2@example.com",
		},
	}
	signIns := []*users.UserSignIn{
		&users.UserSignIn{
			ID:         "1",
			SignInTime: "07/01/2019 10:08:24 AM",
			IP:         "155.187.132.4",
		},
		&users.UserSignIn{
			ID:         "2",
			SignInTime: "07/02/2019 10:08:24 AM",
			IP:         "155.187.132.5",
		},
		&users.UserSignIn{
			ID:         "3",
			SignInTime: "07/03/2019 10:08:24 AM",
			IP:         "155.187.132.6",
		},
	}
	mysql := FakeMySQLConnection{
		db:       fakeUsers,
		signInDB: signIns,
	}
	return &mysql, nil
}

// ErrUserNotFound returns a "user not found error"
var ErrUserNotFound = errors.New("user not found")

//GetByID returns entire row from fake db based on given ID
func (mysql *FakeMySQLConnection) GetByID(id int64) (*users.User, error) {
	for _, user := range mysql.db {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

//GetByEmail returns entire row from fake db based on given email
func (mysql *FakeMySQLConnection) GetByEmail(email string) (*users.User, error) {
	for _, user := range mysql.db {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

//GetByUserName returns entire row from fake db based on given username
func (mysql *FakeMySQLConnection) GetByUserName(username string) (*users.User, error) {
	for _, user := range mysql.db {
		if user.UserName == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

//Insert inserts user into db
func (mysql *FakeMySQLConnection) Insert(user *users.User) (*users.User, error) {
	mysql.db = append(mysql.db, user)
	user.ID = mysql.db[len(mysql.db)-1].ID + 1
	return user, nil
}

// Update updates a user from ID using update struct
func (mysql *FakeMySQLConnection) Update(id int64, updates *users.Updates) (*users.User, error) {
	for _, user := range mysql.db {
		if user.ID == id {
			err := user.ApplyUpdates(updates)
			return user, err
		}
	}
	return nil, ErrUserNotFound
}

// Delete deletes the user at given id from db
func (mysql *FakeMySQLConnection) Delete(id int64) error {
	return nil
}

// InsertSignIn inserts a user log in
func (mysql *FakeMySQLConnection) InsertSignIn(signIn *users.UserSignIn) (*users.UserSignIn, error) {
	mysql.signInDB = append(mysql.signInDB, signIn)
	return signIn, nil
}

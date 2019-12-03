package users

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLConnection makes a connection to the MySQL databse
type MySQLConnection struct {
	Client *sql.DB
}

//GetByID returns entire row based on given ID
func (mysql *MySQLConnection) GetByID(id int64) (*User, error) {
	var user User
	get := "select * from Users where ID = ?"
	res := mysql.Client.QueryRow(get, strconv.FormatInt(id, 10))
	err := res.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
	if err != nil {
		err = ErrUserNotFound
	}
	return &user, err
}

//GetByEmail returns a row based on given email
func (mysql *MySQLConnection) GetByEmail(email string) (*User, error) {
	var user User
	get := "select * from Users where Email = ?"
	res := mysql.Client.QueryRow(get, email)
	err := res.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
	if err != nil {
		err = ErrUserNotFound
	}
	return &user, err
}

//GetByUserName returns the entire row based on given username
func (mysql *MySQLConnection) GetByUserName(username string) (*User, error) {
	var user User

	get := "select * from Users where Username = ?"
	res := mysql.Client.QueryRow(get, username)
	err := res.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
	if err != nil {
		err = ErrUserNotFound
	}
	return &user, err
}

func (mysql *MySQLConnection) Update(id int64, updates *Updates) (*User, error) {
	upd := "update Users set FirstName = ?, LastName = ? where ID = ?"
	mysql.Client.Exec(upd, updates.FirstName, updates.LastName, id)

	return mysql.GetByID(id)
}

//Delete the entire row based on given id
func (mysql *MySQLConnection) Delete(id int64) error {

	_, err := mysql.GetByID(id)
	if err != nil {
		return err
	}
	delete := "delete from Users where ID = ?"
	mysql.Client.Exec(delete, id)
	return nil
}

// Insert inserts data to the user struct
func (mysql *MySQLConnection) Insert(user *User) (*User, error) {
	insq := "insert into Users(Email, PassHash, Username, FirstName, LastName, PhotoURL) values (?,?,?,?,?,?)"
	res, err := mysql.Client.Exec(insq, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		log.Println("res error", err)
	}
	id, _ := res.LastInsertId()
	user.ID = id
	return user, nil
}

type UserSignIn struct {
	ID         int64
	UserID     int64
	SignInTime string
	IP         string
}

func (mysql *MySQLConnection) InsertSignIn(signIn *UserSignIn) (*UserSignIn, error) {
	insq := "insert into UserSignIns(UserID, SignInTime, IP) values (?,?,?)"
	res, err := mysql.Client.Exec(insq, signIn.UserID, signIn.SignInTime, signIn.IP)
	if err != nil {
		log.Println("res error", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		log.Println("Could not insert")
	}
	return signIn, nil
}

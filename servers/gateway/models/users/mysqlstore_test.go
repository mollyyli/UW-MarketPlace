package users

import (
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/DATA-DOG/go-sqlmock"
)

func TestGetByID(t *testing.T) {
	// struct user
	mockUser := &User{
		1,
		"example@example.com",
		[]byte("123456"),
		"aUser",
		"firstName",
		"lastName",
		"https://www.gravatar.com/avatar/user@example.com",
	}

	// mock data
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// defer db.Close()
	row := mock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}).
		AddRow(1, "example@example.com", []byte("123456"), "aUser", "firstName", "lastName", "https://www.gravatar.com/avatar/user@example.com")

	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where ID = ?")).WithArgs(strconv.FormatInt(mockUser.ID, 10)).WillReturnRows(row)
	mockContext := &MySQLConnection{Client: db}

	result, err := mockContext.GetByID(1)
	err = mock.ExpectationsWereMet()
	if err != nil {
		// throw err
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, mockUser) {
		t.Errorf("Wrong user returned")
	}

	// expect error with non-existing ID
	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where ID = ?")).WithArgs(strconv.FormatInt(999, 10)).WillReturnError(ErrUserNotFound)
	_, err = mockContext.GetByID(999)
	if err == nil {
		t.Errorf("Expected error: %v", ErrUserNotFound.Error())
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		// throw err
		t.Errorf(err.Error())
	}
}

func TestGetByEmail(t *testing.T) {
	mockUser := &User{
		1,
		"example@example.com",
		[]byte("123456"),
		"aUser",
		"firstName",
		"lastName",
		"https://www.gravatar.com/avatar/user@example.com",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	row := mock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}).
		AddRow(1, "example@example.com", []byte("123456"), "aUser", "firstName", "lastName", "https://www.gravatar.com/avatar/user@example.com")

	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where Email = ?")).WithArgs(mockUser.Email).WillReturnRows(row)
	mockContext := &MySQLConnection{Client: db}

	result, err := mockContext.GetByEmail(mockUser.Email)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, mockUser) {
		t.Errorf("Wrong user returned")
	}
	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where Email = ?")).WithArgs(mockUser.Email + "asdf").WillReturnError(ErrUserNotFound)
	_, err = mockContext.GetByEmail(mockUser.Email + "asdf")
	if err == nil {
		t.Errorf("Expected error: %v", ErrUserNotFound.Error())
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		// throw err
		t.Errorf(err.Error())
	}
}

func TestGetByUsername(t *testing.T) {
	mockUser := &User{
		1,
		"example@example.com",
		[]byte("123456"),
		"aUser",
		"firstName",
		"lastName",
		"https://www.gravatar.com/avatar/user@example.com",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	row := mock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}).
		AddRow(1, "example@example.com", []byte("123456"), "aUser", "firstName", "lastName", "https://www.gravatar.com/avatar/user@example.com")

	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where Username = ?")).WithArgs(mockUser.UserName).WillReturnRows(row)
	mockContext := &MySQLConnection{Client: db}

	result, err := mockContext.GetByUserName(mockUser.UserName)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(result, mockUser) {
		t.Errorf("Wrong user returned")
	}
	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where Username = ?")).WithArgs(mockUser.UserName + "asdf").WillReturnError(ErrUserNotFound)
	_, err = mockContext.GetByUserName(mockUser.UserName + "asdf")
	if err == nil {
		t.Errorf("Expected error: %v", ErrUserNotFound.Error())
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		// throw err
		t.Errorf(err.Error())
	}
}

func TestInsert(t *testing.T) {
	mockUser := &User{
		0,
		"example@example.com",
		[]byte("123456"),
		"aUser",
		"firstName",
		"lastName",
		"https://www.gravatar.com/avatar/user@example.com",
	}
	mockUser2 := &User{
		0,
		"example1@example.com",
		[]byte("123456"),
		"aUser",
		"firstName",
		"lastName",
		"https://www.gravatar.com/avatar/user@example.com",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"})
	mockContext := &MySQLConnection{Client: db}
	// mock.ExpectExec("insert into Users(Email, PassHash, Username, FirstName, LastName, PhotoURL) values (" + mockUser.Email + ", " + string(mockUser.PassHash) + ", " + mockUser.UserName + ", " + mockUser.FirstName + "," + mockUser.LastName + ", " + mockUser.PhotoURL + ")")
	mock.ExpectExec(regexp.QuoteMeta("insert into Users")).WithArgs(mockUser.Email, mockUser.PassHash, mockUser.UserName, mockUser.FirstName, mockUser.LastName, mockUser.PhotoURL).WillReturnResult(sqlmock.NewResult(1, 1))
	result, err := mockContext.Insert(mockUser)
	mock.ExpectationsWereMet()
	if !reflect.DeepEqual(result, mockUser) {
		t.Errorf("Insert failed")
	}

	mock.ExpectExec(regexp.QuoteMeta("insert into Users")).WithArgs(mockUser2.Email, mockUser2.PassHash, mockUser2.UserName, mockUser2.FirstName, mockUser2.LastName, mockUser2.PhotoURL).WillReturnResult(sqlmock.NewResult(2, 1))
	result, err = mockContext.Insert(mockUser2)
	mock.ExpectationsWereMet()
	if !reflect.DeepEqual(result, mockUser2) {
		t.Errorf("Insert failed")
	}
}

func TestUpdate(t *testing.T) {
	mockUser := &User{
		0,
		"example@example.com",
		[]byte("123456"),
		"aUser",
		"firstName",
		"lastName",
		"https://www.gravatar.com/avatar/user@example.com",
	}
	mockUpdate := &Updates{
		"newfirstname",
		"newlastname",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	row := mock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}).
		AddRow(0, "example@example.com", []byte("123456"), "aUser", "firstName", "lastName", "https://www.gravatar.com/avatar/user@example.com")
	mockContext := &MySQLConnection{Client: db}

	mock.ExpectExec(regexp.QuoteMeta("update Users set")).WithArgs(mockUpdate.FirstName, mockUpdate.LastName, 0).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where ID = ?")).WithArgs(strconv.FormatInt(mockUser.ID, 10)).WillReturnRows(row)

	result, err := mockContext.Update(mockUser.ID, mockUpdate)
	mock.ExpectationsWereMet()
	if !reflect.DeepEqual(result, mockUser) {
		t.Errorf("Update failed")
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	row := mock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}).
		AddRow(5, "example@example.com", []byte("123456"), "aUser", "firstName", "lastName", "https://www.gravatar.com/avatar/user@example.com")
	mockContext := &MySQLConnection{Client: db}
	mock.ExpectQuery(regexp.QuoteMeta("select * from Users where ID = ?")).WithArgs(strconv.FormatInt(5, 10)).WillReturnRows(row)
	mock.ExpectExec(regexp.QuoteMeta("delete from Users where ID = ?")).WithArgs(strconv.FormatInt(5, 10))
	mock.ExpectationsWereMet()
	err = mockContext.Delete(5)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = mockContext.Delete(3)
	if err == nil {
		t.Errorf("Expected error")
	}
}

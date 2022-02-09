package repository

import (
	"database/sql"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/kenanya/fin_coins/account"
	"github.com/stretchr/testify/assert"
)

// var u = &r.UserModel{
// 	ID:    uuid.New().String(),
// 	Name:  "Momo",
// 	Email: "momo@mail.com",
// 	Phone: "08123456789",
// }

var test1 = account.Account{
	ID:        "test-ken999",
	Balance:   50000,
	Currency:  "USD",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var logger log.Logger

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	// var logger log.Loggers
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	db, mock, err := sqlmock.New()
	if err != nil {
		level.Error(logger).Log("error connecting db", err)
		os.Exit(-1)
	}

	return db, mock
}

func TestCreateAccount(t *testing.T) {
	db, mock := NewMock()
	timeNow := time.Now()
	accountRepo := NewAccountRepository(db, logger)

	// query := "INSERT INTO account \\(id, balance, currency, created_at, updated_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)"
	query := `INSERT INTO account (id, balance, currency, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(test1.ID, test1.Balance, test1.Currency, timeNow, timeNow).WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := accountRepo.CreateAccount(test1)
	assert.NoError(t, err)
}

func TestGetAccountByID(t *testing.T) {
	db, mock := NewMock()
	accountRepo := NewAccountRepository(db, logger)

	query := "SELECT id, balance, currency, created_at, updated_at FROM account WHERE id = $1"

	rows := sqlmock.NewRows([]string{"id", "balance", "currency", "created_at", "updated_at"}).
		AddRow(test1.ID, test1.Balance, test1.Currency, test1.CreatedAt, test1.UpdatedAt)

	// mock.ExpectQuery(query).WithArgs(test1.ID).WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(test1.ID).WillReturnRows(rows)

	user, err := accountRepo.GetAccountByID(test1.ID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

// func TestFindByIDError(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &repository{db}
// 	defer func() {
// 		repo.Close()
// 	}()

// 	query := "SELECT id, name, email, phone FROM user WHERE id = \\?"

// 	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"})

// 	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

// 	user, err := repo.FindByID(u.ID)
// 	assert.Empty(t, user)
// 	assert.Error(t, err)
// }

// func TestFind(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &repository{db}
// 	defer func() {
// 		repo.Close()
// 	}()

// 	query := "SELECT id, name, email, phone FROM users"

// 	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
// 		AddRow(u.ID, u.Name, u.Email, u.Phone)

// 	mock.ExpectQuery(query).WillReturnRows(rows)

// 	users, err := repo.Find()
// 	assert.NotEmpty(t, users)
// 	assert.NoError(t, err)
// 	assert.Len(t, users, 1)
// }

// func TestCreateError(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &repository{db}
// 	defer func() {
// 		repo.Close()
// 	}()

// 	query := "INSERT INTO user \\(id, name, email, phone\\) VALUES \\(\\?, \\?, \\?, \\?\\)"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(u.ID, u.Name, u.Email, u.Phone).WillReturnResult(sqlmock.NewResult(0, 0))

// 	err := repo.Create(u)
// 	assert.Error(t, err)
// }

// func TestUpdate(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &repository{db}
// 	defer func() {
// 		repo.Close()
// 	}()

// 	query := "UPDATE users SET name = \\?, email = \\?, phone = \\? WHERE id = \\?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(u.Name, u.Email, u.Phone, u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

// 	err := repo.Update(u)
// 	assert.NoError(t, err)
// }

// func TestUpdateErr(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &repository{db}
// 	defer func() {
// 		repo.Close()
// 	}()

// 	query := "UPDATE user SET name = \\?, email = \\?, phone = \\? WHERE id = \\?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(u.Name, u.Email, u.Phone, u.ID).WillReturnResult(sqlmock.NewResult(0, 0))

// 	err := repo.Update(u)
// 	assert.Error(t, err)
// }

// func TestDelete(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &repository{db}
// 	defer func() {
// 		repo.Close()
// 	}()

// 	query := "DELETE FROM users WHERE id = \\?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

// 	err := repo.Delete(u.ID)
// 	assert.NoError(t, err)
// }

// func TestDeleteError(t *testing.T) {
// 	db, mock := NewMock()
// 	repo := &repository{db}
// 	defer func() {
// 		repo.Close()
// 	}()

// 	query := "DELETE FROM user WHERE id = \\?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 0))

// 	err := repo.Delete(u.ID)
// 	assert.Error(t, err)
// }

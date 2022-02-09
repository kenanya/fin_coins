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
	"github.com/kenanya/fin_coins/payment"
	"github.com/stretchr/testify/assert"
)

var test1 = account.Account{
	ID:        "ken99999",
	Balance:   50000,
	Currency:  "USD",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var test2 = payment.Payment{
	ID:            "x1234",
	AccountID:     "x1234",
	TransactionID: "x1234",
	Amount:        100,
	ToAccount:     "x1234",
	FromAccount:   "x1234",
	Direction:     "x1234",
	CreatedAt:     time.Now(),
}

var logger log.Logger

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
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

	query := `INSERT INTO account (id, balance, currency, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)`
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(test1.ID, test1.Balance, test1.Currency, timeNow, timeNow).WillReturnResult(sqlmock.NewResult(0, 1))

	user, err := accountRepo.CreateAccount(test1)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestGetAccountByID(t *testing.T) {
	db, mock := NewMock()
	accountRepo := NewAccountRepository(db, logger)

	query := "SELECT id, balance, currency, created_at, updated_at FROM account WHERE id = $1"

	rows := sqlmock.NewRows([]string{"id", "balance", "currency", "created_at", "updated_at"}).
		AddRow(test1.ID, test1.Balance, test1.Currency, test1.CreatedAt, test1.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(test1.ID).WillReturnRows(rows)

	user, err := accountRepo.GetAccountByID(test1.ID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestGetAllAccount(t *testing.T) {
	db, mock := NewMock()
	accountRepo := NewAccountRepository(db, logger)

	query := "SELECT id, balance, currency, created_at, updated_at FROM account"

	rows := sqlmock.NewRows([]string{"id", "balance", "currency", "created_at", "updated_at"}).
		AddRow(test1.ID, test1.Balance, test1.Currency, test1.CreatedAt, test1.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	users, err := accountRepo.GetAllAccount()
	assert.NotEmpty(t, users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestGetAllPayment(t *testing.T) {
	db, mock := NewMock()
	paymentRepo := NewPaymentRepository(db, logger)

	query := "SELECT id, account_id, transaction_id, amount, to_account, from_account, direction, created_at FROM payment"

	rows := sqlmock.NewRows([]string{"id", "account_id", "transaction_id", "amount", "to_account", "from_account", "direction", "created_at"}).
		AddRow(test2.ID, test2.AccountID, test2.TransactionID, test2.Amount, test2.ToAccount, test2.FromAccount, test2.Direction, test2.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	users, err := paymentRepo.GetAllPayment()
	assert.NotEmpty(t, users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

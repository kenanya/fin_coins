package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kenanya/fin_coins/account"
	"github.com/kenanya/fin_coins/lib"
	"github.com/kenanya/fin_coins/payment"
	"github.com/pborman/uuid"

	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	ErrRepository = errors.New("unable to handle request")
)

type allRepository struct {
	db     *sql.DB
	logger log.Logger
}

func (repo *allRepository) CreateAccount(accountData account.Account) (account.Account, error) {
	var accountRow = account.Account{}
	timeNow := time.Now()
	sql := `INSERT INTO account (id, balance, currency, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := repo.db.Exec(sql, accountData.ID, accountData.Balance, accountData.Currency, timeNow, timeNow)
	if err != nil {
		return accountRow, err
	} else {
		accountRow = account.Account{
			ID:        accountData.ID,
			Balance:   accountData.Balance,
			Currency:  accountData.Currency,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
	}
	return accountRow, nil
}

func (repo *allRepository) GetAccountByID(id string) (account.Account, error) {
	var accountRow = account.Account{}
	if err := repo.db.QueryRow(
		"SELECT id, balance, currency, created_at, updated_at FROM account WHERE id = $1", id).
		Scan(&accountRow.ID, &accountRow.Balance, &accountRow.Currency, &accountRow.CreatedAt, &accountRow.UpdatedAt); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return accountRow, err
	}

	return accountRow, nil
}

func (repo *allRepository) GetAllAccount() ([]account.Account, error) {

	accounts := []account.Account{}
	rows, err := repo.db.Query(
		`SELECT id, balance, currency, created_at, updated_at FROM account`)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var each account.Account
		if err := rows.Scan(&each.ID, &each.Balance, &each.Currency, &each.CreatedAt, &each.UpdatedAt); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return accounts, err
		}
		accounts = append(accounts, each)
	}
	return accounts, nil
}

func (repo *allRepository) GetAllPayment() ([]payment.Payment, error) {

	payments := []payment.Payment{}
	rows, err := repo.db.Query(
		`SELECT id, account_id, transaction_id, amount, to_account, from_account, direction, created_at FROM payment`)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return payments, err
	}
	defer rows.Close()

	for rows.Next() {
		var each payment.Payment
		if err := rows.Scan(&each.ID, &each.AccountID, &each.TransactionID, &each.Amount, &each.ToAccount, &each.FromAccount, &each.Direction, &each.CreatedAt); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return payments, err
		}
		payments = append(payments, each)
	}
	return payments, nil
}

func (repo *allRepository) GetPaymentByDirection(direction string) ([]payment.Payment, error) {

	payments := []payment.Payment{}
	rows, err := repo.db.Query(
		`SELECT id, account_id, amount, to_account, from_account, direction, created_at 
		FROM payment 
		WHERE direction=$1`, direction)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return payments, err
	}
	defer rows.Close()

	for rows.Next() {
		var each payment.Payment
		if err := rows.Scan(&each.ID, &each.AccountID, &each.Amount, &each.ToAccount, &each.FromAccount, &each.Direction, &each.CreatedAt); err != nil {
			level.Error(repo.logger).Log("err", err.Error())
			return payments, err
		}
		payments = append(payments, each)
	}
	return payments, nil
}

func (repo *allRepository) SendPayment(accountID string, amount float32, toAccount string) error {

	ctx := context.Background()
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	repo.logger.Log("#in begin amount : ", amount)
	//outgoing
	var transactionID = uuid.New()
	sql := `UPDATE account SET balance=balance-$1, updated_at=$2 WHERE id=$3`
	_, err = tx.ExecContext(ctx, sql, amount, time.Now(), accountID)
	if err != nil {
		tx.Rollback()
		level.Error(repo.logger).Log("err", err.Error())
		return err
	}

	var paymentPass = payment.Payment{
		ID:          uuid.New(),
		AccountID:   accountID,
		Amount:      amount,
		ToAccount:   toAccount,
		Direction:   lib.CONS_DIRECTION_OUTGOING,
		FromAccount: "",
	}

	sql = `INSERT INTO payment (id, account_id, transaction_id, amount, to_account, from_account, direction, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err = tx.ExecContext(ctx, sql, paymentPass.ID, paymentPass.AccountID, transactionID, paymentPass.Amount, paymentPass.ToAccount, paymentPass.FromAccount, paymentPass.Direction, time.Now())
	if err != nil {
		tx.Rollback()
		level.Error(repo.logger).Log("err", err.Error())
		return err
	}

	//incoming
	sql = `UPDATE account SET balance=balance+$1, updated_at=$2 WHERE id=$3`
	_, err = tx.ExecContext(ctx, sql, amount, time.Now(), toAccount)
	if err != nil {
		tx.Rollback()
		level.Error(repo.logger).Log("err", err.Error())
		return err
	}

	paymentPass = payment.Payment{
		ID:          uuid.New(),
		AccountID:   toAccount,
		Amount:      amount,
		FromAccount: accountID,
		Direction:   lib.CONS_DIRECTION_INCOMING,
		ToAccount:   "",
	}

	sql = `INSERT INTO payment (id, account_id, transaction_id, amount, from_account, to_account, direction, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err = tx.ExecContext(ctx, sql, paymentPass.ID, paymentPass.AccountID, transactionID, paymentPass.Amount, paymentPass.FromAccount, paymentPass.ToAccount, paymentPass.Direction, time.Now())
	if err != nil {
		tx.Rollback()
		level.Error(repo.logger).Log("err", err.Error())
		return err
	}

	time.Sleep(8 * time.Second)
	// Commit the change if all queries ran successfully
	err = tx.Commit()
	repo.logger.Log("#en tx amount : ", amount)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return err
	}

	return nil
}

func NewAccountRepository(db *sql.DB, logger log.Logger) account.Repository {
	return &allRepository{
		db:     db,
		logger: log.With(logger, "repo", "postgres"),
	}
}

func NewPaymentRepository(db *sql.DB, logger log.Logger) payment.Repository {
	return &allRepository{
		db:     db,
		logger: log.With(logger, "repo", "postgres"),
	}
}

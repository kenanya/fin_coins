package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kenanya/fin_coins/account"
	"github.com/kenanya/fin_coins/payment"

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
	sql := `
			INSERT INTO account (id, balance, currency, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5)`
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

func (repo *allRepository) GetAccountByID(ctx context.Context, id string) (account.Account, error) {
	var accountRow = account.Account{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT id, balance, currency, created_at, updated_at FROM account WHERE id = $1",
		id).
		Scan(
			&accountRow.ID, &accountRow.Balance, &accountRow.Currency, &accountRow.CreatedAt, &accountRow.UpdatedAt,
		); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return accountRow, err
	}

	return accountRow, nil
}

func (repo *allRepository) UpdateBalance(ctx context.Context, amount float32, accountID string) (account.Account, error) {
	var accountRow = account.Account{}
	sql := `UPDATE account SET balance=balance+$1 WHERE id=$2`

	_, err := repo.db.ExecContext(ctx, sql, amount, accountID)
	if err != nil {
		return accountRow, err
	}

	accountRow, err = repo.GetAccountByID(ctx, accountID)
	return accountRow, err
}

func (repo *allRepository) GetAllAcount(ctx context.Context) ([]account.Account, error) {

	accounts := []account.Account{}
	rows, err := repo.db.QueryContext(ctx,
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

type paymentRepository struct {
	db     *sql.DB
	logger log.Logger
}

func (repo *paymentRepository) GetAllPayment(ctx context.Context) ([]payment.Payment, error) {

	payments := []payment.Payment{}
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, account_id, amount, to_account, from_account, direction, created_at FROM payment`)
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

func (repo *paymentRepository) GetPaymentByDirection(ctx context.Context, direction string) ([]payment.Payment, error) {

	payments := []payment.Payment{}
	rows, err := repo.db.QueryContext(ctx,
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

func (repo *paymentRepository) SendPayment(ctx context.Context, payment payment.Payment) error {
	sql := `
			INSERT INTO payment (id, account_id, amount, to_account, from_account, direction, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := repo.db.ExecContext(ctx, sql, payment.ID, payment.AccountID, payment.Amount, payment.ToAccount, payment.FromAccount, payment.Direction, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func NewAccountRepository(db *sql.DB, logger log.Logger) account.Repository {
	return &allRepository{
		db:     db,
		logger: log.With(logger, "repo", "mongodb"),
	}
}

// //Creates and returns an instance
// func NewRepo(db *sql.DB, logger log.Logger) (Repository, error) {
// 	return &repo{
// 	   db:     db,
// 	   logger: log.With(logger, "repo", "mongodb"),
// 	}, nil
//   }

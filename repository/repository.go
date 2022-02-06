package account

import (
	"context"
	"database/sql"
	"errors"
	"fin_coins/account"
	"fin_coins/payment"

	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	ErrRepository = errors.New("unable to handle request")
)

type accountRepository struct {
	db     *sql.DB
	logger log.Logger
}

func (repo *accountRepository) CreateAccount(ctx context.Context, account account.Account) error {
	sql := `
			INSERT INTO account (id, balance, currency, created_at, updated_at)
			VALUES ($1,$2,$3,$4)`
	_, err := repo.db.ExecContext(ctx, sql, account.ID, account.Balance, account.Currency, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (repo *accountRepository) GetAccountByID(ctx context.Context, id string) (account.Account, error) {
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

func (repo *accountRepository) UpdateBalance(ctx context.Context, amount float32, accountID string) error {
	sql := `UPDATE account SET balance=balance+$1 WHERE id=$2`

	_, err := repo.db.ExecContext(ctx, sql, amount, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *accountRepository) GetAllAcount(ctx context.Context) ([]account.Account, error) {

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
		`SELECT id, account_id, amount, to_account, from_account_id, direction, created_at FROM payment`)
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
		`SELECT id, account_id, amount, to_account, from_account_id, direction, created_at 
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
			INSERT INTO payment (id, account_id, amount, to_account, from_account_id, direction, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := repo.db.ExecContext(ctx, sql, payment.ID, payment.AccountID, payment.Amount, payment.ToAccount, payment.FromAccount, payment.Direction, time.Now())
	if err != nil {
		return err
	}
	return nil
}

package repository

// import (
// 	"database/sql"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/go-kit/kit/log"
// 	"github.com/go-kit/log/level"
// 	"github.com/kenanya/fin_coins/account"
// 	"github.com/kenanya/fin_coins/payment"

// 	_ "github.com/lib/pq"
// )

// var (
// 	accountRepo account.Repository
// 	paymentRepo payment.Repository
// )

// const (
// 	defaultPort = "9696"
// 	dbHost      = "localhost"
// 	dbPort      = 5432
// 	dbUser      = "postgres"
// 	dbPassword  = "postgres"
// 	dbName      = "fin_coins"
// 	schemaName  = "wallet"
// )

// func init() {
// 	var logger log.Logger
// 	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
// 	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
// 		dbHost, dbPort, dbUser, dbPassword, dbName, schemaName)
// 	var db *sql.DB
// 	{
// 		var err error
// 		db, err = sql.Open("postgres", psqlInfo)
// 		if err != nil {
// 			level.Error(logger).Log("exit", err)
// 			os.Exit(-1)
// 		}
// 		defer db.Close()
// 	}

// 	accountRepo = NewAccountRepository(db, logger)
// 	paymentRepo = NewPaymentRepository(db, logger)

// }

// func TestCreateAccount(t *testing.T) {
// 	test1 := account.Account{
// 		ID:       "test-ken999",
// 		Balance:  50000,
// 		Currency: "USD",
// 	}
// 	res, err := accountRepo.CreateAccount(test1)
// 	if err != nil {
// 		t.Error("TestCreateAccount failed creating account")
// 	}
// 	fmt.Printf("res.ID : %v", res.ID)
// 	if test1.ID != res.ID {
// 		t.Errorf("TestCreateAccount returned an unexpected result: got %v want %v", res.ID, test1.ID)
// 	}

// }

// func TestDeleteAccount(t *testing.T) {
// 	accountID := "test-ken999"
// 	err := accountRepo.DeleteAccount(accountID)
// 	if err != nil {
// 		t.Error("TestDeleteAccount failed deleting account")
// 	}
// }

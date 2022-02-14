package main

import (
	"database/sql"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/kenanya/fin_coins/account"
	"github.com/kenanya/fin_coins/payment"
	"github.com/kenanya/fin_coins/repository"
	"github.com/pborman/uuid"

	_ "github.com/lib/pq"
)

var logger kitlog.Logger
var testSrv *httptest.Server
var testDb *sql.DB
var globalId string

func initIntegrationTest(entity string) *httptest.Server {

	var (
		dbHost       = envString("DB_HOST", defaultDbHost)
		dbPort       = envString("DB_PORT", defaultDbPort)
		dbUser       = envString("DB_USER", defaultDbUser)
		dbPassword   = envString("DB_PASSWORD", defaultDbPassword)
		dbName       = envString("DB_NAME", defaultDbName)
		dbSchemaName = envString("DB_SCHEMA_NAME", defaultDbSchemaName)
	)

	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "listen", dbPort, "caller", kitlog.DefaultCaller)
	flag.Parse()

	var err error
	testDb, err = OpenDB(logger, dbHost, dbPort, dbUser, dbPassword, dbName, dbSchemaName)
	if err != nil {
		level.Error(logger).Log("failed initialize postgres connection", err)
		os.Exit(-1)
	}

	err = testDb.Ping()
	if err != nil {
		level.Error(logger).Log("failed connect to postgres", err)
		os.Exit(-1)
	}

	var r http.Handler = nil
	if entity == "account" {
		accountRepo := repository.NewAccountRepository(testDb, logger)
		s := account.NewService(accountRepo)
		r = account.MakeHandler(s, logger)
	} else if entity == "payment" {
		accountRepo := repository.NewAccountRepository(testDb, logger)
		paymentRepo := repository.NewPaymentRepository(testDb, logger)
		s := payment.NewService(paymentRepo, accountRepo)
		r = payment.MakeHandler(s, logger)
	}

	return httptest.NewServer(r)
}

func TestIntegrationCreateAccount(t *testing.T) {

	testSrv = initIntegrationTest("account")
	defer testDb.Close()
	if testSrv == nil {
		level.Error(logger).Log("Entity not found")
		os.Exit(-1)
	}

	globalId = "cindy" + uuid.New()

	for _, testcase := range []struct {
		method, url, body string
		want              int
	}{
		{"POST", testSrv.URL + "/account/v1/account", `{"id":"` + globalId + `", "currency":"USD"}`, http.StatusOK},
		{"POST", testSrv.URL + "/account/v1/account", `{"currency":"USD"}`, http.StatusBadRequest},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)
		if testcase.want != resp.StatusCode {
			t.Errorf("%s %s %s: want %d have %d", testcase.method, testcase.url, testcase.body, testcase.want, resp.StatusCode)
		}
	}
}

func TestIntegrationGetAccountByID(t *testing.T) {
	testSrv = initIntegrationTest("account")
	defer testDb.Close()
	if testSrv == nil {
		level.Error(logger).Log("Entity not found")
		os.Exit(-1)
	}

	for _, testcase := range []struct {
		method, url, body string
		want              int
	}{
		{"GET", testSrv.URL + "/account/v1/account/" + globalId, ``, http.StatusOK},
		{"GET", testSrv.URL + "/account/v1/account/" + uuid.New(), ``, http.StatusNotFound},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)
		if testcase.want != resp.StatusCode {
			t.Errorf("%s %s %s: want %d have %d", testcase.method, testcase.url, testcase.body, testcase.want, resp.StatusCode)
		}
	}
}

func TestIntegrationGetAllAccount(t *testing.T) {
	testSrv = initIntegrationTest("account")
	defer testDb.Close()
	if testSrv == nil {
		level.Error(logger).Log("Entity not found")
		os.Exit(-1)
	}

	for _, testcase := range []struct {
		method, url, body string
		want              int
	}{
		{"GET", testSrv.URL + "/account/v1/account", ``, http.StatusOK},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)
		if testcase.want != resp.StatusCode {
			t.Errorf("%s %s %s: want %d have %d", testcase.method, testcase.url, testcase.body, testcase.want, resp.StatusCode)
		}
	}
}

func TestIntegrationSendPayment(t *testing.T) {

	testSrv = initIntegrationTest("payment")
	defer testDb.Close()
	if testSrv == nil {
		level.Error(logger).Log("Entity not found")
		os.Exit(-1)
	}

	globalId = "cindy" + uuid.New()

	for _, testcase := range []struct {
		method, url, body string
		want              int
	}{
		{"POST", testSrv.URL + "/payment/v1/payment", `{"account_id":"mike2167","amount":50,"to_account":"alice456"}`, http.StatusInternalServerError},
		{"POST", testSrv.URL + "/payment/v1/payment", `{"account_id":"bob123","amount":10,"to_account":"alice456"}`, http.StatusOK},
		{"POST", testSrv.URL + "/payment/v1/payment", `{"account_id":"bob123","amount":10000000,"to_account":"alice456"}`, http.StatusInternalServerError},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)
		if testcase.want != resp.StatusCode {
			t.Errorf("%s %s %s: want %d have %d", testcase.method, testcase.url, testcase.body, testcase.want, resp.StatusCode)
		}
	}
}

func TestIntegrationGetAllPayment(t *testing.T) {
	testSrv = initIntegrationTest("payment")
	defer testDb.Close()
	if testSrv == nil {
		level.Error(logger).Log("Entity not found")
		os.Exit(-1)
	}

	for _, testcase := range []struct {
		method, url, body string
		want              int
	}{
		{"GET", testSrv.URL + "/payment/v1/payment", ``, http.StatusOK},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)
		if testcase.want != resp.StatusCode {
			t.Errorf("%s %s %s: want %d have %d", testcase.method, testcase.url, testcase.body, testcase.want, resp.StatusCode)
		}
	}
}

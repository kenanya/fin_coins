package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kenanya/fin_coins/account"
	"github.com/kenanya/fin_coins/payment"
	"github.com/kenanya/fin_coins/repository"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	_ "github.com/lib/pq"
)

const (
	defaultPort         = "9696"
	defaultDbHost       = "localhost"
	defaultDbPort       = "5437"
	defaultDbUser       = "postgres"
	defaultDbPassword   = "postgres"
	defaultDbName       = "postgres"
	defaultDbSchemaName = "wallet"
)

func OpenDB(logger log.Logger, dbHost, dbPort, dbUser, dbPassword, dbName, dbSchemaName string) (*sql.DB, error) {

	intPort, err := strconv.Atoi(dbPort)
	if err != nil {
		return nil, err
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		dbHost, intPort, dbUser, dbPassword, dbName, dbSchemaName)

	// fmt.Println(psqlInfo)
	var db *sql.DB
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, err
}

func main() {
	var (
		addr     = envString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")

		dbHost       = envString("DB_HOST", defaultDbHost)
		dbPort       = envString("DB_PORT", defaultDbPort)
		dbUser       = envString("DB_USER", defaultDbUser)
		dbPassword   = envString("DB_PASSWORD", defaultDbPassword)
		dbName       = envString("DB_NAME", defaultDbName)
		dbSchemaName = envString("DB_SCHEMA_NAME", defaultDbSchemaName)
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	db, err := OpenDB(logger, dbHost, dbPort, dbUser, dbPassword, dbName, dbSchemaName)
	if err != nil {
		level.Error(logger).Log("failed initialize postgres connection", err)
		os.Exit(-1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		level.Error(logger).Log("failed connect to postgres", err)
		os.Exit(-1)
	}
	logger.Log("Successfully connected to DB!")

	var (
		accountRepo = repository.NewAccountRepository(db, logger)
		paymentRepo = repository.NewPaymentRepository(db, logger)
	)

	fieldKeys := []string{"method"}

	var ac account.Service
	ac = account.NewService(accountRepo)
	ac = account.NewLoggingService(log.With(logger, "component", "account"), ac)
	ac = account.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "account_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "account_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		ac,
	)

	var ps payment.Service
	ps = payment.NewService(paymentRepo, accountRepo)
	ps = payment.NewLoggingService(log.With(logger, "component", "payment"), ps)
	ps = payment.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "payment_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "payment_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		ps,
	)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/account/v1/", account.MakeHandler(ac, httpLogger))
	mux.Handle("/payment/v1/", payment.MakeHandler(ps, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

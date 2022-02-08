package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kenanya/fin_coins/account"
	"github.com/kenanya/fin_coins/repository"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	_ "github.com/lib/pq"
)

const (
	defaultPort = "9595"
	// dbsource    = "postgresql://postgres@localhost:5432/ordersdb?sslmode=disable"
	// defaultRoutingServiceURL = "http://localhost:9797"

	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "fin_coins"
	schemaName = "wallet"
)

func main() {
	var (
		addr = envString("PORT", defaultPort)
		// rsurl = envString("ROUTINGSERVICE_URL", defaultRoutingServiceURL)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		// routingServiceURL = flag.String("service.routing", rsurl, "routing service URL")

		// ctx = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, schemaName)
	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		defer db.Close()
	}

	var (
		accountRepo = repository.NewAccountRepository(db, logger)
		// locations      = inmem.NewLocationRepository()
		// voyages        = inmem.NewVoyageRepository()
		// handlingEvents = inmem.NewHandlingEventRepository()
	)

	// // Configure some questionable dependencies.
	// var (
	// 	handlingEventFactory = cargo.HandlingEventFactory{
	// 		CargoRepository:    cargos,
	// 		VoyageRepository:   voyages,
	// 		LocationRepository: locations,
	// 	}
	// 	handlingEventHandler = handling.NewEventHandler(
	// 		inspection.NewService(cargos, handlingEvents, nil),
	// 	)
	// )

	// // Facilitate testing by adding some cargos.
	// storeTestData(cargos)

	fieldKeys := []string{"method"}

	// var rs routing.Service
	// rs = routing.NewProxyingMiddleware(ctx, *routingServiceURL)(rs)

	var ac account.Service
	ac = account.NewService(accountRepo)
	ac = account.NewLoggingService(log.With(logger, "component", "booking"), ac)
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

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/account/v1/", account.MakeHandler(ac, httpLogger))
	// mux.Handle("/tracking/v1/", tracking.MakeHandler(ts, httpLogger))
	// mux.Handle("/handling/v1/", handling.MakeHandler(hs, httpLogger))

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

// func storeTestData(r cargo.Repository) {
// 	test1 := cargo.New("FTL456", cargo.RouteSpecification{
// 		Origin:          location.AUMEL,
// 		Destination:     location.SESTO,
// 		ArrivalDeadline: time.Now().AddDate(0, 0, 7),
// 	})
// 	if err := r.Store(test1); err != nil {
// 		panic(err)
// 	}

// 	test2 := cargo.New("ABC123", cargo.RouteSpecification{
// 		Origin:          location.SESTO,
// 		Destination:     location.CNHKG,
// 		ArrivalDeadline: time.Now().AddDate(0, 0, 14),
// 	})
// 	if err := r.Store(test2); err != nil {
// 		panic(err)
// 	}
// }

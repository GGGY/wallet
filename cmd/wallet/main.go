package main

import (
	"flag"
	"fmt"
	"github.com/GGGY/wallet/internal/wallet/repository"
	"github.com/GGGY/wallet/internal/wallet/service"
	"github.com/GGGY/wallet/internal/wallet/transport"
	httptransport "github.com/GGGY/wallet/internal/wallet/transport/http"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")

		dbName     = envString("WALLET_DB_NAME", "app")
		dbHost     = envString("WALLET_DB_HOST", "localhost")
		dbPort     = envString("WALLET_DB_PORT", "5432")
		dbUser     = envString("WALLET_DB_USER", "app")
		dbPassword = envString("WALLET_DB_PASSWORD", "password")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var db *sqlx.DB
	{
		var dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		)

		db = sqlx.MustOpen("postgres", dsn)
	}

	var s service.Service
	{
		accountRepo := repository.NewAccountPg(db)
		paymentRepo := repository.NewPaymentPg(db)
		s = service.NewService(db, accountRepo, paymentRepo, logger)
	}

	var h http.Handler
	{
		endpoints := transport.MakeEndpoints(s)
		h = httptransport.NewService(endpoints, logger)
	}

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, h)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

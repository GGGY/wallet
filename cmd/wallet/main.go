package main

import (
	"flag"
	"fmt"
	"github.com/GGGY/wallet/internal/wallet/service"
	"github.com/GGGY/wallet/internal/wallet/transport"
	httptransport "github.com/GGGY/wallet/internal/wallet/transport/http"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var s service.Service
	{
		s = service.NewService(nil, nil, logger)
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

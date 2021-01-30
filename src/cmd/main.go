package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/nburunova/currency-loader/src/api"
	"github.com/nburunova/currency-loader/src/currency"
	"github.com/nburunova/currency-loader/src/services/cbr"
	"github.com/nburunova/currency-loader/src/services/log"
	"github.com/nburunova/currency-loader/src/services/signals"
	"github.com/nburunova/currency-loader/src/services/sql_storage"
	"github.com/nburunova/currency-loader/src/services/webserver"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	lifeTimeCheckRatio = 3
)

func main() {
	cfg := parseFlags()
	logger := log.NewLogger("json", "debug")
	transport := &http.Transport{
		IdleConnTimeout:     time.Duration(cfg.http.idleConnTimeout) * time.Second,
		MaxIdleConns:        cfg.http.maxIdleConnectionsPerHost,
		MaxIdleConnsPerHost: cfg.http.maxIdleConnectionsPerHost,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	cbrClient, err := cbr.NewClient(httpClient)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "cannot create Cbr client"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lifeTime := time.Duration(cfg.lifeTime) * time.Second
	storage, err := sql_storage.New(lifeTime, logger)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "cannot create valutes storage"))
	}
	storage.Start(ctx, lifeTime/lifeTimeCheckRatio)

	cbrLoader := currency.NewCbrLoader(cbrClient)
	service := currency.NewService(cbrLoader, storage, logger)
	err = service.Start(ctx, time.Duration(cfg.updatePeriod)*time.Second)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "cannot start valutes service"))
	}

	api, err := api.NewCurrencyAPI(storage, cfg.pageSize, logger)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "cannot start valutes service"))
	}

	commonRouter := newCommonRouter()

	router := webserver.NewRouter(logger)
	router.Mount("/", api.Routers())
	commonRouter.Mount("/", router)

	address := fmt.Sprintf("%v:%v", cfg.http.address, cfg.http.port)
	srv := webserver.NewServer(address, commonRouter)
	signals.BindSignals(logger, srv)
	logger.Info("starting http service...")
	logger.Infof("listening on %s", address)
	if err := srv.Start(); err != nil {
		logger.WithError(err).Fatal()
	}
}

func newCommonRouter() *webserver.Router {
	commonRouter := webserver.NewCommonRouter()
	commonRouter.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "OK")
	})
	return commonRouter
}

type logFlags struct {
	level  string
	format string
}

type httpFlags struct {
	address                   string
	port                      string
	maxIdleConnectionsPerHost int
	idleConnTimeout           int
}

// cliFlags is a union of the fields, which application could parse from CLI args
type cliFlags struct {
	log          logFlags
	http         httpFlags
	updatePeriod int
	lifeTime     int
	pageSize     int
}

// parseFlags maps CLI flags to struct
func parseFlags() *cliFlags {
	cfg := cliFlags{}

	kingpin.Flag("log-level", "Log level.").
		Default("info").
		Envar("LOG_LEVEL").
		EnumVar(&cfg.log.level, "debug", "info", "warning", "error", "fatal", "panic")
	kingpin.Flag("log-format", "Log format.").
		Default("text").
		Envar("LOG_FORMAT").
		EnumVar(&cfg.log.format, "text", "json")

	kingpin.Flag("address", "HTTP service address, default 0.0.0.0.").
		Default("0.0.0.0").
		Envar("ADDRESS").
		StringVar(&cfg.http.address)

	kingpin.Flag("port", "HTTP service port, default 5000.").
		Default("5000").
		Envar("PORT").
		StringVar(&cfg.http.port)

	kingpin.Flag("maxIdleConnectionsPerHost", "maxIdleConnectionsPerHost").
		Default("10").
		Envar("MAXIDLECONNECTIONSPERHOST").
		IntVar(&cfg.http.maxIdleConnectionsPerHost)

	kingpin.Flag("idleConnTimeout", "idleConnTimeout").
		Default("5").
		Envar("IDLECONNECTIONTIMEOUT").
		IntVar(&cfg.http.idleConnTimeout)

	kingpin.Flag("updatePeriod", "currency update period in seconds").
		Default("180").
		Envar("UPDATEPERIOD").
		IntVar(&cfg.updatePeriod)

	kingpin.Flag("lifeTime", "currency save time in seconds").
		Default("1800").
		Envar("LIFETIME").
		IntVar(&cfg.lifeTime)

	kingpin.Flag("pageSize", "api page size").
		Default("10").
		Envar("PAGESIZE").
		IntVar(&cfg.pageSize)

	kingpin.Parse()
	return &cfg
}

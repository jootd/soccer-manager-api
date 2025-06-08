package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/jootd/soccer-manager/app/services/sales-api/handlers"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"github.com/jootd/soccer-manager/business/sdk/v1/debug"
	"github.com/jootd/soccer-manager/foundation/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

var build = "develop"

func main() {
	log, err := logger.New("SALES-API")
	if err != nil {
		fmt.Println("logger init error:", err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	opt := maxprocs.Logger(log.Infof)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	cfg, err := readEnvConfig()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	log.Infow("starting service", "version", build)
	expvar.NewString("build").Set(build)

	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()

	log.Infow("startup", "status", "initializing OT/Zipkin tracing support")

	traceProvider, err := startTracing(
		cfg.Zipkin.ServiceName,
		cfg.Zipkin.ReporterURI,
		cfg.Zipkin.Probability,
	)
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}
	defer traceProvider.Shutdown(context.Background())

	tracer := traceProvider.Tracer("service")

	// =========================================================================
	// Start Debug Service

	log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux(build, log, db)); err != nil {
			log.Errorw("shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =========================================================================

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		DB:       db,
		Tracer:   tracer,
	})

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

type config struct {
	Web struct {
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
		APIHost         string
		DebugHost       string
	}
	DB struct {
		User         string
		Password     string
		Host         string
		Name         string
		MaxIdleConns int
		MaxOpenConns int
		DisableTLS   bool
	}
	Zipkin struct {
		ReporterURI string
		ServiceName string
		Probability float64
	}
}

func readEnvConfig() (*config, error) {
	var cfg config
	var err error

	getEnv := func(key string) (string, error) {
		v := os.Getenv(key)
		if v == "" {
			return "", fmt.Errorf("missing environment variable %q", key)
		}
		return v, nil
	}

	parseDuration := func(key string) (time.Duration, error) {
		v, err := getEnv(key)
		if err != nil {
			return 0, err
		}
		d, err := time.ParseDuration(v)
		if err != nil {
			return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
		}
		return d, nil
	}

	parseInt := func(key string) (int, error) {
		v, err := getEnv(key)
		if err != nil {
			return 0, err
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid integer for %s: %w", key, err)
		}
		return n, nil
	}

	parseBool := func(key string) (bool, error) {
		v, err := getEnv(key)
		if err != nil {
			return false, err
		}
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("invalid boolean for %s: %w", key, err)
		}
		return b, nil
	}

	parseFloat := func(key string) (float64, error) {
		v, err := getEnv(key)
		if err != nil {
			return 0, err
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid float for %s: %w", key, err)
		}
		return f, nil
	}

	// Web timeouts
	if cfg.Web.ReadTimeout, err = parseDuration("WEB_READ_TIMEOUT"); err != nil {
		return nil, err
	}
	if cfg.Web.WriteTimeout, err = parseDuration("WEB_WRITE_TIMEOUT"); err != nil {
		return nil, err
	}
	if cfg.Web.IdleTimeout, err = parseDuration("WEB_IDLE_TIMEOUT"); err != nil {
		return nil, err
	}
	if cfg.Web.ShutdownTimeout, err = parseDuration("WEB_SHUTDOWN_TIMEOUT"); err != nil {
		return nil, err
	}

	// Web hosts
	if cfg.Web.APIHost, err = getEnv("WEB_API_HOST"); err != nil {
		return nil, err
	}
	if cfg.Web.DebugHost, err = getEnv("WEB_DEBUG_HOST"); err != nil {
		return nil, err
	}

	// DB
	if cfg.DB.User, err = getEnv("DB_USER"); err != nil {
		return nil, err
	}
	if cfg.DB.Password, err = getEnv("DB_PASSWORD"); err != nil {
		return nil, err
	}
	if cfg.DB.Host, err = getEnv("DB_HOST"); err != nil {
		return nil, err
	}
	if cfg.DB.Name, err = getEnv("DB_NAME"); err != nil {
		return nil, err
	}
	if cfg.DB.MaxIdleConns, err = parseInt("DB_MAX_IDLE_CONNS"); err != nil {
		return nil, err
	}
	if cfg.DB.MaxOpenConns, err = parseInt("DB_MAX_OPEN_CONNS"); err != nil {
		return nil, err
	}
	if cfg.DB.DisableTLS, err = parseBool("DB_DISABLE_TLS"); err != nil {
		return nil, err
	}

	// Zipkin
	if cfg.Zipkin.ReporterURI, err = getEnv("ZIPKIN_REPORTER_URI"); err != nil {
		return nil, err
	}
	if cfg.Zipkin.ServiceName, err = getEnv("ZIPKIN_SERVICE_NAME"); err != nil {
		return nil, err
	}
	if cfg.Zipkin.Probability, err = parseFloat("ZIPKIN_PROBABILITY"); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func startTracing(serviceName, reporterURI string, probability float64) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(reporterURI)
	if err != nil {
		return nil, fmt.Errorf("creating zipkin exporter: %w", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("exporter", "zipkin"),
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.TraceIDRatioBased(probability)),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

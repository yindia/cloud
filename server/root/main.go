package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"github.com/rs/cors"
	"go.akshayshah.org/connectauth"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	cloudv1connect "task/pkg/gen/cloud/v1/cloudv1connect"
	"task/pkg/x"                                  // Import the x package for env and config
	repository "task/server/repository"           // Import repository package
	interfaces "task/server/repository/interface" // Import repository package
	"task/server/route"                           // Import route package

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	CompressMinByte = 1024 // Minimum byte size for compression
)

// AuthCtx holds user authentication information
type AuthCtx struct {
	Username string
}

// newCORS initializes CORS settings for the server
// It allows all origins and methods, and exposes necessary headers for gRPC-Web
func newCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			return true // Allow all origins
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
	})
}

// main is the entry point of the application.
// It sets up logging and runs the main application logic.
func main() {
	// Initialize structured logging with JSON format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

// run encapsulates the main application logic
// It loads configuration, sets up the server, and handles graceful shutdown
func run() error {
	// Load environment variables
	if err := x.LoadEnv(); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	env, err := x.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	slog.Info("Application started", "config", env)

	// Set up a channel to handle exit signals
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	// Create the repository with DB configuration
	repo, err := repository.GetRepository(env.Database.ToDbConnectionUri(), env.WorkerCount, env.Database.PoolMaxConns)
	if err != nil {
		return fmt.Errorf("failed to initialize database repository: %w", err)
	}

	slog.Info("Database repository initialized", "workerCount", env.WorkerCount)

	// Set up gRPC middleware
	middleware := connectauth.NewMiddleware(GrpcMiddleware)

	// Set up HTTP server
	mux := http.NewServeMux()
	if err := setupHandlers(mux, repo, middleware); err != nil {
		return fmt.Errorf("failed to set up handlers: %w", err)
	}

	// Add Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Initialize HTTP server
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%v", env.ServerPort),
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	// Start the server in a goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		slog.Info("HTTP server starting", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()

	// Wait for exit signal or server error
	select {
	case <-exitChan:
		slog.Info("Shutdown signal received, shutting down server...")
	case err := <-serverErrChan:
		return err
	}

	// Graceful shutdown
	if err := shutdownServer(srv); err != nil {
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}
	slog.Info("HTTP server shut down gracefully")
	return nil
}

// setupHandlers configures the HTTP handlers for the server
// It sets up the gRPC service, health check, and reflection handlers
func setupHandlers(mux *http.ServeMux, repo interfaces.TaskManagmentInterface, middleware *connectauth.Middleware) error {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return fmt.Errorf("failed to create interceptor: %w", err)
	}

	pattern, handler := cloudv1connect.NewTaskManagementServiceHandler(
		route.NewTaskServer(repo),
		connect.WithInterceptors(otelInterceptor),
		connect.WithCompressMinBytes(CompressMinByte),
	)
	mux.Handle(pattern, middleware.Wrap(handler))

	// Health check and reflection handlers
	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(cloudv1connect.TaskManagementServiceName),
	))
	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(cloudv1connect.TaskManagementServiceName),
	))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(cloudv1connect.TaskManagementServiceName),
	))

	slog.Info("Handlers set up successfully", "serviceName", cloudv1connect.TaskManagementServiceName)
	return nil
}

// shutdownServer gracefully shuts down the HTTP server
// It waits for ongoing requests to complete before shutting down
func shutdownServer(srv *http.Server) error {
	slog.Info("Initiating graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}
	slog.Info("Server shutdown completed")
	return nil
}

// GrpcMiddleware is the gRPC middleware used for authentication.
// Currently, it uses a placeholder authentication mechanism.
func GrpcMiddleware(ctx context.Context, req *connectauth.Request) (any, error) {
	// TODO: Implement proper authentication logic
	slog.Warn("Using placeholder authentication", "username", "tqindia")
	return AuthCtx{Username: "tqindia"}, nil
}

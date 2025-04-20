package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/bovinxx/auth-service/internal/closer"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/interceptor"
	"github.com/bovinxx/auth-service/internal/logger"
	"github.com/bovinxx/auth-service/internal/metrics"
	descAccess "github.com/bovinxx/auth-service/pkg/access_v1"
	descAuth "github.com/bovinxx/auth-service/pkg/auth_v1"
	descUser "github.com/bovinxx/auth-service/pkg/user_v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start(ctx context.Context) error {
	err := metrics.Init(ctx)
	if err != nil {
		return err
	}

	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatal("failed to run GRPC server", logger.Err(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatal("failed to run HTTP server", logger.Err(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheusHTTPServer()
		if err != nil {
			logger.Fatal("failed to run Prometheus HTTP server", logger.Err(err))
		}
	}()

	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load("./configs/local.env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.RateLimiterInterceptor,
			interceptor.MetricsInterceptor,
			interceptor.ValidateInterceptor,
			interceptor.AuthzInterceptor(a.serviceProvider.AccessService(ctx))),
	)

	reflection.Register(a.grpcServer)

	descUser.RegisterUserServiceServer(a.grpcServer, a.serviceProvider.UserImplementation(ctx))
	descAuth.RegisterAuthServiceServer(a.grpcServer, a.serviceProvider.AuthImplementation(ctx))
	descAccess.RegisterAccessServiceServer(a.grpcServer, a.serviceProvider.AccessImplementation(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := descUser.RegisterUserServiceHandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runPrometheusHTTPServer() error {
	log.Printf("Prometheus HTTP server is running on %s", a.serviceProvider.PrometheusConfig().Address())

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    a.serviceProvider.PrometheusConfig().Address(),
		Handler: mux,
	}

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

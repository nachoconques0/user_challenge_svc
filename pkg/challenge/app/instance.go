package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/db"
	userAggregate "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/aggregate/user"
	userGrpcController "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/grpc/user"
	userController "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/http/user"
	userService "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/service/user"
	userProto "github.com/nachoconques0/user_challenge_svc/pkg/challenge/proto/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/pubsub"
	grpcServer "github.com/nachoconques0/user_challenge_svc/pkg/challenge/server/grpc"
	httpServer "github.com/nachoconques0/user_challenge_svc/pkg/challenge/server/http"
)

const (
	serviceNEW = iota
	serviceRUNNING
	serviceSTOPPING
	serviceSTOPPED
)

type Instance struct {
	servers []server
	timeout int
	state   int
	mu      sync.Mutex
}

type server interface {
	Run() error
	Stop(context.Context) error
}

func New(opts ...Option) error {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	dbConn, err := setupDatabase(options)
	if err != nil {
		return err
	}

	userSvc, err := setupUserService(dbConn)
	if err != nil {
		return err
	}

	httpSrv, err := setupHTTPServer(userSvc, options.httpPort)
	if err != nil {
		return err
	}

	grpcSrv, err := setupGRPCServer(userSvc, options.gRPCPort)
	if err != nil {
		return err
	}

	i := Instance{
		timeout: 20,
		servers: []server{httpSrv, grpcSrv},
	}

	sigCh := waitForShutdownSignal()
	return i.Run(sigCh)
}

func setupDatabase(opts Options) (*gorm.DB, error) {
	return db.New(opts.dbOptions...)
}

func setupUserService(db *gorm.DB) (userService.Service, error) {
	publisher := pubsub.NewLoggerPublisher(db)
	userAgg, err := userAggregate.New(db, "nontest", publisher)
	if err != nil {
		return nil, err
	}
	return userService.New(userAgg), nil
}

func setupHTTPServer(svc userService.Service, port string) (server, error) {
	srv, err := httpServer.New(httpServer.WithAddress(fmt.Sprintf(":%s", port)))
	if err != nil {
		return nil, err
	}
	router := httpServer.InitHTTPRouter(srv)
	httpServer.InitUserRoutes(router, userController.NewController(svc))
	return srv, nil
}

func setupGRPCServer(svc userService.Service, port string) (server, error) {
	srv := grpcServer.New(port)
	userProto.RegisterUserServiceServer(srv.Server(), userGrpcController.NewController(svc))
	return srv, nil
}

func waitForShutdownSignal() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	return quit
}

func (s *Instance) Run(quit chan os.Signal) error {
	if s.isStopped() {
		return errors.New("instance was already stopped. Can't start it again")
	}

	if s.isNew() {
		log.Info().Msg("Application: starting...")

		s.mu.Lock()
		s.state = serviceRUNNING
		errCh := make(chan error, 1)
		s.runServers(errCh)
		s.mu.Unlock()

		log.Info().Msg("Application: running...")

		select {
		case err := <-errCh:
			return err
		case <-quit:
			s.mu.Lock()
			s.state = serviceSTOPPING
			s.mu.Unlock()

			log.Info().Msgf("Application: stopping in %.0fs...", s.Timeout().Seconds())
			ctx, cancel := context.WithTimeout(context.Background(), s.Timeout())
			defer cancel()
			go s.stopServers(ctx)
			<-ctx.Done()

			s.mu.Lock()
			s.state = serviceSTOPPED
			s.mu.Unlock()
			log.Info().Msg("Application: stopped")
		}
	}
	return nil
}

func (s *Instance) runServers(errCh chan error) {
	for _, srv := range s.servers {
		go func(s server, ch chan error) {
			if err := s.Run(); err != nil {
				ch <- err
			}
		}(srv, errCh)
	}
}

func (s *Instance) stopServers(ctx context.Context) {
	for _, srv := range s.servers {
		_ = srv.Stop(ctx)
	}
}

func (s *Instance) Timeout() time.Duration {
	return time.Duration(s.timeout) * time.Second
}

func (s *Instance) isStopped() bool {
	return s.state == serviceSTOPPED
}
func (s *Instance) isNew() bool {
	return s.state == serviceNEW
}

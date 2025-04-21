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

	"github.com/rs/zerolog/log"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/db"
	userAggregate "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/aggregate/user"
	grpcUserCtrl "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/grpc/user"
	httpUserCtrl "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/http/user"
	pubsubUserCtrl "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/pubsub/user"
	userService "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/service/user"
	userProto "github.com/nachoconques0/user_challenge_svc/pkg/challenge/proto/user"
	simplePubSub "github.com/nachoconques0/user_challenge_svc/pkg/challenge/pubsub/local"
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

	// DB connection
	dbConn, err := db.New(options.dbOptions...)
	if err != nil {
		return err
	}

	// Initialize Bus
	bus := simplePubSub.NewBus(dbConn)

	// Register Subscribers
	err = pubsubUserCtrl.RegisterUserSubscribers(bus)
	if err != nil {
		return err
	}

	// Aggregate
	userAgg, err := userAggregate.New(dbConn, "nontest", bus)
	if err != nil {
		return err
	}

	// Service
	userSvc := userService.New(userAgg)

	// Controller
	httpCtrl := httpUserCtrl.NewController(userSvc)
	grpcCtrl := grpcUserCtrl.NewController(userSvc)

	// HTTP Server
	httpSrv, err := httpServer.New(httpServer.WithAddress(fmt.Sprintf(":%s", options.httpPort)))
	if err != nil {
		return err
	}
	httpRouter := httpServer.InitHTTPRouter(httpSrv)
	httpServer.InitUserRoutes(httpRouter, httpCtrl)

	// gRPC Server
	grpcSrv := grpcServer.New(options.gRPCPort)
	userProto.RegisterUserServiceServer(grpcSrv.Server(), grpcCtrl)

	i := Instance{
		timeout: 20,
		servers: []server{httpSrv, grpcSrv},
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(quitCh)

	return i.Run(quitCh)
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

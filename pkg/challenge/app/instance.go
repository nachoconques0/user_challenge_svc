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

	userAggregate "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/aggregate/user"
	userController "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/http/user"
	userService "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/service"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/server"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/db"
)

// Possible service states
const (
	serviceNEW = iota
	serviceRUNNING
	serviceSTOPPING
	serviceSTOPPED
)

// Instance definitions
type Instance struct {
	servers []server.Server
	timeout int
	state   int
	mu      sync.Mutex
}

func New(opts ...Option) error {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	// Initialize DB
	db, err := db.New(options.dbOptions...)
	if err != nil {
		return err
	}

	// Initialize Aggregate
	userAgg, err := userAggregate.New(db, "nontest")
	if err != nil {
		return err
	}

	// Initialize Service
	userSvc := userService.New(userAgg)

	// Initialize Controllers
	userCtrl := userController.NewController(userSvc)

	// Initialize HTTP server
	httpServer, err := server.New(
		server.WithAddress(fmt.Sprintf(":%s", options.httpPort)),
	)
	if err != nil {
		return err
	}

	httpRouter := server.InitHTTPRouter(httpServer)
	server.InitUserRoutes(httpRouter, userCtrl)

	i := Instance{
		timeout: 20,
		servers: []server.Server{
			httpServer,
		},
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(quitCh)

	return i.Run(quitCh)
}

// Run the service until SIGTERM signal is received. If any failure occurs on Run it will return an error.
// If the application is already running nothing will happen
func (s *Instance) Run(quitCh chan os.Signal) error {
	// If the service was already stopped we cannot start it again
	if s.isStopped() {
		return errors.New("instance was already stopped. Can't start it again")
	}

	// Only run if service is new
	if s.isNew() {
		log.Info().Msg("Application: starting...")

		// Lock the instance while it is being started
		s.mu.Lock()
		s.state = serviceRUNNING

		// Handle startup errors
		errCh := make(chan error, 1)

		// Look at SIGTERM system event
		if quitCh == nil {
			quitCh = make(chan os.Signal, 1)
			signal.Notify(quitCh, syscall.SIGTERM, os.Interrupt)
			defer signal.Stop(quitCh)
		}

		// Run registered servers
		s.runServers(errCh)

		s.mu.Unlock()

		log.Info().Msg("Application: running...")

		// Block waiting for shutdown
		select {
		case err := <-errCh:
			return err
		case <-quitCh:
			s.mu.Lock()
			s.state = serviceSTOPPING
			s.mu.Unlock()

			// Start context with cancellation set for the defined timeout
			log.Info().Msg(fmt.Sprintf("Application: stopping in %.0fs...\n", s.Timeout().Seconds()))
			ctx, cancel := context.WithTimeout(context.Background(), s.Timeout())
			defer cancel()

			// Stop servers and subscribers
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

// Run all the registered servers
func (s *Instance) runServers(errCh chan error) {
	for _, srvr := range s.servers {
		go func(server server.Server, ch chan error) {
			if err := server.Run(); err != nil {
				ch <- err
			}
		}(srvr, errCh)
	}
}

func (s *Instance) stopServers(ctx context.Context) {
	for _, server := range s.servers {
		_ = server.Stop(ctx)
	}

}

func (s *Instance) IsRunning() bool {
	return s.state == serviceRUNNING || s.state == serviceSTOPPING
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

package gracefulshutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type GracefulShutdown struct {
	services []Service

	signal  chan os.Signal
	errChan chan error
}

// New creates a new GracefulShutdown instance.
func New() *GracefulShutdown {
	return &GracefulShutdown{
		services: []Service{},
		signal:   make(chan os.Signal, 1),
		errChan:  make(chan error, 1),
	}
}

func (g *GracefulShutdown) Add(server Service) {
	g.services = append(g.services, server)
}

func (g *GracefulShutdown) CatchSignals() {
	signal.Notify(g.signal, syscall.SIGINT, syscall.SIGTERM)
}

func (g *GracefulShutdown) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, service := range g.services {
		go func(srv Service) {
			log.Printf("Starting service: %s", srv.Name())
			if err := srv.Start(ctx); err != nil {
				log.Printf("Error starting service: %v", err)
				g.errChan <- err
			}
		}(service)
	}

	defer g.shutdown()
	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal, shutting down services...")
	case e := <-g.errChan:
		log.Println("Service start error, shutting down services...")
		return e
	case sig := <-g.signal:
		log.Printf("Received signal: %s, shutting down services...", sig)
	}

	return nil
}

func (g *GracefulShutdown) shutdown() {
	var wg sync.WaitGroup
	for _, service := range g.services {
		wg.Add(1)
		go func(srv Service) {
			log.Printf("Shutting down service: %s", srv.Name())
			if err := srv.Shutdown(); err != nil {
				log.Printf("Error shutting down service: %v", err)
			}
			wg.Done()
		}(service)
	}
	wg.Wait()
	log.Println("All services shut down gracefully.")
}

package main

import (
	"fmt"
	grpc2 "github.com/travix/gotf-example/example-server/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const shutdownDelay = 5 * time.Second

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, FormatTimestamp: formatter})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// create service servers
	shutdown := make(chan bool)
	done := make(chan bool)
	// start servers
	grpc2.NewServer(&grpc2.Servicer{}, shutdown, done)
	handleShutdown(shutdown, done)
	log.Info().Msg("shutdown complete")
}

func formatter(i interface{}) string {
	if i == nil {
		return ""
	}
	if tt, ok := i.(string); ok {
		ts, err := time.ParseInLocation(time.RFC3339, tt, time.Local)
		if err != nil {
			i = tt
		} else {
			i = ts.Local().Format(time.Kitchen)
		}
	}
	return fmt.Sprintf("\u001B[32m[example-server] \x1b[90m%v\x1b[0m", i)
}

// handleShutdown will wait for syscall.SIGINT, syscall.SIGTERM
// Once received will gracefully shutdown server in 10 seconds else will.
func handleShutdown(closed chan<- bool, done <-chan bool) {
	infoLogger := log.Level(zerolog.InfoLevel)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	// received shutdown
	infoLogger.Info().Msgf("received %s(%#v)! shutting down", sig.String(), sig)
	close(closed)
	infoLogger.Info().Msg("waining for grpc to stop")
	select {
	case <-time.After(shutdownDelay):
		return
	case <-done:
	}
}

package grpc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	hmac "github.com/yogeshlonkar/go-grpc-hmac"

	"github.com/rs/zerolog/log"
	hmac "github.com/yogeshlonkar/go-grpc-hmac"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type gRPC struct {
	addr   string
	server *grpc.Server
}

func (s *gRPC) handleShutdown(shutdown <-chan bool, done chan<- bool) {
	<-shutdown
	s.server.GracefulStop()
	log.Info().Msg("gRPC server stopped")
	done <- true
}

func (s *gRPC) start() {
	s.addr = ":50051"
	if os.Getenv("GRPC_PORT") != "" {
		s.addr = fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	}
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	if err = s.server.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve gRPC")
	}
}

func getSecrets(keyId string) (string, error) {
	// same keys are used in terraform scripts
	if keyId == os.Getenv("TF_VAR_key_id") {
		return os.Getenv("TF_VAR_secret_key"), nil
	}
	if os.Getenv("TF_VAR_key_id") == "" {
		return "", errors.New("TF_VAR_key_id is not set")
	}
	return "", nil
}

func (s *gRPC) setup() {
	interceptor := hmac.NewServerInterceptor(getSecrets)
	opts := []grpc.ServerOption{
		interceptor.UnaryInterceptor(),
		grpc.Creds(insecure.NewCredentials()),
	}
	s.server = grpc.NewServer(opts...)
}

// NewServer initializes the gRPC service and the server object.
func NewServer(service *Servicer, shutdown <-chan bool, done chan<- bool) {
	s := &gRPC{}
	s.setup()
	service.RegisterGRPC(s.server)
	go s.start()
	go s.handleShutdown(shutdown, done)
	time.Sleep(1 * time.Second)
	log.Info().Msgf("Listening and serving gRPC on %s", s.addr)
}

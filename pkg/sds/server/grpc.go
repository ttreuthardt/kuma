package server

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	envoy_discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	envoy_server "github.com/envoyproxy/go-control-plane/pkg/server/v2"

	sds_config "github.com/kumahq/kuma/pkg/config/sds"
	"github.com/kumahq/kuma/pkg/core"
	"github.com/kumahq/kuma/pkg/core/runtime/component"
)

const grpcMaxConcurrentStreams = 1000000

var (
	grpcServerLog = core.Log.WithName("sds-server").WithName("grpc")
)

type grpcServer struct {
	server envoy_server.Server
	config sds_config.SdsServerConfig
}

func (s *grpcServer) NeedLeaderElection() bool {
	return false
}

var (
	_ component.Component = &grpcServer{}
)

func (s *grpcServer) Start(stop <-chan struct{}) error {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	useTLS := s.config.TlsCertFile != ""
	if useTLS {
		creds, err := credentials.NewServerTLSFromFile(s.config.TlsCertFile, s.config.TlsKeyFile)
		if err != nil {
			return errors.Wrap(err, "failed to load TLS certificate")
		}
		grpcOptions = append(grpcOptions, grpc.Creds(creds))
	}
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GrpcPort))
	if err != nil {
		return err
	}

	// register services
	envoy_discovery.RegisterSecretDiscoveryServiceServer(grpcServer, s.server)

	errChan := make(chan error)
	go func() {
		defer close(errChan)
		if err = grpcServer.Serve(lis); err != nil {
			grpcServerLog.Error(err, "terminated with an error")
			errChan <- err
		} else {
			grpcServerLog.Info("terminated normally")
		}
	}()
	grpcServerLog.Info("starting", "interface", "0.0.0.0", "port", s.config.GrpcPort, "tls", useTLS)

	select {
	case <-stop:
		grpcServerLog.Info("stopping gracefully")
		grpcServer.GracefulStop()
		return nil
	case err := <-errChan:
		return err
	}
}

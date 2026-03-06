package transportgrpc

import (
	"context"
	"fmt"
	"net"

	"team_service/internal/adapter"
	"team_service/internal/infrastructure/share/settings"
	"team_service/proto/team_service"

	"github.com/thanvuc/go-core-lib/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TeamServer struct {
	logger log.LoggerV2
	config *settings.Server

	adapter *adapter.Dependency
	server  *grpc.Server
}

func NewTeamServer(
	logger log.LoggerV2,
	config *settings.Server,
	adapter *adapter.Dependency,
) *TeamServer {
	return &TeamServer{
		logger:  logger,
		config:  config,
		adapter: adapter,
	}
}

func (s *TeamServer) Start(ctx context.Context) error {

	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", s.config.Host, s.config.TeamPort),
	)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()

	team_service.RegisterGroupServiceServer(
		s.server,
		s.adapter.GroupController,
	)

	team_service.RegisterSprintServiceServer(
		s.server,
		s.adapter.SprintController,
	)

	team_service.RegisterWorkServiceServer(
		s.server,
		s.adapter.WorkController,
	)

	go s.Stop(ctx)

	s.logger.Info(
		fmt.Sprintf(
			"gRPC server listening on %s:%d",
			s.config.Host,
			s.config.TeamPort,
		),
	)

	if err := s.server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		s.logger.Error("failed to serve grpc", log.WithFields(zap.Error(err)))
		return err
	}

	return nil
}

func (s *TeamServer) Stop(ctx context.Context) {
	<-ctx.Done()

	s.logger.Info("shutting down grpc server")

	s.server.GracefulStop()

	s.logger.Info("grpc server stopped")
}

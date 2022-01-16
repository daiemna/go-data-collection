package app

import (
	"fmt"
	"net"

	pb "github.com/daiemna/go-data-collection/internal/services"
	"github.com/daiemna/go-data-collection/pkg/service"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func (t *CliApp) startServer(c *cli.Context) error {
	cassConf := t.Config.Database.ToCassConfig()
	server := service.NewTimeseriesServer(cassConf, "collection server")
	// create listener
	lis, err := net.Listen("tcp", t.Config.Server.HostPort)
	if err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}
	// var opts []grpc.ServerOption
	// opts = append(opts, grpc.)
	// create grpc server
	s := grpc.NewServer()
	pb.RegisterTimeSeriesDataServer(s, server)

	log.Info().Msgf("Server starting!")

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

//InitServer initializes GRPC data server.
func (t *CliApp) InitServer() *cli.Command {
	return &cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Starts the server",
		Action:  t.startServer,
	}
}

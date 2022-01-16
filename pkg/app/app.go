package app

import (
	"context"

	"github.com/daiemna/go-data-collection/internal/config"
	"github.com/daiemna/go-data-collection/pkg/client"
	"github.com/daiemna/go-data-collection/pkg/service"
	cli "github.com/urfave/cli/v2"
)

//CliApp is the central struct to manage the app
type CliApp struct {
	Cli        *cli.App
	Config     *config.Config
	Client     *client.TimeseriesClient
	Server     *service.TimeseriesServer
	cancelFunc context.CancelFunc
}

//NewClientApp creates a new cli client.
func NewClientApp(conf *config.Config, client *client.TimeseriesClient) *CliApp {
	app := NewDefaultCliApp()

	app.Client = client
	app.Config = conf
	return app
}

//NewServerApp creates new server cli app.
func NewServerApp(conf *config.Config, server *service.TimeseriesServer) *CliApp {
	app := NewDefaultCliApp()

	app.Server = server
	app.Config = conf
	return app
}

// NewDefaultCliApp creates new cli app.
func NewDefaultCliApp() *CliApp {
	return &CliApp{
		Cli: cli.NewApp(),
	}
}

// RunWithContext is a wrapper on urfave/cli RunContext function
func (t *CliApp) RunWithContext(ctx context.Context, arguments []string) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	t.cancelFunc = cancelFunc
	return t.Cli.RunContext(ctx, arguments)
}

func (t *CliApp) Close() {
	if t.Server != nil {
		t.Server.Close()
	}
	t.cancelFunc()
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/daiemna/go-data-collection/internal/config"
	"github.com/daiemna/go-data-collection/pkg/app"
	application "github.com/daiemna/go-data-collection/pkg/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// BuildTime gets populated during the build proces
	BuildTime = ""

	//Version gets populated during the build process
	Version = ""
)

const serverAddr = "localhost:50005"
const recordCount = 10
const TSID = "DNA01.temperature"

func setUpLogger(c *cli.Context) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if c.Bool("debug") {
		fmt.Fprintf(c.App.Writer, "Setting the log level to debug\n")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func setUpConfig(c *cli.Context) *config.Config {
	path, err := filepath.Abs(c.String("config"))
	if err != nil {
		log.Fatal().Msgf("error in absolute path: %s", err)
	}
	conf, err := config.NewServerFromFile(path)
	if err != nil {
		log.Fatal().Msgf("Unable to read config from path: %s, %s", path, err)
	}
	return conf
}

func initCliOrPanic(app *app.CliApp, runTimeConf *config.Config) *cli.App {
	cliCtx := cli.NewApp()
	cliCtx.Name = "data-collection-grpc"
	cliCtx.Description = "A service for GRPC data data-collection."
	cliCtx.Version = Version
	cliCtx.Authors = []*cli.Author{
		{
			Name:  "Daiem Ali",
			Email: "daiem.dna@gmail.com",
		},
	}
	cliCtx.EnableBashCompletion = true
	cliCtx.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "buildtime",
			Usage: "time of this build",
		},
		&cli.StringFlag{
			Name:     "config",
			Usage:    "Path to the config file, for e.g. --config=/absolute/path/to/config.yaml",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Set the log level to debug",
		},
	}

	cliCtx.Before = func(c *cli.Context) error {
		if c.IsSet("buildtime") {
			fmt.Fprintf(c.App.Writer, "%v buildtime: %v\n", cliCtx.Name, BuildTime)
			return cli.Exit("", 0)
		}
		// Set up log level based on config
		setUpLogger(c)
		// Set up configuration if not given
		if runTimeConf != nil {
			app.Config = runTimeConf
		} else {
			app.Config = setUpConfig(c)
		}
		return nil
	}
	cliCtx.Commands = []*cli.Command{
		app.InitServer(),
		app.InitClient(),
	}
	sort.Sort(cli.CommandsByName(cliCtx.Commands))
	return cliCtx
}

func main() {
	ctx := context.Background()

	app := application.NewDefaultCliApp()
	app.Cli = initCliOrPanic(app, nil)

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// block on this channel for ctrl-c
		<-c
		app.Close()
		log.Info().Msg("\r- ctrl+C pressed in Terminal. Exiting in 2 seconds.")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()
	log.Info().Msgf("provided flags : %v", os.Args[1:])
	// We pass our custom context, so that
	// we can use this ctx to logout from email server
	// at the end of the program
	err := app.RunWithContext(ctx, os.Args[1:])
	if err != nil {
		log.Error().Msgf("Error in app run: %s", err)
	}
	// clean up here for sure
	app.Close()
	time.Sleep(2 * time.Second)

}

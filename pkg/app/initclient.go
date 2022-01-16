package app

import (
	"math/rand"
	"time"

	"github.com/daiemna/go-data-collection/internal/utils"
	"github.com/daiemna/go-data-collection/pkg/client"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func (t *CliApp) initClient(c *cli.Context) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	errChan := make(chan error)
	conn, err := grpc.Dial(t.Config.Server.HostPort, opts...)
	if err != nil {
		log.Fatal().Msgf("Error in dialing server: %v", err)
	}
	defer conn.Close()
	recordCount := c.Int("record-count")
	TSID := c.String("ts-id")
	sourceId := c.String("source-id")
	rand.Seed(time.Now().UnixNano())
	tsClient := client.NewClient(conn)
	points := make([]client.TimeSeriesPoint, recordCount)
	randFloats := utils.RandFloats(0.0, 1.0, recordCount)
	for i := 0; i < recordCount; i++ {
		points[i].Time = time.Now().Unix() + int64(i)
		points[i].Value = randFloats[i]
	}
	go tsClient.SendPoints(TSID, sourceId, points, errChan)

	for keepReadingErrors := true; keepReadingErrors; {
		select {
		case err, ok := <-errChan:
			if !ok {
				keepReadingErrors = false
				break
			}
			log.Error().Msgf("Error received : %v", err)
		default:
			time.Sleep(time.Millisecond * 1)
		}
	}
	return nil
}

//InitClient initializes GRPC data server.
func (t *CliApp) InitClient() *cli.Command {
	return &cli.Command{
		Name:    "client",
		Aliases: []string{"c"},
		Usage:   "Starts the client",
		Action:  t.initClient,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "record-count",
				Value: 10,
				Usage: "Number of record to send",
			},
			&cli.StringFlag{
				Name:  "ts-id",
				Value: "b34eecb4-5d2a-11ec-b564-cf3495245a9c",
				Usage: "The UUID to be used for sending the data.",
			},
			&cli.StringFlag{
				Name:  "source-id",
				Value: "e04bbf1e-701f-11ec-8fca-8369ee9e1742",
				Usage: "The UUID to be used for sending the data.",
			},
		},
	}
}

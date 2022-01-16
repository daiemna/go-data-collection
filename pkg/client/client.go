package client

import (
	"context"
	"io"
	"sync"

	pb "github.com/daiemna/go-data-collection/internal/services"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

//TimeSeriesPoint can be used to store time series data point.
type TimeSeriesPoint struct {
	Time  int64
	Value float32
}

//TimeseriesClient is client used to send data to server.
type TimeseriesClient struct {
	pb.TimeSeriesDataClient
}

//NewClient creates new TimeseriesClient
func NewClient(conn grpc.ClientConnInterface) *TimeseriesClient {
	return &TimeseriesClient{
		TimeSeriesDataClient: pb.NewTimeSeriesDataClient(conn),
	}
}

// SendPoints sends timeseries data to server for saving.
func (client *TimeseriesClient) SendPoints(tsID, sourceId string, points []TimeSeriesPoint, errChan chan error) {
	log.Info().Msg("Sending data...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	stream, err := client.StreamDataframes(ctx)
	if err != nil {
		errChan <- err
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for res, err := stream.Recv(); err != io.EOF; res, err = stream.Recv() {
			if err != nil {
				// log.Error().Msgf("Error Received: %v", err)
				errChan <- err
			}
			if res != nil {
				log.Info().Msgf("Server response for dataframe `%s` message: %v", res.Dataframeid, res.Response.Msg)
				return
			}
		}
	}()
	log.Debug().Msgf("points: %v", points)

	series := pb.DataSeries{
		Datapointid: tsID,
		Values:      make([]*pb.SeriesPoint, len(points)),
	}
	for i, point := range points {
		series.Values[i] = &pb.SeriesPoint{
			Timestamp: point.Time,
			Value:     point.Value,
		}
	}
	randUUID := uuid.New()
	df := pb.Dataframe{
		Sourceid:    sourceId,
		Dataframeid: randUUID.String(),
		Nseries:     make([]*pb.DataSeries, 1),
	}

	df.Nseries[0] = &series
	log.Debug().Msgf("Sendable data container: %v", df)
	stream.Send(&df)
	stream.CloseSend()
	wg.Wait()
	log.Info().Msg("Done!")
	close(errChan)
}

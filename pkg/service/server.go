package service

import (
	"fmt"
	"io"
	"time"

	"github.com/daiemna/go-data-collection/internal/bigdb"
	"github.com/daiemna/go-data-collection/internal/services"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

//TimeseriesServer is for services to cassandra save.
type TimeseriesServer struct {
	services.UnimplementedTimeSeriesDataServer
	serverName       string
	cassandraConf    *bigdb.CassandraConfig
	cassandraSession *gocql.Session
	batchesChan      chan *gocql.Batch
}

// Close closes the NoSQL DB and batches channel
func (t *TimeseriesServer) Close() {
	t.cassandraSession.Close()
	close(t.batchesChan)
}

// NewTimeseriesServer can be used to initalize a TimeseriesServer.
func NewTimeseriesServer(cassConf *bigdb.CassandraConfig, serverName string) *TimeseriesServer {
	cassClusterConf := bigdb.NewCassandraClusterConfig(cassConf)
	dbSession, err := bigdb.NewCassandraSession(cassClusterConf)
	if err != nil {
		log.Error().Msgf("Error creating cassandra session : %v", err)
		return nil
	}

	tss := TimeseriesServer{
		serverName:       serverName,
		cassandraConf:    cassConf,
		cassandraSession: dbSession,
		batchesChan:      make(chan *gocql.Batch, 10),
	}
	return &tss
}

func (server *TimeseriesServer) processDataSeries(timeseries *services.DataSeries, batch *gocql.Batch) int {
	validRecordCount := 0
	log.Debug().Msgf("timeseries : %v", timeseries)
	for i, point := range timeseries.Values {
		line := fmt.Sprintf("%d;%s;%f", point.Timestamp, timeseries.Datapointid, point.Value)
		if qerr := bigdb.Line2Query(line, server.cassandraConf.Table, batch); qerr != nil {
			log.Error().Msg(fmt.Sprintf("Error in timeseries id `%s` point no `%d`: %v", timeseries.Datapointid, i, qerr))
			continue
		}
		validRecordCount = validRecordCount + 1
	}
	return validRecordCount
}

func (server *TimeseriesServer) processDataRecords(record *services.DataRecord, batch *gocql.Batch) int {
	validRecordCount := 0
	for i, point := range record.Values {
		line := fmt.Sprintf("%d;%s;%f", record.Timestamp, point.Datapointid, point.Value)
		if qerr := bigdb.Line2Query(line, server.cassandraConf.Table, batch); qerr != nil {
			log.Error().Msg(fmt.Sprintf("Error in timeseries id `%s` point no `%d`: %v", point.Datapointid, i, qerr))
			continue
		}
		validRecordCount = validRecordCount + 1
	}
	return validRecordCount
}

// StreamDataSeries handels streaming timeseries data
func (server *TimeseriesServer) StreamDataSeries(stream services.TimeSeriesData_StreamDataSeriesServer) error {
	for timeseries, err := stream.Recv(); err != io.EOF; timeseries, err = stream.Recv() {
		if timeseries == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		batch := gocql.NewBatch(gocql.LoggedBatch)
		validCount := server.processDataSeries(timeseries, batch)
		if validCount > 0 {
			server.batchesChan <- batch
		}
		var msg string
		recordErrorCount := len(timeseries.Values) - validCount
		if recordErrorCount == 0 {
			msg = fmt.Sprintf("timeseries `%s` has been presisted.", timeseries.Datapointid)
		} else {
			msg = fmt.Sprintf("Error in %d points in timeseries %s.", recordErrorCount, timeseries.Datapointid)
		}
		stream.Send(&services.Response{Msg: msg})
	}
	return nil
}

//StreamDataframe handels streaming Batches to Cassandra
func (server *TimeseriesServer) StreamDataframes(stream services.TimeSeriesData_StreamDataframesServer) error {
	for dataframe, err := stream.Recv(); err != io.EOF; dataframe, err = stream.Recv() {
		//Process all `timeseries`
		for _, timeseries := range dataframe.Nseries {
			cassBatch := gocql.NewBatch(gocql.LoggedBatch)
			validCount := server.processDataSeries(timeseries, cassBatch)
			if validCount > 0 {
				server.batchesChan <- cassBatch
			}
			msg := fmt.Sprintf("Error in %d points in timeseries %s.", len(timeseries.Values)-validCount, timeseries.Datapointid)
			err = stream.Send(&services.DataframeResponse{
				Dataframeid: dataframe.Dataframeid,
				Response:    &services.Response{Msg: msg},
			})
			if err != nil {
				log.Debug().Msgf("Error in sending dataframe response: %v", err)
			}
		}

		//Process all `values at time Batch`
		for _, record := range dataframe.Records {
			cassBatch := gocql.NewBatch(gocql.LoggedBatch)
			validCount := server.processDataRecords(record, cassBatch)
			// msg := "No errors occurred!"
			if validCount > 0 {
				server.batchesChan <- cassBatch
			}
			msg := fmt.Sprintf("Error in %d records in dataframe %s.", len(record.Values)-validCount, dataframe.Dataframeid)
			err = stream.Send(&services.DataframeResponse{
				Dataframeid: dataframe.Dataframeid,
				Response:    &services.Response{Msg: msg},
			})
			if err != nil {
				log.Debug().Msgf("Error in sending dataframe response: %v", err)
			}
		}
	}
	return nil
}

// StreamRecords recives values at time t.
func (server *TimeseriesServer) StreamRecords(stream services.TimeSeriesData_StreamRecordsServer) error {
	cassBatch := gocql.NewBatch(gocql.LoggedBatch)
	for record, err := stream.Recv(); err != io.EOF; record, err = stream.Recv() {
		validCount := server.processDataRecords(record, cassBatch)
		if validCount > 0 {
			server.batchesChan <- cassBatch
		}
		msg := fmt.Sprintf("Error in %d records at timestap %v.", len(record.Values)-validCount, record.Timestamp)
		err = stream.Send(&services.Response{
			Msg: msg,
		})
		if err != nil {
			log.Debug().Msgf("Error in sending dataframe response: %v", err)
		}
	}
	return nil
}

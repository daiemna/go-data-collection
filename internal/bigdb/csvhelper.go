package bigdb

import (
	"errors"
	"strconv"

	"github.com/gocql/gocql"
)

type csvLineParser interface {
	fromLine(line []string) error
}

type csvRow struct {
	timestamp int64
	datapoint string
	value     float64
}

type csvRowWithUUID struct {
	csvRow
	datapoint gocql.UUID
}

func (row *csvRow) fromLine(line []string) error {
	utcTime, err := convertToUTCMilli(line[0])
	ts, a2iErr := strconv.ParseInt(line[0], 10, 64)
	if err != nil && a2iErr != nil {
		return errors.New("Failed to parse the timestamp: " + err.Error())
	}
	if err == nil {
		row.timestamp = utcTime
	} else if a2iErr == nil {
		row.timestamp = ts
	}

	val, err := toFloat(line[2])
	if err != nil {
		return errors.New("Illegal value format")
	}
	row.datapoint = line[1]
	row.value = val
	return nil
}

func (row *csvRowWithUUID) fromLine(line []string) error {
	err := row.csvRow.fromLine(line)
	if err != nil {
		return err
	}
	tsidUUIDFormat, err := gocql.ParseUUID(line[1])
	if err != nil {
		return errors.New("Failed to parse the UUID")

	}
	row.datapoint = tsidUUIDFormat
	return nil
}

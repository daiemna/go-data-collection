package bigdb

import (
	"errors"

	"github.com/gocql/gocql"
)

type csvRowWithUUID struct {
	CsvRow
	datapoint gocql.UUID
}

func (row *csvRowWithUUID) getRow() *CsvRow {
	return &CsvRow{
		timestamp: row.timestamp,
		value:     row.value,
		datapoint: row.datapoint.String(),
	}
}

func (row *csvRowWithUUID) fromLine(line []string) error {
	tsidUUIDFormat, err := gocql.ParseUUID(line[1])
	if err != nil {
		return errors.New("failed to parse the UUID")

	}
	row.datapoint = tsidUUIDFormat
	return nil
}

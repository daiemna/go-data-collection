package bigdb

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

//CsvLineParser is interface to parse a csv line to database row.
type CsvLineParser interface {
	//fromLine must be implemented to receive tokenized csv line and parse it.
	fromLine(line []string) error
	//getRow must be implemented to retrive the parsed row.
	getRow() *CsvRow
}

type CsvRow struct {
	timestamp int64
	datapoint string
	value     float64
}

func (row *CsvRow) getRow() *CsvRow {
	return row
}

func (row *CsvRow) fromLine(line []string) error {
	utcTime, err := convertToUTCMilli(line[0])
	ts, toIntErr := strconv.ParseInt(line[0], 10, 64)
	if err != nil && toIntErr != nil {
		return fmt.Errorf("failed to parse the timestamp `%v`: %v", line[0], err.Error())
	}
	if err == nil {
		row.timestamp = utcTime
	} else if toIntErr == nil {
		row.timestamp = ts
	}

	val, err := toFloat(line[2])
	if err != nil {
		return errors.New("illegal value format, value must be float")
	}
	row.datapoint = line[1]
	row.value = val
	return nil
}

//ParseCsvLine attempts to parse csv line tokens to database row,
// which can be retrived by calling getRow().
func ParseCsvLine(line []string) (CsvLineParser, error) {
	csvLineUUID := &csvRowWithUUID{}
	err := csvLineUUID.fromLine(line)
	if err != nil {
		csvLine := &CsvRow{}
		perr := csvLine.fromLine(line)
		if perr != nil {
			log.Debug().Msgf("attempt to parse id to UUID also failed: %v", err)
			return nil, perr
		}
		return csvLine, nil
	}
	err = csvLineUUID.CsvRow.fromLine(line)
	if err != nil {
		return nil, fmt.Errorf("error parsing line: %v", err)
	}
	return csvLineUUID, nil
}

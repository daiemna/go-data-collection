// Recogizer License
//
// Copyright (C) Recogizer Group GmbH - All Rights Reserved.
//
// Unauthorized copying of this file, via any medium is strictly prohibited.
// Proprietary and confidential.

package bigdb

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/araddon/dateparse"
	"github.com/gocql/gocql"
)

// DefaultCassandraPort is default r/w port of cassandra
const DefaultCassandraPort = "9042"

var (
	LocalLoc *time.Location

	// UTC represents a timezone for UTC
	UTC *time.Location
)

// Initialize two different timezones to automatically parse timestamps
func init() {
	// berlinLoc represents a timezone for Berlin
	LocalLoc, _ = time.LoadLocation("Local")

	// UTC represents a timezone for UTC
	UTC, _ = time.LoadLocation("UTC")
}

// Convert a given input to its float representation
//
// It also tries to convert a given input string containing True, False to 1.0, 0.0 else returns error
func toFloat(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if strings.Compare(input, "False") == 0 {
		return 0.0, nil
	}
	if strings.Compare(input, "True") == 0 {
		return 1.0, nil
	}
	return strconv.ParseFloat(strings.TrimSpace(input), 32)
}

// Convert a given timestamp string to its UTC milliseconds representation
func convertToUTCMilli(ts string) (int64, error) {
	parsed, err := dateparse.ParseIn(ts, UTC)
	if err != nil {
		return 0, err
	}
	return parsed.UTC().UnixNano() / int64(time.Millisecond), nil
}

func toUUID4(uuid30 string) (string, error) {
	var newuuid string
	if len(uuid30) < 30 {
		return "", fmt.Errorf("require atleast 30 char string for conversion to UUID4")
	}
	newuuid = uuid30[0:8] + "-" + uuid30[8:12] + "-4" + uuid30[12:15] + "-a" + uuid30[15:18] + "-" + uuid30[18:30]
	if _, err := gocql.ParseUUID(newuuid); err != nil {
		return "", fmt.Errorf("failed to parse the UUID, check provided uuid all charecters must form a valid uuid")
	}
	return newuuid, nil
}

// ExecuteBatch executes a given batch of data and writes it to Cassandra.
// This reads from a channel `batchesChan` where data comes after files are being parsed.
func ExecuteBatch(session *gocql.Session, batchesChan chan *gocql.Batch, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for batch := range batchesChan {
		for {
			if err := session.ExecuteBatch(batch); err != nil {
				log.Info().Msgf("execute batch: %v", err)
				continue
			} else {
				// log.Debug().Msg("Batch executed")
				break
			}
		}
	}
}

// Line2Query converts a csv line to a cql query and add it to batch
func Line2Query(
	scanedText string,
	table string,
	batch *gocql.Batch) error {

	line := strings.Split(scanedText, ";")
	lineLen := len(line)
	if len(line) != 3 {
		return fmt.Errorf("unreadable row, must have 3 tokens but %v was given", lineLen)
	}
	var query string
	// var row csvRowWithUUID
	lineParser, err := ParseCsvLine(line)
	if err != nil {
		return err
	}
	row := lineParser.getRow()
	query = fmt.Sprintf("INSERT INTO %s (tsid, time, value ) VALUES (?, ?, ?)", table)
	batch.Query(query, row.datapoint, row.timestamp, float32(row.value))
	log.Debug().Msgf("Write query: %s", query)
	return nil
}

// IsValidPath checks if file/dir exists
func IsValidPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// ValidateHost checks if the host is valid cassandra host,
// append standrad cassandra port if not give.
func ValidateHost(hostPort string) (string, error) {
	host, port, err := net.SplitHostPort(hostPort)
	log.Debug().Msgf("err in host parsing: %v", err)
	if err != nil {
		log.Debug().Msgf("index of missing port in err : %d", strings.Index(err.Error(), "missing port"))
	}
	log.Debug().Msgf("Host: %s, Port: %s", host, port)
	if err != nil && strings.Contains(err.Error(), "missing port") {
		ip := net.ParseIP(hostPort)
		if ip == nil && hostPort != "localhost" {
			return "", fmt.Errorf("`%s` is not valid hostname, host name can be an IP address or `localhost`", hostPort)
		}
		return hostPort + ":" + DefaultCassandraPort, nil
	} else if err != nil {
		return "", err
	}
	return host + ":" + port, nil
}

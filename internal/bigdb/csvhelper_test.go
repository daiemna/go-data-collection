package bigdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func Test_ParseCsvLine(t *testing.T) {
	assert := assert.New(t)

	testthisdata := func(ts int64, v float64, id string) {
		csvLine := []string{fmt.Sprintf("%d", ts), id, fmt.Sprintf("%f", v)}
		parser, err := ParseCsvLine(csvLine)
		if err != nil {
			assert.Nil(err, fmt.Sprintf("there should be no error in parsing this row: %v", csvLine))
		}
		row := parser.getRow()
		assert.EqualValues(row.datapoint, id, "id is not same")
		assert.EqualValues(row.timestamp, ts, "timestamp is not same")
		assert.EqualValues(row.value, v, "value is not same")
	}

	timestamp := time.Now().UnixMilli()
	value := float64(10.0)
	id := "8"
	testthisdata(timestamp, value, id)
	id = "d0cc4482-7b0f-11ec-941b-cbfab95d687d"
	testthisdata(timestamp, value, id)
}

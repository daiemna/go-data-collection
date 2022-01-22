// Recogizer License
//
// Copyright (C) Recogizer Group GmbH - All Rights Reserved.
//
// Unauthorized copying of this file, via any medium is strictly prohibited.
// Proprietary and confidential.

package bigdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func Test_toFloat(t *testing.T) {
	assert := assert.New(t)

	input := []struct {
		input  string
		output float64
		errMsg string
	}{{
		input:  "False",
		output: 0.0,
		errMsg: "False should have returned 0.0",
	},
		{
			input:  "True",
			output: 1.0,
			errMsg: "True should have returned 1.0",
		},
		{
			input:  "random",
			output: -100.0,
			errMsg: fmt.Sprintf("strconv.ParseFloat: parsing \"random\": invalid syntax"),
		},
	}

	for _, val := range input {
		actual, err := toFloat(val.input)
		if val.input == "random" {
			assert.Equal(err.Error(), val.errMsg)
		} else {
			assert.Equal(val.output, actual, val.errMsg)
		}
	}

}

func Test_toUUID4(t *testing.T) {
	assertions := assert.New(t)
	uuid30 := "fc000c1155541809d6e4cf79534f46"
	uuid4 := "fc000c11-5554-4180-a9d6-e4cf79534f46"
	uuidReturned, err := toUUID4(uuid30)
	assertions.Nil(err, "there should be no err in conversion of uuid")
	assertions.Equal(uuidReturned, uuid4, "uuid conversion is wrong")

	uuid30 = "fc000c1155541809d6e4cg79534f46"
	_, err = toUUID4(uuid30)
	assertions.NotNil(err, "invalid char in uuid30, but conversion did not throw an error")

	uuid30 = "fc000c1155541809d6e4ce79534"
	_, err = toUUID4(uuid30)
	assertions.NotNil(err, "in valid char count in uuid30, but conversion did not throw an error")
}

func Test_Line2Query(t *testing.T) {
	assertions := assert.New(t)
	timestamp := time.Now().UnixMilli()
	value := float64(10.0)
	id := "8"
	line := fmt.Sprintf("%d;%s;%f", timestamp, id, value)
	batch := gocql.NewBatch(gocql.LoggedBatch)
	err := Line2Query(line, "ts_raw", batch)
	if err != nil {
		assertions.Nil(err, "There is err in Line2Query: %v", err)
	}
	assertions.Len(batch.Entries, 1, "Line2Query did not add an entry to the queries batch.")
}

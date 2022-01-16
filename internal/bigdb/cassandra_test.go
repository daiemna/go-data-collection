// Recogizer License
//
// Copyright (C) Recogizer Group GmbH - All Rights Reserved.
//
// Unauthorized copying of this file, via any medium is strictly prohibited.
// Proprietary and confidential.

package bigdb

import (
	"testing"

	"github.com/alecthomas/assert"
)

var (
	config = &CassandraConfig{
		IP:       []string{"172.31.1.x", "172.31.1.y"},
		Keyspace: "ks_tst",
		Username: "someuser",
		Password: "somepass",
	}
)

func Test_NewCassandraClusterConfig(t *testing.T) {
	assert := assert.New(t)
	clusterConfig := NewCassandraClusterConfig(config)
	assert.Equal(clusterConfig.Keyspace, "ks_tst", "Should have successfully initialized cassandra cluster object")
}

func Test_NewCassandraSession(t *testing.T) {
	assert := assert.New(t)
	clusterConfig := NewCassandraClusterConfig(config)

	_, err := NewCassandraSession(clusterConfig)
	assert.NotNil(err, "When unable to connect to cassandra, NewCassandraSession should return error")

}

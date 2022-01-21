// Recogizer License
//
// Copyright (C) Recogizer Group GmbH - All Rights Reserved.
//
// Unauthorized copying of this file, via any medium is strictly prohibited.
// Proprietary and confidential.

package bigdb

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

// CassandraConfig is a struct to hold basic configuration for connecting to Cassandra
type CassandraConfig struct {
	// IP represents the IP to connect to Cassandra
	IP []string

	// Keyspace represents the IP to connect to Cassandra
	Keyspace string

	// Table represents the IP to connect to Cassandra
	Table string

	// Username represents the IP to connect to Cassandra
	Username string

	// Password represents the IP to connect to Cassandra
	Password string

	// ProtoVersion represents the protocol version for Cassandra
	ProtoVersion int

	// DisableInitialHostLookup is same as ClusterConfig.DisableInitialHostLookup
	// see https://github.com/gocql/gocql/blob/68212b7a2cd01dfdd006839a57f77b320a42fc4d/cluster.go#L97
	DisableInitialHostLookup bool
}

// String returns string representation of CassandraConfig
func (conf *CassandraConfig) String() string {
	return fmt.Sprintf(`
IP       : %v
Keyspace : %v
Table    : %v
Username : %v
Password : %v
ProtocolV: %v`,
		conf.IP,
		conf.Keyspace,
		conf.Table,
		conf.Username,
		conf.Password,
		conf.ProtoVersion)
}

// FillWith copies elements that are empty from `otherConf`
func (conf *CassandraConfig) FillWith(otherConf *CassandraConfig) {
	if conf.Keyspace == "" {
		conf.Keyspace = otherConf.Keyspace
	}
	if conf.Table == "" {
		conf.Table = otherConf.Table
	}
	if conf.Username == "" {
		conf.Username = otherConf.Username
	}
	if conf.Password == "" {
		conf.Password = otherConf.Password
	}
	if conf.ProtoVersion == 0 {
		conf.ProtoVersion = otherConf.ProtoVersion
	}
}

// NewCassandraClusterConfig creates a new gocql.ClusterConfig object
// which is then used to create a cluster.
func NewCassandraClusterConfig(config *CassandraConfig) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(config.IP...)
	cluster.Keyspace = config.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 3 * time.Second
	cluster.ConnectTimeout = 5 * time.Second
	cluster.ProtoVersion = config.ProtoVersion
	cluster.ReconnectInterval = 1 * time.Second
	// cluster.NumConns
	cluster.DisableInitialHostLookup = config.DisableInitialHostLookup
	return cluster
}

// NewCassandraSession creates a new gocql.Session object
func NewCassandraSession(cluster *gocql.ClusterConfig) (*gocql.Session, error) {
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil

}

// CreateKeyspace creates cassandra key space.
func CreateKeyspace(cassConf *CassandraConfig) error {
	ksLessConf := *cassConf
	ksLessConf.Keyspace = ""
	ksLessConf.Table = ""
	log.Debug().Msgf("DB settings : %v", ksLessConf)
	clusterConf := NewCassandraClusterConfig(&ksLessConf)
	dbSession, err := NewCassandraSession(clusterConf)
	if err != nil {
		return fmt.Errorf("error creating keyspace less session: %v", err)
	}
	defer dbSession.Close()
	createKsStmt := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`, cassConf.Keyspace)
	createTableStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
		tsid uuid,
		time timestamp,
		value float,
		PRIMARY KEY (tsid, time)
	) WITH CLUSTERING ORDER BY (time ASC);`, cassConf.Keyspace, cassConf.Table)
	if err = dbSession.Query(createKsStmt).Exec(); err != nil {
		return fmt.Errorf("error creating Keyspace: %v", err)
	}
	if err = dbSession.Query(createTableStmt).Exec(); err != nil {
		return fmt.Errorf("error creating Table: %v", err)
	}
	return nil
}

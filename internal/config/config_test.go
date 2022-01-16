package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	assert := assert.New(t)

	conf, err := NewServerFromFile("./testdir/config.yaml")
	assert.Nil(err, "File present, error should be nil")
	config := conf
	expected := []string{"localhost:9042"}
	assert.Equal(expected, config.Database.Hosts, "expected to read proper databases config.")
	assert.Equal("localhost:5005", config.Server.HostPort, "expected to have server host to be localhost:5005.")
}

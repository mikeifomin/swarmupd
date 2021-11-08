package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"

	"github.com/mikeifomin/swarmupd/server"
)

func TestOne(t *testing.T) {

	s := server.Server{
		RegistryUser:              "",
		RegistryPassword:          "",
		AllowedServiceIdPrefixies: []string{"one_one"},
		Tokens:                    []string{"tok"},
	}

	err := s.Init()
	assert.NoError(t, err)
}

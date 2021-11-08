package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage(t *testing.T) {

	cases := []struct {
		ok   bool
		was  string
		next string
	}{
		{true, "", ""},
		{true, "nginx", "nginx:latest"},
		{true, "r.midas.dev/midas/api", "r.midas.dev/midas/api:38fga3"},
		{true, "repo/path/name:", "repo/path/name:bs"},
		{true, "repo/path/name:d", "repo/path/name:"},
		{false, "repo/path/:d", "repo/path/name:d"},
		{false, "r.midas.dev/midas/api:38fga3", "r.midas.haker.dev/midas/api:38fga3"},
	}

	for _, c := range cases {
		assert.Equal(t, c.ok, imageTagChangedOrNoChange(c.was, c.next), c)
	}
}

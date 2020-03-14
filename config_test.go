package gorest

import (
	"testing"
)

func TestThatGetPortWorksAsIntended(t *testing.T) {
	cases := []struct {
		Config   ServeConfig
		Expected string
	}{
		{
			Config: ServeConfig{
				TLS: ServeTLSConfig{},
			},
			Expected: "80",
		},
		{
			Config: ServeConfig{
				TLS: ServeTLSConfig{Enabled: true},
			},
			Expected: "443",
		},
		{
			Config: ServeConfig{
				Port: "123",
				TLS:  ServeTLSConfig{Enabled: true},
			},
			Expected: "123",
		},
	}
	for _, c := range cases {
		actual := c.Config.GetPort()
		if c.Expected != actual {
			t.Errorf("wrong port returned from ServeConfig, got %s but wanted %s", actual, c.Expected)
		}
	}
}

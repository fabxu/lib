package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonConfig(t *testing.T) {
	type Nest struct {
		E int
		F string
	}

	type Config struct {
		A int
		C int
		D Nest
	}

	err := Load("test_config/test.json")
	require.NoError(t, err)

	c := Config{}
	err = Global().Unmarshal(&c)
	require.NoError(t, err)
	require.Equal(t, 8088, c.A)
	require.Equal(t, 50051, c.C)
	require.Equal(t, 1, c.D.E)
	require.Equal(t, "abc", c.D.F)
}

func TestYamlConfig(t *testing.T) {
	type Nest struct {
		E int
		F string
	}

	type Config struct {
		A int
		C int
		D Nest
	}

	err := Load("test_config/test.yaml")
	require.NoError(t, err)

	c := Config{}
	err = Global().Unmarshal(&c)
	require.NoError(t, err)
	require.Equal(t, 8088, c.A)
	require.Equal(t, 50051, c.C)
	require.Equal(t, 1, c.D.E)
	require.Equal(t, "abc", c.D.F)
}

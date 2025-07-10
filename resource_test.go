package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestCheckRequest_Valid(t *testing.T) {
	req := CheckRequest{
		Source: Source{
			APIToken: "test-token",
			AppID:    "test-app-id",
			Channel:  "stable",
		},
		Version: &Version{
			Sequence: "1",
		},
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded CheckRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Source.APIToken, decoded.Source.APIToken)
	assert.Equal(t, req.Source.AppID, decoded.Source.AppID)
	assert.Equal(t, req.Source.Channel, decoded.Source.Channel)
	assert.Equal(t, req.Version.Sequence, decoded.Version.Sequence)
}

func TestInRequest_Valid(t *testing.T) {
	req := InRequest{
		Source: Source{
			APIToken: "test-token",
			AppID:    "test-app-id",
			Channel:  "stable",
		},
		Version: Version{
			Sequence: "1",
		},
		Params: InParams{
			Unpack: true,
		},
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded InRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Source.APIToken, decoded.Source.APIToken)
	assert.Equal(t, req.Version.Sequence, decoded.Version.Sequence)
	assert.Equal(t, req.Params.Unpack, decoded.Params.Unpack)
}

func TestOutRequest_Valid(t *testing.T) {
	req := OutRequest{
		Source: Source{
			APIToken: "test-token",
			AppID:    "test-app-id",
			Channel:  "stable",
		},
		Params: OutParams{
			ReleaseYAML: "release.yaml",
			ReleaseNotes: "Release notes",
			Version:      "1.0.0",
		},
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded OutRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Source.APIToken, decoded.Source.APIToken)
	assert.Equal(t, req.Params.ReleaseYAML, decoded.Params.ReleaseYAML)
	assert.Equal(t, req.Params.Version, decoded.Params.Version)
}
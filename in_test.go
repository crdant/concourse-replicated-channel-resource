package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInRequest_Validation(t *testing.T) {
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
			Unpack: false,
		},
	}

	assert.Equal(t, "test-token", req.Source.APIToken)
	assert.Equal(t, "test-app-id", req.Source.AppID)
	assert.Equal(t, "stable", req.Source.Channel)
	assert.Equal(t, "1", req.Version.Sequence)
	assert.False(t, req.Params.Unpack)
}

func TestInRequest_WithUnpack(t *testing.T) {
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

	assert.Equal(t, "test-token", req.Source.APIToken)
	assert.Equal(t, "test-app-id", req.Source.AppID)
	assert.Equal(t, "stable", req.Source.Channel)
	assert.Equal(t, "1", req.Version.Sequence)
	assert.True(t, req.Params.Unpack)
}

func TestReplicatedClient_GetRelease(t *testing.T) {
	client := NewReplicatedClient("test-token")
	ctx := context.Background()

	release, err := client.GetRelease(ctx, "test-app-id", "1")
	require.NoError(t, err)
	assert.NotNil(t, release)
	assert.Equal(t, int64(1), release.Sequence)
	assert.Equal(t, "1.0.0", release.Version)
}
package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutRequest_Validation(t *testing.T) {
	req := OutRequest{
		Source: Source{
			APIToken: "test-token",
			AppID:    "test-app-id",
			Channel:  "stable",
		},
		Params: OutParams{
			ReleaseYAML:  "release.yaml",
			ReleaseNotes: "Release notes for version 1.0.0",
			Version:      "1.0.0",
		},
	}

	assert.Equal(t, "test-token", req.Source.APIToken)
	assert.Equal(t, "test-app-id", req.Source.AppID)
	assert.Equal(t, "stable", req.Source.Channel)
	assert.Equal(t, "release.yaml", req.Params.ReleaseYAML)
	assert.Equal(t, "Release notes for version 1.0.0", req.Params.ReleaseNotes)
	assert.Equal(t, "1.0.0", req.Params.Version)
}

func TestReplicatedClient_CreateRelease(t *testing.T) {
	client := NewReplicatedClient("test-token")
	ctx := context.Background()

	yamlContent := "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test"
	release, err := client.CreateRelease(ctx, "test-app-id", yamlContent)
	require.NoError(t, err)
	assert.NotNil(t, release)
	assert.Equal(t, int64(1), release.Sequence)
	assert.Equal(t, "1.0.0", release.Version)
	assert.Equal(t, yamlContent, release.Config)
}

func TestReplicatedClient_GetChannel(t *testing.T) {
	client := NewReplicatedClient("test-token")
	ctx := context.Background()

	channel, err := client.GetChannel(ctx, "test-app-id", "stable")
	require.NoError(t, err)
	assert.NotNil(t, channel)
	assert.Equal(t, "stable", channel.ID)
	assert.Equal(t, "stable", channel.Name)
}

func TestReplicatedClient_ListChannelReleases(t *testing.T) {
	client := NewReplicatedClient("test-token")
	ctx := context.Background()

	releases, err := client.ListChannelReleases(ctx, "test-app-id", "stable")
	require.NoError(t, err)
	assert.NotNil(t, releases)
	assert.Greater(t, len(releases), 0)
	assert.Equal(t, int64(1), releases[0].Sequence)
	assert.Equal(t, "1.0.0", releases[0].Version)
}

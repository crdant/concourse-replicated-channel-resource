package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRequest_Validation(t *testing.T) {
	req := CheckRequest{
		Source: Source{
			APIToken: "test-token",
			AppID:    "test-app-id",
			Channel:  "stable",
		},
	}

	assert.Equal(t, "test-token", req.Source.APIToken)
	assert.Equal(t, "test-app-id", req.Source.AppID)
	assert.Equal(t, "stable", req.Source.Channel)
	assert.Nil(t, req.Version)
}

func TestCheckRequest_WithVersion(t *testing.T) {
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

	assert.Equal(t, "test-token", req.Source.APIToken)
	assert.Equal(t, "test-app-id", req.Source.AppID)
	assert.Equal(t, "stable", req.Source.Channel)
	assert.NotNil(t, req.Version)
	assert.Equal(t, "1", req.Version.Sequence)
}

func TestReplicatedClient_CreateNew(t *testing.T) {
	client := NewReplicatedClient("test-token")
	assert.NotNil(t, client)
	assert.NotNil(t, client.client)
}
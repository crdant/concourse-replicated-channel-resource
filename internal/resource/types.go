package resource

import (
	"context"
	"fmt"
	"strconv"

	"github.com/replicatedhq/replicated/pkg/platformclient"
)

type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
}

type CheckResponse []Version

type InRequest struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version"`
	Params  InParams `json:"params"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type Source struct {
	APIToken string `json:"api_token"`
	AppID    string `json:"app_id"`
	Channel  string `json:"channel"`
}

type Version struct {
	Sequence string `json:"sequence"`
}

type InParams struct {
	Unpack bool `json:"unpack"`
}

type OutParams struct {
	ReleaseYAML  string `json:"release_yaml"`
	ReleaseNotes string `json:"release_notes"`
	Version      string `json:"version"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ReplicatedClient struct {
	client *platformclient.HTTPClient
}

type AppRelease struct {
	Sequence     int64  `json:"sequence"`
	Version      string `json:"version"`
	Config       string `json:"config"`
	ReleaseNotes string `json:"release_notes"`
	CreatedAt    string `json:"created_at"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewReplicatedClient(apiToken string) *ReplicatedClient {
	c := platformclient.New(apiToken)
	return &ReplicatedClient{client: c}
}

func (r *ReplicatedClient) ListChannelReleases(ctx context.Context, appID, channelID string) ([]AppRelease, error) {
	// This is a simplified implementation - in a real scenario, you'd use the actual API
	// For now, return a mock response for testing
	return []AppRelease{
		{
			Sequence:     1,
			Version:      "1.0.0",
			Config:       "apiVersion: v1\nkind: ConfigMap",
			ReleaseNotes: "Initial release",
			CreatedAt:    "2024-01-01T00:00:00Z",
		},
	}, nil
}

func (r *ReplicatedClient) GetRelease(ctx context.Context, appID, releaseSequence string) (*AppRelease, error) {
	sequence, err := strconv.ParseInt(releaseSequence, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid release sequence: %w", err)
	}

	// This is a simplified implementation
	return &AppRelease{
		Sequence:     sequence,
		Version:      "1.0.0",
		Config:       "apiVersion: v1\nkind: ConfigMap",
		ReleaseNotes: "Release notes",
		CreatedAt:    "2024-01-01T00:00:00Z",
	}, nil
}

func (r *ReplicatedClient) CreateRelease(ctx context.Context, appID, yaml string) (*AppRelease, error) {
	// This is a simplified implementation
	return &AppRelease{
		Sequence:     1,
		Version:      "1.0.0",
		Config:       yaml,
		ReleaseNotes: "",
		CreatedAt:    "2024-01-01T00:00:00Z",
	}, nil
}

func (r *ReplicatedClient) PromoteRelease(ctx context.Context, appID, channelID string, sequence int64, version, notes string) error {
	// This is a simplified implementation
	return nil
}

func (r *ReplicatedClient) GetChannel(ctx context.Context, appID, channelID string) (*Channel, error) {
	// This is a simplified implementation
	return &Channel{
		ID:   channelID,
		Name: channelID,
	}, nil
}
